+++
author = ["Paul Dix"]
date = "2014-11-12T00:00:00-08:00"
title = "Why InfluxDB is written in Go"
series = ["Birthday Bash 2014"]
+++

[InfluxDB is an open source time series database](http://influxdb.com) written in Go. One of the important distinctions between Influx and some other time series solutions is that it doesn't require any other software to install and run. This is one of the many wins that Influx gets from choosing Go as its implementation language.

While the first commit to InfluxDB was just over a year ago, our decision to use Go can be traced back to November 2012. At the time, we were working on a SaaS product that had a time series API at its core. This was written in Scala using Cassandra as the underlying data storage and Redis as a quick index. As a developer I had grown tired of the complexity of Scala. I had read a little about Go and was interested in using it, particularly after the 1.0 release earlier in the year. The simplicity of the language appealed to me as a great strength over my experience working with Scala.

I was also interested in using Go for the time series API because it would be a good fit for deploying on customer hardware. I didn't think we'd be able to deploy a complex architecture with many moving parts and third party dependencies behind the firewall. Having a single binary that didn't require anything else would be the way to go, which made Go a great fit. Its simplicity of deployment, copying a file to a server, can't be beat.

I hadn't actually worked with Go before, but I did a quick spike over a weekend to test performance on using Go with LevelDB as a storage engine. My code wasn't idiomatic, but the performance numbers looked great so I began the great rewrite of our API. We deployed in late January of 2013 and ran that code in production for months handling billions of requests without issues.

Fast forward to September of 2013. We had decided to pursue building out our time series API as an open source project. Since the previous project was SaaS, the architecture was a bit more complex than what we wanted in the open source offering. We decided to start with a fresh code base and incorporate some of the lessons we had learned on the user facing API design.

Since it was to be a new project, we thought again about whether Go would be the right choice. The only other contenders were either C or C++. Both of those would allow us to have a single binary and give excellent performance. However, because our experience with Go so far we knew it would be easier to work with and we'd be much more productive. That is, we'd get more done faster.

The only concern we had at the time was if the garbage collector would cause us problems. While most installations of Influx are running fine without issues from GC pauses, there are certain cases where it's an open issue for us. However, with the coming [improvements to the GC in 1.5](https://docs.google.com/document/d/16Y4IsnNRCN43Mx0NZc5YXZLovrHvvLhK_h0KN8woTO4/edit), it looks like the Go team is again doing the hard work for us.

The major advantage to Go that I'd like to close this post out with, but certainly not that last advantage, is the community. The growth of the Go community has been astounding. There are so many incredibly smart and enthusiastic members in this community doing great things with this new language. The InfluxDB team feels lucky and privileged to have stumbled upon it almost by accident 2 years ago.

Happy birthday Go, we're looking forward to an even more exciting 5 years ahead.