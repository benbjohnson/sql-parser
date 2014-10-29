+++
title = "Go Advent Day 17 - Pond: a New RSS+Atom Syncing Protocol"
date = 2013-12-17T06:40:42Z
author = ["Arturo Vergara"]
series = ["Advent 2013"]
+++

## The Problem

We're standing on the verge of a new era of data ownership and privacy, with decentralization and cryptography taking center stage on the technical side of things. The series of events that led us to this point has been taking place since long before our suspicions about governments and corporations invading our privacy were confirmed.

In the past decade major players in the "cloud" industry have emerged, and most users trust them--or used to trust them--blindly with their information. If you use the services of one of these companies, and it fails, all your data is gone. We have already seen what happens with one of the most illustrative examples in recent times: the death of Google Reader.

In response to this, former Google Reader users migrated over to new services, which benefitted from its death. These new services are great alternatives and offer some neat APIs, but sadly, we were left with one big problem: inconsistency and fragmentation.

There is no standard protocol for synchronising your RSS+Atom subscriptions across devices--not to mention the great inconsistent mess that RSS already is. Without a consistent way to sync, clients must be specifically built for each of these services, and that's painful and slow for users and developers alike.

## The Solution

[Pond](https://github.com/ArturoVM/pond) is a new protocol that aims to solve the problems of inconsistency and portability, by defining [a set of RESTful HTTP APIs](https://github.com/ArturoVM/pond/blob/master/api_doc.md) that send and receive data in a standardized format--namely, a set of JSON schemas (and OPML for importing and exporting subscriptions). Its main design goals are comfort, ease of use, minimalism and elegance. It is worth noting that authentication is a work in progress.

Anybody can implement Pond using their technology stack of choice and, as long as any particular implementation satisfies all the API endpoints, anything can be thrown on top of it--for example, [the Tent protocol](https://tent.io). The Pond reference server in particular, was written in Go.

## Why Go

Go, and the Go toolset were obvious choices for this project for a number of reasons:

First, when I started this project, I was already very comfortable writing Go code. I declare myself a fan of Go; it has always been a delight to use.

Second, I wanted to provide users with a way to easily deploy the reference implementation server (which I aim to make as high quality as possible, so it can also be used in production, and it's not just there as a learning tool), and Go, with its statically linked binaries and wide range of target platforms and architectures was a perfect fit.

Third, there were concurrency and performance as factors to consider when writing the reference implementation, both of which are important in server code. Go is built for speed, and has an _amazing_ concurrency model built right into it.

## An Amazing Learning Experience

In the process of writing the reference server, I learned a lot of things about Go, despite me being already fairly familiar with it. It never ceases to amaze me, with neat idiomatic expressions, super clean code, and impressive performance feats. Plus, I absolutely _love_ the docs (and `godoc`).

## In Closing

I'd like to ask you for a small favor:

If you are a software developer and care about data ownership, privacy and decentralization, please help Pond grow.

You can contribute a lot just by using the reference server and submitting bugs and feature requests. Or, if you write Go, you can contribute by solving issues, and diving into the code and getting your hands dirty.

I, and--I think I can say this confidently--every other RSS+Atom developer and user, will thank you.
