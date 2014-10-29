+++
title = "Go Advent Day 13 - Go web services with Tiger Tonic"
date = 2013-12-13T06:40:42Z
author = ["Richard Crowley"]
series = ["Advent 2013"]
+++

## Welcome

Go is unique among mainstream programming languages in that its standard library web server is not a complete afterthought.  The Go language is well-suited for engineering complex networked services and Go's standard library recognizes that many (if not most) of those services communicate via HTTP.  Sprinkle some [Google scale](https://groups.google.com/forum/#!msg/golang-nuts/BNUNbKSypE0/E4qSfpx9qI8J) on it and your web applications and services can really hit the ground running.

The standard library sets the tone but it's far from the end of the story of how to effectively build web services in Go.  As most Go stories do, ours begins with an interface:

    type Handler interface {
        ServeHTTP(ResponseWriter, *Request)
    }

`http.Handler` may not look like much but its simplicity is the key to its universality.  The standard library provides a number of powerful implementations of this interface but most serious web applications and services are quick to find shortcomings, differences of opinion, and higher-level abstractions that require us to implement `http.Handler` ourselves.

### Introducing Tiger Tonic

[Tiger Tonic](https://github.com/rcrowley/go-tigertonic) packages several `http.Handler` implementations that make engineering web services easier.  They strive to be orthogonal like the features of the Go language itself.  Because everything is an `http.Handler` it's possible and easy to use Tiger Tonic in conjunction with your own or anyone else's handlers.

Tiger Tonic eschews HTML templates, JavaScript and CSS asset pipelines, cookies, and the like to remain squarely focused on building JSON web services.  Why JSON?  It's become the lingua franca for serialization in web services everywhere.  Why the focus on web services?  The authors' love [Dropwizard](http://dropwizard.codahale.com) and missed all it brought to the Java community so they worked to reproduce that feeling in Go.

### Multiplexing requests

Most web services respond differently to requests for different URLs so they reach for the standard `http.ServeMux`.  Then these web services start handling requests like `POST /games/{id}/bet` (an example from [Betable](https://developers.betable.com)) and `http.ServeMux` gets in the way.  [tigertonic.TrieServeMux](http://godoc.org/github.com/rcrowley/go-tigertonic#TrieServeMux) is here to help.

    mux := tigertonic.NewTrieServeMux()
    mux.Handle("POST", "/games/{id}/bet", handler)

`tigertonic.TrieServeMux` multiplexes requests differently than the standard `http.ServeMux`.  URL components surrounded by braces like `{id}` above are wildcards added to `r.URL.Query` and otherwise the request must match the method and path of the matching handler exactly.  In the author's experience this is the least surprising thing to do.  Plus it gives the framework the opportunity to respond 404 and 405 to requests it cannot fulfill.  See?  Your web service is already a good HTTP citizen!

And as cool as that is it's even cooler that your web service is still an `http.Handler` that you can `http.ListenAndServe` without a second thought.

### JSON in, JSON out

So now the question is how we respond to all the requests we're now multiplexing.  This is where we start to specialize on web services.  Once again, the standard library has giant shoulders on which to stand.  Import `encoding/json` and you're off to the races, right?  Well:

    func handlerFunc(w http.ResponseWriter, r *http.Request) {
        err := json.NewEncoder(w).Encode(&MyResponse{/*...*/})
        if nil != err {
            fmt.Fprintln(w, err)
        }
    }

That's a bit clumsy.  [tigertonic.Marshaled](http://godoc.org/github.com/rcrowley/go-tigertonic#Marshaled) allows you to construct handlers from functions that automatically deserialize request bodies according to the type of the final argument and serialize response bodies from the error or response object returned.

    var handler http.Handler = tigertonic.Marshaled(func(
       url.URL, http.Header, *MyRequest,
    ) (int, http.Header, *MyResponse, error) {
       return http.StatusOK, nil, &MyResponse{/*...*/}, nil
    })

This proves to be a powerful higher-level abstraction that removes JSON serialization from the web service programmer's long list of concerns.  Just as before, handlers returned by `tigertonic.Marshaled` understand how to respond 400, 406, and 415 like Roy Fielding intended.

### Testing

Perhaps the best feature of this new abstraction is its effects on testing web services.  No disrespect to the `net/http/httptest` package or `httptest.ResponseRecorder` but Go web services should test their responses not their serializations.  For example:

    func TestAdvent(t *testing.T) {
        s, _, rs, err := handler(
            mocking.URL(mux, "POST", "/games/ID/bet"),
            mocking.Header(nil),
            nil,
        )
        if nil != err {
            t.Fatal(err)
        }
        if http.StatusOK != s {
            t.Fatal(s)
        }
        if "win" != rs.Outcome { // Merry Christmas!
            t.Fatal(rs)
        }
    }

That's far more precise than matching strings with regular expressions and far less verbose than deserializing the JSON response.

For those keeping score at home: yes, `tigertonic.Marshaled` also returns an `http.Handler`.

### Extra batteries

Though `tigertonic.TrieServeMux` and `tigertonic.Marshaled` are the main attractions, Tiger Tonic packages a number of other useful handlers that make building web services with Go easier:

- [tigertonic.TrieServeMux](http://godoc.org/github.com/rcrowley/go-tigertonic#TrieServeMux.HandleNamespace) has a `HandleNamespace` method to match and remove prefixes from requested URLs.
- [tigertonic.HostServeMux](http://godoc.org/github.com/rcrowley/go-tigertonic#HostServeMux) supports virtual hosting of many domains in a single Go process.
- [tigertonic.First](http://godoc.org/github.com/rcrowley/go-tigertonic#First) and [tigertonic.If](http://godoc.org/github.com/rcrowley/go-tigertonic#If) enable handler chaining a la Rack or WSGI middleware.
- [tigertonic.Counted](http://godoc.org/github.com/rcrowley/go-tigertonic#Counted) and [tigertonic.Timed](http://godoc.org/github.com/rcrowley/go-tigertonic#Timed) emit metrics about all your requests via [go-metrics](https://github.com/rcrowley/go-metrics).
- [tigertonic.WithContext](http://godoc.org/github.com/rcrowley/go-tigertonic#WithContext) and [tigertonic.Context](http://godoc.org/github.com/rcrowley/go-tigertonic#Context) add support for strongly-typed per-request context.
- [tigertonic.HTTPBasicAuth](http://godoc.org/github.com/rcrowley/go-tigertonic#HTTPBasicAuth) is a specialization of `tigertonic.If` that conditionally handles requests if an acceptable `Authorization` header is present.
- [tigertonic.CORSBuilder](http://godoc.org/github.com/rcrowley/go-tigertonic#CORSBuilder) and [tigertonic.CORSHandler](http://godoc.org/github.com/rcrowley/go-tigertonic#CORSHandler) facilitate setting the basic CORS response headers.
- [tigertonic.Server](http://godoc.org/github.com/rcrowley/go-tigertonic#Server) has `CA` and `TLS` methods that simplify listening for TLS connections.
- [tigertonic.Configure](http://godoc.org/github.com/rcrowley/go-tigertonic#Configure), in conjunction with method values, makes it easy to read configuration files.

### Go forth

[Tiger Tonic](https://github.com/rcrowley/go-tigertonic) is available on GitHub an includes a complete [example](https://github.com/rcrowley/go-tigertonic/tree/master/example) that covers all the handlers included.

`http.Handler` is the common currency for Go web frameworks of all shapes and sizes.  The handlers in the `tigertonic` package are meant to make engineering web services more correct, more efficient, and more testable but because they're handlers they're at your service and up to the challenge just as long as you stick to the humble little `http.Handler` interface.
