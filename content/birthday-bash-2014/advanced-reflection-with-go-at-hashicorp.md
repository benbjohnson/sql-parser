+++
author = ["Mitchell Hashimoto"]
date = "2014-11-26T00:00:00-08:00"
title = "Advanced Reflection with Go at HashiCorp"
series = ["Birthday Bash 2014"]
+++

# Advanced Reflection with Go at HashiCorp

[HashiCorp](https://www.hashicorp.com) builds a diverse set of
popular DevOps tools written in Go:
[Packer](https://www.packer.io),
[Serf](https://www.consul.io),
[Consul](https://www.consul.io),
and [Terraform](https://www.terraform.io).
While these tools range from desktop software to highly scalable distributed
systems, their internals all have one thing in common: they use reflection,
and a lot of it. In this post, I'll share the libraries and techniques we
use at HashiCorp to get the most out of reflection, all while being safe
and efficient.

Reflection is a powerful tool that can enable some beautiful
functionality, but the downside is that it isn't particularly efficent
and a misuse of the standard reflection library almost always results in
a panic. Because of this, HashiCorp has developed a set of core libraries
on top of the standard `reflect` package in order to safely use reflection
throughout our tools.

The three open source libraries we use heavily are:
[mapstructure](https://github.com/mitchellh/mapstructure),
[copystructure](https://github.com/mitchellh/copystructure), and
[reflectwalk](https://github.com/mitchellh/reflectwalk). We'll go over
each of these libraries: what they do, when they should be used, and where
they're used in HashiCorp tools.

## mapstructure

[mapstructure](https://github.com/mitchellh/mapstructure) is a library
for decoding generic map values to Go structures. In other words, given
any `interface{}` and a pointer to another Go value, mapstructure does
its best to decode the value to into the destination Go value.
Most commonly, mapstructure is used to take a `map[string]interface{}`
and decode it into a Go struct.

**Why?** There are a lot of times, especially with APIs and configuration,
where you receive a value in some format, and the final struct that value
should be decoded to is variable. Go's standard decoding libraries for
JSON, XML, etc. do not handle this case particularly well. mapstructure
gives you a way to cleanly take any of those formats and decode it into
a specific structure determined at runtime.

An example explains this best:

```go
// A structure we want to decode into.
type Person struct {
	Name   string
	Age    int
	Emails []string
	Extra  map[string]string
}

// This input can come from anywhere, but typically comes from
// something like decoding JSON where we're not quite sure of the
// struct initially.
input := map[string]interface{}{
	"name":   "Mitchell",
	"age":    91,
	"emails": []string{"one", "two", "three"},
	"extra": map[string]string{
		"twitter": "mitchellh",
	},
}

var result Person
if err := mapstructure.Decode(input, &result); err != nil {
	panic(err)
}

fmt.Printf("%+v", result)
```

mapstructure exposes a very developer-friendly API, and protects you from
all the edge cases of reflection that result in panics, instead returning
an `error` when types mismatch or invalid input is given.

We've found that mapstructure is one of those libraries where at first you look
at it and say: "I'm not sure why I would ever need that." But then within
a few months it quickly becomes a library where you say: "How did I ever
write tools without mapstructure?"

We use mapstructure in every one of our tools at the minimum for
configuration but many times for APIs as well.

## copystructure

[copystructure](https://github.com/mitchellh/copystructure) is a library for
performing deep copies of values in Go.

This library usually doesn't need a "why?" as the utility is clear: sometimes
you need a deep copy of a structure. Perhaps you're trying to return a modified
map without modifying the original value, or you're trying to duplicate a
structure for multiple function calls. Whatever the case, deep copying
is something that comes up from time to time.

Usage is simple:

```go
input := map[string]interface{}{
	"bob": map[string]interface{}{
		"name":   "bob",
		"emails": []string{"a", "b"},
	},
	"jane": map[string]interface{}{
		"name": "jane",
	},
}

dup, err := copystructure.Copy(input)
if err != nil {
	panic(err)
}

fmt.Printf("%#v", dup)
```

Simple, easy to use API, powered by `reflect`.

We use this within [Terraform](https://www.terraform.io) in order to
duplicate full config objects (which contain many nested maps, slices,
etc.).

## reflectwalk

[reflectwalk](https://github.com/mitchellh/reflectwalk) uses reflection
to "walk" any value in Go, calling callbacks along the way for each component
of the structure. In other terms, it turns reflection into something similar
to the visitor pattern.

**Why?** We've found that a callback-based approached to reflection is a powerful
and natural way to implement a lot of reflection-powered code. In fact, it
is the underlying library for implementing `copystructure`.

Usage:

```go
var walker MyWalker

value := getComplexValue()
if err := reflectwalk.Walk(value, walker); err != nil {
	panic(err)
}
```

In this case, `MyWalker` is a struct that implements one or more of the
interfaces in copystructure. Depending on the interfaces implemented,
copystructure will invoke various callbacks on the walker for different
events while traversing the `value`.

We use this within [Terraform](https://www.terraform.io) in order to
implement our [interpolations](https://terraform.io/docs/configuration/interpolation.html).
Interpolations work by reflectwalking the configuration structure,
finding all string values, and then parsing them for interpolations.

## Great Power

With great power comes great responsibility. Reflection can enable
some really great use cases that have historically been difficult to
achieve with a static language like Go. On the other hand, reflection
is quite slow (relative to knowing the types of ahead of time) and introduce
a layer of complexity at runtime that can result in crashes if you're not
careful. However, we've developed these libraries in order to make reflection
easier to use where it makes sense, and we make heavy use of them in every
one of our projects.

It is a great testimony to Go that it is able to create performant static
binaries while also being flexible enough to enable some dynamic behavior.
We love Go at [HashiCorp](https://www.hashicorp.com) and look forward to
our long future with it.
