+++
title = "Go Advent Day 22 - A Journey Into NSQ"
date = 2013-12-22T06:40:42Z
author = ["Matt Reiferson"]
series = ["Advent 2013"]
+++

## Introduction

([NSQ](http://bitly.github.io/nsq) is a realtime distributed messaging platform. It's designed
to serve as the backbone of a modern infrastructure composed of loosely connected services running
on many computers.

This post describes the internal architecture of NSQ with an emphasis on Go, focusing on
performance optimization, stability, and robustness for high throughput network servers.

Arguably NSQ would not exist if it were not for the timing of our adoption of Go at
[bitly](https://bitly.com/). There is a strong symmetry between the functionality of NSQ and
features provided in the language. Naturally,
[language influences thought](http://en.wikipedia.org/wiki/Linguistic_relativity) and this is no
exception.

In retrospect the choice to use Go has paid off tenfold. The excitement around the language and the
willingness of the community to provide feedback, support the project, and contribute patches has
been a tremendous help.

## Overview

NSQ is composed of 3 daemons:

- [nsqd](http://bitly.github.io/nsq/components/nsqd.html) is the daemon that receives, queues, and delivers messages to clients.
- [nsqlookupd](http://bitly.github.io/nsq/components/nsqlookupd.html) is the daemon that manages topology information and provides an eventually consistent discovery service.
- [nsqadmin](http://bitly.github.io/nsq/components/nsqadmin.html) is a web UI to introspect the cluster in realtime (and perform various administrative tasks).

Data flow in NSQ is modeled as a tree of _streams_ and _consumers_. A *topic* is a distinct
stream of data. A *channel* is a logical grouping of consumers subscribed to a given *topic*.

![](https://f.cloud.github.com/assets/187441/1700696/f1434dc8-6029-11e3-8a66-18ca4ea10aca.gif)

A single *nsqd* can have many topics and each topic can have many channels. A channel
receives a _copy_ of all the messages for the topic, enabling _multicast_ style delivery while each
message on a channel is _distributed_ amongst its subscribers, enabling load-balancing.

These primitives form a powerful framework for expressing a variety of [simple and complex topologies](http://bitly.github.io/nsq/deployment/topology_patterns.html).

For more information about the design of NSQ see the [design doc](http://bitly.github.io/nsq/overview/design.html).

## Topics and Channels

Topics and channels, the core primitives of NSQ, best exemplify how the design of the system
translates seamlessly to the features of Go.

Go's channels (henceforth referred to as "go-chan" for disambiguation) are a natural way to express
queues, thus an NSQ topic/channel, at its core, is just a _buffered_ go-chan of `Message` pointers.
The size of the buffer is equal to the `--mem-queue-size` configuration parameter.

After reading data off the wire, the act of publishing a message to a topic involves:

- instantiation of a `Message` struct (and allocation of the message body `[]byte`)
- read-lock to get the `Topic`
- read-lock to check for the ability to publish
- send on a buffered go-chan

To get messages from a topic to its channels the topic cannot rely on typical go-chan receive
semantics, because multiple goroutines receiving on a go-chan would _distribute_ the messages
while the desired end result is to _copy_ each message to every channel (goroutine).

Instead, each topic maintains 3 primary goroutines. The first one, called `router`, is responsible
for reading newly published messages off the incoming go-chan and storing them in a queue (memory
or disk).

The second one, called `messagePump`, is responsible for copying and pushing messages to channels as
described above.

The third is responsible for `DiskQueue` IO and will be discussed later.

Channels are a _little_ more complicated but share the underlying goal of exposing a _single_ input
and _single_ output go-chan (to abstract away the fact that, internally, messages might be in
memory or on disk):

![](https://f.cloud.github.com/assets/187441/1698990/682fc358-5f76-11e3-9b05-3d5baba67f13.png)


Additionally, each channel maintains 2 time-ordered priority queues responsible for deferred and
in-flight message timeouts (and 2 accompanying goroutines for monitoring them).

Parallelization is improved by managing a _per-channel_ data structure, rather than relying on the
Go runtime's _global_ timer scheduler.

*Note:* Internally, the Go runtime uses a single priority queue and goroutine to manage timers.
This supports (but is not limited to) the entirety of the `time` package. It normally obviates the
need for a _user-land_ time-ordered priority queue but it's important to keep in mind that it's a
_single_ data structure with a _single_ lock, potentially impacting `GOMAXPROCS`>`1` performance.
See [runtime/time.goc](http://golang.org/src/pkg/runtime/time.goc?s=1684:1787#L83).

## Backend / DiskQueue

One of NSQ's design goals is to bound the number of messages kept in memory. It does this by
transparently writing message overflow to disk via `DiskQueue` (which owns the _third_ primary
goroutine for a topic or channel).

Since the memory queue is just a go-chan, it's trivial to route messages to memory first, if
possible, then fallback to disk:

    for msg := range c.incomingMsgChan {
        select {
        case c.memoryMsgChan <- msg:
        default:
            err := WriteMessageToBackend(&msgBuf, msg, c.backend)
            if err != nil {
                // ... handle errors ...
            }
        }
    }


Taking advantage of Go's `select` statement allows this functionality to be expressed in just a few
lines of code: the `default` case above only executes if `memoryMsgChan` is full.

NSQ also has the concept of *ephemeral* channels. Ephemeral channels _discard_ message overflow
(rather than write to disk) and disappear when they no longer have clients subscribed.  This is
a perfect use case for Go's interfaces. Topics and channels have a struct member declared
as a `Backend` _interface_ rather than a concrete type. Normal topics and channels use a
`DiskQueue` while ephemeral channels stub in a `DummyBackendQueue`, which implements a no-op
`Backend`.

## Reducing GC Pressure

In any garbage collected environment you're subject to the tension between throughput (doing useful
work), latency (responsiveness), and resident set size (footprint).

As of Go 1.2, the GC is mark-and-sweep (parallel), non-generational, non-compacting, stop-the-world
and mostly precise . It's _mostly_ precise because the remainder of the work wasn't completed in
time (it's slated for Go 1.3).

The Go GC will certainly continue to improve, but the universal truth is: _the_less_garbage_you_create_the_less_time_you'll_collect_.

First, it's important to understand how the GC is behaving _under_real_workloads_. To this end,
*nsqd* publishes GC stats in [statsd](https://github.com/etsy/statsd/) format (alongside other internal metrics).
*nsqadmin* displays graphs of these metrics, giving you insight into the GC's impact in both
frequency and duration:

![](https://f.cloud.github.com/assets/187441/1699828/8df666c6-5fc8-11e3-95e6-360b07d3609d.png)

In order to actually _reduce_ garbage you need to know where it's being generated. Once again the
Go toolchain provides the answers:

- Use the [testing](http://golang.org/pkg/testing/) package and `go`test`-benchmem` to benchmark hot code paths. It profiles the number of allocations per iteration (and benchmark runs can be compared with [benchcmp](http://golang.org/misc/benchcmp)).
- Build using `go`build`-gcflags`-m`, which outputs the result of [escape analysis](http://en.wikipedia.org/wiki/Escape_analysis).

With that in mind, the following optimizations proved useful for *nsqd*:

- Avoid `[]byte` to `string` conversions.
- Re-use buffers or objects (and someday possibly [sync.Pool](https://groups.google.com/forum/#!topic/golang-dev/kJ_R6vYVYHU) aka [issue 4720](https://code.google.com/p/go/issues/detail?id=4720)).
- Pre-allocate slices (specify capacity in `make`) and always know the number and size of items over the wire.
- Apply sane limits to various configurable dials (such as message size).
- Avoid boxing (use of `interface{}`) or unnecessary wrapper types (like a `struct` for a "multiple value" go-chan).
- Avoid the use of `defer` in hot code paths (it allocates).

### TCP Protocol

The [NSQ TCP protocol](http://bitly.github.io/nsq/clients/tcp_protocol_spec.html) is a shining example of a section where these GC optimization
concepts are utilized to great effect.

The protocol is structured with length prefixed frames, making it straightforward and performant to
encode and decode:

    [x][x][x][x][x][x][x][x][x][x][x][x]...
    |  (int32) ||  (int32) || (binary)
    |  4-byte  ||  4-byte  || N-byte
    ------------------------------------...
        size      frame ID     data

Since the exact type and size of a frame's components are known ahead of time, we can avoid the
[encoding/binary](http://golang.org/pkg/encoding/binary/) package's convenience [Read()](http://golang.org/pkg/encoding/binary/#Read) and
[Write()](http://golang.org/pkg/encoding/binary/#Write) wrappers (and their extraneous interface lookups and conversions) and
instead call the appropriate [binary.BigEndian](http://golang.org/pkg/encoding/binary/#ByteOrder) methods directly.

To reduce socket IO syscalls, client `net.Conn` are wrapped with [bufio.Reader](http://golang.org/pkg/bufio/#Reader) and
[bufio.Writer](http://golang.org/pkg/bufio/#Writer). The `Reader` exposes [ReadSlice()](http://golang.org/pkg/bufio/#Reader.ReadSlice), which reuses its
internal buffer. This nearly eliminates allocations while reading off the socket, greatly reducing
GC pressure. This is possible because the data associated with most commands does not escape (in
the edge cases where this is not true, the data is _explicitly_ copied).

At an even lower level, a `MessageID` is declared as `[16]byte` to be able to use it as a `map`
key (slices cannot be used as map keys). However, since data read from the socket is stored as
`[]byte`, rather than produce garbage by allocating `string` keys, and to avoid a copy from the
slice to the backing array of the `MessageID`, the `unsafe` package is used to cast the slice
directly to a `MessageID`:

    id := *(*nsq.MessageID)(unsafe.Pointer(&msgID))

*Note:* _This_is_a_hack_. It wouldn't be necessary if this was optimized by the compiler and
[Issue 3512](https://code.google.com/p/go/issues/detail?id=3512) is open to potentially resolve this. It's also worth reading through
[issue 5376](https://code.google.com/p/go/issues/detail?id=5376), which talks about the possibility of a "const like" `byte` type that
could be used interchangeably where `string` is accepted, _without_ allocating and copying.

Similarly, the Go standard library only provides numeric conversion methods on a `string`. In order
to avoid `string` allocations, *nsqd* uses a [custom base 10 conversion method](https://github.com/bitly/nsq/blob/master/util/byte_base10.go#L7-L27)
that operates directly on a `[]byte`.

These may seem like micro-optimizations but the TCP protocol contains some of the *hottest* code
paths. In aggregate, at the rate of tens of thousands of messages per second, they have a
significant impact on the number of allocations and overhead:

    benchmark                    old ns/op    new ns/op    delta
    BenchmarkProtocolV2Data           3575         1963  -45.09%

    benchmark                    old ns/op    new ns/op    delta
    BenchmarkProtocolV2Sub256        57964        14568  -74.87%
    BenchmarkProtocolV2Sub512        58212        16193  -72.18%
    BenchmarkProtocolV2Sub1k         58549        19490  -66.71%
    BenchmarkProtocolV2Sub2k         63430        27840  -56.11%

    benchmark                   old allocs   new allocs    delta
    BenchmarkProtocolV2Sub256           56           39  -30.36%
    BenchmarkProtocolV2Sub512           56           39  -30.36%
    BenchmarkProtocolV2Sub1k            56           39  -30.36%
    BenchmarkProtocolV2Sub2k            58           42  -27.59%

## HTTP

NSQ's HTTP API is built on top of Go's [net/http](http://golang.org/pkg/net/http/) package. Because it's _just_ HTTP, it
can be leveraged in almost any modern programming environment without special client libraries.

Its simplicity belies its power, as one of the most interesting aspects of Go's HTTP tool-chest
is the wide range of debugging capabilities it supports. The [net/http/pprof](http://golang.org/pkg/net/http/pprof)
package integrates directly with the native HTTP server, exposing endpoints to retrieve CPU, heap,
goroutine, and OS thread profiles. These can be targeted directly from the `go` tool:

    $ go tool pprof http://127.0.0.1:4151/debug/pprof/profile

This is a tremendously valuable for debugging and profiling a _running_ process!

In addition, a `/stats` endpoint returns a slew of metrics in either JSON or pretty-printed text,
making it easy for an administrator to introspect from the command line in realtime:

    $ watch -n 0.5 'curl -s http://127.0.0.1:4151/stats | grep -v connected'

This produces continuous output like:
 
![](https://f.cloud.github.com/assets/187441/1780159/cbbfc3da-6865-11e3-831a-7a6177b66e70.png)


Finally, Go 1.2 brought [measurable HTTP performance gains](https://github.com/davecheney/autobench/blob/master/linux-amd64-x220-go1.1.2-vs-go1.2.txt#L156-L181). It's always nice when
recompiling against the latest version of Go provides a free performance boost!

## Dependencies

Coming from other ecosystems, Go's philosophy (or lack thereof) on managing dependencies takes a
little time to get used to.

NSQ evolved from being a single giant repo, with _relative imports_ and little to no separation
between internal packages, to fully embracing the recommended best practices with respect to
structure and dependency management.

There are two main schools of thought:

- *Vendoring*: copy dependencies at the correct revision into your application's repo and modify your import paths to reference the local copy.
- *Virtual*Env*: list the revisions of dependencies you require and at build time, produce a pristine `GOPATH` environment containing those pinned dependencies.

*Note:* This really only applies to _binary_ packages as it doesn't make sense for an importable
package to make intermediate decisions as to which version of a dependency to use.

NSQ uses [godep](https://github.com/kr/godep) to provide support for (2) above.

It works by recording your dependencies in a [Godeps](https://github.com/bitly/nsq/blob/master/Godeps) file, which it later uses to
construct a `GOPATH` environment. To build, it wraps and executes the standard Go toolchain
inside that environment. The `Godeps` file is just JSON and can be edited by hand.

It even supports `go`get` like semantics. For example, to produce a reliable build of NSQ:

    $ godep get github.com/bitly/nsq/...

## Testing

Go provides solid built-in support for writing tests and benchmarks and, because Go makes
it so easy to model concurrent operations, it's trivial to stand up a full-fledged instance of
*nsqd* inside your test environment.

However, there was one aspect of the initial implementation that became problematic for testing:
global state. The most obvious offender was the use of a global variable that held the reference to
the instance of *nsqd* at runtime, i.e. `var`nsqd`*NSQd`.

Certain tests would inadvertently mask this global variable in their local scope by using
short-form variable assignment, i.e. `nsqd`:=`NewNSQd(...)`. This meant that the global reference
did not point to the instance that was currently running, breaking tests.

To resolve this, a `Context` struct is passed around that contains configuration metadata and a
reference to the parent *nsqd*. All references to global state were replaced with this local
`Context`, allowing children (topics, channels, protocol handlers, etc.) to safely access this data
and making it more reliable to test.

## Robustness

A system that isn't robust in the face of changing network conditions or unexpected events is a
system that will not perform well in a distributed production environment.

NSQ is designed and implemented in a way that allows the system to tolerate failure and behave in a
consistent, predictable, and unsurprising way.

The overarching philosophy is to fail fast, treat errors as fatal, and provide a means to debug any
issues that do occur.

But, in order to _react_ you need to be able to _detect_ exceptional conditions...

### Heartbeats and Timeouts

The NSQ TCP protocol is push oriented. After connection, handshake, and subscription the consumer
is placed in a `RDY` state of `0`. When the consumer is ready to receive messages it updates that
`RDY` state to the number of messages it is willing to accept. NSQ client libraries continually
manage this behind the scenes, resulting in a flow-controlled stream of messages.

Periodically, *nsqd* will send a heartbeat over the connection. The client can configure the
interval between heartbeats but *nsqd* expects a response before it sends the next one.

The combination of application level heartbeats and `RDY` state avoids [head-of-line blocking](http://en.wikipedia.org/wiki/Head-of-line_blocking), which can otherwise render heartbeats useless (i.e. if a consumer is
behind in processing message flow the OS's receive buffer will fill up, blocking heartbeats).

To guarantee progress, all network IO is bound with deadlines relative to the configured heartbeat
interval. This means that you can literally unplug the network connection between *nsqd* and a
consumer and it will detect and properly handle the error.

When a fatal error is detected the client connection is forcibly closed. In-flight messages are
timed out and re-queued for delivery to another consumer. Finally, the error is logged and various
internal metrics are incremented.

### Managing Goroutines

It's surprisingly easy to _start_ goroutines. Unfortunately, it isn't quite as easy to orchestrate
their cleanup. Avoiding deadlocks is also challenging. Most often this boils down to an ordering
problem, where a goroutine receiving on a go-chan exits _before_ the upstream goroutines sending on
it.

Why care at all though? It's simple, an orphaned goroutine is a _memory_leak_.

To further complicate things, a typical *nsqd* process has _many_  active goroutines. Internally,
message "ownership" changes often. To be able to shutdown cleanly, it's incredibly important to
account for all _intraprocess_ messages.

Although there aren't any magic bullets, the following techniques make it a little easier to
manage...

### WaitGroups

The [sync](http://golang.org/pkg/sync/) package provides [sync.WaitGroup](http://golang.org/pkg/sync/#WaitGroup), which can be used to
perform accounting of how many goroutines are live (and provide a means to wait on their exit).

To reduce the typical boilerplate, *nsqd* uses this wrapper:

    type WaitGroupWrapper struct {
        sync.WaitGroup
    }

    func (w *WaitGroupWrapper) Wrap(cb func()) {
        w.Add(1)
        go func() {
            cb()
            w.Done()
        }()
    }

    // can be used as follows:
    wg := WaitGroupWrapper{}
    wg.Wrap(func() { n.idPump() })
    // ...
    wg.Wait()


### Exit Signaling

The easiest way to trigger an event in multiple child goroutines is to provide a single go-chan
that you close when ready. All pending receives on that go-chan will activate, rather than having
to send a separate signal to each goroutine.

    func work() {
        exitChan := make(chan int)
        go task1(exitChan)
        go task2(exitChan)
        time.Sleep(5 * time.Second)
        close(exitChan)
    }

    func task1(exitChan chan int) {
        <-exitChan
        log.Printf("task1 exiting")
    }

    func task2(exitChan chan int) {
        <-exitChan
        log.Printf("task2 exiting")
    }


### Synchronizing Exit

It was quite difficult to implement a reliable, deadlock free, exit path that accounted for all
in-flight messages. A few tips:

- Ideally the goroutine responsible for sending on a go-chan should also be responsible for closing it.
- If messages cannot be lost, ensure that pertinent go-chans are emptied (especially unbuffered ones!) to guarantee senders can make progress.
- Alternatively, if a message is no longer relevant, sends on a single go-chan should be converted to a `select` with the addition of an exit signal (as discussed above) to guarantee progress.

The general order should be:

- Stop accepting new connections (close listeners)
- Signal exit to child goroutines (see above)
- Wait on `WaitGroup` for goroutine exit (see above)
- Recover buffered data
- Flush anything left to disk

### Logging

Finally, the most important tool at your disposal is to _log_the_entrance_and_exit_of_your_goroutines!_. 
It makes it _infinitely_ easier to identify the culprit in the case of deadlocks or leaks.

*nsqd* log lines include information to correlate goroutines with their siblings (and parent),
such as the client's remote address or the topic/channel name.

The logs are verbose, but not verbose to the point where the log is overwhelming. There's a fine
line, but *nsqd* leans towards the side of having _more_ information in the logs when a fault
occurs rather than trying to reduce chattiness at the expense of usefulness.

## Summary

Our journey has come to an end...

If you have any questions don't hesitate to reach out on twitter [@imsnakes](https://twitter.com/imsnakes).

Special thanks to [@mccutchen](https://twitter.com/mccutchen), [@danielhfrank](https://twitter.com/danielhfrank),
[@ploxiln](https://twitter.com/ploxiln) and [@elubow](https://twitter.com/elubow) for reviewing this post.

Finally, shoutout to [@jehiah](https://twitter.com/jehiah), co-author of NSQ.
