+++
author = ["Quinn Slack"]
title = "Go at Sourcegraph - Serving Terabytes of Git Data, Tracing App Performance, and Caching HTTP Resources"
date = "2014-11-28"
series = ["Birthday Bash 2014"]
+++

[Sourcegraph](https://sourcegraph.com) is a code search and review
application that supports and analyzes code in multiple languages: Go,
Python, Java, Ruby, JavaScript, and soon more. Even though we have
experts in each language on our team, Sourcegraph's core has been
written in Go since day one, and we've chosen Go for each new project
and system we've built. We'll run through all of the major open-source
systems and projects we've built within Sourcegraph in Go.

At a high level, Sourcegraph has 2 parts. The first is
[Sourcegraph.com](https://sourcegraph.com), the application that users
see, whose [architecture and code patterns we presented at Google I/O
2014](https://sourcegraph.com/blog/google-io-2014-building-sourcegraph-a-large-scale-code-search-engine-in-go).
The second is [srclib](https://srclib.org), our multi-language source
code analysis engine, which is completely open source. Both of these
parts are built using a number of libraries and systems that we've also
released as open source, but we're just going to cover those that
Sourcegraph.com uses, since those are the most broadly useful.

Sourcegraph.com, our main application, manages fetching and updating
several terabytes of VCS (git/hg) data, scheduling builds of projects,
integrating with external APIs, storing user data, and serving the web
app. Here's what we've built with Go to make this possible.


## Storing and serving VCS (git/hg) data: [go-vcs](https://sourcegraph.com/sourcegraph/go-vcs) and [vcsstore](https://sourcegraph.com/sourcegraph/vcsstore)

To access and fetch git and hg repositories in Go, we wrote [go-vcs](https://sourcegraph.com/sourcegraph/go-vcs). It
provides a common Repository interface that has 4 implementations: git
using [git2go](https://github.com/libgit2/git2go) (a native Go git
library that uses [libgit2](https://libgit2.github.com/)), git by
shelling out to the "git" command, hg using
[hgo](https://github.com/knieriem/hgo) (a native Go hg library), and hg
by shelling out to the "hg" command. There's an extensive test suite
(using Go's testing package) that tests that the behavior of each
implementation is identical.

<script type="text/javascript" src="https://sourcegraph.com/sourcegraph.com/sourcegraph/go-vcs/.GoPackage/sourcegraph.com/sourcegraph/go-vcs/vcs/.def/Repository/.sourcebox.js"></script>

In addition to providing VCS-specific methods such as [GetCommit](https://sourcegraph.com/sourcegraph.com/sourcegraph/go-vcs@master/.GoPackage/sourcegraph.com/sourcegraph/go-vcs/vcs/.def/Repository/GetCommit),
[ResolveBranch](https://sourcegraph.com/sourcegraph.com/sourcegraph/go-vcs@master/.GoPackage/sourcegraph.com/sourcegraph/go-vcs/vcs/git/.def/Repository/ResolveBranch), [Diff](https://sourcegraph.com/sourcegraph.com/sourcegraph/go-vcs@master/.GoPackage/sourcegraph.com/sourcegraph/go-vcs/vcs/git/.def/Repository/Diff), and [Commits](https://sourcegraph.com/sourcegraph.com/sourcegraph/go-vcs@master/.GoPackage/sourcegraph.com/sourcegraph/go-vcs/vcs/.def/Repository/Commits) (to get a list), the [vcs.Repository interface](https://sourcegraph.com/sourcegraph.com/sourcegraph/go-vcs@master/.GoPackage/sourcegraph.com/sourcegraph/go-vcs/vcs/.def/Repository) can return a virtual FileSystem that can access files and
directories as of a given revision. This [FileSystem](https://sourcegraph.com/sourcegraph.com/sourcegraph/go-vcs@master/.GoPackage/sourcegraph.com/sourcegraph/go-vcs/vcs/.def/Repository/FileSystem) has the standard
Open, Stat, ReadDir, etc., methods, which means it works with other
libraries that expect this [standard
interface](https://godoc.org/golang.org/x/tools/godoc/vfs#FileSystem).
It also lets us use
[mapfs](https://godoc.org/golang.org/x/tools/godoc/vfs/mapfs) to test
it.

Here's an example of using it to show a file at a specific revision:

<script type="text/javascript" src="https://sourcegraph.com/sourcegraph.com/sourcegraph/go-vcs@9378eed42e5ef190c80983efee7a0a6223aea7ec/.tree/cmd/go-vcs/go-vcs.go/.sourcebox.js?StartLine=94&EndLine=129"></script>

To scale [go-vcs](https://sourcegraph.com/sourcegraph/go-vcs) to work on hundreds of thousands of repositories and
terabytes of data, we built[vcsstore](https://sourcegraph.com/sourcegraph/vcsstore). It has an HTTP server, which
provides [HTTP handlers](https://sourcegraph.com/sourcegraph.com/sourcegraph/vcsstore@master/.GoPackage/sourcegraph.com/sourcegraph/vcsstore/server/.def/NewHandler) and an HTTP API to access data from any stored
repository (at URL paths like
`/git/https/github.com/user/repo/.branches/mybranch`), and an [API
client](https://sourcegraph.com/sourcegraph.com/sourcegraph/vcsstore@master/.GoPackage/sourcegraph.com/sourcegraph/vcsstore/vcsclient/.def/New), which implements the same [vcs.Repository interface](https://sourcegraph.com/sourcegraph.com/sourcegraph/go-vcs@master/.GoPackage/sourcegraph.com/sourcegraph/go-vcs/vcs/.def/Repository) interface with methods
that issue HTTP requests to the [vcsstore](https://sourcegraph.com/sourcegraph/vcsstore) server. This means our code can
access remote repositories over HTTP as though they were local
repositories. By setting HTTP cache headers on the server and using a
caching HTTP transport, we get nearly automatic caching of VCS data,
which makes our app fast.


## Web application integration testing using Selenium via [go-selenium](https://sourcegraph.com/sourcegraph/go-selenium)

As a large application with a lot of separate services, Sourcegraph can
fail in a lot of ways. To let us sleep a bit more easily at night, we
have a large suite of [Selenium](http://www.seleniumhq.org/) browser integration tests (among other
tests). We adapted an existing library to build [go-selenium](https://sourcegraph.com/sourcegraph/go-selenium), a Go
library that makes it easy to drive a Web browser and programmatically
perform front-end actions. Before building [go-selenium](https://sourcegraph.com/sourcegraph/go-selenium), we considered
using Ruby, Python, or JavaScript Selenium libraries, but we chose
[go-selenium](https://sourcegraph.com/sourcegraph/go-selenium) so that our integration tests can use other logic we've
written in Go to trigger backend user actions.

While Go's explicit error return values are nice in most cases, they can
lead to verbose test code (if you check every error, even ones unrelated
to the test at hand) or flaky test code (if you ignore error returns).
To solve this problem, [go-selenium](https://sourcegraph.com/sourcegraph/go-selenium) provides wrapper types WebDriverT
and WebElementT intended for use by test code, which combine a web driver
or element and a `*testing.T` and call `t.Fatal` if a method from the
underlying driver or element returns an error. This lets test authors
omit error checks but still report granular test errors.

Here's what a test case looks like:

<script type="text/javascript" src="https://sourcegraph.com/sourcegraph.com/sourcegraph/go-selenium/.GoPackage/sourcegraph.com/sourcegraph/go-selenium/.def/TestClick/.sourcebox.js"></script>

Test helpers in Go present another problem, though: if a helper function
calls `t.Log` (or anything that calls it, such as `t.Fatal`, `t.Error`, etc.),
the message is associated with the file and line in the helper function,
not in the test case that called the helper. We made a [quick hack to
show the test case's file and
line](https://twitter.com/sqs/status/532683226905468928), which helps us
identify the source of test failures better:

<script type="text/javascript" src="https://sourcegraph.com/sourcegraph.com/sourcegraph/go-selenium/.GoPackage/sourcegraph.com/sourcegraph/go-selenium/.def/fatalf/.sourcebox.js"></script>


## Fast HTTP caching with [httpcache](https://github.com/gregjones/httpcache), [multicache](https://sourcegraph.com/sourcegraph/multicache), and [s3cache](https://sourcegraph.com/sourcegraph/s3cache)

Our front-end app hits our HTTP API to fetch all the data it needs, so
we can use standard HTTP caching techniques to cache data. (We've found
this to be far simpler than if we had an application-specific cache, in
Redis for example, and had to reinvent caching and eviction semantics
and behavior.) At first our app just used Greg Jones'
[httpcache](https://github.com/gregjones/httpcache), which provides a
[caching HTTP transport](https://sourcegraph.com/github.com/gregjones/httpcache@master/.GoPackage/github.com/gregjones/httpcache/.def/Transport/RoundTrip) that writes to memory and a local disk. The
beauty of net/http and Go's trust in interfaces shines here; it was
super easy to drop in this caching transport, and the rest of our
application logic didn't need to change.

But as we grew to multiple app servers and performed frequent redeploys,
the hit rate of our servers' memory and disk caches declined because
each server's cache was separate and was purged on each deploy.

Thankfully, it was easy to extend the [httpcache.Cache interface](https://sourcegraph.com/github.com/gregjones/httpcache/.GoPackage/github.com/gregjones/httpcache/.def/Cache):

<script type="text/javascript" src="https://sourcegraph.com/github.com/gregjones/httpcache/.GoPackage/github.com/gregjones/httpcache/.def/Cache/.sourcebox.js"></script>

We first built [s3cache](https://sourcegraph.com/sourcegraph/s3cache), which implements the same [Cache
interface](https://godoc.org/github.com/gregjones/httpcache#Cache) that
httpcache expects and accesses [Amazon S3](http://aws.amazon.com/s3).
This meant our servers all shared the same cache.

However, the HTTPS
latency to/from S3 was a significant overhead, so we built [multicache](https://sourcegraph.com/sourcegraph/multicache),
which let us specify cache policies such as "when reading, try the
in-memory cache first, then disk, then S3" and "when writing, return
after the item has been written to memory, but continue writing to disk
and S3 asynchronously." This gives us near-RAM speeds for most
frequently used cache entries but with the high hit rates of using a
remote persisted cache.

## Distributed application tracing with [apptrace](https://sourcegraph.com/sourcegraph/apptrace)

To improve performance in a web app that hits multiple services to serve
each request, it's useful to see the timings for each action, no matter
which host it occurred on or how deep in the call stack it is. Just
looking at the top-level page generation time isn't enough. It's also
important to see metadata like HTTP cache headers to see what's actually
occurring and why things are slow.

We took ideas from [Google's Dapper
paper](http://research.google.com/pubs/pub36356.html), implementation
tips from [Twitter's Zipkin](https://twitter.github.io/zipkin/), and
code from Coda Hale's [lunk](https://github.com/codahale/lunk) to create
a distributed application tracing system in Go called [apptrace](https://sourcegraph.com/sourcegraph/apptrace). In
keeping with the principles of simple distributed tracing in the Dapper
paper, we get near-total visibility into the performance of our
distributed services by instrumenting two external call points: external
HTTP API calls by using an [HTTP transport](https://sourcegraph.com/sourcegraph.com/sourcegraph/apptrace@master/.GoPackage/sourcegraph.com/sourcegraph/apptrace/httptrace/.def/Transport) that records to [apptrace](https://sourcegraph.com/sourcegraph/apptrace), and
[SQL queries](https://sourcegraph.com/sourcegraph/apptrace@master/.tree/sqltrace/sql.go) by wrapping
[modl.SqlExecutor](https://godoc.org/github.com/jmoiron/modl#SqlExecutor).


## Why Go?

We've been able to release these projects as open source in part
because of Go's easy composition. Interfaces make it easy to improve
our app by providing better, faster implementations of an interface,
without affecting the interface's contract. We often develop these
improved implementations in separate repositories so they can't
introduce complex interdependencies into our app. This is what
occurred with our VCS data storage and our HTTP cache: we started out
with a simple concrete implementation in our app's main codebase and
then developed improved implementations of the same interface in
external repositories. Once done, open-sourcing these repositories is
a no-brainer because they're already standalone projects.

Also, Go makes it far easier to create separate projects than any other
language we're familiar with. All it takes is a .go file in a directory.
Other languages require package description files, complex directory
structures, install/setup scripts, etc.

We think all of this means that Go's open-source ecosystem is far more
mature than languages of comparable age and popularity. That, combined
with a beautifully designed and implemented standard library, makes Go a
joy to work with.

From all of us at Sourcegraph, we wish Go a very happy birthday!
