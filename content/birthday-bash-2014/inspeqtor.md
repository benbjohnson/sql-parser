+++
author = ["Mike Perham"]
date = "2014-11-21T00:00:00-06:00"
title = "Inspeqtor"
series = ["Birthday Bash 2014"]

+++

When I decided to build [Inspeqtor](http://contribsys.com/inspeqtor) ([source](https://github.com/mperham/inspeqtor)), I had a fundamental choice: what language should I build it in? I’ve worked in Ruby for the last 8 years so it was a natural choice: “use the tool you know best” is never a bad choice when solving your own problem.

<img alt="Inspeqtor Logo"
     src="/postimages/why-go/inspeqtor.jpg"
     width=450 height=138
     style="padding: 10px"/>

However I’m not building something for myself: I’m building a product that will be used by thousands of others. Since Inspeqtor is an infrastructure monitoring tool, it needs to run 24/7 efficiently and reliably. I also want everyone to be able to install and use Inspeqtor with the absolute bare minimum of hassle. This means minimizing third-party dependencies like the Ruby VM itself, gems, etc. In the end, I needed to select the right tool for the job, not just the tool I know.

Ultimately these requirements led me to two options: Go or Rust. Both can build binaries which have no runtime dependencies at all. Both are reliable and fast. Rust has a huge amount of potential but version 1.0 is still not ready yet — they’re still working out syntax details — and learning their memory ownership model has a definite learning curve. I chose Go because of its maturity and the simplicity of the language.

I want to call out two features that make Go so nice to work with:

### Easy concurrency

Modern languages must have a better concurrency story than “thin wrapper around the POSIX thread API”. Go’s goroutines and channels are a simple but powerful abstraction and easier to use safely than traditional threads. [Here Inspeqtor gathers the current metrics in parallel](https://github.com/mperham/inspeqtor/blob/master/inspeqtor.go#L281) for the entities it is monitoring. I’m still trying to figure out best practices for handling errors and ensuring timeouts in goroutines. Google’s [Context](http://blog.golang.org/context) pattern looks like a strong contender to solve that problem.

### Full development workflow

To paraphrase the poet John Donne: “No programming language is an island”. Go ships with tools for testing, code profiling, documentation, cross compiling and syntax formatting. There was only one thing I felt that Go should provide that it doesn’t: an assert package for test code. [Here’s some example code](https://github.com/mperham/inspeqtor/blob/master/events_test.go#L25) which uses a 3rd party assert package I like. Adding if blocks here for every check would really hurt the readability of the code.

So far I’m very pleased with how well Go has worked out for Inspeqtor. I’ve been running Inspeqtor on my server for about two months now and never once has it gone over 10MB of RAM or any significant CPU usage. Here’s Inspeqtor monitoring itself:


    $ sudo inspeqtorctl status
    Inspeqtor Pro 0.5.0, uptime: 54h14m14s, pid: 11645

    Service: inspeqtor [Up/11645]
      cpu:system                  0.1%
      cpu:total_system            0.0%
      cpu:total_user              0.1%
      cpu:user                    0.1%       90%
      memory:rss                  8.97m      10m


Every language has strengths and weaknesses. For this purpose, Go has worked out great.

PS C and C++ were never considered. I vowed in 1997 to never manage my own memory again and not for a single day have I ever regretted that decision.
