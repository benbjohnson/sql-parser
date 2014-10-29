+++
title = "Go Advent Day 20 - Go in Academia: Emulating Wireless Networks"
date = 2013-12-20T06:40:42Z
author = ["Song Gao"]
series = ["Advent 2013"]
+++


## How I Came to This

_TL;DR;_ — We didn't want to use simulators and using real device for experiments became infeasible, so I decided to build an emulator for 802.11-like networks. The first version was in C++ and it didn't work well. I rewrote the simulator in Go and it gives twice throughput than C++ version and that's what we are using now.

I am in a wireless networking research group in Auburn University. We were working on a project in which we run [OLSR](https://en.wikipedia.org/wiki/Optimized_Link_State_Routing_Protocol) (an mobile Ad-hoc routing daemon) on mobile devices, and used data dissemination algorithms to periodically broadcast data into a group of 12 nodes.

To measure how OLSR performs in real world environment, and to benchmark our algorithms, we did various tests on laptops and Android devices. The most ridiculous-looking one was conducted in the basement of our department building, where 8 of us walked in a L shaped hallway back-and-forth, each holding a laptop, with the other 4 placed statically in corners.

We did these real device based experiments because we were doing systematic engineering and we needed to study behaviors of implementations that people would be using in real world. Simulations wouldn't work because everything has to be written in the simulator's script language (Tcl in case of ns-2), in the context of the simulation framework, and that means the protocol running in the simulator is not the one people can run on a laptop or a mobile device.

You might wonder why simulation can't be made as close to real-world programs as possible. Just look at the [OLSR module for ns-2](http://um-olsr.cvs.sourceforge.net/viewvc/um-olsr/um-olsr/) and the de facto standard implementation [OLSRd](http://olsr.org/git/). Most of advances happen in `OLSRd` instead of the ns-2 module. In addition, applications and OS kernel can make a big difference between real device and simulators.

But real device based experiments are apparently not scalable, especially when we want to consider mobility of wireless nodes. So I started to build an emulator. I was never great at C++ but for some reason I naturally picked C++ for this. I used `boost::asio` for asynchronous networking. It ended up with a lot of crazy callbacks written in weird C++11 syntax, and there was a race that I couldn't think of a way to solve efficiently.

So I looked into Go, and rewrote everything I had written in C++ in Go. I took half the time, wrote half number of lines of code, but the maximum throughput it could handle was twice as much as the C++ version! It worked and it worked better. That was how I started my first Go project, Squirrel Land ([squirrel-land/squirrels](https://github.com/squirrel-land/squirrels))

## So What Is Squirrel Land Exactly?

The goal of Squirrel Land is to emulate wireless network behaviours on Ethernet infrastructure. In other words, networking applications running on different Linux hosts connected to each other through Ethernet (physical or virtual) should see similar networking performance as in a mobile wireless network, e.g. 802.11 Ad-hoc network.

Conceptually, we want to replace the physical layer and part of link layer in TCP/IP stack with software that we have control on. See figure below:

    <div class="row">
        <img class="col-sm-offset-2 col-sm-8" src="day-20-squirrel/concept.png"></img>
    </div>


Squirrel Land produces two executables, `squirrel-worker` and `squirrel-master`. `squirrel-worker` runs on each host where a network application is being tested. Let's call such hosts worker hosts. `squirrel-master` runs only one instance on a special host that has powerful processing power and is reachable by all worker hosts.

When executed, `squirrel-worker` creates a [TAP](https://www.kernel.org/doc/Documentation/networking/tuntap.txt) interface, e.g. `tap0`, on that host, assigns an IP address to the interface and sets up routing for it. Since TAP interfaces work in link layer, `squirrel-worker` can capture every link layer frame that any process on the host sends through `tap0`. Upon getting a frame, `squirrel-worker` forwards it to the `squirrel-master`.

`squirrel-master` can apply an interference model, and based on virtual locations of each emulated "mobile node", decides whether the frame should be deliverable or not. If deliverable, it forwards it to the corresponding `squirrel-worker`, the second `squirrel-worker` writes the frame into its TAP interface, and the OS would parse the frame and deliver the data to corresponding application.

Following figure is an example of how a piece of data is sent from `OLSRd` process in one node to `OLSRd` process in another node:

    <div class="row">
        <img class="col-sm-offset-2 col-sm-8" src="day-20-squirrel/data_flow.png"></img>
    </div>


`squirrel-worker` is really just a relay between TAP device and a TCP stream to `squirrel-master`. On the other hand, `squirrel-master` can easily be the bottle neck, because every single frame in the entire system needs to be handled by `squirrel-master`.

## Implementation Details

### Re-using Buffers
GC is great but when you have a new 1500 bytes large slice allocated every 40 microseconds (assuming of 300Mbps traffic with 1500 MTU), there's a lot of pressure on GC.

To reduce this overhead, I use a circular buffer ([buffered.go](https://github.com/squirrel-land/squirrels/blob/dev/common/buffered.go), [link.go](https://github.com/squirrel-land/squirrels/blob/dev/common/link.go)) for all frames being handled on `squirrel-master`. To be more specific, there's a `chan` called `owner` buffering all available byte slices. Whenever a frame comes in, a byte slice is taken from `owner` and used to hold the frame. When the frame is no longer useful, it's returned into the `owner`.

Here's simplified code that illustrates it:

    type Token struct {
        Data  []byte
        owner chan<- *Token
    }

    // done with the byte slice; returning it for re-use
    func (t *Token) Return() {
        t.owner <- t
    }

    type BufferManager struct {
        buffer chan *Token
    }

    func NewBufferManager(size int) *BufferManager {
        buffer := make(chan *Token, size)
        ret := &BufferManager{buffer}
        for i := 0; i < size; i++ {
            buffer <- &Token{Data: []byte{}, owner: buffer}
        }
        return ret
    }

    // get a re-usable byte slice
    func (b *BufferManager) GetToken() *Token {
        return <-b.buffer
    }

I have a separate repo showing a simple benchmark about this technique. You are welcomed to [check out the results](https://github.com/songgao/bufferManager.go). Also, CloudFlare explaines how they recycle memory buffers in Go in [an awesome blog post](http://blog.cloudflare.com/recycling-memory-buffers-in-go).

### Plugin System

Go does not have dynamic linking, instead everything is built into a single binary, which is a great because it's much easier to deploy. However, this removes the ability to dynamically load modules.

We needed a way to allow other researchers to contribute modules into Squirrel Land. These modules include ones that describe mobility patterns of mobile nodes, and ones that models wireless transmission properties, such as interference on wireless links and [DCF](https://en.wikipedia.org/wiki/Distributed_coordination_function) in 802.11.

I ended up creating a separate repo ([squirrel-land/models](https://github.com/squirrel-land/models/) for all such modules. Each module is in a directory and is built into a Go package. Then in [`constructors.go`](https://github.com/squirrel-land/models/), constructors of each module is mapped to a string that can be used in configuration file. In this way, new modules could be easily integrated into Squirrel Land without much coupling with `squirrel-land/squirrels`.

I'm not very proud of this approach. If you have a better idea, please let me know :)

### TUN/TAP Driver

[`tuntap`](https://code.google.com/p/tuntap/) was the only one I found. It didn't quite meet my needs because all I wanted was a simple way to read data from TUN/TAP into my existing byte slices. In fact I just wanted something like `os.File`. So I implemented another one, [songgao/water]([https://github.com/songgao/water/).

`tuntap` was a really good example to start from. There's not much to talk about here and I just want to say I was impressed by Go's way of handling system calls. It feels so close to C, but it's not done through a C binding layer like CGO. There's a lot of great stuff in [`syscall`](http://golang.org/pkg/syscall/) package. If you are thinking about building some OS-related stuff, it's definitely worth looking into.

## Try It!

Squirrel Land only works on Linux for now. If you've got 20 minutes and have docker on your Linux machine, let's try something here.

### Get the Squirrel Land master

    # On Linux host
    go get github.com/squirrel-land/squirrels/squirrel-master

This would clone `squirrel-land/squirrels` into your $GOPATH, along with its dependency `squirrel-land/models`. But it only builds and installs `squirrel-master`, which is the binary that you would run on your Linux host.

### Build `squirrel-worker` and (optionally) make it available through HTTP

    # On Linux host
    cd $GOPATH/src/github.com/squirrel-land/squirrels/squirrel-worker
    go build

    # following is optional; you can use whatever way you feel comfortable
    # with to transfer file into containers
    mkdir /tmp/blahblah
    cp squirrel-worker /tmp/blahblah/squirrel-worker
    cd /tmp/blahblah
    python2.7 -m SimpleHTTPServer 9999

### Create LXC containers and get `squirrel-worker`

    docker run -i -t -name="squirrel_test_1" ubuntu /bin/bash
    
    # now you are in container

    apt-get update && apt-get install wget iperf
    # assuming 172.17.42.1 is the gateway of your docker bridge
    (cd /usr/local/bin && wget 172.17.42.1:9999/squirrel-worker && chmod +x squirrel-worker)

Now open up two new shells and create another two such containers, and we'll have three containers ready for `squirrel-worker`.

### Run `squirrel-master`
    
On linux, create `master.conf.json` with following content:

    {
        "ListenAddress":                ":1234",
        "Network":                      "10.0.0.0/24",

        "MobilityManager":              "InteractivePositions",
        "MobilityManagerParameters":    {
            "laddr": ":8765"
        },

        "September":                    "September2nd",
        "SeptemberParameters":          {
            "LowestZeroPacketDeliveryDistance": 120000,
            "InterferenceRange":                250000
        }
    }

This configuration file uses `September2nd` as wireless model, which considers interference and wireless transmission range. `InteractivePositions` spins up a web server that allows you to set position of nodes dynamically in browser. To run the master with this configuration file:

    squirrel-master -c master.conf.json

### Run `squirrel-worker` in containers

In first container:

    squirrel-worker -m 172.17.42.1:1234 -i 1 -t tap0 &
    
In second container:

    squirrel-worker -m 172.17.42.1:1234 -i 2 -t tap0 &
    
In third container:

    squirrel-worker -m 172.17.42.1:1234 -i 3 -t tap0 &
    
Again, this assumes your Linux host address in docker bridge is `172.17.42.1`. Change it accordingly if it's different. Now you should see some output from `squirrel-master` on Linux host.

If you type `ip`addr` and `ip`route` now in any of these three containers, you'll see a `tap0` interface with `10.0.0.*` IP address, and an entry for `10.0.0.0/24` in routing table.

Open [`http://localhost:8765/`](http://localhost:8765/) in browser on your Linux host. It should show a grid with three colorful dots crowded in top left corner. Drag the three dots to center, and place them not two far away from each other.

    <div class="row">
        <a href="day-20-squirrel/grid_and_nodes.png"><img class="col-sm-offset-3 col-sm-6" src="day-20-squirrel/grid_and_nodes.png"></img></a>
    </div>


### Run some `iperf` tests in containers

In first container:
    
    iperf -sui1

In second container:

    iperf -uc 10.0.0.1 -b 50M -t 100

Now in first container `iperf` should be printing bandwidth usage (about 7 Mbits/sec) if your node 1 and node 2 are not too far away from each other. Try moving node 2 away from node 1 (in browser) and see how bandwidth usage changes.

Now let's add third node's traffic. In third container:

    iperf -uc 10.0.0.1 -b 50M -t 100

The first container would show bandwidth usage for both node 2 and node 3. They both drop to around 3 Mbits/sec due to interference. If you are interested, try to use that to produce [hidden station problem](https://en.wikipedia.org/wiki/Hidden_node_problem).

Please be aware that the web page is still quite buggy by the time this is written. It sends a HTTP request whenever you move a node. If it seems unresponsive, just refresh the page :)

## Contribution

Squirrel Land is GPL licensed. It's still in early stage. (There's no test coverage yet please don't hate it.) Please feel free to [drop me a line](https://song.gao.io/), send PRs or submit issues on [GitHub](https://github.com/squirrel-land/squirrels). If you are doing research in computer networking, please consider improve or submit new [models](https://github.com/squirrel-land/models).

-- 

This article is also posted at: [http://blog.song.gao.io/posts/2013/go-in-academia-emulating-wireless-networks/](http://blog.song.gao.io/posts/2013/go-in-academia-emulating-wireless-networks/)
