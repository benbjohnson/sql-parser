+++
title = "Go Advent Day 5 - An introduction to beego"
date = 2013-12-05T06:40:42Z
author = ["Jiahua Chen"]
series = ["Advent 2013"]
+++


##  An introduction to beego

Beego is an open-source, high-performance and lightweight application framework for the Go programming language. It supports a RESTful router, MVC design, session, cache intelligent routing, thread-safe map and many more features that you can check out [here](http://beego.me).

This post will give you an overview and get you started with the beego framework.

### Overview

The goal of beego is to help you build and develop Go applications effectively in the Go way. Beego integrates features belonging to Go and great mechanisms from other frameworks; in other words, it's not a translated framework, it's native and designed only for Go.

Beego also uses a modular design and thus gives you the freedom of choosing which modules you want to use in your applications and leave alone useless ones for you.

The following figure shows 8 major modules of beego:

![](/postimages/day-05-beego/beego_arch.png)

And here is the classic organization of projects that are based on beego:

	├── conf
	│   └── app.conf   
	├── controllers
	│   ├── admin
	│   └── default.go
	├── main.go
	├── models
	│   └── models.go
	├── static
	│   ├── css
	│   ├── ico
	│   ├── img
	│   └── js
	└── views
	    ├── admin
	    └── index.tpl

When you develop with beego, the [Bee](https://github.com/beego/bee) tool will give great convenience features, like hot compile for your code. 

## Getting started

### Installation

Beego contains sample applications to help you learn and use beego App framework.

You will need a functioning Go 1.1 installation for this to work.

Beego is a "go get" able Go project:

       go get github.com/astaxie/beego

Or through [gopm](https://github.com/gpmgo/gopm) by executing 
	
	gopm get -gopath github.com/astaxie/beego

You may also need the `bee` tool for developing: 

	go get github.com/beego/bee 

or 

	gopm bin -dir github.com/beego/bee $GOPATH/bin

For convenience, you should add `$GOPATH/bin` into your `$PATH` variable.

### Your first beego application

Want to quickly setup an application see if it works?

	$ cd $GOPATH/src
	$ bee new hello
	$ cd hello
	$ bee run

These commands will help you:

1. Install beego into your $GOPATH.

2. Install Bee tool in your computer.

3. Create a new application called "hello".

4. Start hot compile.

Once it's running, open a browser and point it to [http://localhost:8080/](http://localhost:8080/).

### Simple example

The following example prints string "Hello world" to your browser, it shows how easy to build a web application with beego.

	package main
	
	import (
		"github.com/astaxie/beego"
	)
	
	type MainController struct {
		beego.Controller
	}
	
	func (this *MainController) Get() {
		this.Ctx.WriteString("hello world")
	}
	
	func main() {
		beego.Router("/", &MainController{})
		beego.Run()
	}

Save file as `hello.go`, build and run it:

	$ go build -o hello hello.go
	$ ./hello

Open address [http://localhost:8080/](http://localhost:8080/) in your browser and you will see "hello world".

What is happening in this example?

1. We import package `github.com/astaxie/beego`. As we know that Go initialize packages and runs init() function in every package ([more details](https://github.com/Unknwon/build-web-application-with-golang_EN/blob/master/eBook/02.3.md#main-function-and-init-function)), so beego initializes the BeeApp application at this time.

2. Define controller. We define a struct called `MainController` with a anonymous field `beego.Controller`, so the `MainController` has all the methods that a `beego.Controller` has.

3. Define RESTful methods. Once we use anonymous combination, `MainController` has already had `Get`, `Post`, `Delete`, `Put` and other methods, these methods will be called when user sends corresponding request, like `Post` method is for requests that are using POST method. Therefore, after we overloaded `Get` method in `MainController`, all GET requests will use `Get` method in `MainController` instead of in `beego.Controller`.

4. Define main function. All applications in Go use main function as entry point like C does.

5. Register routers, it tells beego which controller is responsibility for specific requests. Here we register `/` for `MainController`, so all requests in `/` will be handed to `MainController`. Be aware that the first argument is the path and the second one is pointer of controller that you want to register.

6. Run application in port 8080 as default, press `Ctrl+c` to exit.

### Going further

The [official website](http://beego.me) of beego contains all the information you may need, and the [documentation](http://beego.me/docs) page is made as an [open source project](https://github.com/beego/beedoc) so that you can open issues for requesting more details or pull request and be a part of the [community](http://beego.me/community). To follow up news of beego, keep your eyes on our [blog](http://beego.me/blog).


