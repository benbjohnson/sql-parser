+++
author = ["Derek Collison"]
date = "2014-11-15T00:00:00-06:00"
title = "How Continuum ended up being written in Go"
series = ["Birthday Bash 2014"]
+++

<img alt="Continuum Logo"
     src="/postimages/apcera/continuum-logo.png"
     width=128 height=128
     style="float:left;"/>
In March of 2012, I had just left VMware and the project I had founded, architected and built, Cloud Foundry. PaaS then was still very new, as was a distributed system built in Ruby. Many Go advocates these days come from the Ruby world, which was a surprise to Go authors who believed many would come from C or C++ worlds. Go was built inside of Google from an amazing cast of authors who were looking to solve problems with the current build and link process for large C++ applications. Inside Google (I worked at Google from 2003-2009), Java applications were not much better in terms of build times and binary sizes.

As [Apcera](https://www.apcera.com) was taking shape in April of 2012, I was exploring different languages, knowing that Apcera's new system, now named [Continuum](https://www.apcera.com/continuum), would most likely not be written in Ruby. Ruby was, and still is a great language, but it presented challenges for large-scale distributed systems with high update cycles. Moreover, Ruby tends to be a bit meta, the language actually encourages you to do so, which presents its own challenges when reviewing and understanding code at a later date. I had been playing with Go since the 0.56 days in the spring of 2011, and was considering moving the [NATS](https://nats.io) messaging system towards Go. NATS was written originally in Ruby and was the control plane for systems like Cloud Foundry.

With Go, I liked what I saw so much that I predicted the rise of Go as a language for cloud systems and tooling. I believe this prediction has largely come true, and that Go as a major force in the Internet of Things will also come true.

![](/postimages/apcera/go-prediction.jpg)

I started on a prototype Go client for NATS. Actually, my first Go program was to test whether or not Go stacks were real and bypassed the garbage collector all together. At the time, this was the most important thing about the new language for me. I was tired of trying to re-architect programs from a memory model perspective to appease the garbage collectors of the world. Go’s stacks were real, which meant you could keep your memory model, and just shift it from the heap to the stack! Also, the stack would auto-promote to the heap as needed, which was a problem I had previously worked hard to solve in the C language.


    // My First Go program!

    func stack() {
		var a [22]int
		as := a[:0]

		for i := 0; i < 22; i++ {
			as = append(as, 22)
		}
    }

    func stackPromotion() {
		var a [10]int
		as := a[:0]
		ocap := cap(as)

		println("slice starts at 10, appending 100 times")
		for i := 0; i < 100; i++ {
			as = append(as, i)
			diff := cap(as) - ocap
			ocap = cap(as)
			if diff > 0 {
				println("slice cap expanded by:", diff)
			}
		}
    }



Another thing I loved about Go were static compiles. One of my big pet peeves was the inability to eliminate runtime dependencies for deployments. Dealing with these in the past had left some scars that I did not want to re-live. The builds were fast, I mean really fast. I do believe that the Ruby converts choose Go for 3 main reasons:

* **Blazing Fast Compiles**  It is as fast to develop and run code in Go as it is in Ruby and other interpreted languages.
* **Static Typing** Dynamic typing is great to start, but can cause unknown side effects later down the road. Go is a static language, but presents the users with easy-to-use syntax with the inference-aware compiler and interfaces, which present very similar to duck typing.
* **Performance** Even in the early days, Go showed promise of fast performance against the interpreted languages, and its current performance stacks up well against the leaders in the space, C++ and Java.

As performance goes, [NATS](https://nats.io) in Ruby did around 150k msgs./sec. When I moved the server to Go, but maintained the same architecture, it jumped to 700k msgs./sec. When I really started to take advantage of what Go had to offer, I pushed it to 6M msgs./sec. on a Macbook Air! I was fortunate to present at the very first GopherCon this year, and detailed my journey to get to those performance numbers in my talk at GopherCon on [High-Performance Systems in Go](http://www.slideshare.net/derekcollison/gophercon-2014)

Apcera’s choice to build Continuum in Go was simple and straightforward, the team has never looked back. We are proud that we were one of the firsts companies to fully embrace Go. 99% of Continuum is written in Go, so if you want to work for a cool company in San Francisco building distributed cloud systems in Go, email us at [jobs@apcera.com!](mailto:jobs@apcera.com?subject=I%20Want%20to%20Program%20in%20Go%20at%20Apcera)

What do I love most about Go today? Many things come to mind, from the tool chain, to the broad standard library, and the pace of innovation of Go's releases. But, the one thing that stands out most is the *Go Community*.

**Happy Birthday Go!,** we have all made an excellent choice.






