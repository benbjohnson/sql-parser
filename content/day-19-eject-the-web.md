+++
title = "Go Advent Day 19 - Eject the Web"
date = 2013-12-19T06:40:42Z
author = ["Yasuhiro Matsumoto"]
series = ["Advent 2013"]
+++

_Editors Note:_ Yasuhiro is not a native English speaker, so during the editing of this post is was necessary to make some minor corrections. We felt that it was very important however, that the Author's original phrasing and intent be preserved as much as possible.

## Introduction

As you already know, Go is a good programming language to write web applications.

There are already many packages providing routers, MVC, and Sinatra-like frameworks. However, Go is also a good programming language to write command-line applications. I have written many small programs written in Go, and I believe Go is the best for this task.

## OS Compatibility

Most of Go's packages are beautiful for beginners to see, because Go hides the operating system's peculiar matters as much as possible. For example: file system differences, character encodings, socket I/O, concurrency. You'll notice this is a feature of Go.

For example, `path/filepath` provides a way to manipulate file paths in an OS independent manner. This code will work fine for both POSIX Operating Systems and Windows.

    package main

    import (
        "fmt"
        "log"
        "os"
        "path/filepath"
    )

    func main() {
        pwd, err := os.Getwd()
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(filepath.ToSlash(
            filepath.Join(append([]string{pwd}, os.Args[1:]...)...)))
    }


`path/filepath` treats file paths as simple UTF-8 strings. All of the APIs for file systems are doing the proper encoding/decoding for their environment. `os.Args` is encoded to UTF-8 even though it's on Windows. So Windows user can write programs without entering the hell of multi-byte encoding.

In Japanese, Windows applications use Shift_JIS (MS932) file system encoding. Shift_JIS is a double-byte character-set in which some characters contain `0x5C` (a.k.a backslash) in the trailing-byte. So many Windows application developers meet a hell of multi-byte encoding.

    C:\path\contains\backslash\char>gcc -o unix_like_app.c
    
    C:\path\contains\backslash\char>unix_like_app.exe

![](/postimages/day-19-eject-the-web/image1.png)

but Go can make them happy.

### Concurrency Patterns

I speculate that Go will provide a new style of UI coding.

Traditionally it is difficult to handle both key-press events and input timeouts at once. However, if you separate key event handling into a goroutine and pass events to a channel external to the goroutine, it may make complex UI coding simpler.
	
    package main

    import (
        "fmt"
        "github.com/mattn/go-runewidth"
        "github.com/nsf/termbox-go"
        "log"
        "time"
    )

    func print_tb(x, y int, msg string) {
        for _, c := range []rune(msg) {
            termbox.SetCell(x, y, c, termbox.ColorWhite, termbox.ColorDefault)
            x += runewidth.RuneWidth(c)
        }
        termbox.Flush()
    }

    func main() {
        err := termbox.Init()
        if err != nil {
            log.Fatal(err)
        }
        termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

        event := make(chan termbox.Event)
        go func() {
            for {
                // Post events to channel
                event <- termbox.PollEvent()
            }
        }()

        print_tb(1, 1, "Hit any key")
    loop:
        for {
            // Poll key event or timeout
            select {
            case ev := <-event:
                print_tb(1, 2, fmt.Sprintf("Key typed: %v", ev.Ch))
                break loop
            case <-time.After(5 * time.Second):
                print_tb(1, 2, "Timeout")
                break loop
            }
        }
        close(event)
        time.Sleep(1 * time.Second)
        termbox.Close()
    }


If you create a channel for the events, you will be able to handle other events in the `select` easily.

- Receiving new tweet
- Redraw window event
- Key events

When typing keys, in most cases, internal conditions will be changed and the UI will need to be redrawn. But the rendering may be expensive to call for each of key press,and we do not want to block other operations. A solution is to use a goroutine and a timer.

    dirty := false
    timer := time.AfterFunc(0, func() {
        if dirty {
            redraw_full()
        } else {
            redraw_part()
        }
    })
    for {
        select {
        case ev := <-event:
            // handle terminal event
            switch ev.Type {
            case termbox.EventKey:
                // handle key event
                switch ev.Key {
                case termbox.KeyEnter:
                    // update internal condition
                    update_condition(ev)
                    dirty = false
                    // redraw immediately
                    timer.Reset(1 * time.Microsecond)
                default:
                    // update candidate
                    update_candidate(ev)
                    dirty = true
                    // redrwa in later
                    timer.Reset(200 * time.Microsecond)
                }
            }
            // handle another events
        }
    }

The `timer` makes a delay for redrawing screen. If the next event requires forcing a screen update, our redraw goroutine will be run immediately.

Of course, the above code works fine for Windows. You don't need to do things the hard way. Writing pthread codes? Writing multiple event polling? No. Do it the easy way by using Go.

### But... I must write OS Specific Code

For example, you want to eject the CD-ROM tray in Go. Then you can use [go-eject](https://github.com/mattn/go-eject).

As you can see, Go can separate operating system specific code using a file naming convention.

- [`eject_windows.go`](https://github.com/mattn/go-eject/blob/master/eject_windows.go)
- [`eject_linux.go`](https://github.com/mattn/go-eject/blob/master/eject_linux.go)

Another way to separate specific OS code is to add magic into the top of the file.

    // +build !windows

This mean "Build this code unless compiling for Windows".

Go provides an easy way to write bindings. If you already have C code, and you want to use it, you can use CGO. For example, `eject_linux.go` uses CGO.

.code day-19-eject-the-web/eject_linux.go

If you don't want to mix C code in Go code, you can also separate C code to another file. The Go tool will build it perfectly. 

For Windows, you can call functions in DLLs dynamically.

    package eject

    /*
       #include <fcntl.h>
       #include <linux/cdrom.h>
       #include <sys/ioctl.h>
       #include <sys/stat.h>
       #include <sys/types.h>
       #include <unistd.h>
       #include <errno.h>

       static int
       _eject(int f) {
         int r = 0;
         int fd = open("/dev/cdrom", O_RDONLY | O_NONBLOCK);
         if (fd == -1) {
           r = errno;
         } else {
           int e = CDROMEJECT;
           if (f == 0) {
             if (ioctl(fd, CDROM_DRIVE_STATUS, 0) == CDS_TRAY_OPEN)
               e = CDROMCLOSETRAY;
           } else if (f == 1)
             e = CDROMEJECT;
           else
             e = CDROMCLOSETRAY;
           if (ioctl(fd, e, 0) < 0) {
             r = errno;
           }
           close(fd);
         }
         return r;
       }
    */
    import "C"
    import "errors"
    import "syscall"

    func Eject() error {
        if r := C._eject(0); r != 0 {
            return errors.New(syscall.Errno(r).Error())
        }
        return nil
    }


Note that you're better off using the wide string APIs on Windows.

As I mentioned before, Go is a very nice language for Windows programmers. If you write a package wrapping OS specific issue in one place, most of users doesn't need to worry about compatibility. i.e. many users can eject CD-ROM with by using `eject.Eject()`, regardless of their operating system.

### Let's Go Web

Finally, you can write a full application using [go-eject](https://github.com/mattn/go-eject) and [web.go](https://github.com/hoisie/web) to eject your cd from the Web.

    package main

    import (
        "encoding/json"
        "github.com/hoisie/web"
        eject "github.com/mattn/go-eject"
    )

    type result struct {
        Error interface{} `json:"error"`
    }

    func main() {
        web.Get("/eject", func(ctx *web.Context) {
            ctx.ContentType("application/json")
            if err := eject.Eject(); err != nil {
                json.NewEncoder(ctx).Encode(&result{err.Error()})
            } else {
                json.NewEncoder(ctx).Encode(&result{nil})
            }
        })
        web.Run(":8080")
    }


Add in `static/index.html`

    <!DOCTYPE html>
    <html>
    <head>
    <meta charset="UTF-8">
    <title>Eject the Web</title>
    <script src="http://code.jquery.com/jquery-2.0.3.min.js"></script>
    <script>
    $(function() {
      $('#eject-the-web').click(function() {
        $.ajax("/eject", function(res) {
    	  if (res.error) {
            alert(res.error);
          }
        });
      });
    })
    </script>
    </head>
    <body>
      <input type="button" value="Eject the Web">    
    </body>
    </html>

Click "Eject the Web". Have fun.

![](/postimages/day-19-eject-the-web/push-the-button.png)
