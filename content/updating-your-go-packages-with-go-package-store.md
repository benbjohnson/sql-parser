+++
author = ["Dmitri Shuralyov"]
date = "2014-11-19T00:00:00-06:00"
title = "Updating your Go packages with Go Package Store"
series = ["Birthday Bash 2014"]
+++

[Go Package Store](https://github.com/shurcooL/Go-Package-Store#go-package-store) is an app that displays updates for the Go packages in your GOPATH. Why another way to update Go packages when you can already just do `go get -u`, you might think. In true Go tradition, Go Package Store doesn't try to replace what already exists. Instead, it uses composition to augment it. In the end, Go Package Store simply uses the `os/exec` package to execute `go get -u` for you (which is why it's safe to run). But its goal is to make that experience more delightful and informative.

![Go Package Store Screenshot](/postimages/updating-your-go-packages-with-go-package-store/go-package-store.png)

Go Package Store runs locally and scans your GOPATH for all Go packages that are under a version control system. It checks if there's a newer version available remotely, meaning that running `go get -u` on that package would have some effect.

Then it tries to present the update with as much information about the change as possible. This typically means including the commit messages for every new commit.

If some Go package has a dirty working tree or non-default branch checked out, it simply skips it, as running `go get -u` on it would fail anyway.

### History

Go Package Store has an interesting history of how it came to be, and there are good reasons why it targets Go packages only.

I mentioned the `os/exec` package earlier, which is used to start the `go get -u` process. When I first discovered Go, just over 2 years ago, I was still working in C++, dealing with Visual Studio projects for Windows with Makefiles for Linux, modifying header files by hand whenever any function signature changed. I also needed the ability to launch a process, and I was pretty disappointed that there was no standard cross-platform library for it. I looked at boost, and a few separate libraries, but ended up writing hacky code that was not cross-platform at all.

It was then that I ran into Go, and decided to look at its standard library. I found `os/exec` and was able to do in 4 lines of Go code cleanly what I couldn't with 40+ of my C++ code.

The second thing that I was intrigued by was seeing statements like... `import "github.com/some/package"`. Coming from C++, importing a new library would often involve manual steps of downloading said library, figuring out the peculiarities and installation steps. But with Go, you could just import a package directly from the internet, I thought at the time.

While in reality it's not as direct, which is a good thing as it saves on bandwidth and allows your Go code to compile when you're on a plane without internet. But the ability to make any native Go package available with a single `go get` command and not having to do any manual steps was transformative and extremely valuable. It means you can tackle arbitrarily difficult problems, reuse any number of existing high quality and well maintained libraries, all while still being easy to install for anyone using Go.

It was around that time that I started playing with Go more and more, in the process [creating many small individual packages](https://twitter.com/shurcooL/status/478413714572312576) for my needs. Because I was making so many separate repos, and they were changing often as I was learning more about Go, I found it hard to keep track as to which repos still had uncommitted changes. I wrote a little bash script that iterated over all the folders, performing a `git status` on each one. That eventually lead to a full featured program [`gostatus`](https://github.com/shurcooL/gostatus), which tells you the status of many Go packages.

`gostatus` would simply display a `+` symbol for any Go packages that had updates. After looking at many of those `+`s and blindly doing `go get -u` on them, I became more curious as to what was changing behind the scenes.

It was then that everything needed for Go Package Store to exist came together. I already had a robust way to detect which Go packages had updates. I had played with the GitHub API and knew how to use it to fetch a list of commits with their descriptions. I only needed a way to present all that information, and the terminal wasn't very fitting. But I had played with `html/template` package and got the idea to simply use the browser to present the information.

Go made it very easy to put all the pieces together. But more importantly, Go and its simple but effective design choice for remote import paths and `go get` made it feasible for Go Package Store to exist in the first place.

### Present

One recent addition to Go Package Store was clickable links for each commit.

![Commit Links](/postimages/updating-your-go-packages-with-go-package-store/go-package-store-commit-links.png)

Sometimes, when seeing an interesting commit message, I found myself wanting to learn more about it, or look at the actual code change (perhaps to review it, or decide if I want to update right away). Being able to get to a commit with just one click has been very helpful for that.

This is just another example of how Go Package Store can evolve and become better. With that said, there's currently a PR to [add detailed change support](https://github.com/shurcooL/Go-Package-Store/pull/25) for all Go packages hosted on azul3d.org. It already supports GitHub and code.google.com repositories. But if you run into any issues or want to help make it better, Go Package Store [is open source](https://github.com/shurcooL/Go-Package-Store), so feel free to contribute!

### Conclusion

Go is very young but it's growing up fast. In the few years that I've used it so far, it has been a very joyful experience. It motivates and inspires me to work harder to write better code and make improvements to various open source Go packages, because I feel the foundation on which we build is very solid.

It's really great to see more and more functionality become available as pure Go packages that are easy to `go get`. It makes it easier to reuse work done by others and not spend time on reinventing the wheel.

Lots of talented and hardworking people have created a very rich collection of Go packages that are powering many Go programs today. With Go Package Store, I wanted to make it easier for you to see the improvements, bug fixes and other changes to the Go packages you use, so that updating your Go packages can be a more delightful experience.

Happy Birthday Go! Let's continue to move forward and build beautiful, simple things together!
