+++
title = "Go Advent Day 6 - Service Discovery with etcd"
date = 2013-12-06T06:40:42Z
author = ["Andrew Bonventre"]
series = ["Advent 2013"]
+++

## Introduction

At [Poptip](https://poptip.com/), our first foray into Go was a small but critical service that required extremely high throughput for a non-trivial amount of text processing. Skeptical at first, I remember having a conversation with a [friend](http://robert.sesek.com/) who was raving about how much he had been enjoying Go, and noticed that some other [very smart people](http://blog.golang.org/building-stathat-with-go) had chosen to bet their entire companies on the language. After writing a few benchmarks, we were very happy with the results and confident to move forward with the project.

One of the many things that we didn’t want to overcomplicate early on in the process was deployment, so we made it dead simple. Once a commit to master is made, our continuous integration server pulls, compiles, and runs the tests. If everything is green, it scps the binary and a run script to a set of hosts listed in a text file checked into master. Once that’s complete, it executes the run script which will kill the existing process and start up the new binary. The time from push to live in production was less than 2 minutes, which made rollbacks simple: just revert the change and push.

This process was fine for a single service, but we’ve since moved the majority of our stack to Go, and it’s become suboptimal. Many services share the same config flags, which causes a headache when we rotate a database credential but forget to update one of the service’s config. Additionally, our infrastructure has become mature enough that autoscaling has become a crucial need in order for us to handle traffic spikes while minimizing cost overhead. Updating a text file in the repo just wasn’t cutting it anymore.

Enter [etcd](https://github.com/coreos/etcd#etcd), “[a] highly-available key value store for shared configuration and service discovery.”

The README has wonderful documentation on the basics of the service, but I’d like to go into some specific features that address some of the problems outlined above: host management and storing configuration parameters.

### Host Management

Let’s get started by firing up an etcd process we can use to play around with.

You can either [compile the binary from source](https://github.com/coreos/etcd#building), or [download the latest release from GitHub](https://github.com/coreos/etcd/releases/). I’ll wait here while you get things set up.

While etcd is highly available — a crucial feature for Poptip’s mission-critical systems — we won’t go into the details of multi-node clusters or failover scenarios. Let’s just start up a single node called dot:

	 $ ./etcd -data-dir dot -name dot
	 [etcd] Dec  4 16:28:39.451 INFO      | etcd server [name dot, listen on 127.0.0.1:4001, advertised url http://127.0.0.1:4001]
	 [etcd] Dec  4 16:28:39.451 INFO      | raft server [name dot, listen on 127.0.0.1:7001, advertised url http://127.0.0.1:7001]

The reasoning behind the name dot is left as an exercise for the reader, but it isn’t important for these set of exercises. You can name it whatever you like.

Let’s say we have a set of web frontend servers that we want to keep track of. Upon startup, they can make a simple curl request to register themselves within a directory in etcd. Let’s add two servers within the frontends/ directory:

	 $ curl -L http://127.0.0.1:4001/v2/keys/frontends/fe1 -XPUT -d value=10.0.1.1
	 {"action":"set","node":{"key":"/frontends/fe1","value":"10.0.1.1","modifiedIndex":2,"createdIndex":2}}
	 $ curl -L http://127.0.0.1:4001/v2/keys/frontends/fe2 -XPUT -d value=10.0.1.2
	 {"action":"set","node":{"key":"/frontends/fe2","value":"10.0.1.2","modifiedIndex":4,"createdIndex":4}}

Now we can list all the frontends within that directory:

	 $ curl -L http://127.0.0.1:4001/v2/keys/frontends/
	 {"action":"get","node":{"key":"/frontends","dir":true,"nodes":[{"key":"/frontends/fe1","value":"10.0.1.1","modifiedIndex":3,"createdIndex":3},{"key":"/frontends/fe2","value":"10.0.1.2","modifiedIndex":4,"createdIndex":4}],"modifiedIndex":2,"createdIndex":2}}

You’ll notice at this point that each call returns a very nice JSON-encoded response. This makes integrating with clients extremely easy (especially if you’re using [encoding/json](http://golang.org/pkg/encoding/json)). That said, there is a client library built for Go by the same [team](http://coreos.com/) that built etcd called [go-etcd](https://github.com/coreos/go-etcd).

This is all well and good, but what about if a host or process dies? They’ll be this zombie entry in there forever, right? Not with key TTLs. Let’s set a TTL of 5 seconds for one of our web frontends:

	 $ curl -L http://127.0.0.1:4001/v2/keys/frontends/fe2 -XPUT -d value=10.0.1.2 -d ttl=5
	 {"action":"set","node":{"key":"/frontends/fe2","prevValue":"10.0.1.2","value":"10.0.1.2","expiration":"2013-12-04T16:56:54.123531985-05:00","ttl":5,"modifiedIndex":5,"createdIndex":5}}

Now we wait for at least 5 seconds and query the directory again:

	 $ curl -L http://127.0.0.1:4001/v2/keys/frontends/
	 {"action":"get","node":{"key":"/frontends","dir":true,"nodes":[{"key":"/frontends/fe1","value":"10.0.1.1","modifiedIndex":3,"createdIndex":3}],"modifiedIndex":2,"createdIndex":2}}

As expected, the frontend with the short TTL was removed. This is useful because you can set your hosts to periodically send a “heartbeat” request that will ensure that all values within the frontends directory are up to date. The request is simply to set the key to the same value and TTL, therefore extending the lifetime that it will remain in the store. If there is no heartbeat request after the TTL has elapsed, then the entry is removed and it can be assumed that the machine or process is not available. This is especially helpful for deployment, but also can be used to update load balancers automatically when a machine is added or removed.

### Shared Configuration

We’ve gone though the basics of keys and directory storage. In addition to them being very useful for machine management, configuration values can be stored so that your run scripts don’t require too many flags, environment variables can be left alone, and config files don’t need to be managed on each box. A simple way of working with this setup would be for the relevant processes to query for the keys that it needs on startup, but that requires the binary to be restarted in order for new config values to take effect. What if we want to respond to changes immediately without having to restart anything? That’s where etcd’s watch feature comes into play.

So far, we’ve been using curl, but let’s get our hands dirty with a bit of Go code and the [go-etcd](https://github.com/coreos/go-etcd) client library. Make sure your server is still running. We’re going to need it.

Let’s re-add the second frontend server that was removed during our TTL example:

 	$ curl -L http://127.0.0.1:4001/v2/keys/frontends/fe2 -XPUT -d value=10.0.1.2

The following code will connect and list out each frontend stored within the frontends/ directory:

	package main

	import (
		"github.com/coreos/go-etcd/etcd"
		"log"
	)

	func main() {
		client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
		resp, err := client.Get("frontends/", false, false)
		if err != nil {
			log.Fatal(err)
		}
		for _, n := range resp.Node.Nodes {
			log.Printf("%s: %s\n", n.Key, n.Value)
		}
	}


That’s great, you say, but what about the configuration stuff I was talking about? Let’s say I have a set of credentials that more than one service uses to connect to one of our datastores. We set that as a key in our etcd instance for shared usage:

 	$ curl -L http://127.0.0.1:4001/v2/keys/creds -XPUT -d value='dbname=naughtylist host=ec2-123-73-145-214.northpole.compute-1.amazonaws.com port=6212 user=saintnick password=ilovemrsclaus sslmode=require'

Now we need a program to watch for updates when those credentials were to change:

	package main

	import (
		"github.com/coreos/go-etcd/etcd"
		"log"
	)

	func main() {
		client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
		resp, err := client.Get("creds", false, false)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Current creds: %s: %s\n", resp.Node.Key, resp.Node.Value)
		watchChan := make(chan *etcd.Response)
		go client.Watch("/creds", 0, false, watchChan, nil)
		log.Println("Waiting for an update...")
		r := <-watchChan
		log.Printf("Got updated creds: %s: %s\n", r.Node.Key, r.Node.Value)
	}


Run the above program...

	 $ go run watch.go 
	 2013/12/04 18:37:23 /creds: dbname=naughtylist host=ec2-123-73-145-214.northpole.compute-1.amazonaws.com port=6212 user=saintnick  password=ilovemrsclaus sslmode=require
	 2013/12/04 18:37:23 Waiting for an update...

and it will wait for a change to the /creds key, so when you change it...

 	$ curl -L http://127.0.0.1:4001/v2/keys/creds -XPUT -d value='dbname=naughtylist host=ec2-123-73-145-214.northpole.compute-1.amazonaws.com port=6212 user=saintnick password=iadoremrsclaus sslmode=require'

It will print the updated credentials value:

 	2013/12/04 18:37:39 Got updated creds: /creds: dbname=naughtylist host=ec2-123-73-145-214.northpole.compute-1.amazonaws.com port=6212 user=saintnick password=iadoremrsclaus sslmode=require

And that’s it! You can use this functionality to update any clients to reconnect using the new database credentials however you see fit.

### Conclusion

We’ve only scratched the surface of how you can make your life easier using etcd with some basic examples, and I hope that they helped to demonstrate how powerful it can be for you. I find it to be one of the best-designed tools precisely because of this simplicity, and I hope you do too. Also, did I mention it’s also completely written in Go?

Happy hacking and if you’d like to see more posts like this, get @ me on [Twitter](https://twitter.com/andybons)!
