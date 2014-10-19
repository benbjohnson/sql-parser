+++
title = "Automatically Deploy A Revel Web Application"
date = 2014-07-09T06:40:42Z
author = ["Brian Ketelsen"]
tags = ["revel","deployment","git"]
+++

# Introduction

The websites that power GopherAcademy and GopherCon are written using [Revel](http://robfig.github.io/revel/), which is a very nice framework for building web applications.  Go has a great built-in HTTP server, but there are times when you don't want to roll-your-own web framework.  Revel is great if you're looking for a batteries-included approach to web development in Go.

I come from a Ruby and Rails background, and one of my favorite parts of the Rails ecosystem is [Capistrano](https://github.com/capistrano/capistrano/wiki).  Deploying a Rails app is just a simple "cap deploy" command once you've set things up correctly.  I missed that with Revel.  I know that I could have used Capistrano to deploy our Revel applications; capistrano isn't limited in what you can use it to deploy.  But I wanted an opportunity to learn new technologies and push my Linux/Git/Go learning a little bit.  When Jeff Lindsay released [Dokku](https://github.com/progrium/dokku) I had my inspiration.  I wanted to create an automated deployment system like Dokku, but without using Docker. 

If you haven't checked out Dokku yet, you should.  It's a tiny amount of Bash scripting that wraps the fabulous [Docker](http://www.docker.io) deployment system.  With Dokku and Docker, you can create a Heroku-like deployment system on your own Ubunut server.  Dokku adds some hooks in the git workflow that intercepts your `git push` and deploys your app automatically.

I cloned the Dokku source and found that the meat of the git interception happens in [gitreceive](https://github.com/progrium/gitreceive).  Armed with this new knowledge, I set off to make something smaller than Dokku that would accept a `git push` and automatically deploy my Revel web applications.

## Install gitreceive

The first step in making this happen is to install gitreceive.  Follow the instructions on the [gitreceive](https://github.com/progrium/gitreceive) repository home page to download and initialize gitreceive on your server.  This article presumes that you're running Ubuntu 13+ like I am.

After you run _sudo_gitreceive_init_, gitreceive creates a "git" user on your server with a home directory at /home/git.  It also creates a sample _receiver_ script that shows an exmample of what you can do with the script.  I copied this script and saved it as _receiver.original_ so I could reference it later.

## Install NginX

The next step in the deployment methodology is to have NginX proxy all of my Revel web applications.
	
	sudo apt-get install nginx

That was easy.  I've decided that I'm going to deploy all my web applications at _/var/www/$APPNAME/_ using a capistrano-like setup where each release is in its own folder and the root has a symlink to the most current release.

## Create the new _receiver_ script

To deploy a Revel application, you have a few options.  I develop on a Mac but deploy to Linux, so I knew I wanted to push the code to the server and use the Revel tools to package and deploy the web apps.  See the [Revel Deployment Guide](http://robfig.github.io/revel/manual/deployment.html) for more information on your options. 

Here's how my new receiver script looks:

	#!/bin/bash
	echo "Removing previous directory"
	rm -rf /home/git/tmp/src/$1
	echo "Creating new package directory"
	mkdir -p /home/git/tmp/src/$1 && cat | tar -x -C /home/git/tmp/src/$1
	export PATH=$PATH:/home/git/tmp/bin:/usr/local/go/bin
	export GOPATH=/home/git/tmp
	echo "Packaging application $1"
	sed -i 's/BOGUSPASS/myMysqlRootPassword/g' /home/git/tmp/src/$1/conf/app.conf
	revel package $1
	mkdir -p /var/www/$1/releases/$2
	tar -zxvf /home/git/$1/$1.tar.gz -C /var/www/$1/releases/$2
	rm /var/www/$1/current 
	ln -s /var/www/$1/releases/$2 /var/www/$1/current
	chmod +x /var/www/$1/current/run.sh
	sudo restart $1

Let's walk through it line by line and explain what is going on.

	#!/bin/bash

This is a bash script.  We have to put the shebang line in to mark it as such.

	echo "Removing previous directory"

Gitreceive is kind enough to pass stdout from your git push session back to the client.  It's really nice to have this feedback as the app is being assembled and deployed.

	rm -rf /home/git/tmp/src/$1

I created a $GOPATH in _/home/git/tmp_ so that I could compile the Revel binary and have any packages installed that the web apps depend on.  Here, I'm removing the last version of the application before the _gitreceive_ script puts the new code there.  You'll need to have the revel binary somewhere in your path for this process to work.  I installed it as the _git_ user in _/home/git/tmp/bin_ so it would be available to these scripts.

	echo "Creating new package directory"

More feedback for the remote user.

	mkdir -p /home/git/tmp/src/$1 && cat | tar -x -C /home/git/tmp/src/$1

This is where the _gitreceive_ magic happens.  First I create a directory with the name of the repository that's being pushed.  Then pipe the output of the _gitreceive_ script (which is a tar'd version of the repository files you're pushing with git) to the tar command, un-tarring them to _/home/git/tmp/src/$1_ where $1 represents the name of the repository.

	export PATH=$PATH:/home/git/tmp/bin:/usr/local/go/bin
	export GOPATH=/home/git/tmp

These two lines setup a temporary $GOPATH and working environment for the compilation of the Revel web app.

	echo "Packaging application $1"

More feedback.

	sed -i 's/BOGUSPASS/myMysqlRootPassword/g' /home/git/tmp/src/$1/conf/app.conf

Here, I remove a placeholder password in my Revel app's app.conf file and replace it with my real mysql password.  Now I can keep the app's source in a public repo on Github without compromising my server.

	revel package $1

The _revel package_ command takes the Revel web application, compiles it, and creates a tar.gz file with the binary, assets and a shell script to run the app.

	mkdir -p /var/www/$1/releases/$2

This creates a directory under _/var/www_ with the application name, and the git SHA like this: /var/www/gopheracademy/releases/biglonggitshahere.

	tar -zxvf /home/git/$1/$1.tar.gz -C /var/www/$1/releases/$2

Now untar the Revel package into the previously created directory.

	rm /var/www/$1/current 
	ln -s /var/www/$1/releases/$2 /var/www/$1/current

These two commands remove the previously created _current_ symlink and replace it with a symlink to the newest deployment.

	chmod +x /var/www/$1/current/run.sh

Mark the Revel run script as executable.

	sudo restart $1

Restart the _Upstart_ service that runs the web application.

## UPDATE

Rob Figueiredo emailed me with a few comments about the post.  The best piece of advice he gave me was to ditch the `revel package $1` command in favor of `revel build $1 /var/www/$1/releases/$2` which shaves many seconds off the deployment and avoids the creation of the tar.gz file completely.  Awesome!  He also said that Revel supports environment variables in the app.conf file so I could export ${DBPASS} and use ${DBPASS} in my mysql connection script.  I tried this and couldn't make it work, more than likely it's a problem with the setuid stanza in the Upstart script and how Upstart handles environment variables.  I'll keep trying with that piece.  Thanks for the great advice Rob.  And of course, thanks for Revel, we love it.  The new script looks like this:

	#!/bin/bash
	echo "Removing previous directory"
	rm -rf /home/git/tmp/src/$1
	echo "Creating new package directory"
	mkdir -p /home/git/tmp/src/$1 && cat | tar -x -C /home/git/tmp/src/$1
	export PATH=$PATH:/home/git/tmp/bin:/usr/local/go/bin
	export GOPATH=/home/git/tmp
	echo "Packaging application $1"
	sed -i 's/BOGUSPASS/MyGoodPass/g' /home/git/tmp/src/$1/conf/app.conf
	revel build $1 /var/www/$1/releases/$2
	rm /var/www/$1/current 
	ln -s /var/www/$1/releases/$2 /var/www/$1/current
	chmod +x /var/www/$1/current/run.sh
	sudo restart $1

## NginX Configuration

There could be any number of Revel web apps deployed using this method.  We have two running on this server, so I thought it would be convenient to put the NginX configuration files right in each application's repositories.  Each application has an nginx.conf file in the root of the repo that looks like this:
	
	upstream ga { server 127.0.0.1:9001; }
	server { 
  		listen      80;
  		server_name $hostname;
  		location    / {
    		proxy_pass  http://ga;
    		proxy_http_version 1.1;
    		proxy_set_header Upgrade \$http_upgrade;
    		proxy_set_header Connection "upgrade"; 
  		}
	}

This one is for GopherAcademy, it declares an upstream host on port 9001, which matches the Revel app's configuration file setting for port.

In _/etc/nginx/conf.d/_ I added a file called autodeploy.conf with this line:

	include /var/www/*/current/src/*/nginx.conf;

Now NginX will automatically include configurations for any revel application deployed using this gitreceive script.

## Upstart

The last step in this process is to create an Upstart script to run the shell script that Revel creates.  Here's the one for Gopher Academy, it lives in _/etc/init/_ and it's called gopheracademy.conf:

	description "Gopheracademy Website"

	start on (local-filesystems and net-device-up IFACE!=lo)

	kill signal TERM
	kill timeout 60

	respawn
	respawn limit 10 5

	setgid mydeploymentuser
	setuid mydeploymentuser

	# oom score -999
	#console log

	script
  		/var/www/gopheracademy/current/run.sh
	end script
	
The important thing to note is that this Upstart script runs my Revel app as `mydeploymentuser` instead of running it as root.   To start the web app from the command line type the following:

	start gopheracademy

That's it.  The service will start on server reboot, and it's restarted as the last line of the gitreceive script mentioned above.

## Setting up your git repository

To enable auto-deployment, you have to enable SSH keys for a git user.  The gitreceive documentation shows this example:

	cat ~/.ssh/id_rsa.pub | ssh you@yourserver.com "sudo gitreceive upload-key progrium"

where `progrium` is the user attached to the key.  In my case it's the user `bketelsen` because that's who my key belongs to.  

Now add a git remote to your repository:

	git remote add production bketelsen@myserver.ip.address:gopheracademy

This line is important, so let's break it down.  We're creating a remote git repository reference called _production_, with the git user _bketelsen_ (which has to match the git user you created with the upload-key command above) at my Ubuntu deployment server.  The repository name _gopheracademy_ corresponds to the $1 variable in the gitreceive scripts.  So to emphasize that point again, the :gopheracademy piece of this command will be concatenated into the _/var/www/$1_ lines as _/var/www/gopheracademy_.

## Profit

After adding the git remote, make changes to your local Revel application, commit them, then deploy using the following command:

	git push production master

This means "deploy the master branch to the remote named production".  You can watch the output of your _gitreceive_ command from your command prompt as you deploy.

I suspect you'll end up with some permissions issues as you follow along.  The simplest thing to do is to make your _/var/www_ directory owned by the _git_ user, and make the user in your Upstart script's _setuid_ stanza the _git_ user as well.

I hope you enjoyed this tutorial.  It was a fun project project and I learned quite a bit about Revel packaging and Bash scripting along the way.

Here's some sample output of an actual deployment of the Gopher Academy website:

	Ix:gopheracademy bketelsen$ git push production master
	Counting objects: 19, done.
	Delta compression using up to 4 threads.
	Compressing objects: 100% (10/10), done.
	Writing objects: 100% (10/10), 946 bytes, done.
	Total 10 (delta 8), reused 0 (delta 0)
	remote: Removing previous directory
	remote: Creating new package directory
	remote: Packaging application gopheracademy
	remote: ~
	remote: ~ revel! http://robfig.github.com/revel
	remote: ~
	remote: 2013/07/09 15:52:44 revel.go:267: Loaded module static
	remote: 2013/07/09 15:52:44 build.go:72: Exec: [/usr/local/go/bin/go build -tags  -o /home/git/tmp/bin/gopheracademy gopheracademy/app/tmp]
	remote: Your archive is ready: gopheracademy.tar.gz
	remote: gopheracademy
	remote: run.bat
	remote: run.sh
	remote: src/github.com/robfig/revel/conf/mime-types.conf
	remote: src/github.com/robfig/revel/modules/static/app/controllers/static.go
	remote: src/github.com/robfig/revel/modules/testrunner/app/controllers/testrunner.go
	remote: src/github.com/robfig/revel/modules/testrunner/app/plugin.go
	remote: src/github.com/robfig/revel/modules/testrunner/app/views/TestRunner/FailureDetail.html
	remote: src/github.com/robfig/revel/modules/testrunner/app/views/TestRunner/Index.html
	remote: src/github.com/robfig/revel/modules/testrunner/app/views/TestRunner/SuiteResult.html
	remote: src/github.com/robfig/revel/modules/testrunner/conf/routes
	remote: src/github.com/robfig/revel/modules/testrunner/public/css/bootstrap.css
	remote: src/github.com/robfig/revel/modules/testrunner/public/images/favicon.png
	remote: src/github.com/robfig/revel/modules/testrunner/public/js/jquery-1.9.1.min.js
	remote: src/github.com/robfig/revel/templates/errors/403.html
	remote: src/github.com/robfig/revel/templates/errors/403.json
	remote: src/github.com/robfig/revel/templates/errors/403.txt
	remote: src/github.com/robfig/revel/templates/errors/403.xml
	remote: src/github.com/robfig/revel/templates/errors/404-dev.html
	remote: src/github.com/robfig/revel/templates/errors/404.html
	remote: src/github.com/robfig/revel/templates/errors/404.json
	remote: src/github.com/robfig/revel/templates/errors/404.txt
	remote: src/github.com/robfig/revel/templates/errors/404.xml
	remote: src/github.com/robfig/revel/templates/errors/500-dev.html
	remote: src/github.com/robfig/revel/templates/errors/500.html
	remote: src/github.com/robfig/revel/templates/errors/500.json
	remote: src/github.com/robfig/revel/templates/errors/500.txt
	remote: src/github.com/robfig/revel/templates/errors/500.xml
	remote: src/gopheracademy/README.md
	remote: src/gopheracademy/app/content/building-stathat-with-go.article
	remote: src/gopheracademy/app/content/building-stathat-with-go_stathat_architecture.png
	remote: src/gopheracademy/app/content/building-stathat-with-go_weather.png
	remote: src/gopheracademy/app/controllers/app.go
	remote: src/gopheracademy/app/controllers/db.go
	remote: src/gopheracademy/app/controllers/init.go
	remote: src/gopheracademy/app/controllers/jobs.go
	remote: src/gopheracademy/app/controllers/newsletter.go
	remote: src/gopheracademy/app/models/jobs.go
	remote: src/gopheracademy/app/routes/routes.go
	remote: src/gopheracademy/app/tmp/main.go
	remote: src/gopheracademy/app/views/Application/Index.html
	remote: src/gopheracademy/app/views/Jobs/Confirm.html
	remote: src/gopheracademy/app/views/Jobs/Find.html
	remote: src/gopheracademy/app/views/Jobs/HandlePostSubmit.html
	remote: src/gopheracademy/app/views/Jobs/Index.html
	remote: src/gopheracademy/app/views/Jobs/Post.html
	remote: src/gopheracademy/app/views/Jobs/Show.html
	remote: src/gopheracademy/app/views/errors/404.html
	remote: src/gopheracademy/app/views/errors/500.html
	remote: src/gopheracademy/app/views/footer.html
	remote: src/gopheracademy/app/views/header.html
	remote: src/gopheracademy/conf/app.conf
	remote: src/gopheracademy/conf/routes
	remote: src/gopheracademy/db/dbconf.yml
	remote: src/gopheracademy/db/migrations/001_jobs.sql
	remote: src/gopheracademy/messages/sample.en
	remote: src/gopheracademy/nginx.conf
	remote: src/gopheracademy/public/css/bootstrap-responsive.css
	remote: src/gopheracademy/public/css/bootstrap-responsive.min.css
	remote: src/gopheracademy/public/css/bootstrap.css
	remote: src/gopheracademy/public/css/bootstrap.min.css
	remote: src/gopheracademy/public/images/building-stathat-with-go_stathat_architecture.png
	remote: src/gopheracademy/public/images/building-stathat-with-go_weather.png
	remote: src/gopheracademy/public/images/favicon.png
	remote: src/gopheracademy/public/images/project.png
	remote: src/gopheracademy/public/img/glyphicons-halflings-white.png
	remote: src/gopheracademy/public/img/glyphicons-halflings.png
	remote: src/gopheracademy/public/js/bootstrap.js
	remote: src/gopheracademy/public/js/bootstrap.min.js
	remote: src/gopheracademy/tests/apptest.go
	remote: gopheracademy start/running, process 21270
	To git@my.server.ip:gopheracademy
	   dd71f2b..5117a28  master -> master

## UPDATE

With the new _receiver_ script changes that Rob Figueiredo pointed out, the output of a deployment is much smaller, and significantly faster to boot.  Here's what it looks like now:

	brians-air:gopheracademy bketelsen$ git push production master
	Counting objects: 7, done.
	Delta compression using up to 4 threads.
	Compressing objects: 100% (4/4), done.
	Writing objects: 100% (4/4), 401 bytes, done.
	Total 4 (delta 2), reused 0 (delta 0)
	remote: Removing previous directory
	remote: Creating new package directory
	remote: Packaging application gopheracademy
	remote: ~
	remote: ~ revel! http://robfig.github.com/revel
	remote: ~
	remote: 2013/07/10 01:22:13 revel.go:267: Loaded module static
	remote: 2013/07/10 01:22:13 build.go:72: Exec: [/usr/local/go/bin/go build -tags  -o /home/git/tmp/bin/gopheracademy gopheracademy/app/tmp]
	remote: gopheracademy start/running, process 22938
	To git@my.server.ip.address:gopheracademy
	   c3bf1da..2806965  master -> master


I love automating things!  This entire process was _HEAVILY_ inspired by Jeff Lindsay's Dokku scripts.  Many thanks to him for creating such an inspiring project.

