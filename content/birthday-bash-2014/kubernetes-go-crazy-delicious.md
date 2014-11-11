+++
author = ["Joe Beda"]
date = "2014-11-11T00:00:00-08:00"
title = "Kubernetes + Go = Crazy Delicious"
series = ["Birthday Bash 2014"]
+++

### What is Kubernetes? And what kind of name is that?

<img alt="Kubernetes Logo"
     src="/postimages/kubernetes-go-crazy-delicious/kubernetes-logo-256x248-2x.png"
     width=128 height=124
     style="float:left; padding: 10px"/>
[Kubernetes](http://kubernetes.io) is a container cluster management system.
Modeled after Google's internal systems, Kubernetes (or k8s for short) allows
users to schedule the running of Docker containers over a cluster of machines.
It is a toolset for starting, tracking and finding what work you have running
and where it is running.

In fact, Kubernetes has been off to such a great start, we've created an
official Google Cloud Platform product powered by Kubernetes: [Google Container
Engine](https://cloud.google.com/container-engine/).  The speed that we've been
able to do this is in no small part due Go and its community.

Honestly, the name is the best we could get through the Google trademark
lawyers.  It also is ancient greek for "helmsman or pilot of a ship".  As
Kubernetes part of the Docker ecosystem we find this name fitting.

### Go is Goldilocks for systems software

Kubernetes is an all new code base. While it is heavily influenced by systems we
have been developing and running for a long time (10+ years) at Google, we
started with a blank slate.  As we made this decision, Go was the only logical
choice.

We considered writing Kubernetes in the other main languages that we use at
Google: C/C++, Java and Python. While each of these languages has its upsides,
none of them hits the sweet spot like Go.

Go is neither too high level nor too low level.

C/C++ is great. We love C and C++.  Much of the similar software at Google is
written in C or C++. Many of the main Kubernetes developers grew up
professionally and feel most at home with it.  However, we wanted to ensure that
we had a vibrant community of contributors and a set of standard and extended
libraries that were easy to install and access.  C/C++ can be a little daunting
for more casual users.

Java is great.  We love Java.  Java has a great tool situation with some amazing
IDEs.  Refactoring stuff in Java is a breeze.  But we wanted to make it easy to
install on a wide variety of platforms.  The heavy runtime download and install
for Java made it less attractive along this dimension.

Python is great.  We love Python.  It is so easy to get something up and running
quickly with Python.  But the dynamic typing of Python presents challenges for
system software.  Strong typing eliminates a whole class of errors and lets us
concentrate on building the system we wanted.

### A lot to love

As many of us were new to Go, we found a lot to love.

* **Great Libraries.**  There are a great set of system libraries that come out
  of the box with Go.  In addition, there exist high quality libraries for
  pretty much every thing we needed.
* **Fast tools.**  Building and testing is fast.  We are getting addicted to the
  speed of development.
* **A KISS culture.**  Code in Go isn't overly complex.  People don't create
  FactoryFactory objects.  Tracing through the code for the part of the system
  you are interested in is usually pretty easy.
* **Cool extended tools.**  We love gofmt, easy code coverage, race detection
  coverage and go vet.
* **Built in concurrency.**  Building distributed systems in Go is helped
  tremendously by being able to fan out and collect network calls easily.
* **Anonymous functions.**  Something with the feel of C with more advanced
  features like anonymous functions is a great combo.
* **Garbage Collection.**  We all know how to clean up after our selves but it
  is _so nice_ to not have to worry about it.
* **Type safety.**  Any time you are parsing untrusted bits of the internet,
  having a language that protects you from many simple buffer overflow bugs is a
  huge step up.

### Upward and Onward

As the project and the team has grown—both within Google proper and with great
partners like RedHat—Go has scaled with us.  We've been able to make great
progress by relying on all of the qualities above.  The patterns and tools in Go
have encouraged us to make well factored and reusable code that will give us a
great degree of flexibility and velocity.

The Kubernetes team is having a blast bringing Google's ideas around cluster
management to a wider audience.  Building and supporting an open source project
and community has led to a more diverse set of "co-workers" with more diverse
points of view and has resulted in a stronger product.  I'm excited to see where
Go and Kubernetes can go together.

