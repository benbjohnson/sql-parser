+++
author = ["Ian Rose"]
date = "2014-11-17T00:00:00-06:00"
title = "Going fast at startups"
series = ["Birthday Bash 2014"]

+++

<img alt="FullStory Logo"
     src="/postimages/go-at-fullstory/fs_rect.png"
     style="float:left; padding-right: 10px"/>[FullStory](https://www.fullstory.com/) is a tool for understanding website visitors in a whole new way.  An in-page script captures everything that happens during a user's online session, including the entire DOM and every mutation. Through this novel approach, you can reconstruct and play back every session in high fidelity. Since we are capturing directly at the DOM level, this also allows us to make interactions and page elements super searchable and analyzable.

Our previous product was built as a single, monolithic Java app.  Despite our team's expertise with Java (several of the founders were the creators of [GWT](http://www.gwtproject.org/)) we still struggled with concurrency bugs, application server/framework headaches, and just generally how to grow a flexible codebase that we could iterate on quickly yet safely.  When we started to design FullStory, we knew that we wanted something different.

### Go is great for startups

Startups like ours experience frequent and disruptive change.  A successful team must quickly develop "minimum viable" products and features, test them with customers, and iterate.  Go makes this process possible through:

1. **Static typing.**  For many teams, the need for speed motivates their use of dynamically-typed languages, but we find this to be a bit myopic.  Although the *initial* development might be a little faster, we find iterating and growing these systems to be exponentially harder (both slower and more error prone) as the system grows.  Large-scale refactorings (which are quite common early on) are much easier and safer with a static type checker in hand.
1. **Explicit error handling.**  Initially, code often implements only the simplest of error handling (e.g. a server returning a 500 status for any error).  As the code matures, "smarter" error handling is often desired; for example, returning more useful error messages or status codes, or returning partial results where possible.  Finding where these improvements should be implemented is much easier when the error handling is already explicit, as required in Go programs.  In languages that use exceptions, finding exactly where to add error handling (and how to do it in a safe way) can be arduous.
1. **Emphasis on readability and simplicity.** "Debugging is twice as hard as writing the code in the first place. Therefore, if you write the code as cleverly as possible, you are, by definition, not smart enough to debug it."[1]  Go concretely pushes you to write code that is simpler and easier to understand and debug.  Channels encourage the use of local state instead of global state plus locking.  Goroutines encourage straight line procedural code instead of asynchronous calls and callback hell.  Returning errors explicitly encourage handling errors locally instead of somewhere up the stack in an exception handler.

### Go is great for both frontend and backend servers.

FullStory is built on Google App Engine for our frontend (web) servers and Google Compute Engine for our backend services.  This hybrid approach gives us the autoscaling flexibility of App Engine plus the ability to run stateful or resource-intensive services on GCE.  We use Go in both tiers, which we have found surprisingly effective.  This approach allows greater code sharing and reuse, as well as reducing the mental friction of moving between areas of the codebase.

**What do we look for in a frontend language?**

1. **Low memory baseline.**  Go's small stack sizes and minimal per-object memory overhead allows servers to serve many concurrent requests, even when running on the relatively resource-constrained App Engine VMs (max 1GB RAM)
1. **Libraries built around "IO glue".**  Instead of raw computation, frontend servers often spend most of their time [de]serializing request payloads and shuttling bytes back and forth to other services (e.g. backend servers or external APIs).  Go greatly simplifies this kind of work by standardizing IO around io.Reader and io.Writer.  Together with powerful concurrency primitives and convenient (yet concise) support for JSON encoding, most of our frontend handlers require little code and are very easy to understand.
1. **Powerful templating system.**  Although somewhat complex, Go's templating system is undoubtedly powerful and flexible.  Especially when you just want to throw together a quick internal-only admin page, html/template will get you there very quickly.

**What do we look for in a backend language?**

1. **Strong out-of-the-box profiling.**  Alongside our own Go services, we run several Java-based open source services.  The difference in ease of profiling is stark.  Dropping http/pprof in a server is an amazingly simple way to enable *remote* (this is crucial for us) cpu and memory profiling.  We have solved several bugs just by looking at the call graph generated by pprof's 'web' command.
1. **Scatter/gather is really easy.**  Backend services often perform a "scatter/gather" pattern where they query several other services (concurrently) and then aggregate the results together.  Channels + goroutines makes this very easy to implement, even when adding timeouts, partial results, retries, etc.
1. **Statically-linked binaries.**  DLL hell.  MethodNotFoundException.  Incompatible shared libraries.  Anyone with experience deploying a live service has felt these pains.  Compiling to statically-linked binaries neatly sweeps away whole classes of deployment headaches.

All in all, we have been very happy with our choice of Go.  We have seen a strong uptick in interest and community involvement from big web companies, resulting in some [really](https://github.com/facebookgo) [nice](https://github.com/googlecloudplatform/kubernetes) [libraries](https://github.com/dropbox/godropbox).  We're very excited for what the future holds for both Go and FullStory!

[1] Brian W. Kernighan and P. J. Plauger in "The Elements of Programming Style".
