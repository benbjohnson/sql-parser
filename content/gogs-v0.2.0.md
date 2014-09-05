+++
title = "Gogs: Binary Deployment: The Right Way to Deploy"
date = 2014-04-01T06:40:42Z
author = ["Jiahua Chen"]
+++

# Gogs: When you're deploying a binary, you're doing it Right.

This post is published corresponding to the [Gogs - Go Git Service](https://github.com/gogits/gogs) `v0.2.0` release.

First, please let me speak for the develop team to thank all of our friends who are supporting us on GitHub. As you may know, `v0.2.0` is the first public release of Gogs, and the community has contributed over 650 stars to this project on GitHub in just one week.

As for me, I want to give the most sincere gratitude to all members of develop team, everybody worked extremely hard for the first public release. Our united and tacit understanding is the key to the success of this project.

## Project Overview

We want to be clear with this, the first public release, to help you understand more about why we launched this project, how we develop it, and its current develop status.

### Purpose

In the area of self-hosted Git services, there are some very succeesful products running all around the world, so why did we decide to develop Gogs? How can we compete with them(e.g. gitlab)? What's the view of us in terms of recreating the wheel?

Officially, our purpose is to take the advantage of Go building everything into one binary to setup your own self hosted Git Service. But in my words, I just don't like the products that exist right now. Because they need a lot of installation steps and dependencies, or they come from a single developer with no anticipated prospect, or they aren't cross-platform. And Gogs certainly supports a rich number of operating systems because of Go's cross-platform reach.

We have at least two competitors in the long term, gitlab and GitHub Enterprise. So what does Gogs have to compete with them? I think right now we have only one advantage: the binary deployment. Setting up your git service simply will be our big competitive advantage.

Many people just hate so-called "creating the wheel repeatedly" without any reason, then they miss the opportunity to be creative. You need to see it in several perspectives. First of all, it's a great thing for learning and practicing skills, especially try something you've never got a chance to try, like we chose [martini](http://martini.codegangsta.io/) as our basic framework. Secondly, current services do not make us feel satisfied, and we won't give up the chance to improve it.

### Develop Team

Gogs is the product of the internet age, and its development team is also uncoupled through the internet. All 5 members were known each other for a while by developing Go and we've never met face-to-face. We use IM(QQ) to discuss, use [Trello](https://trello.com/b/uxAoeLUl/gogs-go-git-service) to assigning tasks, and use GitHub to work together. Four of them are taking time after work to join the project except me(as a student), it's remarkable. 

Now, please let me introduce our develop team:

- [@Unknown](https://github.com/Unknwon)：Project manager, back-end developer
- [@lunny](https://github.com/lunny)：Back-end Git and database developer
- [@fuxiaohei](https://github.com/fuxiaohei)：Front-end developer
- [@slene](https://github.com/slene)：Front and back-end developer
- [@skyblue](https://github.com/shxsun)：Back-end developer

### Current Development Status

We've released the first public version with many features but it only has very low version number `v0.2.0` and status `alpha`. I have to say that we do not recommend and Gogs isn't ready for enterprise usage. So does Gogs just something look nice and useless? Of course not! The core of Git hosting is the Git data, Gogs is pretty stable there, so every time new version is released, you only need to do `copy-paste` binary upgrade with no damage to your Git data.

## Getting Started

All documentation is in the [GitHub Wiki](https://github.com/gogits/gogs/wiki]) which helps to give you a deeper understanding of Gogs.

- For people who just want to have a try, after installing the [Prerequirements](https://github.com/gogits/gogs/wiki/Prerequirements), then [Install from binary](https://github.com/gogits/gogs/wiki/Install-from-binary) and you're all set.
- If you want to [Install from source](https://github.com/gogits/gogs/wiki/Install-from-source), it's also quite easy once you have a Go environment.

To be more convenient, Gogs also provides an installer page(go to URL:`/install`) to help you configure for the first-time run.

## Notice

We're glad you choose Gogs, and here we have some notices for you:

- Avoid unnecessary exceptions by looking at our [Known Issues](https://github.com/gogits/gogs/wiki/Known-Issues).
- See [FAQs](https://github.com/gogits/gogs/wiki/FAQs) and [Troubleshooting](https://github.com/gogits/gogs/wiki/Troubleshooting) before you [Open an Issue](https://github.com/gogits/gogs/issues/new) on GitHub.

## Future Plans

- In general, we'll release new mirror(+0.1) version every month in average.
- For non-English users, we plan to support i18n in `v0.5.0`. Keep eyes on [Issue #9](https://github.com/gogits/gogs/issues/9).
- Optimizing CPU and memory usage is big part of future releases.
- To keep tracking our news, follow us on [Twitter](https://twitter.com/gogitservice).
- Since we distribute Gogs in binary, [gobuild.io](http://gobuild.io) is our compile and download service provider.

Thank you for taking time read this post, Gogs will grow much better with your support and feedback!

