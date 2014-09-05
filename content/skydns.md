+++
title = "SkyDNS (Or The Long Road to Skynet)"
date = 2013-10-10T06:40:42Z
author = ["Brian Ketelsen"]
+++

# SkyDNS and Skynet

This article is in two sections.  The first is the announcement of SkyDNS, a new tool to help manage service discovery and announcement.  The second part of the article 
is a bit of a back story about why we needed this tool and how we got here.  If you're the impatient type, you can read the annoucement and description of SkyDNS and skip the 
rest.

## SkyDNS

Today we're releasing SkyDNS as an open source project on [github](https://github.com/skynetservices/skydns).  SkyDNS is our attempt to solve the problem of finding 
the services that are running in a large environment.  SkyDNS acts as a DNS service that only returns [SRV records](http://en.wikipedia.org/wiki/SRV_record).  Services announce 
their availability by sending a POST with a small JSON payload.  Each service has a [Time to Live](http://en.wikipedia.org/wiki/Time_to_live) that allows SkyDNS 
to expire records for services that haven't updated their availability within the TTL window.  Services can send a periodic POST to SkyDNS updating their TTL to keep 
them in the pool.

Now when services die, hosts die, datacenters die, or meteors fall, your clients won't be wasting valuable time making requests to services that aren't running or aren't 
reachable.

When a client is looking for a particular service, it issues a DNS request to SkyDNS with the particulars of the service needed.  SkyDNS supports several interesting 
parameters in the SRV record to facilitate returning the service that is best situated to serve the client.

### Environment 

SkyDNS understands the last parameter in the DNS query as an Environment.  You can put any value here, but for our purpose, words like `development`, `staging`,
`qa`, `production`, and `integration` make the most sense.  The environment parameter allows you to segregate your services so that you can run a single SkyDNS 
service cluster but serve your entire development environment.

### Service

The Service parameter is the name of the service that is running.  Yours might be called addressvalid if you're running an address validation service, or it could be called 
indexentries if you're running a service that indexes documents in the background.  The service name is the unique name that identifies WHAT your service does.

### Version 

Services can be versioned to allow for variations in the request/response data without breaking backwards compatibility.  The version parameter allows you to request a 
specific version of a service.

### Region 

We run services in multiple data centers, so SkyDNS understands locality of a service.  Services register themselves as running in a specific region, which allows clients 
to request a service running in the same region.  If you're running the same service in multiple regions, SkyDNS will use the priority field in the SRV record to return services 
in the requested region first, with the same priority.  Services running in different regions are returned at a lower priority so that the client can use them in case they're needed 
but prefer the local services first.

### Host 

SkyDNS understands the host that a service is running on so that if a client prefers it, the client can request a service on the same host.

### UUID 

Services are required to have a unique identifier, which allows multiple services at the same version to be running on the same host.  If you really wanted to request 
a specific instance of a service you can specify the UUID in the request and you'll only receive that specific instance.

### Wildcards

The only field that is required in the DNS request is the "Environment" field.  Any missing fields are interpreted as "any", and you can specify "any" at any point along 
the parameter list.

To make a request for a service running in the production environment named `addressvalid` you would issue a dns request that looks like this:

     dig @localhost addressvalid.production SRV

This request would return a DNS SRV response that included all instances of the  `addressvalid` service that are running in any region, on any host, at any version level.  You can specify 
some of the parameters and leave others as wildcards.  To request version 1 of the `addressvalid` service you could make this request:

     dig @localhost 1.addressvalid.production SRV 

Similarly, you can request any version, but in a specific region with this request:

     dig @localhost east.any.addressvalid.production SRV

This wildcard flexibility allows you to version your services and serve the requests from the closest available server without too much effort.

### Responses

SkyDNS returns valid SRV records that list the services that are running matching the query.  Here's an example response:

`dig @localhost testservice.production SRV`

   ;; QUESTION SECTION:
   ;testservice.production.    IN  SRV

   ;; ANSWER SECTION:
   testservice.production. 615   IN  SRV 10 20 80   web1.site.com.
   testservice.production. 3966  IN  SRV 10 20 8080 web2.site.com.
   testservice.production. 3615  IN  SRV 10 20 9000 server24.
   testservice.production. 3972  IN  SRV 10 20 80   web3.site.com.
   testservice.production. 3976  IN  SRV 10 20 80   web4.site.com.

   ;; Query time: 0 msec
   ;; SERVER: 127.0.0.1#53(127.0.0.1)
   ;; WHEN: Thu Oct 10 11:47:22 EDT 2013
   ;; MSG SIZE  rcvd: 310

The response shows the service name, environment, TTL, priority, weight, service port, and host name.  With that response your client can construct 
a request to the service that is provided.

Get the full details in the [README](https://github.com/skynetservices/skydns/blob/master/README.md) file.


### Clustered Operations

SkyDNS was built to enable highly available services, so it only makes sense that SkyDNS itself is highly available.  We built SkyDNS on the super 
easy to use [goraft](https://github.com/goraft/raft) implementation of the RAFT protocol.  SkyDNS servers communicate using the Raft protocol and 
must achieve consensus before any records are committed.  Only one SkyDNS server runs as "master" at any given time so requests made to SkyDNS servers that aren't 
serving as "master" are redirected to the current master.  When a SkyDNS server becomes unavailable for any reason, the other servers elect a new master and operations 
continue without interruption.  When you bring a failed node back online, it automatically retrieves missing data from the other servers in the cluster.

Using Raft allowed us to make SkyDNS highly available with minimum fuss.  Kudos to Ben Johnson for making such an easy to use Raft implementation.  The new [ETCD](https://github.com/coreos/etcd) project uses the same implementation of Raft.  It's pretty solid!

SkyDNS was written by Erik St. Martin with help from Brian Ketelsen.  It's available today on Github: [SkyDNS](https://github.com/skynetservices/skydns), so 
go fork it, make some changes, kick the tires and let us know what you think!  We think it's written in a way that will allow it to be used outside of our 
specific use case, and it might enable some fun projects.  We can envision load balancers, reverse proxies, queue workers and more using SkyDNS to 
make sure that a client only attempts to call a service that is actively running.  Using SkyDNS removes the need for complicated round-robin load balancing in your 
service layer.

### Enjoy!

Please let us know what you think about SkyDNS.  Read on for some more information about how we got here, and the mis-steps we made along the way.


## Skynet

Sometime in 2011 I came up with the idea for Skynet. I was trying to figure out ways to make it easier to scale our big
Ruby on Rails application, and I had been playing with Go for several months.  Keith Rarick and Blake Mizerany introduced Doozer to the 
Go mailing list, and I downloaded it and started playing around.  One of the features that I really liked was "watching" a path for changes.
I had never played with Zookeeper before, so this concept was completely new to me.  Immediately the idea for skynet came to me.  What if we could write services that registered with Doozer using paths for service discovery, and we relied on the node disappearing when
the service terminated to remove that service's availability from the registry?  [Skynet](https://github.com/skynetservices/skynet) was born. I played with it a lot for a few months, and even did a StrangeLoop talk on it.  But it didn't really have a heart and soul until Erik joined the Skynet team.

## Skynet prime

When Erik St. Martin joined the team at work we changed a lot of things.  Originally it used Go's RPC for communication between Skynet services, but 
we realized quickly that it would be painful to port Go's binary encoding protocol Gob to other languages, and we had a direct need for the framework 
to be polyglot.  Instead we evaluated JSON, MessagePack, Protocol Buffers, and BSON and chose to use BSON/RPC because the encoding was a little more 
strongly typed.  Skynet originally had ingress services called initiators that accepted inbound traffic from any source (ftp/http/filesystem) and inserted 
the requests into the Skynet mesh.  Initiators turned out to be kind of useless, so we dropped the idea.  Similarly, the original Skynet had the concept of 
routing, where requests could be routed through a series of services that each had the same request/response parameters, with each service modifying the 
response along the way.  This was too complicated, so we killed that complexity too.

We spent some time brainstorming around the problem we were actually trying to solve with Skynet.  The biggest problem we had was that it isn't easy 
building highly available services in a complex environment, and a lot of the code is boilerplate.  We knew we wanted to start to move away from a 
monolithic Rails application and separate concerns for quicker development, and we knew that we'd want to be able to make use of newly added services in the 
cluster quickly, hopefully without any configuration changes.  Skynet became a platform for creating services that were easy to create, easy to call, and 
had most of the HA batteries included.  A new developer on the team could easily use Skynet to create services that solved business problems rather than 
writing the same HA/metrics/discovery/RPC logic over and over.  Skynet evolved into a simple Service and Client interface that defined a contract between the client 
and the server allowing it to be used in polyglot environments with minimal effort.

The skynet that remained worked nicely.  Sure, it had it's warts to be fixed in future releases, but it worked fine.  Unfortunately, Doozer was painful to install, 
bootstrap, maintain, and keep running.  Skynet depending on Doozer was slowing adoption in the community.  The problem was compounded when Doozer was more or less 
abandoned by Heroku.  We volunteered with others to help maintain the project on Github, but without a benign dictator the project languishes from a lack of attention.  
We decided to try Zookeeper instead.  This was the beginning of Skynet2.

## Skynet2

[Skynet2](https://github.com/skynetservices/skynet2) was supposed to be a temporary fork of Skynet while we moved off of Doozer and on to Zookeeper.
Of course we had to change a million other small things along the way, addressing many of the design decisions and bugs in the original due to it's rapid 
prototyping.  We still haven't bothered to merge Skynet2 back into the original Skynet repository, mostly because there were plenty of people using Skynet 
happily with Doozer, and Zookeeper wasn't much better than Doozer in the long run.  Sure it was easier to install and keep running, 
but from a client perspective we lost a lot of the nice functionality we enjoyed in Doozer, like seeing the actual change that caused a watch to fire. 
We were forced to watch thousands of keys, spawning thousands of goroutines, creating a maintenance nightmare. 
Skynet2 works, and is also in production in a few places, but it's still not ideal.  Zookeeper, like Doozer, just isn't the tool to announce service 
availability in a rapidly changing environment.  Zookeeper is a perfectly good tool, but the way Skynet registrations worked, it required us to do way too 
many crazy hacks to apply our concept of registration.

## SkyDNS

Then a thought came to me on my morning commute.  Why are we working so hard to solve availability and discovery when there are tools that
exist to do this?  Isn't that what DNS is for?  By the time I got to work I had a loose idea of what would eventually evolve into SkyDNS. 
Building upon the amazing [DNS](https://github.com/miekg/dns) library that Miek Gieben wrote, we would use DNS SRV records to announce services 
and a simple JSON post to the DNS server to announce availability and update TTL on the DNS records.

Conceptually this sounded like a good approach, so we put it up on a whiteboard and worked out the details.  SkyDNS would serve only SRV records 
for the services that register with it.  Skynet uses the concept of Regions (which are really availability zones), Services,
Versions (so you can change services without breaking compatibility), and Environments (production, staging, etc).  Using these concepts, we could construct 
a DNS server that responded to SRV queries with the services that were running that matched the needs of the requesting client.  
We'd use the SRV record's priority to tell the client which services were closest to it, and the weight to tell the clients which services were under the lowest load.  
Services use a simple JSON API to register and update their TTL, but clients use DNS requests to find services.

After Erik fixed all of my poor assumptions, he talked me into believing that we could use Raft to make SkyDNS resilient.  I thought Raft would be hard to implement, 
but of course it wasn't.  Erik wrote the whole thing in less than a week.  Now we have a DNS server that accepts registrations from services, and responds to DNS queries from clients.  It only returns DNS entries for services that haven't outlived their TTL, or extended their TTL by updating it before it expires.

Of course, we are going to use it to power Skynet, and it's the first step in our plan to completely revamp and simplify the next version.

- Fix availability and discovery with SkyDNS
- Replace complicated BSON/RPC with JSON/REST
- Make our metrics collector more generic so it can be released and used by anyone with nearly any metrics reporting service 
- Use rsyslog for all logging, supporting cool searchable syslog daemons like [ELSA](https://code.google.com/p/enterprise-log-search-and-archive/)

The end goal for Skynet is the ability to hand the library to a new developer so that she can write services without having to think about 
the hard stuff.  Resilience, availability, metrics, logging, and all the other heavy lifting is already done. 

In my continuing effort to build muscle memory in Vim, this post was created entirely from the command prompt.  No GUI's were harmed in the writing of this text.

