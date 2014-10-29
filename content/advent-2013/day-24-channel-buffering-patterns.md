+++
title = "Go Advent Day 24 - Channel Buffering Patterns"
date = 2013-12-24T06:40:42Z
author = ["Caleb Spare"]
series = ["Advent 2013"]
+++

## Introduction

One common method for message processing in Go is to receive from an input channel and write to an output
channel, often using intermediate channels for message transformation and filtering. Multiple goroutines may
perform these functions concurrently and independently, making for code that is easily parallelized and
tested.

Message buffering is one kind of transformation that is sometimes useful in these systems. Some programs don't
need to process each message immediately, and can more efficiently process several messages at once. Other
programs receive bursty input but are able to coalesce groups of messages.

An example of the latter case is software built with an inotify/kqueue library such as
[github.com/howeyc/fsnotify](https://github.com/howeyc/fsnotify). This library provides the user with a
channel of file events. Messages from fsnotify often come in bursts; for instance, some text editors save
files by writing out a temporary file and moving it into place -- a create event, one or more modifications,
and a rename.

I observed such behavior when writing [github.com/cespare/reflex](https://github.com/cespare/reflex), a tool
for running commands when files change, and the code below is similar to
[real code in that program](https://github.com/cespare/reflex/blob/master/workers.go).

In this post I will explore how to handle bursty messages in channel pipelines.

## Producing messages

We'll represent messages with `Event` s, which are just `int` s. Multiple `Event` s are allowed to be merged
into a single event using addition. (In an fsnotify-based program, for example, you might be able to
coalesce multiple events by preserving only the unique file names, or only the fact that some file was changed
at all.) Our producer intermittently emits a burst of random-valued `Event` s.

	type Event int

	func NewEvent() Event                   { return 0 }
	func (e Event) Merge(other Event) Event { return e + other }

	func produce(out chan<- Event) {
		for {
			delay := time.Duration(rand.Intn(5)+1) * time.Second
			nMessages := rand.Intn(10) + 1
			time.Sleep(delay)
			for i := 0; i < nMessages; i++ {
				out <- Event(rand.Intn(10))
			}
		}
	}

If we print out all the messages created by `produce`, we can see that they come out in small batches.
[See the runnable playground example](http://play.golang.org/p/aHsHNFIfGP) (you may have to scroll down to
see the output).

## Batching by time

Let's suppose we are willing to wait for a short time after receiving a message to see if other messages
arrive that we can coalesce with the first. One simple way of doing this is simply to send every N
milliseconds and merge all received messages in the meantime. This is easy with a
[`time.Ticker`](http://golang.org/pkg/time/#Ticker):

	func coalesce500(in <-chan Event, out chan<- Event) {
		ticker := time.NewTicker(500 * time.Millisecond)
		event := NewEvent()
		for {
			select {
			case <-ticker.C:
				out <- event
				event = NewEvent()
			case e := <-in:
				event = event.Merge(e)
			}
		}
	}

[See the full runnable code here](http://play.golang.org/p/yelQeZiyly). This works, but it has a few
problems:

- If no messages are arriving, every 500ms `coalesce500` sends an empty message.
- If a message arrives shortly before the ticker fires, we'll send it by itself without waiting for more messages. (More generally, a burst falling across the ticker will create two messages, rather than one.)
- If the receiver doesn't handle the message immediately, the coalescing goroutine will be blocked, which in turns blocks the upstream producer. (In real code, that producer may be constructed to discard any messages it cannot immediately send.)

Let's tackle the third problem first.

## Handling slow receivers

One obvious way to avoid blocking on a slow receiver is to change the output channel to be buffered. This only
helps until the buffer fills up, though, and furthermore we won't merge multiple buffered messages,
potentially causing a downstream burst as well. A better way of handling the slow receiver is to implement the
desired semantics directly: keep coalescing messages until the receiver is ready to accept more input.

	func coalesceSlow(in <-chan Event, out chan<- Event) {
		event := NewEvent()
		for e := range in {
			event = event.Merge(e)
		loop:
			for {
				select {
				case e := <-in:
					event = event.Merge(e)
				case out <- event:
					event = NewEvent()
					break loop
				}
			}
		}
	}

You can [run the code with a slow receiver](http://play.golang.org/p/bYZpWB974A) and see that this has the
desired effect: we only create a single coalesced message while the receiver blocks. `coalesceSlow` still
suffers from the problem as `coalesce500` that it does not wait for a follow-up burst after receiving a
message. Further, it introduces a bit of ugly code: we receive from `in` in two different places (the range
and inside the select) and we have to bounce between the inner and outer loops depending on whether we have
anything to send.

## Putting it together

To review, we would like for our final version to do the following:

- Wait for 500ms after the receipt of a message to send
- Wait until the receiver is ready before sending
- Keep coalescing incoming messages while waiting
- Only have single channel receive and send statements instead of duplicates

Our state machine will need to be a little more complex to meet these goals. To solve the empty message
problem from `coalesce500`, we'll start a timer running only after receiving an event. Then we'll coalesce
messages until the timer fires. Finally, we'll continue coalescing messages until the output channel is ready,
after which we can send and go back to the first state again.

Here is a diagram to show these transitions:

![](/postimages/day-24-channel-buffering-patterns/channel-buffering.png)

Here is the code:

	func coalesce(in <-chan Event, out chan<- Event) {
		event := NewEvent()
		timer := time.NewTimer(0)

		var timerCh <-chan time.Time
		var outCh chan<- Event

		for {
			select {
			case e := <-in:
				event = event.Merge(e)
				if timerCh == nil {
					timer.Reset(500 * time.Millisecond)
					timerCh = timer.C
				}
			case <-timerCh:
				outCh = out
				timerCh = nil
			case outCh <- event:
				event = NewEvent()
				outCh = nil
			}
		}
	}

There are a few techniques used here that deserve explanation. Most notably, we're exploiting the fact that
any cases inside a select which receive from or send to a nil channel are skipped (read more about this
in the [language spec](http://golang.org/ref/spec#Select_statements)). We use this in the second and third
cases: `timerCh` and `outCh` are toggled between nil and the `timer.C` / `out` channels to selectively
remove them from the select.

In the initial state (labeled "ready" in the graph), both `timerCh` and `outCh` are nil, so we simply block on
`in`, waiting for input. After we receive a message, the timer is set for 500ms in the future and `timerCh` is
set to the timer's output channel; further messages are coalesced while we wait for the timer to fire (the
second state, "waiting"). Once this happens, we toggle `outCh` to `out`, so that we can detect when the output
channel is ready to receive (this is the last state, "ready to send"). Once we have successfully sent, we
return to the initial state, `outCh` and `timerCh` having previously been reset to nil.

The `timer` variable initialized with `time.NewTimer(0)` fires immediately, but we don't start listening to
its output channel until after we have `Reset` it. (There is no way to create a `time.Timer` without
scheduling it.)

If you [run this on the playground](http://play.golang.org/p/r4nH9M38Jt) you can verify that it behaves as
advertised.

## The end

If you have questions, you can find me on [Google+](https://plus.google.com/+CalebSpare/) or
[Twitter](https://twitter.com/kingfishr). I am also cespare on IRC; if you haven't checked it out, we have a
great community at #go-nuts on Freenode.

Special thanks to Tommi Virtanen and Dominik Honnef for the excellent diagram and valuable feedback as I was
editing this post.
