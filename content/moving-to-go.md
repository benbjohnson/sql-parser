+++
title = "Moving to Go: A Pragmatic Guide"
linktitle = "moving to go"
date = 2014-01-27T06:40:42Z
author = ["Paddy Foran"]
+++

# Moving to Go

You’ve read [all](http://blog.iron.io/2013/03/how-we-went-from-30-servers-to-2-go.html) [the](http://blog.golang.org/go-at-heroku) [blog](http://word.bitly.com/post/29550171827/go-go-gadget) [posts](https://airbrake.io/blog/status/planned-airbrake-migration-love-go-love-riak) about how great Go is. You’ve lost patience with your monolithic framework of choice—Ruby on Rails, Django, etc. You’re ready to take the leap and switch to Go.

Well, what now?

That’s exactly the position we find ourselves in at [DramaFever](http://www.dramafever.com). Our site is built on Django, and it just isn’t scaling to keep up with our rapidly growing traffic. We had read great things about Go, and some of our engineers are big proponents of the language ([Dan Worth](http://www.danworth.com) runs the [Go Philly meetup](http://www.meetup.com/GoLangPhilly/)), so we decided to take the plunge and start migrating things to Go. I want to talk a bit about how we’re doing it, because it raises some interesting challenges.

## Don’t Say Goodbye Just Yet

It’s tempting to say “Yeah! Let’s throw out all our old, legacy code and rewrite everything from scratch in Go!” And what could possibly go wrong?

[Everything](http://www.joelonsoftware.com/articles/fog0000000069.html). Everything could possibly go wrong.

Most businesses can’t afford to stop all forward development to rewrite everything from scratch. They need features and bug fixes on a regular basis, or the business stagnates and dies. And the wonderful thing about rewriting everything from scratch is that it takes much, much longer than you expect it to. Always.

So we decided to apply our development into two directions: maintaining and enhancing our current Django application while [slowly migrating](http://programmingisterrible.com/post/73023853878/getting-away-with-rewriting-code-from-scratch) things to Go.

When I say slowly, I mean _slowly_. One piece at a time. We’ll be maintaining our Django application for years to come, offloading its responsibilities one at a time to Go micro-services. Each Go service does one thing, and only one thing. The Go services communicate with each other using message queues and brokers (right now SQS, but we’ll be using [NSQ](http://bitly.github.io/nsq) soon) and APIs. They communicate with our legacy Django application in the same exact way.

Breaking our monolithic application into a bunch of services preserves our ability to migrate piecemeal. Each service is ignorant of and indifferent to the programming language the other services are written in. They all speak JSON rather than `gob` or `pickle`. Each service is self-contained.

This raises some interesting problems.

## Integrating With Django

It’s all well and good to say we’re integrating with Django, but what does that really _mean_?

### Exposing Business Data

Obviously, these services are going to need access to at least a subset of the data we’re storing in Django. Things like user profiles need to be available.

We discussed two options for approaching this: we could have Go services share the same database the Django app is using, or we could use messaging and API endpoints to allow Go services to mirror the data into their own long-term cache. Having the services use the same database as Django is tempting, because it requires far fewer moving parts and is a lot simpler to implement. The down-side is that these two separate pieces of code—the Django app and the Go services—are then very tightly coupled. If either has a requirement change that forces the database schema to change, suddenly both need to be updated to account for it, in lockstep.

To make our services independent and self-contained, we opted to create a new data store for them and mirror the information. This means examining how the services are going to require that information—will they be loading lots of records at once? One at a time? At runtime or in the background?—and tailoring the API to those access patterns, adding new endpoints as they become necessary. This also means that caching and messaging need to be implemented in the appropriate places, so your services stay in sync with minimal lag, without generating a ton of API traffic as the services poll for changes.

### Templating

There’s an even trickier piece to this puzzle: our new Go services have user-facing elements. How do we make two services, written in two different languages, serve pages that look like they come from the same site? The same header and footer, the same styles and JavaScript, the whole shebang.

The ideal solution would be “Ah, you just need another service, one whose job is to render user-facing elements!” Which would be awesome, but we’re transitioning slowly, and Django isn’t really built to work that way. Also, rewriting each and every one of our pages isn’t a “small” transition. If we’re trying to use small increments in our transition, what other options do we have?

We could use Django as a proxy, routing requests through Django, making a sub-request to the Go service, obtaining the information, and injecting that into our Django template the same way we inject results from a database. But that means that the feature is now split across Django and Go; the presentation of it needs to be handled in Django, and the business logic has to be handled in Go. It also has some performance implications, and overall was just more of a hack than we were comfortable with.

The other option, the one we selected, is to treat the template as nothing more than text. We developed a [Django management command](https://docs.djangoproject.com/en/dev/howto/custom-management-commands/) that renders every possible permutation of our base template (different languages, different user types, etc.) and uploads them to S3. The Operations team will run this command as part of our deploy process, which allows us to keep our templates in sync. The Go services then download and cache these templates, and use them the basis for a Go template that can be rendered with the Go templating engine.

To achieve this, we had to manually go through our Django template and discern what information it needed to render, and we had to either inject a Go variable there, so Go could replace it at runtime, or we needed to render another set of template variations to account for every possible permutation of that value. For example, translations aren’t something we can just inject at runtime from the Go service, so we have to generate a different version of the base template for each language we support. As more variables like this add up, it leads to a combinatorial explosion.

## Refactor in Favor of Simplicity

Some of our changes aren’t new feature additions, they’re refactoring the way we handle existing things. Authentication, for example, can be a headache. We refactored things so that every request gets a header—applied inside our server stack—providing the user’s ID, so we no longer have to authenticate on every service. Instead, each service can just check for that header. Because this is a global expectation, our security models can be designed around it, so that header can always be trusted. In this way, our services become much simpler and have fewer dependencies.

By utilizing micro-services this way, we’re significantly lowering the contextual information a developer needs to keep in their brain while working on our codebase. Code with a single focus allows developers to focus, too. We’re really excited about our transition to Go, but with millions of requests a day being served, we can’t jump ship all at once. By moving slowly, our users haven’t even noticed that we’re rewriting code, our business has continued to grow, and our product continues to improve—but we’re still eliminating technical debt with extreme prejudice.

## Raindrops on Roses & Whiskers on Kittens

We’ve fallen a bit in love with Go here at DramaFever. Our Operations team loves the simple deploy procedure—they just download the binary our Jenkins server produces as an artifact. Our engineers love its flexibility and clarity. The accountants love the lower server expenses. Everyone wins!

These are a few of our favourite things.

- Interfaces make testing easy

Testing our Django application is a major pain. You need to generate fixtures, load them in specific orders, keep track of them, and figure out mismatches between your development environment’s data and the fixture data. Plus, they take forever to populate the database with.

Comparatively, Go is elegant in its database testing, thanks to interfaces. Our Go database interactions all run through interfaces. So we’ll have a special database type, something like this:

	type Database interface {
		SaveThisThing(thingToSave *Thing) error
		GetThing(thingId string) (*Thing, error)
		DeleteThing(thingId string) error
	}

Then we can call this interface from all our business logic. In production, the interface is filled by something like this:

	type SQL sql.DB

	func (sql *SQL) SaveThisThing(thingToSave *Thing) error {
		// generate and execute sQL statements here
		// return any errors
	}

But when we’re testing, we can stub the database interactions out:

	type testDB struct{}

	func (db testDB) SaveThisThing(thingToSave *Thing) error {
		// return an error or success, as your test requires
	}

Our databases are set up in a read-slave configuration, and there are some interesting race conditions that can come up due to the replication lag between writing to master and that write reaching a slave. With these interfaces, we can use goroutines and the `time.After` function to manually create situations in which race conditions would occur, then test against them. It’s trivial to test our Go code for bugs related to replication lag.

Testing nirvana.

- go test -cover is addicting

Because testing in Django is so slow and fraught, our test coverage for it sometimes falls short of our ideal test coverage. With the [built-in test coverage tool](http://blog.golang.org/cover), testing has been gamified. Without modifying our code, we can see exactly how much of our code has tests associated with it. Another command will show us exactly what our tests are missing. And watching that output climb towards 100% is strangely motivational. When a project finally hits 100% test coverage, we celebrate with animated GIFs in the IRC channel (we’re on `#dramafever`, come say hi!).

- Taking a Go approach

One of the problems with switching to Go was the loss of Django’s ORM. And while Go ORMs exist, the entire idea of an ORM seems somehow out of place in Go. Andrew Gerrand [semi-famously](http://go-lang.cat-v.org/quotes) said “In Go, the code does exactly what it says on the page.” ORMs hide a lot of what the program is doing, so the cost and complexity of a function is less clear than in other Go code.

That being said, ORMs save a huge amount of time, and a lot of our engineers are used to working with them. Transitioning to writing SQL by hand would be difficult for us. We wanted a happy medium, so we [built one](https://github.com/DramaFever/pan). While the project is still undergoing development, it tries to take a more idiomatic approach to the problem of writing SQL. Rather than writing an ORM that hides all the SQL, we built an abstraction layer on top of SQL that makes writing queries easier for us. We tried to take the same pragmatic approach to abstraction you see in the language design of Go: enough abstraction to ensure developers don’t get bogged down in verbosity, but not so much abstraction that the actual work being done is obscured. So far, the Go way hasn’t led us astray.

## Invested in Go

We’re investing heavily in Go, from our use of it in these micro-services to our sponsorship of GopherCon. We’re [always](http://gopheracademy.com/jobs/show/60) looking for talented gophers to join our team, so if these problems seem interesting to you, definitely get in touch.

And if you’re considering moving to Go, but aren’t sure how to begin, definitely share your experiences as you make the transition. Hopefully we’ve given you a starting point; we’re very happy with how it has gone for us.
