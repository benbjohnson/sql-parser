+++
title = "Go Advent Day 11 - Build a Christmas List with Martini"
date = 2013-12-11T06:40:42Z
author = ["Jeremy Saenz"]
series = ["Advent 2013"]
+++

## Introduction
[Martini](http://github.com/codegangsta/martini) is a Go package for web server development that has gained quite a bit of popularity over the last month. Martini was written to help make web development in Go a convenient, expressive, and DRY (pun intended) process.

As of this writing Martini has *161* watchers, *2316* Stars, and *153* Forks on Github. There is a ton of weekly activity around both the [martini](http://github.com/codegangsta/martini) and [martini-contrib](http://github.com/codegangsta/martini-contrib) repositories.

If you haven't already be sure to check out the [Video Demo](http://martini.codegangsta.io/#demo).

## Hello world
Martini makes web development classy, just take a gander at the following code: 

	package main

	import "github.com/codegangsta/martini"

	func main() {
		m := martini.Classic()

		m.Get("/", func() string {
			return "Merry Christmas!"
		})

		m.Run()
	}


The Martini API was obsessively designed to make HTTP servers easy to write and easy to read. A [`martini.Classic()`](http://godoc.org/github.com/codegangsta/martini#Classic) contains a good set of base functionality like logging, error recovery, routing, and static file serving. If you do not need any of the base functionality it is just as easy to instantiate a blank canvas with [`martini.New()`](http://godoc.org/github.com/codegangsta/martini#New).

Handling HTTP requests is extremely intuitive. A [`Handler`](http://godoc.org/github.com/codegangsta/martini#Handler) in Martini is any callable function. If your function returns something, Martini will write it out to the http response body. In the same vein, a handler function can return a `(int,`string)` and Martini will write out a response code as well as a body.

The return handling is fun, but most of Martini's power comes from [Middleware](https://github.com/codegangsta/martini#middleware-handlers) and [Services](https://github.com/codegangsta/martini#services). Let's dive a little deeper and create our first real web app in Martini!

## A Go Advent Christmas List

In the spirit of Christmas, I decided that it would only be appropriate to create a Christmas Wish List web app in Martini. By the end of this tutorial we will have a functioning we app complete with *HTML*template*rendering*, *Form*parsing*, and *MongoDB*integration*. Oh yeah, did I mention this will be only *50*lines*of*code?*

### Rendering HTML Templates
Our Christmas wish list only needs one endpoint, `/wishes`. When a browser visits the `/wishes` endpoint we should provide them with a list of wishes that already exists as well as a form to create a new wish.

To render out the HTML we will use an external package from martini-contrib called `render`.

	package main

	import (
		"github.com/codegangsta/martini"
		"github.com/codegangsta/martini-contrib/render"
	)

	func main() {
		m := martini.Classic()
		m.Use(render.Renderer())

		m.Get("/wishes", func(r render.Render) {
			r.HTML(200, "list", nil)
		})

		m.Run()
	}


With minimal amount of code we were able to integrate some powerful functionality. Calling `m.Use(render.Renderer())` adds the [`render.Renderer()`](http://godoc.org/github.com/codegangsta/martini-contrib/render#Renderer) to our HTTP stack as middleware. When a HTTP request comes in, Martini will pass it through the applications middleware layer for processing. In this case, the `render.Renderer()` middleware provides a [`render.Render`](http://godoc.org/github.com/codegangsta/martini-contrib/render#Render) interface for use to access via our handler functions argument list.

The `render.Render` interface allows us to easily render HTML templates using the Go standard library's `html/template` package. By default, templates are read from the `templates/` directory with a `.tmpl` file extension.

Here is our wishlist template, located at `templates/list.tmpl`:

	<html>
	  <head>
	    <link rel="stylesheet" href="http://yui.yahooapis.com/pure/0.3.0/pure-nr-min.css" />
	  </head>

	  <body style="margin: 20px;">
	    <h1>Wishes</h1>
	    {{range .}}
	      <div> {{.Name}} - {{.Description}}</div>
	    {{ end }}

	    <h1>Add a wish</h1>
	    <form action="/wishes" method="POST" class="pure-form">
	      <input type="text" name="name" placeholder="name" />
	      <input type="text" name="description" placeholder="description" />
	      <input type="submit" value="submit" class="pure-button pure-button-primary"/>
	    </form>
	  </body>
	</html>


We have a very basic form that POSTs to the `/wishes` endpoint and a loop over the target object, which will eventually be a list of wishes. We will cover both of in more detail, but first let's get MongoDB set up.

### Hooking up MongoDB
Now that we learned how to use a middleware with `render.Renderer`, let's build our own! For this example we will use the wonderful `labix.org/v2/mgo` package as a driver for Mongo. This is what the code looks like for our `DB` middleware:

	// DB Returns a martini.Handler
	func DB() martini.Handler {
		session, err := mgo.Dial("mongodb://localhost")
		if err != nil {
			panic(err)
		}

		return func(c martini.Context) {
			s := session.Clone()
			c.Map(s.DB("advent"))
			defer s.Close()
			c.Next()
		}
	}

When `DB()` is called we initialize a Mongo session on localhost. `DB()` returns a [martini.Handler](http://godoc.org/github.com/codegangsta/martini#Handler) which will be called on every request. We simply clone the session for every request and make sure it is closed once the request is done being processed. The important bit is the call to `c.Map`. This maps an instance of `*mgo.Database` to our request context. This allows all subsequent handler functions to specify a `*mgo.Database` as an argument and get it injected.

Now that we have a `DB()` middleware we can add it to our middleware stack like so:

	m.Use(DB())

Having a database isn't that useful unless we have data to store. So we will create a Wish struct that we can serialize and deserialize into our database:

	type Wish struct {
		Name        string `form:"name"`
		Description string `form:"description"`
	}

You may have noticed the `form:"name"` tags on each field. These will be utilized a bit later when we need to parse out our HTML POST form.

For convenience we will throw in a simple `GetAll` method to retrieve all of the `Wish` objects our of our database:

	// GetAll returns all Wishes in the database
	func GetAll(db *mgo.Database) []Wish {
		var wishlist []Wish
		db.C("wishes").Find(nil).All(&wishlist)
		return wishlist
	}

### Listing the Wishes

Now that all of the heavy lifting is done, we can modify our `GET`/wishes` handler to use the our new `*mgo.Database` service.

	m.Get("/wishes", func(r render.Render, db *mgo.Database) {
		r.HTML(200, "list", GetAll(db))
	})

We are now telling the `templates/list.tmpl` template to render with the result of `GetAll(db)` as the target. Looking back at our HTML you can see that we loop over our target and render a name and description for every `Wish`. 

We get the `*mgo.Database` type injected into our argument list because it was mapped by the handler that we returned from `DB()`. Since we called `session.Copy` in our `DB()` handler the request receives it's own database connection. This is a very powerful pattern in Martini as it allows us to forget about writing all of the boilerplate for setting and tearing down a session for each new request.

Time to get some wishes in our list!

### Creating new Wishes

Our HTML form POSTs to the `/wishes` endpoint, so we can write a new route to handle that form. In this example we will make use of the `github.com/codegangsta/martini-contrib/binding` package, which gives us some awesome utilities for binding form data to our `Wish` struct:

	m.Post("/wishes", binding.Form(Wish{}), func(wish Wish, r render.Render, db *mgo.Database) {
		db.C("wishes").Insert(wish)
		r.HTML(200, "list", GetAll(db))
	})

The call to `binding.Form(Wish{})` will parse out form data when the request comes in. It will bind the data to the struct and map it to the request context so we can get it injected into our next handler function. We then `Insert` the wish into the database and render out our `template/list.tmpl` our newly updated list of wishes.

### 50 Lines of Go Later...

Putting all of this together, we have a whopping 50 lines of Go code! For everything it does, this wish list app is surprisingly succinct and elegant:

	package main

	import (
		"github.com/codegangsta/martini"
		"github.com/codegangsta/martini-contrib/binding"
		"github.com/codegangsta/martini-contrib/render"
		"labix.org/v2/mgo"
	)

	type Wish struct {
		Name        string `form:"name"`
		Description string `form:"description"`
	}

	// DB Returns a martini.Handler
	func DB() martini.Handler {
		session, err := mgo.Dial("mongodb://localhost")
		if err != nil {
			panic(err)
		}

		return func(c martini.Context) {
			s := session.Clone()
			c.Map(s.DB("advent"))
			defer s.Close()
			c.Next()
		}
	}

	// GetAll returns all Wishes in the database
	func GetAll(db *mgo.Database) []Wish {
		var wishlist []Wish
		db.C("wishes").Find(nil).All(&wishlist)
		return wishlist
	}

	
	func main() {
		m := martini.Classic()
		m.Use(render.Renderer())
		//START3 OMIT
		m.Use(DB())
	
		m.Get("/wishes", func(r render.Render, db *mgo.Database) {
			r.HTML(200, "list", GetAll(db))
		})
		
		m.Post("/wishes", binding.Form(Wish{}), func(wish Wish, r render.Render, db *mgo.Database) {
			db.C("wishes").Insert(wish)
			r.HTML(200, "list", GetAll(db))
		})

		m.Run()
	}


[Github Source](https://github.com/codegangsta/go-advent-martini)

I hope this example gives a little more of an in depth look at how to build real applications in Martini. If you are interested in learning more about Martini check our the links below. 

## Community
Although Martini is young, the community surrounding it is vibrant. Here are some resources to help you get involved:

- [Martini Site](http://martini.codegangsta.io)
- [Github](https://github.com/codegangsta/martini)
- [Useful Martini Components](https://github.com/codegangsta/martini-contrib)
- [Mailing List](https://groups.google.com/forum/#!forum/martini-go)
- [Build a RESTful API with Martini](http://0value.com/build-a-restful-API-with-Martini)

## Go Build Stuff!
The Beauty of Martini is that it embraces the diversity of the web. Martini gives you the flexibility to build API's, render HTML, and serve rich content over HTTP. Make your Martini the way you like it!

2014 is going to be rad for the Go community. Let us together build some valuable web components for the Go ecosystem!
