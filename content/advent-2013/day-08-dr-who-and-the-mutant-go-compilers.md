+++
title = "Go Advent Day 8 - Doctor Who and the mutant Go compilers"
date = 2013-12-08T06:40:42Z
author = ["Elliott Stoneham"]
series = ["Advent 2013"]
+++


At the end of 1963 the UK was in the grip of [Beatlemania](http://en.wikipedia.org/wiki/Beatlemania), but for an impressionable 7-year old like me, it was the arrival of [Doctor Who](http://en.wikipedia.org/wiki/Doctor_Who) that really fired my imagination. The highlight of my week was peeping out from behind the furniture in our living-room on a Saturday evening, alternately terrified and transfixed by the exciting new program on our small black-and-white television set. During the rest of the week, my friends and I would play "[Daleks](http://en.wikipedia.org/wiki/Dalek)”. We would take turns either to mechanically move towards the inferior humans shouting “EXTERMINATE, EXTERMINATE", or recoil in mock-horror and run away back to our trusty [TARDIS](http://en.wikipedia.org/wiki/TARDIS) to escape. Happy days.

Half a century later, I found myself tweeting on the parallels between my life-long favourite sci-fi show and the joys of being a gopher: 

<blockquote class="twitter-tweet" lang="en" align="center"><p>The Go language is like Dr Who&#39;s TARDIS: A small idiosyncratic exterior, conceals large-scale quality engineering. <a href="https://twitter.com/search?q=%23golang&amp;src=hash">#golang</a></p>&mdash; Elliott Stoneham (@ElliottStoneham) <a href="https://twitter.com/ElliottStoneham/statuses/403929226371289088">November 22, 2013</a></blockquote>
<script async src="//platform.twitter.com/widgets.js" charset="utf-8"></script>


The TARDIS is a good metaphor for the [Go language specification](http://golang.org/ref/spec). For most languages this document would be unreadably dry, but not for Go. In fact I think it is the essence of why the language will be successful in the long term - it contains the smallest possible number of features necessary to provide the functionality required. Everything else is in the extensive libraries. 

By resisting feature-bloat in the core language, Go is very easy to read. [C. A. R. Hoare](http://en.wikipedia.org/wiki/Tony_Hoare) highlighted why this is a key implementation issue: 

> ”There are two ways of constructing a software design: One way is to make it so simple that there are obviously no deficiencies, and the other way is to make it so complicated that there are no obvious deficiencies. The first method is far more difficult.”

Go enables the first method, even for very large code-sets. That is quality engineering.

But Go's small feature set also has another big advantage, it makes the core language relatively easy to implement. So, like the TARDIS, the Go language specification can travel to other (programming) worlds.

### The Time Lords of Golang 

Rob Pike has recently [blogged](http://blog.golang.org/cover) that: “From the beginning of the project, Go was designed with tools in mind.”

There are now a dazzling array of tools, yet the team are still enthusiastic to build more. One current prototype tool and it’s libraries have inspired my Go programming efforts over the last 8 months: [Oracle](http://godoc.org/code.google.com/p/go.tools/cmd/oracle), a tool for answering questions about Go source code, analyses the entire code-set of a Go application. 

As a by-product of creating the Oracle tool, the Golang team have already created much of what is required to implement a Go compiler in Go.

The [design document](http://golang.org/s/oracle-design) for the oracle tool states that: 

> “As a prerequisite to pointer analysis, the program must first be converted from typed syntax trees into a simpler, more explicit intermediate representation (IR), as used by a compiler.  We use a high-level static single-assignment (SSA) form IR in which the elements of all Go programs can be expressed using only about 30 basic instructions.”

It is this SSA form IR created by [go.tools/ssa](http://godoc.org/code.google.com/p/go.tools/ssa) that holds such promise for the future.

For testing purposes, [go.tools/ssa/interp](http://godoc.org/code.google.com/p/go.tools/ssa/interp) provides an interpreter for these basic instructions; which you can either access through the [ssadump](https://code.google.com/p/go/source/browse/?repo=tools#hg%2Fcmd%2Fssadump) command or through your own short Go program, as in this [gist](https://gist.github.com/elliott5/7578605).

But the more intriguing prospect is to translate these basic instructions into executable code.

### Materialising on alien worlds

The future I envisage, is one in which your Go code can be automatically translated into many other programming languages and intermediate forms.
 
Imagine how useful it would be to be able to write a set of code once in Go and have it available to be used from many other languages, after automatic translation. It would be possible to create portable library packages. 

Or what if you want to run a chunk of Go code in some environment not targeted by the existing compilers. Maybe using JavaScript or Flash on a web page; or inside a mobile app; or as part of some large established existing non-Go application. Here is a route to doing that.

This is not pie-in-the-sky:

- [Haxe](http://haxe.org) is a successful open-source language with a vibrant user community that already takes this multi-target approach. 

- [GopherJS](https://github.com/neelance/gopherjs), a transpiler from Go to JavaScript, already implements much of the Go language specification. It does not use the go.tools/ssa library itself, but does rely on the same underlying libraries in the go.tools sub-repository.

- And I have already successfully translated some Go (via SSA and Haxe) into JavaScript, C++, Java, C#, PHP and Flash as an [experiment](https://speakerdeck.com/elliott5/ssa?slide=29) (including a single-threaded emulation of goroutines and channels).

But the target of the translation need not be a programming language, it could be another intermediate representation. The [llgo](https://github.com/axw/llgo) project is translating the Go SSA form into LLVM bitcode, in order to exploit that compiler architecture’s comprehensive back end. And there is no reason why, for example, Java bytecode could not be similarly targeted by some future project.

More speculatively, if the SSA form of the Go program were further analysed and re-worked, it might be possible for Go programs to run on parallel heterogeneous platforms. Maybe using [OpenCL](http://en.wikipedia.org/wiki/OpenCL), perhaps with the addition of a new Go library package to express data-parallelism. This approach might lead to even faster Go execution times for some applications.

### Sonic screwdriver required

Aside from the prototype nature of parts of the current go.tools code-base, which will improve as it matures, the two biggest long-term issues are speed and standard library portability.

Speed is an issue for two reasons:

- Firstly, the target environment will not have all of the language features of Go, so some level of emulation will always be required, potentially slowing down execution. The inspiring [Emscripten](http://emscripten.org) LLVM to JavaScript compiler project has an excellent [technical paper](https://github.com/kripken/emscripten/blob/master/docs/paper.pdf?raw=true) which describes the types of issues involved.

- Secondly, the SSA form of a Go program is intentionally un-optimised, to allow for easy analysis. Thus if the Go SSA is naively translated into some other form, the speed that can be achieved is dependent on the optimisation abilities of the compiler or interpreter of that other form. So in the medium term, it would probably be helpful to write an SSA-to-SSA optimiser of some sort.
 
The excellent Go standard library packages, which naturally were not written with the sort of extreme portability that I am advocating in mind, also present some challenges:

- Some execution environments may not fully support unsafe pointers, or have the same data sizes and alignments as the main compilers.

- Many of the library packages make use of runtime assembler/C functions for speed or to access low-level functionality, which need to be emulated in the target environment, where possible.

- Links to operating system functionality is not be possible in some execution environments, where that underlying functionality does not exist.

- Finally, the size and interconnectedness of library packages may also be a practical issue on the client side. For example, a test involving the frequently imported `unicode` package generated un-minified JavaScript files of 1.7M using my experimental (Go->SSA->Haxe->JS) set-up (reduced down to 492K after using the [Closure compiler](https://developers.google.com/closure/)) and 597K using [GopherJs](https://github.com/neelance/gopherjs).    

In consequence, not all of the Go standard libraries will be either available or desirable in every Go target execution environment. 

However, by contrast, I believe there is no imperative programming environment in which the core Go language specification itself could not be implemented. The only exception being the “System considerations” section of the specification regarding “Package unsafe” and “Size and alignment guarantees”. 
 
### A trip in the TARDIS

My personal plan is to continue my [adventures with go.tools/ssa](https://speakerdeck.com/elliott5/ssa). Using the Go language specification as my TARDIS, I want to see just how far in computing time and space I can travel. 

Let’s use the TARDIS together right now. Since Go has just had its fourth open-source birthday, let’s travel forward in time by another four years.

Picture the festive scene on Christmas morning 2017, with the whole family gathered around. An eager young gopher opens their presents, while their parent looks adoringly on and asks:

- Has Father Christmas brought you any nice new toys? 

- Really? A whole set of new Go compilers. How exciting! 

- Oh, so you’ll be able to use Go for almost every project now. How wonderful!

- OK, OK, you go off and play with your new Golang toys straight away, if you must. (sigh)

Happy Christmas!
