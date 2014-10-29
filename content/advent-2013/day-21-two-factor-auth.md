+++
title = "Go Advent Day 21 - Two Factor Authentication in Go"
date = 2013-12-21T06:40:42Z
author = ["Damian Gryski"]
series = ["Advent 2013"]
+++

## Introduction

Every day we hear about another site getting hacked and more passwords being
leaked.  Bad passwords and password reuse are two of the biggest problems with
the human side of computer security.

Two-factor Authentication (2FA) is an attempt to improve things.  Passwords
alone are "something you know", and if the database is compromised then
somebody else can know it too.  This is why password reuse is a problem.

With 2FA, passwords are augmented with "something you have".  This used to be a
key fob, but today mobile phone apps are more common.  When you log in, you
enter your username, password, and a single-use security token from the key fob.
A successful login requires validation of both the password and the token.

## Google Authenticator

One of the most commonly used security token generators is
[Google Authenticator](https://code.google.com/p/google-authenticator/),
which uses one-time passwords specified in
[RFC6238](http://tools.ietf.org/html/rfc6238) and
[RFC4226](https://tools.ietf.org/html/rfc4226).  There are clients available
for a wide range of devices, and many 2FA apps support generating Google
Authenticator tokens in addition to their own. A number of large websites use
Authenticator tokens for 2FA: Amazon Web Services, Dreamhost, Dropbox,
Facebook, GitHub, Linode, and of course Google.

There are several Go libraries available to generate the corresponding tokens,
including one by me which also implements the same validation logic as Google's
PAM module: [dgoogauth](https://github.com/dgryski/dgoogauth), christened back
when I wasn't sure how to name things. Kyle Isom also has a complete HOTP-based
authentication solution and demo app at
[github.com/gokyle/hotp](https://github.com/gokyle/hotp)

The most common way to configure applications using Google Authenticator is
with QR codes.  These can be generated with
[Russ Cox's QR code library](http://godoc.org/code.google.com/p/rsc/qr)

## YubiKey

[Yubico's YubiKey](http://www.yubico.com/) hardware token generates one-time
passwords. The key fob plugs into a USB port and acts as a keyboard which types
the password directly into the input field when the button is pressed.

If you just want to validate tokens provided by your users, there is a client
library for YubiKey's cloud-based authentication service:
[github.com/dgryski/go-yubicloud](https://github.com/dgryski/go-yubicloud)

If you are more interested in building your own authentication service, you'll
want to look at [Conformal's YubiKey library](https://github.com/conformal/yubikey)
for parsing tokens.  Yubico's cloud validation service is open source and
written in PHP.  I've implemented a
[simplified clone in Go](https://github.com/dgryski/go-yubiauth).

## Duo Security

[Duo Security](https://www.duosecurity.com/) provides three different
authentication factors, all connected to a user's smartphone.  The user can be
prompted with a push notification from the app, be provided with a single-use
token via SMS, and even have the phone generate a token if it is off-line.

Duo provides a demo web app which uses a iframes, JavaScript, and a high-level
API to authenticate users with Duo's 2FA. I have
[ported it to Go](https://github.com/dgryski/go-duoweb).

There is also a lower-level REST API for which there are currently no bindings:
write some and become famous!

## Other Providers

Obviously this list isn't exhaustive.  I've covered the major players, but
there are many more.  Some already have client libraries (for example
[Authy](https://www.authy.com/):
[github.com/dcu/authygo](https://github.com/dcu/authygo)), but most do not.
Luckily, writing clients is generally "easy" given Go's great `net/http` package.

There are also solutions targeted at enterprise security and require
interfacing with a RADIUS or LDAP server.  Integration with these becomes more
likely as Go's prominence as an infrastructure language increases.

## Questions Needing Answers

Adding a complete two-factor authentication solution is more complicated than
just using passwords.  There are a number of tricky areas that need to be
figured out.  I can't cover all the issues here, but I'll list some to give you
an idea of some of the issues you will face:

- Are you going to do true 2FA, or have you just replaced a password ("something you know") with a token ("something you have")?
- Is two factor authentication required or optional?  How frequently do you need to enter a token?
- Do you need to add support for backup or scratch codes that allow entry without the token?
- How are you going to save state?  Hashed passwords are essentially static.  For 2FA, you need to store information about each attempt.  You must verify that each token is used only once.  You need to ensure that the time or counters are always increasing.
- How are you going to prevent race conditions?  It shouldn't be possible for multiple simultaneous login attempts with the same token to succeed.
- How are you going to provision and store the shared secrets?
- What happens if a user loses their phone? How do you validate a new device?  How do you revoke the old one?  What do you do in the meantime?
- Do you want an on-premise solution or is a cloud provider sufficient?

If you're using a cloud service, they should take care of some of the technical
issues, but that still leaves a number of business decisions that need to be
figured out with the CSO.

## Final Thoughts

Adding two-factor authentication can make your web app more secure, and there
are great libraries for the most popular solutions.

If you're dealing with sensitive information, you owe it to your users to
integrate more than just passwords-based logins.
