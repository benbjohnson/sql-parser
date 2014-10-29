+++
title = "Go Advent Day 4 - GoConvey"
date = 2013-12-04T06:40:42Z
author = ["Matthew Holt", "Michael Whatcott"]
series = ["Advent 2013"]
+++


## Introduction

One highly appealing aspect about Go is built-in testing with `go test`. From one who once eschewed test-driven development, I now wholly encourage it. [Testing is fundamental](http://golang.org/pkg/testing/) to writing Go code, and Go 1.2's [new test coverage tools](http://blog.golang.org/cover) make TDD more compelling than ever.

## Introducing GoConvey

[GoConvey](http://smartystreets.github.io/goconvey) is a new project that makes testing even better in Go. It consists of (1) a framework for writing behavioral-style tests, and (2) a web UI which reports test results in real-time. Both are optional, depending on your own workflow.

![](http://d79i1fxsrar4t.cloudfront.net/goconvey.co/gc-1-light.png)

The web interface works with any tests that run with go test. In other words, it fully supports native Go tests using the standard `testing` package, and it should play nicely with other frameworks like gocheck, Testify, and Ginkgo.

The browser page updates automatically when Go files are changed. You can also enable desktop notifications so you don't have to leave your editor:

![](http://d79i1fxsrar4t.cloudfront.net/goconvey.co/gc-7-light.png)

At its core, though, GoConvey is a testing framework. Its simple DSL is designed so that code is not only tested, but also documented clearly and unambiguously:

	package main

	import (
		. "github.com/smartystreets/goconvey/convey"
		"testing"
	)

	func TestIntegerExample(t *testing.T) {
		Convey("Subject: Integer incrementation and decrementation", t, func() {
			var x int

			Convey("Given a starting integer value", func() {
				x = 42

				Convey("When incremented", func() {
					x++

					Convey("The value should be greater by one", func() {
						So(x, ShouldEqual, 43)
					})
					Convey("The value should NOT be what it used to be", func() {
						So(x, ShouldNotEqual, 42)
					})
				})
			})
		})
	}

The web UI comes with a tool to stub this out for you without having to worry about Go syntax.

If you prefer the terminal over a browser, just run go test as usual:

![](http://d79i1fxsrar4t.cloudfront.net/goconvey.co/terminal.png)

## Assertions

GoConvey provides _positive_ assertions, which is different from traditional Go tests. Quick:tell me what the expected results are, not what they aren't: 

	func TestBoring(t *tesing.T) {

		if foo < bar || baz != "abc" {
			t.Error("Test failed...")
		}

	}

(Are you right? Keep reading.) This statement tells us what the code _shouldn’t_ do, rather than telling us what it _should_ do. It would be much easier to understand if we replace it with two clear `So` assertions:

	func TestSmart(t *testing.T) {
		
		So(foo, ShouldBeGreaterThanOrEqualTo, bar)
		So(baz, ShouldEqual, "abc")
		
	}

Did you read the original test correctly? Most people forget that `foo` can _equal_ `bar` and still pass.

## Summary

GoConvey is a different way of thinking about testing in Go. GoConvey is verbose, whereas Go convention is concise and even abbreviated. But Go purists may still find the web UI appealing while sticking to the standard testing package. (And let’s remind ourselves that the test files and packages, with all their verbosity, are not compiled into the production binaries.)

The framework and the web UI are independent; if the framework isn't your thing, try the web UI (and vice versa). This is designed to be a drop-in convenience around existing tests.

If you would like to contribute new features or fixes, please be our guest. Open an issue on [GitHub](https://github.com/smartystreets/goconvey) to discuss it, and remember: "a pull request, with test, is best."
