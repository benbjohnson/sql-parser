+++
author = ["Jason Moiron"]
date = "2014-11-27T00:00:00-06:00"
title = "Go at Datadog"
series = ["Birthday Bash 2014"]
+++

## Go at Datadog

In the last year, Go has started to supplant parts of our intake pipeline at Datadog that were previously written in Python. The embrace of Go by users of dynamic languages is [well documented](http://blog.iron.io/2013/03/how-we-went-from-30-servers-to-2-go.html), generally focusing on the performance and memory usage benefits of switching to a compiled language that lacks heavy per-object overhead.

We could tell the same story.  [Datadog](http://datadoghq.com) is a <abbr title="Software as a Service">SaaS</abbr> cloud monitoring and metrics platform, and because of an increasing focus on horizontal scaling in the cloud, our users have lots of servers, and their servers send us *lots* of data.  Our largest Go clusters are smaller than our smallest Python clusters. We've [written](https://www.datadoghq.com/2014/04/go-performance-tales/) about ways in which we were able to get more performance out of our Go applications.

Now that Go has turned 5, I don't want to *bury the lead* by focusing on performance.  While performance has been important, the focus of Go as a language and the story of its success at Datadog is also one of its simplicity and the impact this has on the programs we write in it.  In "[Less is Exponentially More](http://commandcenter.blogspot.com/2012/06/less-is-exponentially-more.html)", Rob Pike claims:

> Go was designed to help write big programs, written and maintained by big teams.

This might [sound suspicious](http://docs.spring.io/spring-framework/docs/2.5.6/api/org/springframework/jmx/support/WebLogicJndiMBeanServerFactoryBean.html), but the approach that Go takes is to make the language as simple as possible,  and this scales *down* extremely well to small teams writing small programs, even if the amount data is still big.

This simplicity manifests itself in the way that our software is constructed.  Because it's fairly easy to read even sophisticated Go code, new developers have been able to produce high quality Go code quickly.  This goes well beyond modifying a routine in a pre-existing project:  most of our Go programs represent an authors *first* Go code in production.  The speed of compilation allows new developers to iterate quickly and test new ideas, while the tool chain (virtually everything in the Go tool) makes it easy to build, test, visualize and profile real Go programs.  This feedback loop cycles surprisingly quickly; we've had daemons doing heavy lifting on our full data intake within a month of conception written by authors who had no prior experience with Go.

Compared to C/C++, having a garbage collector and managed memory means we've yet to suffer *any* memory access or corruption errors in production.  Compared to Python, having an affinity with the way data is laid out in memory has allowed us to write more efficient programs almost by default.  Despite lacking some higher order features of other languages, Go's feature set makes premature abstraction painful, reducing unnecessary complexity, while still allowing us to recognize and take advantage of patterns that *actually* manifest in real code.  It sounds like a "my greatest weakness is a strength" sort of cop out, but the real cost of abstraction notoriously underestimated, and some the biggest warts in our larger code base are failed or half-realized abstractions.

Its concurrency model has massively simplified both development and deployment.  Channels let us set up pipelines within programs that mirror our architecture at large.  Each program can be built by composing the required components, with the decision for what level of concurrency to use localized to that workload.  Adjusting these levels is straightforward both in concept and execution:  it's the same as adding new machines to a cluster on different sides of a queue.  The runtime's effective parallelism has massively simplified deployment and operations:  although our cluster sizes are not so significantly different in node count, in order to maximize their use of multi-core systems, our largest Python clusters run over 700 processes, whereas our largest Go cluster runs 6 (one per node). Things like balancing and locality suddenly get much, much easier.

Finally, the stability of the Go language itself has meant that we've had no problems justifying the addition of new Go consumers to our pipeline and the gradual replacement of old Python processors with Go ones.  A stable language has made for a stable ecosystem upon which to build.  Despite being only 5 years old in total, in a little over a month's time on December 14th 2014, half of its lifetime will have been post-1.0 release and [compatibility guarantee](https://golang.org/doc/go1compat).

We at Datadog do not shy away from using a language where it is well suited.  We still rely heavily on Python's excellent numerical processing libraries for a lot of analysis and its customizability on the client side, JavaScript for our front-end user experience, Ruby for our internal devops, and more.  Due to its performance, simplicity, and safety, Go has earned a growing place in our infrastructure.  We are really excited about where future releases are going to take Go, and where Go is going to take our infrastructure and our product.  Here's to 5 more years!
