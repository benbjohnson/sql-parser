+++
title = "Go Advent Day 14 - gobrew"
date = 2013-12-14T06:40:42Z
author = ["Craig Wickesser"]
series = ["Advent 2013"]
+++

## What is gobrew?

Simply put, gobrew lets you easily switch between multiple versions of go.  It is based on [rbenv](https://github.com/sstephenson/rbenv) and [pyenv](https://github.com/yyuu/pyenv).

## Why gobrew?

Often times you'll be developing against one version of Go when another version is released (or perhaps a release candidate is made available). Instead of fighting to manage multiple versions and changing your $PATH repeatedly, you can use one simple tool to manage your installations of Go.

## Get gobrewing

1. Clone the repository

    	git clone git://github.com/grobins2/gobrew.git ~/.gobrew

2. Update your shell config

    	export PATH="$HOME/.gobrew/bin:$PATH"
    	eval "$(gobrew init -)"

3. List available versions

    	gobrew list

4. Install Go

    	gobrew install 1.2

note: there's a [bug](https://github.com/grobins2/gobrew/issues/7) in gobrew which doesn't allow go version 1.2 to be installed on Mac OS X.

## Contribute

Give [gobrew](https://github.com/grobins2/gobrew) a shot, fork it or submit issues and let's make it the best Go version manager around. If you're
behind a proxy server, check out [gobrew::fw](https://github.com/mindscratch/gobrew_fw) as well.

