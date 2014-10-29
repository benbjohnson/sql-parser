+++
title = "Go Advent Day 7 - A Router for GOV.UK"
date = 2013-12-07T06:40:42Z
author = ["James Stewart"]
series = ["Advent 2013"]
+++


## Introduction

When we set out to build [GOV.UK](https://www.gov.uk/), the new home for UK Government information and services, we decided up front that we wanted an architecture that would allow us to build very focussed applications that did one thing well. We didn't have a clear idea of how our product or our teams would develop and we wanted to keep our options open. And we wanted to encourage a culture of experimentation where people could easily plug in whichever HTTP-fluent tool helped them most effectively meet whichever user need they were working on.

## Flexibility costs
Unfortunately that flexibility comes with costs and we quickly found ourselves tangled in the problem of how to dispatch incoming requests to the right application when nothing but the full path indicated which application that was.

For example https://www.gov.uk/bank-holidays goes to a "calendars" application, while https://www.gov.uk/calculate-your-maternity-pay is served by "smart answers" and something needs to make sure those applications receive those requests and can return responses to their users.

Over the past couple of years we've adopted a few approaches - initially a ruby-based proxy, then some varnish configuration, then a scala based dynamic reverse proxy, then some varnish configuration, and now, finally a solution we believe is going to stick. Since you're reading this here you'll not be surprised to here it's written in Go.

## Sharing the details
Quite a few of us had been playing with Go for a while before this project (I for one kept showing up having spent an evening reimplementing some part of our stack or other) but this was the first time we actually followed through with thorough testing, monitoring, a deployment mechanism and all those vital details.

We've just begun a series of posts on our technology team's blog going into the details of that router. We'll be covering how we started, how we decided it was right for us, how we tested it and how we deployed it. You can find those posts at https://gdstechnology.blog.gov.uk/tag/router/
