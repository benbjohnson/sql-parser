+++
title = "Go Advent Day 10 - Beyond Static Binaries"
date = 2013-12-10T06:40:42Z
author = ["Joseph Anthony Pasquale Holsten"]
series = ["Advent 2013"]
+++

## Introduction

Today we're going to go against the general theme of the Go Advent Calendar and introduce No New Hotness™. That's because today is all about why folks in IT operations <3 go.

Fortunately, we've had a couple teasers of the ops perspective with discussions on [environment variable configs](http://blog.gopheracademy.com/day-03-building-a-twelve-factor-app-in-go) and [service discovery](http://blog.gopheracademy.com/day-06-service-discovery-with-etcd]).

## Get in Touch with Your Inner Sysadmin

Since you may not be a natural born sysadmin, let's try and get you in the mood. Have you ever:

- set up a home server?
- had the hard drive fail?
- had an update go wrong?
- had a dependency conflict when installing?
- had it become mysteriously slow?

If you remember that feeling, you know what every day feels like in IT operations. Just like you, we hate that lost, frustrated, “Isn't this just supposed to work?” sensation. Once you've spent years waging war against the robots, you learn the ops mantra:

.blockquote boring is good

All the cool languages have [incredible type systems](http://www.haskell.org/), or can [automatically scale horizontally](http://nodejs.org), or can even [metaprogram their own syntax](http://racket-lang.org). Go barely qualifies as object-oriented, has the same concurrency system used by Pike almost 25 years ago, and doesn't even have preprocessor macros. With so much boring, an admin knows they'll be able to reason about how and why code is behaving the way it does in production.

## Making the Daily Grind Delightful

By now you should be feeling the urge to wear a fedora, if not start trolling about package management systems in [`alt.sysadmin.recovery`](https://groups.google.com/forum/#!forum/alt.sysadmin.recovery). Don't worry, these sensations will pass. Let's just get through a good day of getting a Go application into production. Assuming your devs already have the code ready (ha!), it should be your job to build & deploy.

## Building

Before you've got anything to deploy, you need to build it. Unfortunately, unlike many stacks we haven't locked down a set of standard tools or phases for building. But you've got to start somewhere, and [Travis CI has a reasonable starting point](http://about.travis-ci.org/docs/user/languages/go/):

    $ go get -d -v . # download dependencies
    $ go build -v .  # build top-level binaries
    $ go test -v     # run tests

But while those phases are critical, they aren't the only hammer for each nail.

## Dependencies

Presumably you know all about [```go get```](http://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies) and its way of tracking packages from git, hg, svn and bzr. If you work at Google, or maintain every single one of your dependant libraries, you'll be comfortable running on HEAD from each of your dependencies. But most of us don't live in that world.

The most low-tech approach is to use source-control submodules, and all of the above version control systems support them out of the box. Of course, I've yet to find someone (who is not a version control hacker) who actually likes submodules. Even if you do love submodules, they still don't support dependencies stored in other version control systems.

For maintaining a `$GOPATH` of dependencies, [`godep`](https://github.com/kr/godep) supports locking to revisions in any version control system and has a delightfully boring config file:

    {
        "ImportPath": "github.com/kr/hk",
        "GoVersion": "go1.1.2",
        "Deps": [
            {
                "ImportPath": "code.google.com/p/go-netrc/netrc",
                "Rev": "28676070ab99"
            },
            {
                "ImportPath": "github.com/kr/binarydist",
                "Rev": "3380ade90f8b0dfa3e363fd7d7e941fa857d0d13"
            }
        ]
    }

While working with transitive dependencies isn't perfect, it's got support for importing an existing set of dependencies and works around many of the magical solutions of other package management tools.

## Automation

While plenty of folks are happy with [```go build```](http://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies), by far the most popular way to deal with more complicated automation situations is boring old `make`. I promise your admin loves traditional `Makefile`s over `autotools`.

Presumably you'll use this to wrap your favorite testing tool, from [```go test```](http://golang.org/cmd/go/#hdr-Test_packages) to a Test-Anything tool like [`bats`](https://github.com/sstephenson/bats]).

## Deployment

Now you've got binaries built and tested, it's time to get them onto your servers (hopefully staging before production). You'll have to drop the binary on the servers, as well as relevant configuration. Then you've got to roll forward from the last release. Finally, you may even want to support basic signals to allow for config reloads, gradual bleed-off and zero-downtime restarts.

## Installation

If you wanted boring, this phase is the epitome of it. With Go binaries statically compiled, all you've got to is drop the files. Whether you love `scp` or tarball artifacts or system packages, they're all delightfully boring!

## Local Configuration

While you deploy out binaries across all your servers, you may not desire each server or cluster to perform exactly the same. You may need logging turned up, or features flippers enabled, or certain parameters specifically tuned. And here you get more options than you could ask for!

Ultimately, you've got three main philosophies for local config: file, flag and env var.

File is perhaps most traditional, with the most common being the standard [`encoding/json`](http://golang.org/pkg/encoding/json/) package (used by e.g. Docker, Packer, Lumberjack & Camlistore).

Command line flags are the style used for Google services, easily accessed by the [`flag`](http://golang.org/pkg/flag/) package (also used by NSQ, Doozerd & Etcd).

Environment variables are most popular amongst twelve-factor apps and folks fond of `daemontools` or `runit` process management, and are well benefited by reading Hightower's discussion of [envconfig](http://blog.gopheracademy.com/day-03-building-a-twelve-factor-app-in-go).

But of course, sometimes a boring old config file is not enough. Sometimes you need something more powerful than a simple tree of config values. Sometimes, only a domain specific language that can describe full-powered plugins is needed. For that a number of projects have taken to using Lua, with implementations like [golua](https://github.com/xenith-studios/golua), [luar](https://github.com/stevedonovan/luar), and [heka's lua sandbox](https://github.com/mozilla-services/heka/tree/master/sandbox/lua).

## Cross-system Configuration

Of course, local system configuration is not nearly as important as cross-system config in today's distributed systems. Especially in systems that need [blue-green deployments](http://martinfowler.com/bliki/BlueGreenDeployment.html), or database master election, or even if you want to easily remove memcached servers; it's important to be able to push configs out across clusters of services.

[Paxos](http://research.microsoft.com/en-us/um/people/lamport/pubs/lamport-paxos.pdf) and [Raft](https://ramcloud.stanford.edu/wiki/download/attachments/11370504/raft.pdf) are distributed consensus protocols implemented by [zookeeper](http://zookeeper.apache.org), [doozerd](https://github.com/ha/doozerd) and [etcd](https://github.com/coreos/etcd) (which you should read more about in [Bonventre's Go Advent discussion of service discover](http://blog.gopheracademy.com/day-06-service-discovery-with-etcd)). These provide the strongest consistency possible for coordinating config changes across services.

For less stringent requirements, [Serf](https://github.com/hashicorp/serf) is designed to let services discover their peers.

## Signals

Once the new binary in dropped in place along with new configuration, it's time to migrate to the new release or configuration. While it's easy enough to just terminate the existing process and start a new one, that means that any existing requests will be dropped on the floor. And while this might be the [Worse is Better](http://www.jwz.org/doc/worse-is-better.html) approach, many protocols don't love this approach.

Fortunately, it's easy to add a signal handler so you can stop accepting new requests or reload a configuration:

    package main

    import (
    	"fmt"
    	"os"
    	"os/signal"
    	"time"
    )

    func main() {
    	c := make(chan os.Signal, 1)
    	signal.Notify(c, os.Interrupt, os.Kill)
    	go func() {
    		<-c
    		fmt.Println("Received Interrupt")
    		os.Exit(1)
    	}()

    	for {
    		time.Sleep(10 * time.Second)
    	}
    }


## Zero Downtime Restarts

One of the most important uses of signals is to have a restart that drops no requests. Ruby and Python users may have experience with the (g)unicorn web server, which has the same feature. For this same approach [goagain](https://github.com/rcrowley/goagain) along with [manners](https://github.com/braintree/manners) handle passing open socket connections between processes and cleanly closing down resources.

Of course, sometimes it's not as easy to modify an application to make it work with zero downtime restarts as it is to hide it behind a proxy that does the job for you. If that's the case, [socketmaster](https://github.com/zimbatm/socketmaster) can do the job for you. And even better, it doesn't even require you application to be in Go, so it's a great tool to have in your tool-belt no matter the application.

## Service management

Of course, while there's a great deal going for services written in Go, it's not perfect yet. Due to [how go's thread scheduler works](https://code.google.com/p/go/issues/detail?id=227), it's not possible for a process written in go to use traditional [fork(2)](http://man.cx/fork(2)) to daemonize.

As such, it's best to just use a modern service manager (like [upstart](http://upstart.ubuntu.com), [systemd](http://freedesktop.org/wiki/Software/systemd/), [smf](https://en.wikipedia.org/wiki/Service_Management_Facility) or [runit](http://smarden.org/runit/)) which can handle services that don't daemonize themselves. Of course, each of these also has at least a dozen reasons to use them over traditional System V or BSD init.

## Congrats!

With all that, you've got a lovely process running that you'll have no trouble building and deploying in future. Of course, nothing's perfect and there's still room for improvement. But things should be sufficiently boring here that you can return back to interesting things.

And that's what IT ops is all about!
