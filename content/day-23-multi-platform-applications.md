+++
title = "Go Advent Day 23 - Multi-Platform Applications: Architecture and Cross-Compilation"
date = 2013-12-23T06:40:42Z
author = ["Mitchell Hashimoto"]
series = ["Advent 2013"]
+++

## Introduction

While Go is touted for its utility on the server side and in networked
environments, Go is incredibly powerful as a client-side (desktop)
application language as well.

An often unknown feature of Go is that it is more or less completely
portable: you can compile your Go code to run on any other operating
system Go supports from the comfort of your own operating system. In
addition to this, Go has build constraints to control which files are compiled
under what conditions, allowing you to write OS-specific code and still have
your application compile.

In this post, I'm going to show you how to effectively write Go code
that runs on multiple platforms by taking advantage of all the features
Go has to offer in this area.

### Interfaces, interfaces, interfaces!

The key to easily supporting alternate operating systems is to make
copious use of interfaces in places where operating system specific behavior
might exist. One example is network connections: don't use a `net.TCPConn`
if a `net.Conn` will suffice (or, for that matter, an `io.ReadWriteCloser`).

This is important because if your application makes use of client/server
communication, maybe you can use Unix domain sockets on BSD or Linux. Of
course, Unix domain sockets don't exist on Windows. If you're just using
a `net.Conn`, it doesn't matter.

On Windows, maybe you can use [named pipes](http://msdn.microsoft.com/en-us/library/windows/desktop/aa365590(v=vs.85).aspx)
(which Go doesn't support out of the box at the moment). Again, you'll
need the `net.Conn` or `io.ReadWriteCloser` abstraction to make this work
properly.

Of course, interface usage could be a whole blog post on its own, so without
going into detail: don't go overboard on interfaces. Use them where they make
sense. But one of the places where they make sense is separating complex
behaviors from their actual implementation, and distilling their core
operations (the functions of the interface).

### Checking the OS at run time

While perhaps obvious, Go provides standard library functions for determining
the running operating system at run time. This lets you easily change minor
behaviors based on the OS. For major behaviors, build constrains and most
likely interfaces should be used. Both are covered in their own sections.

From the `runtime` package, you can use `runtime.GOOS` to determine the
operating system that the binary was compiled for (which should be the only
operating system it is running on). This allows you to do basic switches:

    func RootDrive() string {
        if runtime.GOOS == "windows" {
            return "C:/"
        } else {
            return "/"
        }
    }

If used too much, it can get confusing what code runs and where, so try to
limit the OS-specific branches. Or, if the behavior is large enough, pull
it out into an interface. Or, finally, maybe a build constraint makes sense.

### Build Constraints

A build constraint is a line comment in a Go file that lists the conditions
under which a file should be included in a package. One of the possible
constraints is target operating system, allowing you to selectively include
and exclude certain files based on the target operating system of the package.

This shouldn't be used to exclude any and all code used
for only a single operating system. It should only be used to include or
exclude code that will only compile on a single operating system.

For example: don't use build constraints two implement to versions of a
function `ConvertSlashes` that converts file path slashes to the proper
direction based on the operating system. Instead, use an if statement
on `runtime.GOOS`. This makes finding bugs much easier, writing tests easier,
and keeps cognitive overhead to a minimum when determining what code
is running where.

Build constraint syntax and available options are covered exhaustively in
the [go build docs](http://golang.org/pkg/go/build/), but a brief example
is shown here.

The most common idiom with code in OS-constrained files is to export the
callable function in the platform-independent version of the file, and use
a private function in the platform-dependent files. For example, in
[Packer](http:/www.packer.io), we have a function that returns the path
to where we should put the configuration file. On Unix, this is as a dot file
in the home folder. On Windows, this is in the application settings folder,
which can only be determined by making some DLL calls. We use build constraints
to control whether the DLL calls are included in the package, since they
won't compile on Unix.

To do this, we first have a non-constrained file
[configfile.go](https://github.com/mitchellh/packer/blob/7c9c7afd828b63211d0586e5ed032d21e20d042a/configfile.go). This exposes a public function `ConfigFile`. As you can
see, though, it simply calls a private function `configFile`. This private
function is then implemented in OS-constrained files:
[configfile_unix.go](https://github.com/mitchellh/packer/blob/7c9c7afd828b63211d0586e5ed032d21e20d042a/configfile_unix.go) and
[configfile_windows.go](https://github.com/mitchellh/packer/blob/7c9c7afd828b63211d0586e5ed032d21e20d042a/configfile_windows.go).

If you look in those files, you'll see the build tags at the top of the
file controlling when they are compiled. In the Windows file, you'll see
we make use of standard functions that are only available when compiling for Windows,
such as `syscall.MustLoadDLL`.

Build constraints are extremely useful in separating out platform-specific
code, but can also make following the direction of your code confusing,
so use them carefully.

Also, take note that build constraints can also be used in test files! This
lets you write platform-specific tests if you need to. You should do this if
you have platform-specific implementations of functions.

### Cross Compilation

The final piece of the puzzle is being able to compile your Go application
for multiple operating systems with ease. For many other languages, you either
have to compile directly on the target operating system, or you must follow
an almost impossibly complex process to build a cross-compilation toolchain.
With Go, cross compilation is available right out of the box.

To make cross-compilation a little bit nicer, I recommend using
[Gox](https://github.com/mitchellh/gox). Gox mimics `go`build` in usage
but will compile for multiple platforms in parallel.

But to see the raw bits of how cross compilation in Go works, I
recommend reading [Dave Cheney's blog post on it](http://dave.cheney.net/2013/07/09/an-introduction-to-cross-compilation-with-go-1-1).

To install Gox, just `go`get` it:

    go get github.com/mitchellh/gox

Once it is installed, you'll have to build the toolchains so you can cross
compile. You only need to do this once per Go version:

    gox -build-toolchain

That will take some time, but once it is complete, you're ready to cross
compile! Just run `gox` (just like you would just run `go`build`) in the
directory of your application, and it'll build the application for every
platform that your version of Go supports! The output should look like the
following:

    $ gox
    Number of parallel builds: 4

    -->      darwin/386: github.com/mitchellh/gox
    -->    darwin/amd64: github.com/mitchellh/gox
    -->       linux/386: github.com/mitchellh/gox
    -->     linux/amd64: github.com/mitchellh/gox
    -->       linux/arm: github.com/mitchellh/gox
    -->     freebsd/386: github.com/mitchellh/gox
    -->   freebsd/amd64: github.com/mitchellh/gox
    -->     openbsd/386: github.com/mitchellh/gox
    -->   openbsd/amd64: github.com/mitchellh/gox
    -->     windows/386: github.com/mitchellh/gox
    -->   windows/amd64: github.com/mitchellh/gox
    -->     freebsd/arm: github.com/mitchellh/gox
    -->      netbsd/386: github.com/mitchellh/gox
    -->    netbsd/amd64: github.com/mitchellh/gox
    -->      netbsd/arm: github.com/mitchellh/gox
    -->       plan9/386: github.com/mitchellh/gox

In my case, Gox parallelized the builds 4-ways because my computer has
4 cores. Your parallelization factor might be different, but the end result
is the same: your application is cross-compiled!

If you inspect some of the files made, you can prove to yourself that
they're for other platforms. For example, when cross-compiling Gox itself:

    $ file gox_windows_386.exe
    gox_windows_386.exe: PE32 executable for MS Windows (console) Intel 80386 32-bit

    $ file gox_plan9_386
    gox_plan9_386: Plan 9 executable, Intel 386

    $ file gox_openbsd_amd64
    gox_openbsd_amd64: ELF 64-bit LSB executable, x86-64, version 1 (OpenBSD), statically linked, for OpenBSD, not stripped

You can also limit the platforms you want to build for by using the `-os`
and `-arch` flags. See the
[Gox README](https://github.com/mitchellh/gox/blob/master/README.md) for
full documentation on usage.

### Final Thoughts

In my opinion, the ease of cross-compilation with Go rivals the portability
of languages like Java, while still allowing you to easily dig into the
platform-specific bits if you need to, and without having to ask your end users
to install a large runtime.

And compared to C, a highly portable language for sure, getting the same
complex application to compile across multiple platforms is much simpler, and
getting the cross-compilation toolchain functioning is incredibly easier.

As an ending note, I often get asked of my thoughts on using Go for GUI-based
desktop applications. Personally, I think GUIs never feel quite right unless
they're written specifically for their target platforms. Therefore, I recommend
using the native language/toolkit for building the GUI application, but
putting all the complex logic into a Go applicatin for portability. This works
for almost all types of applications except perhaps games. This lets you
test and share the complex logic of an application across multiple platforms,
while getting a native look for your application as well. This itself could
be a blog post so I'll just end there!
