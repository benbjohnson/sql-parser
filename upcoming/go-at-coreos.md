+++
author = ["Kelsey Hightower"]
date = "2014-11-25T00:00:00-08:00"
title = "Go at CoreOS"
series = ["Birthday Bash 2014"]
+++

# Go at CoreOS

When we launched the CoreOS project we knew from the very beginning that everything we built would be written in Go. This was not to make a fashion statement, but rather Go happened to be the perfect platform for reaching our goals -- to build products that make distributed computing as easy as installing a Linux distro.

It’s been almost 10 years since a new Linux distro hit the scene, and during that time Python was the language of choice. Python had good support for C interop, a large standard library, and was quickly becoming the standard for building user space system tools .

10 years was a long time ago.

CoreOS was designed during the Cloud era with the goal of moving the industry beyond the Cloud deep into distributed computing based on Linux containers and multi core machines. We wanted to build Google’s infrastructure as an open source project that anyone could download and run on their gear.

We needed a language that would enable us to build applications quickly and not get in our way or become a major bottleneck in terms of performance down the road.

Another major factor that made Go an excellent choice was the ability to produce standalone binaries. While this does not sound like a big deal it was huge for us. We wanted to keep CoreOS, our Linux distribution, extremely lightweight and only ship the bits required for running Linux containers. That meant no runtimes such as Ruby, Python, or Perl and their related dependencies would ship by default. That also meant no package manager. As a result our base image checks in around 140MB, and get this, it’s fully bootable!

Finally we knew choosing Go would give us the ability to onboard new developers with very little fuss. Over the last couple of years we have observed that our new hires are able to hit the ground running, even those with little to no Go experience. This is largely because of Go’s simple syntax and focus on straightforward solutions to general programming problems. You’ll often hear people talk about the benefits of gofmt ending tons of bikeshedding; trust me that’s a real thing. We just tell new developers to run `gofmt` and done.

So far Go has been awesome to work with, but to be honest it has not been a perfect experience. One of the sore spots early on was managing dependencies across projects. I’m pretty sure we tried everything from make files to building our own [dependency management tool](https://github.com/coreos/third_party.go), but these days we’ve pretty much [settled on using godep](https://coreos.com/blog/godep-for-end-user-go-projects), which I highely recommend.

## CoreOS Projects using Go

While Go is a great language, it does not mean very much if you are not building and shipping stuff. While CoreOS has [tons of projects](https://github.com/coreos) written in Go, I would like to focus on 3 of our major projects: etcd, fleet, and flannel.

### etcd

etcd is at the heart of CoreOS -- it’s our highly-available key value store for shared configuration and service discovery. Many of our projects leverage etcd for leader election, lock service, and/or configuration. etcd has also been adopted by other projects such as Cloud Foundry, SkyDNS, and Kubernetes for similar use cases.

etcd makes extensive use of the Go standard library and a few third party libraries including [gogoprotobuf/proto](https://code.google.com/p/gogoprotobuf) and [net/context](https://godoc.org/golang.org/x/net/context), which provides everything we needed to implement the raft protocol, once or twice, and our new persistent datastore backed by a WAL with CRC checksums for data integrity. To the surprise of many, our new [raft implementation](http://godoc.org/github.com/coreos/etcd/raft) is pretty lightweight and easy to navigate. 

When it comes to performance, so far so good, Go’s GC has not been an issue.

### fleet

fleet ties together [systemd](http://coreos.com/using-coreos/systemd) and [etcd](https://github.com/coreos/etcd) into a distributed init system. Think of it as an extension of systemd that operates at the cluster level instead of the machine level.

fleet is interesting because it plays a big role in a CoreOS cluster; it’s the default solution for scheduling long running jobs (applications) and provides a lightweight machine database for tracking cluster inventory. fleet runs on every machine, which means it’s absolutely critical that fleet does not consume more resources than necessary. To our delight fleet only consumes about 15MB of RAM and a fraction of the CPU during normal operation.

### flannel

One of the newest projects in the CoreOS stack is flannel, an etcd backed overlay network fabric for Linux containers.  When we set out to build flannel, time to market was absolutely critical, and doing what startups do, we put one of our newest engineers on the project. He had little Go experience, but to be fair, he is a pretty bad-ass C/C++ developer with great low level skills. He was able to build and ultimately ship his vision for flannel while learning Go. With the help of code reviews the project came out pretty nice, idiomatic Go and all. The language was not a blocker and allowed us to make a new team member effective almost immediately.

flannel is also one of the projects where we really got to leverage Go’s ability to handle low system stuff. When flannel initially shipped we only had UDP encapsulation, and if you know anything about networking that means we had to take a huge performance hit by doing everything in user-space.

In a short time period we were able to add native VXLAN support to flannel which allows us to use in-kernel VXLAN to encapsulate the packets between containers. This resulted in a dramatic performance boost, and we did not have to refactor the project to do it, or switch to another language or runtime.

## Conclusion

Go has been ***and continues to be*** the go to programming language for CoreOS. We have been very fortunate to live in the sweet spot for the language. Our focus is on building simple tools to solve a large variety of infrastructure automation tasks, and high performance components for distributed systems. So far Go has not let us down and we are anxiously looking forward to all the new stuff we plan on shipping next year. Using Go of course.

