+++
author = ["Matt Holt"]
date = "2014-11-18T00:00:00-06:00"
title = "Building Street Address Autocomplete with Go"
series = ["Birthday Bash 2014"]

+++

Almost two years ago, [SmartyStreets began an exodus from the .NET platform](http://blog.jonathanoliver.com/why-i-left-dot-net/). What would take its place? Go. Before moving our entire code base from .NET to a language none of us knew, we decided to write a completely new product in Go: a service to help users enter their addresses while they're still typing. We had 3 months.

Since there are over 300,000,000 addresses designated by the USPS, our first challenge was to figure out how to serve up relevant data to a user based both on their location and what they had already typed. It had to be fast since the user would still be typing.

We wanted new suggestions at every keystroke. Since the average professional can type [at least 50 words per minute](http://en.wikipedia.org/wiki/Words_per_minute), that amounts to about 3-4 keystrokes every second, which means results need to come back in a quarter of a second, otherwise things appear broken and unresponsive.

Initially, we thought about using a massive prefix tree ("trie") to store all the data because it's easy to build and traverse, but then we did some math: the results stopped us in our tracks.

The components of typical US street addresses are ordered from most specific to least specific. This means that most of the repetition found in 300,000,000+ street addresses is at the end, not the beginning. This spells bad news for prefix trees, which use lots of pointers. Assuming a 64-bit system and that most addresses contain about 25 ASCII characters, a typical trie in the worst case uses at least 300,000,000 x 25 x 8 = 60,000,000,000 bytes = 60 GB in characters and pointers alone. Yikes! Radix trees and other more optimized data structures would bring down the size, but because of the unfortunate construction of street addresses and our use case, we had to devise a better data structure. (For what it's worth, there is now [an implementation of MA-FSA in Go](https://github.com/smartystreets/mafsa) which will change things for us in the future.)

This is where Go really shined. Where most languages like Python, Ruby, and PHP [bloat even simple data structures](http://nikic.github.io/2011/12/12/How-big-are-PHP-arrays-really-Hint-BIG.html) with high-level features to make them more convenient and powerful, Go keeps the pedal to the metal, giving us C-like performance and space optimization.

It's important to know about [Go's numeric types](https://golang.org/ref/spec#Numeric_types). When loading lots of data, using plain int where you only need int8 is a serious error. On 64-bit systems, you're wasting 7 bytes of space for every 1-byte integer. Instead of using `[][]int` which wasted about 2/3 of the space we allocated, we made our own structure:

    type Optimized struct {
        field1  int8
        field2  int32
        field3  int16
        field4  int8
        field5  int16
        field6  int8
    }

and then used `[]Optimized` instead of `[][]int`.

But we can do even better. Struct fields are sequential in memory, and a type's values are aligned to word boundaries. While you would expect one `Optimized` struct to use 11 bytes, `unsafe.Sizeof(Optimized{})` [shows 16 bytes](http://play.golang.org/p/I0RVsUxDuy). That's because field2 must be aligned on a 32-bit word boundary, so field1 has 3 bytes of padding after it. There's another byte of padding after field4 to make field5 start on a 16-bit boundary. Reordering the fields like this:

    type Packed struct {
        field2  int32
        field3  int16
        field5  int16
        field1  int8
        field4  int8
        field6  int8
    }

eliminates those 4 padding bytes and the size goes down to 12 (1 padding byte at the end). With our data, this saved another ~25 MB of RAM.

Thanks to simple optimizations like this and some tricks we learned about street addresses, the data structure we finally settled on fits entirely in memory on low-priced, low-powered 32-bit servers. That's a huge win: no disk I/O.

What's the usual trade-off when you optimize memory space? CPU. Now that we had the data in memory, we needed to account for the user's location and typing habits. Their city would be guessed from their IP address using a standard geo-IP database, which we load into memory from a text file. Lookups are performed by a binary search, and the results are cached between keystrokes so a geolocation search only has to happen once per user.

Some users enter directionals (North, South, ...) and other's don't. Some type suffixes (Street, Avenue, Blvd, ...) as abbreviations, some spell them out, others don't type them at all. Worse yet, some street names contain words that look like suffixes or directionals. These issues and others multiply the burden of parsing the user's partial input.

Yep, we need lots of CPU time, and we need it immediately.

Thankfully, Go makes it easy to hog the processors. We spin up a dozen or more goroutines for every keystroke, and each goroutine follows a different path, reporting its findings to an aggregator which collects and orders the results closest to what the user typed. Channels coordinate the execution of these lightweight threads, and thanks to Go's scheduler, we get lots of time on the clock.

To ensure that every HTTP response is dispatched quickly enough, we use time.After() in a `select` to halt all the goroutines and send back whatever results we've got:

    select {
    // ...
    case <-time.After(tolerance * time.Millisecond):
        // halt processing, send response
    }

Now things are looking good. We fit the address data, geolocation database, and the program itself into memory on a small 32-bit system. Responses are guaranteed within a split-second, and users find it helpful to be offered suggestions while they type.

Building this data set isn't exactly easy, but it doesn't require massive cloud infrastructure either. The address data is compiled from a highly-compressed, proprietary binary database that takes up about 6 GB on disk. The build process, also written in Go, takes about 15-20 minutes on a modern Mac, maxing out the CPU and all 16 GB of RAM. The finished result is a beautiful, small, optimized data set that can be loaded into memory on the production servers in just a minute.

We've been maintaining our autocomplete service for over a year now. People love it. It has never crashed or panicked in production, mainly due to Go's convention for explicit error handling, which makes it easy to find and handle potential pitfalls. At time of publication, this little 8 MB application has served over 35 million precious keystrokes.

We don't use Go for everything, but it's good at a lot of things. Go has been a huge boon to our success at SmartyStreets, and our happiness as developers. We can write more competitive software more quickly. We aren't constrained to a particular development platform and we don't need an IDE. Deployment is a breeze. I look forward to seeing where Go takes us.

(Special thanks to Damian Gryski for reviewing this post.)
