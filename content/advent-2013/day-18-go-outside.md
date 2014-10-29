+++
title = "Go Advent Day 18 - Go Outside"
date = 2013-12-18T06:40:42Z
author = ["Tony Wilson"]
series = ["Advent 2013"]
+++


## Introduction

[Outside](http://github.com/tHinqa/outside) is a Go package to dynamically link to and execute functions in Windows Dynamic Link Libraries and Linux Shared Libraries.
Its current status is 'prerelease' with only 32-bit register size implemented and tested so far.
Also, some functionality is *very* experimental and will probably change a lot before release 0.1.

I only came to explore Go as a viable general purpose or glue language in the 2nd quarter of 2013.
Not so much because of its (at the time) fairly low profile and newness but more because I had no real need to change from my home-grown developer-base-of-one Prolog system.
So, realising the amount of effort it would take to roadworthy the beast for general consumption (typical WORN code - Write Once Read Never), I decided to rather tag along with Go for future works.

My needs almost always involve interacting with external (hence `Outside`) functions in proprietary or open-source Windows `.dll` or Linux `.so` files.
`Cgo` is a perfectly capable solution but I did find the reliance on `gcc` a bit detracting from the scripting look and feel of pure Go.

Tackling the (initially) rather bewildering [reflect](http://golang.org/pkg/reflect) package turned out to be the reward I needed.

[syscall](http://golang.org/pkg/syscall) on Windows already has ```[Must]LoadDLL```, `[Must]FindProc` and `Call` to load libraries, find entry-points and call procedures.
Note: Because these functions are Windows-specific, they don't show up in documentation on [golang.org](http://golang.org).
See godoc and [syscall](http://golang.org/pkg/syscall) on reading documentation for other systems.
Linux provides more limited (6 arguments max) [Syscall](http://golang.org/pkg/syscall/#Syscall) functions, but neither dynamic loading nor lookup. Also, the neater `Call` is missing. In `outside` similar Linux functionality has been implemented using the `dl` library and some C code.

All that was missing (on Windows) was a method of "black-boxing" the otherwise messy [unsafe.Pointer](http://golang.org/pkg/unsafe/#Pointer) conversions and type coercion necessary to issue a `Call` (or its underlying `Syscall`).
For this we have (drum roll please) the magical [reflect.MakeFunc](http://golang.org/pkg/reflect/#MakeFunc).
Within a MakeFunc'd function we can transform the input arguments, do the call and manipulate the return in any way desired.

### Enough already, show me an example

Suppose we want to 'stylize' an image using ImageMagick (like [William Kennedy did on day 9](day-09-building-a-weather-app-using-go))

Starting with a rather dental portrayal of our `favourite` rodent...

![](/postimages/day-18-go-outside/i.png)

we want to 'weather' him a bit...

![](/postimages/day-18-go-outside/o.jpeg)
 
(Ok, I'm from the sunny south of the equator, so sleet is a bit tongue-in-cheek - and I never miss a chance to take a dig at Microsoft for forcing the spelling `favorite` down the international throat)

### Start with the your choice of prototypes

Identify the C prototypes for the functions we will need.

      MagickWand *NewMagickWand(void);
      MagickWand *DestroyMagickWand(MagickWand *wand);
      MagickBooleanType MagickReadImage(MagickWand *wand, const char *filename);
      MagickBooleanType MagickSketchImage(MagickWand *wand,
          const double radius,const double sigma, const double angle);
      MagickBooleanType MagickEqualizeImage(MagickWand *wand)
      MagickBooleanType MagickWriteImage(MagickWand *wand, const char *filename);

### Set up the Go function variables and types

Do the usual usual type-name to name-type transposition, name simplification and casing that we all love.

      type Wand struct{}

      var New func() *Wand
      var Destroy func(m *Wand) *Wand
      var Read func(m *Wand, filename string) bool
      var Sketch func(m *Wand,
          radius, sigma, angle float64) bool
      var Equalize func(m *Wand) bool
      var Write func(m *Wand, filename string) bool

### Relate the Go functions to the dll entry-points

      var allApis = outside.Apis{
          {"NewMagickWand", &New},
          {"DestroyMagickWand", &Destroy},
          {"MagickReadImage", &Read},
          {"MagickSketchImage", &Sketch},
          {"MagickEqualizeImage", &Equalize},
          {"MagickWriteImage", &Write},
      }

### Link it all together

This code does all the setup needed.
The DLL or SO is loaded and the entry-points are resolved in a soft fashion - that is, there is no hard failure at this point for any missing names, they are merely logged.
The bool 2nd argument to `AddDllApis` determines whether `string` types are converted to char string `*char` (false), or UTF-16 string `*uint16` (true) types.
The Linux `dll` name is definitely distribution dependent (I used the current Mint Debian Linux).
You can get pre-built binary releases from [ImageMagick](http://www.imagemagick.org/script/binary-releases.php)
For a `.dll` or `.so` to be loaded dynamically, its directory needs to be in the `path` when the program is run (not necessary at compile time).
After init, functions are immediately executable.
A call to an unresolved entry-point will panic.

      func init() {
          var dll string
          if runtime.GOOS == "windows" {
              dll = "CORE_RL_wand_.dll"
          }
          if runtime.GOOS == "linux" {
              dll = "libMagickWand.so.5"
          }
          outside.AddDllApis(dll, false, allApis)
      }

### Add your own Go flavour

Although you can use the functions directly, there is no harm in a bit of Go idiomatic embelishment.

      func (m *Wand) Read(filename string) *Wand {
          Read(m, filename)
          return m
      }
      func (m *Wand) Sketch(radius, sigma, angle float64) *Wand {
          Sketch(m, radius, sigma, angle)
          return m
      }
      func (m *Wand) Equalize() *Wand {
          Equalize(m)
          return m
      }
      func (m *Wand) Write(filename string) *Wand {
          Write(m, filename)
          return m
      }

### Formulate the main plan of action

      func main() {
          m := New()
          defer Destroy(m)
          m.Read("i.png").Sketch(0, 30, 60).Equalize().Write("o.jpeg")
      }

### Wait a minute - what about error handling?

Since we're treating the bool return sequence as a single entity we could 

      type Bool bool
      var ... func(...) (Bool, error)

with the optional `Error` method on `Bool`

      func (ok Bool) Error func(e error) (bool, error) {
          if ok {
              return ok, nil
          }
          if e != nil {
               panic("Failed with: " + e.Error())
          }
          panic("Failed with unknown error")
      }

### Throw in the Go housekeeping header and run

      package main

      import (
          "github.com/tHinqa/outside"
          "runtime"
      )

[Complete code](day-18-go-outside/example.go)

Happy holidays and Go safely.
