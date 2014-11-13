+++
author = ["Tom Maiaroto"]
date = "2014-11-13T00:00:00-06:00"
title = "To be Concurrent or Not to be? Sometimes Both"
series = ["Birthday Bash 2014"]

+++

Go really makes concurrency easy. That said, there's still some things to watch out for and not every 3rd party package you find out there is ready to have "go" put in front of its functions.

Likewise, it's important to know when you actually need concurrency. Just because Go makes concurrency easy doesn’t mean we always need to use it. Sometimes our applications only have small needs for it.

Fortunately with Go, it's not an all or nothing type situation.

### Concurrency for Data Mining
In my case, the application was data mining various APIs with some web crawling and event logging. [Social Harvest](http://www.socialharvest.io) is a social media analytics platform re-built in Go and it does a lot of data mining and analysis.

Previously Social Harvest was built using PHP and then later on parts were re-built with Node.js. Each of these languages has a different way of handling the task and it's really good to compare them. These are all very popular web languages and they all get the job done.

So the requirements were to take unstructured data from various sources across the web (from rate limited APIs), analyze it, and then store and log it for a flexible workflow. Data comes in and goes out to literally anywhere. It needs to do this as fast and as efficiently as possible. It must scale.

### How PHP Handled it
cURL was used to handle the HTTP requests in a procedural fashion. While it gets the job done, it’s slow. Believe it or not, PHP does have some extensions for multithreading. One of which is called pthreads. Trying to work with this extension in order to concurrently make multiple HTTP requests quickly became messy and proved to be unreliable.

On the bright side, API rate limits were not an issue. Still, it ended up being a slower data mining application and that quickly lead to the desire for a different language.

Dirty data was a problem because the internet is full of it. PHP knows this and is very accommodating because it's not strongly typed. However, when storing the data into a database that was not schemaless, there were a few annoyances. Consistent and reliable data became a problem.

These factors ultimately became too much. Especially when math was involved. It compounded all of the shortcomings and resulted in a lot of sanity checks throughout the code. This lead to a more difficult to maintain codebase as well as an even slower application.

### How Node.js Handled it
Callback hell. I'll just put that out there immediately. That said, Node.js was a good bit better than PHP for this task.

Like PHP, it's not strongly typed so it worked well with dirty data but had the same limitations when it came to storage and analysis.

Unlike PHP, Node.js is asynchronous in nature and this allowed data to be gathered faster. However, this came at the cost of system resources. A lot more RAM was required by the application now.

Telling Node.js when to garbage collect is possible, but not simple. However, that did allow me to keep some performance issues at bay. Still, the application now required more system resources than the PHP version.

Rate limits became a problem as well. The application now needed to artificially add in pauses with the timer and manage a bunch of wait groups. This wasn't too sloppy to read, but it wasn't ideal to manage.

### How Go Handled it
Go was like PHP in that the first write of everything ended up being very procedural; however, the codebase was more maintainable.

Go is also strongly typed and that initially was something to get used to coming from PHP and Node.js. When working with dirty data it felt like there was more work to be done. Truthfully, there was - up front. It paid off in the long run though when it came to storing and sending that data off elsewhere.

Go makes things very easy with structs and its marshalling system. Converting data to or from JSON or CSV is a breeze. It took the unpredictable data and made it predictable which made application development much easier and faster in the long run.

After the application was gathering data, concurrency was added. It was easy to add it to the existing application. When using other languages, more thought needs to be put into how you design your application when working asynchronously. It seems to be a do or do not. With Go, it's much easier to do either or both! [Andrew Gerrand has a great talk on this](http://blog.golang.org/two-recent-go-talks).

When it came to concurrency, the typical pain points were quickly identified with Go's tooling. Profiling (pprof) helps track down goroutines that may be misbehaving. This is not quite as easy with Node.js. It's far more difficult to get on top of a memory leak with Node.js than it is with Go.

The rate limiting by various APIs was not an issue because those function calls simply ran without the "go" keyword in front of them. In two characters the application adapted to handle this requirement. No messy timers. No callback hell.

Sure, there are wait groups in Go. Though I found myself reaching for them less than with Node.js. There's also channels which really made some things easier too. Go is well equipped to handle various concurrency related tasks, but on the surface it's all neatly tucked away like a Swiss Army knife.

Last, but not least there was a great performance increase. What required a server with 4GB of RAM with Node.js, ended up running just fine on a server with 512MB of RAM. Of course having a very easy to deploy binary is also quite nice. I even joked about running it on a Rasberry Pi.

How would you like that? A data mining tool that would typically cost you hundreds (and even thousands) of dollars per month as a service...Instead now sitting on your desk at work running on a device that costs less than $100? Gathering and analyzing millions of messages per month.

Go makes all of this possible. From the adaptable and easy to maintain codebase to the efficient and portable application.

Don't get me wrong, the other languages are very capable, but Go seemed to handle data mining and analysis best for Social Harvest.
