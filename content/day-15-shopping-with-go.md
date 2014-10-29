+++
title = "Go Advent Day 15 - Accelerating ecommerce with Go"
date = 2013-12-15T06:40:42Z
author = ["Shane M. Hansen"]
series = ["Advent 2013"]
+++

## Welcome

### Writing an ecommerce site in Go

Go adoption in the enterprise is increasing since the 1.0 release. Large respected tech companies have been using Go to build interesting back end services like [etcd](https://github.com/coreos/etcd), specialized content delivery systems like [dl.google.com](http://www.oscon.com/oscon2013/public/schedule/detail/28669), and mobile optimization services like [Moovweb](http://www.moovweb.com/).

[Steals.com](http://steals.com), a boutique quality daily deal site for children and women, was preparing to launch its new retail site, with the goal of engaging customers in a Pinterest-like product presentation.

Due to past experience working with the search architects behind Backcountry.com, Best Buy, and Walmart, we (the Steals.com technology team) knew that search is a great opportunity to build a section of an ecommerce site using a separate stack such as Solr, Elasticsearch, or Go.

One of the deficiencies we wanted to rectify is that traditional commerce search focuses on drilling down and eliminating options. For example, when I'm looking at HP Laptops, I don't want to have to hit the back button to see related products, like Asus Laptops. We wanted to focus on never eliminating options, instead grouping and sorting our entire product list based on user input. We aren't aware of a search engine built around this paradigm.

### Why Go?

We wanted the instant performance gains we'd get by building on top of Go's HTTP stack rather than the more full featured, but creaky, legacy framework our current cart system uses. Our technology team decided to build [Shop.steals.com](http://shop.steals.com/) purely in Go using an in-memory inverted index of our entire product catalog. Customers could then instantly reshuffle products based on categories or brands they found interesting.

Ultimately we chose Go because our Director of Technology was a [Go author](http://golang.org/AUTHORS#L308), we wanted to get better performance for "free". We knew that the stdlib support for [JSON](http://golang.org/pkg/encoding/json/), [SQL](http://golang.org/pkg/database/sql), and [HTTP](http://golang.org/pkg/net/http) would be able to satisfy our needs. We also wanted to take advantage of strong typing and the easy built-in testing framework to catch errors and speed the development process.

Executive buy-in at Steals.com was easy. The front end engineer for the project, [Sam Wecker](https://github.com/swecker), was also excited to build a site in Go. Steals.com is a fan of Google technologies in general, and recognized the value of Go in our architecture as we grow both from a systems standpoint and a recruiting standpoint.

Risk was mitigated by testing Go out on a new site that didn't need to share as much existing business logic and also by sharing templates with the rest of the deal sites using the [Mustache](http://mustache.github.io/) templating language with the [go mustache bindings](https://github.com/hoisie/mustache).

### Architecture of a Go search engine

We made use of [database/sql](http://golang.org/pkg/database/sql) and [encoding/json](http://golang.org/pkg/encoding/json) to pull product content from our core systems in the form of feeds. We used [github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx) for our database interaction.

We also contributed back a patch to the sqlx project for mapping multiple tables in a query to multiple embedded structs. For a simple example:

    query := "select product.*, inventory.* from products join inventory USING (product_id)"
    type SearchProduct struct {
        data.Product
        data.Inventory
    }
    product := []SearchProduct{}
    session.Select(&product, query)

We found this patch extremely useful because it allows us to easily compose our existing data structs (product, inventory, images, etc) into new structs for our product extraction feeds.

We also developed a small utility for introspecting databases to auto-generate models that work with sqlx using existing database schemas [github.com/stealnetwork/goper](https://github.com/stealnetwork/goper)

The resulting data is sucked up into our server on start up and added to an in-memory inverted index for brand, category, etc. We have plans to load balance the individual instances, but haven't seen the performance need yet. We also achieve almost 0-downtime restarts by binding to the HTTP port, indexing content, and then serving requests.

    listener, err := net.Listen("tcp", ":80")
    // index feeds, register HTTP handlers, etc
    http.Serve(listener, nil)

Additionally we developed [gosp](https://github.com/stealnetwork/gosp), a strongly typed templating language that compiles to Go. The performance is great, but catching all syntax and type errors in your template at compile time is the real productivity boon. Gosp compiles to a function that closes over template arguments, returning a function that itself takes an `io.Writer`. For example:

    func ShowProduct(resp http.ResponseWriter, req *http.Request) {
        // fetch product data from request	

        // compose the Product page and write it to resp
        templates.Master(templates.Product(product))(resp)
    }

### Testing

We found no need to extend or enhance Go's testing features. Our templating library compiles down to functions that close over template variables and take a simple `io.Writer` as an argument, so we were trivially able to test our app using [bytes.Buffer](http://golang.org/pkg/bytes/#Buffer).

### Results

We think [Shop.steals.com](http://shop.steals.com) is one of the first ecommerce sites built on Go. By leveraging the stdlib and builtin Go data structures we built an extremely fast site for displaying products to our customers in an innovative way. We built a few libraries along the way and contributed to another open source library. The entire experience was pleasant and Go has been completely stable for us.
