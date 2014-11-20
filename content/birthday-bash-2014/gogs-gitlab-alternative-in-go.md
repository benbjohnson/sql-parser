+++
author = ["Jiahua Chen"]
date = "2014-11-20T00:00:00-06:00"
title = "Gogs: GitLab alternative in Go"
series = ["Birthday Bash 2014"]

+++

### What is Gogs and why we make it?

<img alt="Gogs Logo"
     src="/postimages/gogs-gitlab-alternative-in-go/Logo_gopher.png"
     width=128 height=128
     style="float:left; padding: 10px"/>
[Gogs](http://gogs.io) is a painless self-hosted Git Service written in Go. It aims to make the easiest, fastest and most painless way to set up a self-hosted Git service. With Go, this can be done in independent binary distribution across **ALL platforms** that Go supports, including Linux, Mac OS X, and Windows. 

### Why we choose Go?

As a strong type and compiled system programming language, Go has significant ability to catch errors at compile time to reduce possibility of runtime errors, and it is extremely useful when we want to make changes in the project.

In fact, we use Go compiler as the error checker to check whether types are mismatched, or function/method calls are giving correct arguments number and type. Meanwhile, Go compiler only requires few seconds to compile a project that has over 40k lines of code, so we can get feedback from compiler very quickly.

To be more awesome, Go has a tool called "gofmt" to keep code style unified, this helps community contributors understand team members' code better and faster, and vice versa.

### Why don't just use GitLab?

GitLab is an awesome product, I cannot deny that, it has tons of features and integrations, especially huge ecosystem build upon it. However, everything just don't make sense if user cannot even install on their machine or get service running correctly and easily. 

This is where Gogs comes with huge advantages.

Go allows us compile code into one single binary, and thus there is rarely no requirements for distributing and deploying Gogs, and makes Gogs to be a cross-platform application without much additional work. Even when users want to install Gogs from source code, the universal command `go get github.com/gogits/gogs` handles most of work for them. This is how Go helps Gogs deploy easy and make users happy.

#### Low resource usage and requirements

Not surprisingly, Gogs has significantly lower resource usage and requirements than GitLab, many users converts from GitLab to Gogs for this reason:

![](/postimages/gogs-gitlab-alternative-in-go/twitter-screenshot-1.png)

![](/postimages/gogs-gitlab-alternative-in-go/twitter-screenshot-2.png)

#### Easy to upgrade

For Gogs, our goal is not just to make it easy to install, but easy to upgrade as well. 

- For binary deployment, simply download and unzip new files into old directory.
- For source code builds, simply `git pull origin master` and rebuild Gogs are basically all you need to do.

![](/postimages/gogs-gitlab-alternative-in-go/twitter-screenshot-3.png)

### Future

It is true that Gogs is a new star and lacks some features, most importantly, its ecosystem. But it is what we're hardly working on, to develop more features users need, and more integrations users love.

In the end, I want to thank you, Go, love you at first sight, and happy 5th birthday! 
