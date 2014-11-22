+++
author = ["Clayton Coleman"]
date = "2014-11-22T00:00:00-06:00"
title = "OpenShift 3 and Go: Teaching Old Dogs New Tricks"
series = ["Birthday Bash 2014"]

+++

The first commit to [OpenShift](https://openshift.github.io) (the Platform as a Service that is so hipster that we were doing containers even before it was cool) was four years ago.  From day one it's been about making a platform that helps developers and operations move their applications into a cloudy future with the tools and technologies that they are already familiar with.  As developers, we love working on brand new things: things written in the newest languages, the hottest stacks, and the trendiest databases. But our operators and those who use the software we write know that new doesn't always mean better - it can also mean buggy, slow, and prone to major failure right when you need it to work.  And while we want to write Twelve Factor Apps, we also have to deal with It Takes Twelve People To Agree That Changing This One Java EJB Annotation Is Safe Apps.  

Building OpenShift for developers and ops alike has forced us to internalize that spectrum - how do we pick the right technologies to support those extreme contrasts?


In the beginning...
-------------------

There was Ruby.  And lo, you could use Rails to build web applications fast.  The original OpenShift team came from a diverse background.  Ops guys with years of experience running Linux clusters in enterprises, Java guys from JEE collaborative software shops, Pythonistas who cut their teeth on lambdas and iterators, and a bunch of user experience guys who swore to never write another line of JavaScript again (hint: they had to).  The only thing they could all agree on was that the most important thing in starting any new project was not having to write anything that looked like what they did before.  The next most important thing they eventually agreed on was to build something real as fast as they could, because software no one uses doesn't matter.

Ruby and Rails was the obvious choice at the time - it let us create the first version of OpenShift in [six months](http://web.archive.org/web/20110504163157/http://openshift.redhat.com/app/), and helped us launch on-premise OpenShift Enterprise 1.0 a year after that.  Ruby is easy to learn and script (important when you're pulling together the foundational building blocks of Linux), and to many on the team was a refreshing change from the static typed languages of yore.  Its ecosystem was rapidly growing and solving problems that let us focus on our product and users.

But it is not without its warts.  We use RubyGems.  Lots and lots of Ruby Gems.  Packaging and installing those gems on systems takes time and care, and every package we ship is a package that has to be watched for security issues, updated, checked for mismatches, conflicts with customer environments, etc.  The larger we grew the more we valued speed (and size) of deployment and more rigid boundaries between components.  Our client tool `rhc` was also written in Ruby, and on Windows it's still difficult to get a Ruby runtime, Git, and SSH to all play nicely.


In the middle...
----------------

We're container hipsters.  OpenShift has always isolated user applications with a combination of Unix user security (the flannel shirt of containers), SELinux mandatory access control (the slightly ungroomed beard), and good old-fashioned ops know-how (eclectic music taste?).  Looking to the future we knew we wanted to take advantage of the latest features in the kernel to isolate applications even further - to let them own their own network interfaces, install their own packages, etc.  We were putting together our first designs and prototypes... when suddenly a wild Docker appeared.

Docker made containers easy - so easy it scarcely needs introduction now - and beyond the obvious parallels to our own work (we've been all-in on Docker for a while) it introduced us to Go.

Go caught our attention.  Go was fast.  Go was simple.  Go was clean shaven and well dressed - the kind of language you wouldn't be afraid to introduce to your mom.  Go is a peacemaker - it won't let you argue about 2 space indentation or 4, tabs or spaces, or whether you can have naked if statements (don't worry, you can't).

Our experimentation with Docker started with the geard project - building a simple orchestration agent to prototype the next generation of OpenShift.  We wanted to stay close to Docker, be easy to install, and fast to load and execute, and trying Go gave us a chance to do all of those.  

Our confidence with Go grew along with our appreciation of the simplicity of the language.  It helped a few of us actually learn what [SOLID](http://en.wikipedia.org/wiki/SOLID_(object-oriented_design)) was, it reminded us how to build easily testable code by composing interfaces [ed: but monkey patching is so easy…] , and it taught - nay, forced - us how to write code that everyone else could read.

When it came time to make the jump to our next generation platform - based on Docker and the [Kubernetes cluster manager](www.kubernetes.io) from Google - we were ready to Go all-in.

Well, mostly.


dlopen\(\"something_else\", RTLD_NOW\)
--------------------------------------

Statically compiled binaries are awesome for distribution and deployment (OpenShift 3 includes an entire clustering system and client [in a single binary!](https://github.com/openshift/origin/releases)), but lack a dynamic language's ability to load code at runtime.  Our operators and customers often need to deeply tweak the behavior of the system, and where previously we could just let Rails load a few more extensions from disk, we now need to be more cognizant of building a composable system and designing configuration and customization that can be tweaked without compilation.  Valuable tools like New Relic depend strongly on the dynamic nature of Rails to work their magic - integrating and decorating the system the same way in Go is a lot more work.

Moving back to a compiled language has also limited the ability to peek at and debug (and sometimes tweak) the code of a running system.  Our operations team has used their familiarity with Ruby quite a few times to work around a development introduced issue in customer environments until fixes could be delivered.  As we move forward, more of our function will be inaccessible and we will have to compensate with richer APIs and better runtime debugging to address that deficit.


ENOMORECOMPLAINING
------------------

What we love:

* **Static binaries and easy cross compilation**

  We can run one command and generate a client and server binary for Linux, Windows, and Mac that includes our server, our client, our database, an agent, our admin tools, etc, etc, etc.  It's 22 megabytes.  It just works everywhere.  Thank you, Go.  Thank you.

* **Write code that everyone can read**

  Idiomatic Go (after a while) all starts to look the same, and that’s a good thing.  Reviewing large amounts of Go code comes easier to many of us - there’s less syntax to keep in your head and fewer ways to express the same concepts.  Metaprogramming is a lot of fun but having to debug the awesomely elegant language extension you cooked up six months ago when you have a critical bug today is a lot less fun.  Boring… but predictable.

* **Compose, don’t inherit**

  Having to learn to write object-oriented code without using inheritance was a shock. But Go makes it easy to declare interfaces on the fly (interfaces match any object that has the same methods defined) and easy to compose interfaces with embedding.  In a statically typed language, unit testing is heavily dependent on composition of interfaces, which has the side benefit of letting you spend more time looking at the code in between the interfaces and how they transform their inputs and outputs.  It's not often you can stop using a fundamental feature of (most) programming languages and realize you don't miss it.

* **A small language that compiles fast makes for a happy developer**

  The Go language is small, compiles really fast, and as a result it lets your mind focus on the actual problem and less on the tool you are using to solve it.  Code, test, debug cycles are so quick that you forget you are not working with an interpreted language.  Looking at our code, you see less boilerplate and more business logic.  Type marshalling, concurrency, and defer style control flow keeps the code clean and compact.  We miss ternary operators, but we will survive.

* **The Golang community**

  We had to get changes into the golang project to support a user namespaces feature for Docker.  Once code review started, things moved very fast, and the reviewers were friendly and worked to help get the fix in.  Whether working on the core language, or many of the extended libraries, the Go community is open by default.  Even though Go is still young, we’ve been able to build on top of a great (and surprisingly extensive) set of libraries and tools.


The undiscovered country
----------------------------------

When we started working with Google on [Kubernetes](https://kubernetes.io), we had a moment of tension.  The decision to build the next version of OpenShift around the Kubernetes vision was a no brainer - but were we ready to move our entire team and codebase to Go?

Turns out… yes.  The test of Go’s strength as an engineering language was how fast the team became productive.  We threw 30 engineers at Go, and they were delivering code in weeks.  A few months later it was the new normal.  The big board of “things I hate about Go” has moved from mostly serious to mostly joking.  Go isn’t perfect, and it isn’t magic.  It’s quick to learn, but does take some time to master.  But the rough edges don’t stop real work from getting done (and it’s sooooo fast).

Kubernetes and OpenShift intend to make it easier to build, deploy, and run any application - Twelve Factor or Twelve Levels of Abstraction, microservice or monolith, greenfield or brownfield, hipster or graybeard - in the private and public clouds of today and tomorrow.  Go is helping us build that future on a solid, no-nonsense foundation.

Thanks, and Happy Birthday Go! 
