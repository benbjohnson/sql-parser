+++
author = ["Brian Ketelsen"]
date = "2014-11-10T00:00:00-08:00"
title = "Introducing the Go Birthday Bash 2014"
series = ["Birthday Bash 2014"]
+++

### The Go Birthday Bash

Today is a [special day](https://blog.golang.org/5years) for Go enthusiasts across the globe.  We're celebrating the 5th birthday of a programming language, and perhaps just as importantly: a community.  From the beginning, many knew Go was special.  With a [heritage](http://plan9.bell-labs.com/plan9/) [befitting](http://en.wikipedia.org/wiki/Newsqueak) [nobility](http://en.wikipedia.org/wiki/Communicating_sequential_processes), it was clear that Go was intended for big things.

This month we are celebrating Go by inviting guests to post articles about how Go has made a difference in their business, in their projects, and even in their lives.  We're going to hear from the people behind projects and companies big and small; but one common thread joins them together: their love of the Go programming language.

### Go: My Love Story
I'm going to get the ball rolling with my own story.  I read the [blog post](http://google-opensource.blogspot.com.au/2009/11/hey-ho-lets-go.html) on November 10, 2009 that introduced Go.  I downloaded the compiler, I wrote a few simple applications.  Then I forgot about it.  Our beloved Gopher was a little too young in those early days to keep my attention.  It wasn't until the summer and fall of 2010 when I started to pay serious attention.  Around that time weekly releases started coming out and a community started to evolve.  Truly as much as I enjoyed what the language, I enjoyed the community too.  The people on the mailing list treated newcomers well, and had an abundance of patience for a future Ruby expat such as me.  More often than not, my questions were answered clearly and concisely by friendly Gophers willing to help; often with pieces of my horrible code rewritten for me.

I began to tell coworkers and friends about Go.  This language is going somewhere.  This language is going to help us solve the problems we're facing right now in our business.  My ramblings often fell on deaf ears, though.  Not to be deterred, I began writing small utility applications in Go.  You know, the ones that are sometimes throw-away applications, not line-of-business critical.  One of the first ones I wrote was a program to parse a giant CSV file and insert the records into several different databases. 

The Ruby script that we had written to do this task frequently took a whole weekend to finish.  My first Go attempt finished in just two hours.  I spent another two hours running queries and counts against the databases just to make sure I hadn't skipped two-thirds of the records accidentally.  I was astonished at how easy it had been to break the tasks into smaller pieces and feed them through a configurable number of channels that wrote records to the database.  I went from zero to concurrent programming in just a few hours.

Fast-forward to 2013, and I'm writing nothing but Go. I've converted friends and neighbors over to my newfound religion, and I was feeling a little left out when RailsConf 2014 was announced.  Go needed a conference.  Erik St. Martin and I had lamented the lack of Go conferences several times before.  I asked: "Why don't _we_ run a Go conference?"

In much the same way that an expecting first-time parent has no idea what they're getting in to, Erik and I plunged head first into organizing the first GopherCon.  We had help and support from Gophers all across the globe.  When we opened up ticket sales, Mitchell Hashimoto (of [hashicorp](http://hashicorp.com)) bought the first TWO tickets almost immediately after sales went live.  We spent the next several months sweating bullets.  Were people going to buy tickets to go to a conference for a new programming language run by a couple of guys who had never before run a conference?  In February of 2014 we sold out our first batch of 500 tickets.  We juggled budgets, increased room sizes, begged and pleaded with the hotel and were able to make another 250 tickets available.  Those sold out within weeks.  That's when the real pressure began for us.  We knew we had to put on a conference that would make a Gopher proud.  

You will probably be surprised when I tell you that Erik and I didn't pull it off.  True, by all accounts GopherCon 2014 was a huge success.  We sold out 750 tickets for an unknown conference.  We had amazing speakers giving talks on amazing topics.  We had big name sponsors and little Go startup sponsors.  But it wasn't Erik and I that pulled off this incredible event.
 
It was a whole community.  Weeks and months before the first day of the conference we had offers for help pouring in from around the globe.  And it was a whole village that raised our Gopher and made it the event it was.  We've thanked them all before, but it's worth doing again.  We couldn't have pulled off GopherCon 2014 without the help of all the volunteers who put in countless hours doing boring tasks like designing websites, stuffing goodie bags, sorting badges, staffing check-in booths, and stacking coffee mugs.

The feeling at the end of GopherCon was one I'll always remember.  I was relieved that things went off without the public perception of a hitch.  But more importantly, I felt like the Go community took that opportunity in Denver to come together as a single entity.  It was our time to celebrate our love for Go with friends we've never met.  I hope it doesn't sound like hubris when I say that GopherCon 2014 felt like a pivotal moment in Go's history.  We certainly can't take credit for the ecosystem, the community, and the language.  But it seemed like those days were Go's coming of age.

Today Erik and I use Go for every corner of our development stacks, using
tools like [Kubernetes](http://kubernetes.io), [docker](http://www.docker.com), [CoreOS](http://coreos.com), 
and [etcd](https://github.com/coreos/etcd).  A prime example is my startup [XOR Data Exchange](http://xor.exchange) where we use that stack combined with [crypt](https://github.com/xordataexchange/crypt) and
[viper](https://github.com/spf13/viper) for encrypted configuration
storage.  We serve our API with
[net/http](http://golang.org/pkg/net/http/).  We vendor our dependencies
with [godep](https://github.com/tools/godep), and we love every minute of
our development process using [vim-go](https://github.com/fatih/vim-go).

### GopherCon 2015 - Announcement
So how about those announcements?

GopherCon 2015 will be a two day event on the 8th and 9th of July in Denver Colorado.

On the 7th of July we'll have workshops for Gophers of all skill levels who are interested in getting started a day early.  These workshops will have an additional cost, and space will be limited.

July 8th and 9th will be the official conference, with two full days of speakers in a single-track schedule.

We'll wrap up the festivities on the 10th with an optional "Hack Day".  Last year's hack day was a resounding success, so we'll get bigger rooms and plan for more people this time.  We planned for a small percentage of last year's group to stay the extra day, but ended up scrambling at the last minute to feed an additional 300 people who stayed for the _hours_ of lightning talks and group projects.  This year there will be distinct space available for the lightning talks and for groups who want to hack together on group projects.  

We're stepping up our game this year when it comes to the venue, as well.  GopherCon 2015 will be at the Colorado Convention Center in beautiful Downtown Denver.  The Colorado Convention Center has the room to accomodate a large group for a single-track conference, so we're going to double the number of tickets available.  So how large is large? This year, you can join 1500 of your Go friends in the mountains for four days of learning, camraderie, and fun.

We'll be opening up the Call for Proposals soon, and making a limited number of early-bird tickets available at discounted rates.  We'll put up a new website with all the details, too.  This year we'll have more hotels available with discounted room blocks, so it should be easier than last year to find a comfortable place to stay close to the event.  There are a number of nice hotels right across the street from the Convention Center that have agreed to discounted GopherCon rates.  Those details and more will be coming within a few weeks.

While GopherCon 2015 will be bigger and better than last year, we're striving to keep that feeling of community.  We want GopherCon to feel like coming home for the Gophers who travel from around the world to celebrate the language, the project, and the people that make The Go Community such a strong one.

Until the GopherCon site is refreshed you can get up-to-date information by
following @gopheracademy and @gophercon on Twitter.

See you there.



