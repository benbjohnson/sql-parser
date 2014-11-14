+++
author = ["Dan Carley", "Kushal Pisavadia"]
date = "2014-11-14T00:00:00-06:00"
title = "Using Go in government"
series = ["Birthday Bash 2014"]
+++

When the UK [Government Digital Service (GDS)](https://gds.blog.gov.uk/) started working on
[GOV.UK](https://www.gov.uk/), much of it was
[built in Ruby](https://gds.blog.gov.uk/govuk-launch-colophon/). Since
then, we’ve used a number of different programming languages across
government including Java, Clojure, Scala, Python and Javascript. More
recently, we’ve turned to Go for some projects.

This is a brief experience report. It’s about how we’ve used Go and
what we feel would be useful to know for others considering it. If
you’re more interested in reading a case study delving into the
details of what we’ve done with Go, we’ve posted on our blog about our
[router](https://gdstechnology.blog.gov.uk/2013/12/05/building-a-new-router-for-gov-uk/),
[crawler](https://gdstechnology.blog.gov.uk/2014/08/27/taking-another-look-at-gov-uks-disaster-recovery/),
and
[CDN acceptance test](https://gdstechnology.blog.gov.uk/2014/10/01/cdn-acceptance-testing/)
projects.

## What made Go a viable option?

As an organisation we feel that learning and experimenting with new
technologies exposes us to different approaches and broadens our
thinking. In the case of modern programming languages they solve
problems in different ways.

We’d heard a lot of good things about Go. Its successful use at Google
for their internal systems and our knowledge of the calibre of the
team were both bonuses. However, concurrency, runtime speed, and
resource usage were important qualities for the
[first project that we prototyped in Go](https://gdstechnology.blog.gov.uk/2013/12/09/choosing-go-for-a-new-project/).
These weren’t all satisfied by some of the languages that we were
already using.

## Easy to pick up

Go has a simple language specification. This has proved valuable in
getting interest from colleagues that have no prior Go experience,
from peer reviewing code and to later contributing. Yet at no point
have we felt particularly constrained by that simplicity. When you
want to customise things, interfaces and composition make it easy and
reliable to do so.

Go’s standard library was touted as being good and has proved to be
excellent. It has a wide breadth of packages for common tasks. These
include interacting with file systems, HTTP services and building
command-line tools, through to working with JSON data and formatting
dates and times. The standard library has a depth that hints at
experienced and well-considered design. RFC standards are adhered to
and useful functions are provided as building blocks for problems you
might be working on.

The standard library has
[wonderful documentation](http://golang.org/pkg/), which is also a
great source of learning for new and seasoned Go programmers. There’s
even an
[excellent guide to writing Go code idiomatically](http://golang.org/doc/effective_go.html),
and a
[tool that formats your code correctly](http://golang.org/cmd/gofmt/).

## Easy to deploy

Over the last few years we’ve learnt about the different Ruby models
of deployment (like Unicorn workers) and
[built tooling to help us](https://github.com/gds-operations/unicornherder).
We have a culture of
[releasing regularly and releasing often](https://gds.blog.gov.uk/2012/11/02/regular-releases-reduce-risk/).
Any technology that made deployment easy was going to do well in our
environment and Go happened to shine in this instance.

Go has no special runtime requirements. A single binary is compiled
and transferred to a remote machine. There’s no extra runtime
dependency resolution (such as `bundle install` in Ruby) required on
the other end. And, restarting the service is fast in comparison to
Rails where it can take a number of seconds before you get feedback.

## Easy to use

Teams tend to decide which languages work, individuals don't. Our
usage of Go has increased over the past year and there are certain
characteristics of Go that have enabled this and made it easier to
work with.

It's been easy to get people interested in Go, from sysadmins who
claim they can't code through to developers who are picking it up as
their second language. There's a lot of momentum behind Go and what
the maintainers are trying to do with it. Specifically, the
[version stability promise for 1.x releases](https://golang.org/doc/go1compat)
is important to us. Having backwards compatible releases meant that we
could be confident working with the language over a longer period of
time and not have to worry about recompiling source code during minor
or patch releases.

Having the `go` tool cover the majority of project lifecycle tasks has
made getting to grips with the language a lot easier. Similarly, the
C-like syntax has reduced the barrier for many who have had trouble
with other language idioms.

If you follow the statement
["Make it work. Make it right. Make it fast"](http://c2.com/cgi/wiki?MakeItWorkMakeItRightMakeItFast)
then using Go means what you write is often fast enough by default.
The runtime is quick and improving on every release and the standard
library comes well equipped. This has meant that we can concentrate on
characteristics of our software that are more important to us as a
team: clarity and readability.

## Where we’re going next

For GDS to fully embrace Go there are certain problems we need to
solve. One of these problems is management of versioned dependencies.
For some of our core systems we need to guarantee the deployed
versions of code and their respective dependencies. The language
maintainers have
[publicly endorsed vendoring](http://golang.org/doc/faq#get_version).
We’re looking at using
[gom and Godep as possible solutions](https://github.com/alphagov/styleguides/blob/master/go.md#external-dependencies)
to this problem to be more developer-friendly.

It doesn’t look like our usage of Go is going to decrease any time
soon. You can read more about our experiences of Go and other
technologies on the
[GDS Technology blog](https://gdstechnology.blog.gov.uk/).
