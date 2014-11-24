+++
author = ["Matt Aimonetti"]
date = "2014-11-24T00:00:00-08:00"
title = "How Go Helped Shape Splice's engineering culture"
series = ["Birthday Bash 2014"]
+++


Go is a trendy programming language, but let’s be honest: the Go language doesn’t have anything new that wasn’t first implemented somewhere else. As a matter of fact, Go is a pretty boring programming language. Why would anyone pay attention to a typed compiled language that doesn’t have generics, doesn’t let you do metaprogramming, doesn’t support inheritance and, as some said, feels like we are taken back to the 70s? Why not just write C?

5 years after its release, an increasingly large number of companies of all sizes are joining the Go community. As Co-Founder and CTO of [Splice](https://splice.com), a cloud solution for music creators built from the ground up with Go at its core,I’d like to take a step back and explain how and why Go has made a difference for us.

I won’t cover the obvious Go advantages:

* concurrency
* easy deployment
* performance

A lot of very smart people have already covered these advantages in depth. While Go does these things very well (and keeps on improving), other languages provide similar features. What I’d like to reflect on are the consequences that using Go has had on our code, team and culture. More specifically 5
 Go side effects that affect Splice engineering culture the most:

* Code standardization
* Code culture
* Less is more / simplicity
* Maintainability / explicitness
* Composition / modularity


### Code standardization.

I don’t think I’ve ever worked in a company where engineers weren’t having [long discussions about code styles and conventions](http://lightgray.bikeshed.com/). As companies grow, they usually explicitly or implicitly define their styleguides (see [GitHub’s](https://github.com/styleguide/), and [NASA/JPL’s](http://lars-lab.jpl.nasa.gov/JPL_Coding_Standard_C.pdf)). Getting there takes energy and expertise, not to mention long meetings, exhausting arguments and defending the guidelines each time they are challenged by new developers. The Go language is very opinionated when it comes to not spending time arguing about these “details”. The philosophy of simplicity is so strong that the code conventions feel like part of the language design.

> Formatting issues are the most contentious but the least consequential. People can adapt to different formatting styles but it's better if they don't have to, and less time is devoted to the topic if everyone adheres to the same style. The problem is how to approach this Utopia without a long prescriptive style guide.” - [Effective Go](https://golang.org/doc/effective_go.html#formatting)


Go ships with a few tools to **make sure your team focuses on "real" issues** and doesn’t waste time arguing about indentation characters or naming. [Gofmt](https://golang.org/cmd/gofmt/) (pronounced “go feumt”) uses the built in AST parser to rewrite your code. By default it follows the formatting conventions defined by the Go team, but can also be adapted to your own rules. [Go vet and golint](https://blog.splice.com/going-extra-mile-golint-go-vet/) are two other tools provided by Go to help you stick to conventions and find potential errors. _Golint_ is a [linter](http://en.wikipedia.org/wiki/Lint_(software)) which highlights style issues while _go vet_ focuses on code correctness. These tools are defined and used by the language creators and at Google. When I run these tools (automatically when saving or in a git pre-commit hook), I feel like I’m getting a quick code review from Rob “bot” Pike before submitting my code to someone else for a more in-depth review. It feels good to know that I am following conventions that are defined and enforced by people with so much experience and practice.
At Splice we’ve adapted to these idioms where possible, and refer to [Effective Go](https://golang.org/doc/effective_go.html) when we aren’t quite sure. Our developers spend much less time talking about “how” to do something and more time discussing "why" we are doing it. We spend time talking about logic, consequences and interfaces. In other words, most engineering discussions at Splice are focused on strategies more than tactics.


### Code culture

When you’re building a startup, you often hear that you don’t have time to do things “right”. You need to iterate quickly. You hear Zuckerberg’s infamous ["move fast, break things"](http://www.aghanomics.com/wp-content/uploads/2013/02/ZuckPic1.png) quote a dozen time a day and very often quality is not a primary concern. The other quote you hear often is Donald Knuth’s ["premature optimization is the root of all evil"](http://c2.com/cgi/wiki?PrematureOptimization) quote (usually totally missing the context but that’s besides the point).
It’s true that a good part of what we do when working on a startup is to learn about the problem, not about the implementation. We make assumptions and validate them. But does it really mean that you have to sacrifice quality, and what is quality?
Most startups take the normal prototypal exploration approach, often using high level technology with shared/community components. Prototypes are great but they are by definition meant to be thrown away and implemented properly. The challenge is that it’s way too tempting to build on top of prototypes, after all, “premature optimization is the root of all evil” right? This approach results in startups building their products on top of code that was never meant to be used as a foundation. They keep doing that for a little while and one of two scenarios play out: the startup is successful or, in the statistically more likely scenario, the startup is a flop, which is why the prototype approach makes so much sense.

I strongly believe that **programmers are misguided when they assume that speed of execution and quality are on two opposite ends of a spectrum**. This misconception is probably due to our field of practice being too young and lacking hindsight. Another part of the problem might be the lack of attention/education put into software architecture design. Too many “engineers” are looking for a collection of libraries they can put together to achieve the result they are after. They are focused on the "how" (short sighted quick fix), not the "why" (deeper understanding of the consequences), and end up very often disconnected from the business goals because they are looking for a "quick" solution. A more experienced engineer will quickly evaluate the consequences each potential solution will have on the system as a whole. It doesn’t have to be costly to architect code so the effect of future inevitable change of directions will be isolated. It certainly takes a longer for a novice than an expert to design such systems, but an expert can design code that will be able to handle changes of requirements without collapsing or requiring a full rewrite.
Between the three co-creators of Go: Ken Thompson (unix, B (the language right before C), UTF-8, regexp, plan 9), Robert Griesemer (v8 JS engine, Java’s HotSpot VM) and Rob Pike (UTF-8, plan 9), encompasses a significant amount of of design expertise. They spent years boiling the language down to the simplest, most pragmatic form. They argued back and forth until they were all in agreement which resulted in a small language with a set of very distinct values. Values that were chosen because they help developers focus on one thing: producing value!

### Less is more

Go’s take on simplicity is by far what  had [the largest effect on Splice](https://blog.splice.com/golang-improved-simplicity-reduced-maintenance/). In our code reviews, if someone doesn’t understand the intent of the code right away, we know it’s a red flag. Though it is not catalogued as such in [Martin Fowler's refactoring book](http://martinfowler.com/books/refactoring.html), at Splice we have agreed that clever is a code smell. The funny thing is that it’s often way harder to write simple code than complex code.

> Perfection is achieved, not when there is nothing more to add, but when there is nothing left to take away. - Antoine de Saint-Exupéry

It’s true that sometimes taking this minimalistic approach does result in writing a few more lines of code, and sometimes even that dreaded enemy: duplication.

> Duplication is far cheaper than the wrong abstraction. - [Sandi
> Metz](http://www.confreaks.com/videos/3358-railsconf-all-the-little-things)

### Maintainability

At Splice we are proud to be multilingual company with code written mainly in Go, JS, Ruby, Objective-C and C#. For the past 10 years of so, I spent most of my programming time writing code in Ruby. Ruby, very much like Python is a great language that I enjoy using. Something I didn’t expect and didn’t realize until we’ve had shipped Go code in production for a while is that we spend far less time maintaining/fixing production code than what I was used to in Ruby. I often joke that **our Go code requires 82% less maintenance than the same code in Ruby**. There are some obvious reasons for that. Go is a simpler, typed language that is compiled. Go’s amazingly fast compiler catches a lot of small typos we often make when using dynamic languages. These typos are usually found when writing tests or later on at runtime. But I believe that most of maintenance reduction is due to the fact that our Go code is just simpler, less abstracted than code in other languages. The surface-area being smaller, the opportunity for bugs to creep in is smaller. This might be the reason why most critical code out there is still written in C. 
As you know, debugging an issue usually means spending 90% of the time looking for the root cause, 5% thinking about a fix, 5% fixing the bug. **With explicit and simple code, finding the root cause is easier and therefore maintenance is reduced drastically**. A “boring” / straightforward code base means fewer surprises and therefore fewer bugs.

### Composition

Composition isn’t a new concept for any of us at Splice, however Go reshaped the way we looked at it. Without inheritance, Go’s approach to Object Oriented programming is quite different than what we are used to. You can still define methods on types (somewhat similar to instance methods) but your "classes" don’t inherit from each other ([you can embed types in each other though](http://www.golangbootcamp.com/book/types#sec-struct_composition)). Instead, Go relies on [composition via interfaces](https://talks.golang.org/2012/splash.article#TOC_15.). Interfaces don’t include implementations, they are just a way to define behavioral contract. Any “instance”/type can implement one or multiple interfaces and functions and methods can require their input or output to implement a specific interface. These limitations have had an interesting effect on our code: we spend a lot more time talking about how the code interacts, what interfaces we need and why. For instance, we have a storage interface and 2 storage types that implement this interface: AWS S3 and file storage. We can very easily switch between file and S3 storage. If we need to, we could add a Google Cloud Storage type implementing the same interface, and switching storage providers becomes trivial. Of course, this isn’t a new concept and you can and should take this approach in whatever language you use. But Go is pushing developers to think about making the right choices early on and makes bad design decisions harder than usual. **At the end of the day our code is very modular and is composed to multiple simple pieces coming together nicely**.

Go’s power is not in small implementation details or features, but in its opinionated and holistic approach to software design. The concrete application of Go’s philosophy directly resulted in better engineers as well as simpler, more flexible and more maintainable code.
