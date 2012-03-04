GoPaste

Introduction
============

This project contains two pieces:
- A server for accepting and serving pastes
- A client for sending files and command output to the server

By default, the client will send pastes to my gopaste server.

Installation
------------

To install, run the following command:

    go get -u github.com/kylelemons/gopaste/gp

Then, the `gp` command will be installed into `${GOPATH}/bin`.
The rest of this document assumes that your `PATH` contains this
directory, so that typing `gp` will execute the binary.

Usage
=====

The `gc` command is simple to use.  It consumes standard input
or a file, depending on the command-line arguments.

    Usage of gp:
      -f    <file>    The name of a file to read (standard input if not provided)
      -name <name>    The name of the paste (use filename or MD5 sum if not provided)

Pasting files
-------------

If you want to paste a file, use the following command:

    gp -f <filename>

The name of the file will automatically be included in the paste.
If the filename is a common one (say, `main.go`) you will want to
specify a name for the file, which can be done like this:

    gp -f <filename> --name <more_descriptive_name>

Pasting output
--------------

If you want to paste the output of a command, pipe it into gp.

    go build | gp

This will generate a hashed name that's not very readable,
so you may want to give it a name:

    go build | gp --name <more_descriptive_name>

GoPaste Service
===============

I run a GoPaste server.  Please do not abuse it.
Pastes are currently limited by the server to 1MB in size,
so don't try to send larger files or you'll just waste your outgoing bandwidth :-).

Auto-linking
------------

The primary use of this command is for linking pastes on the #go-nuts IRC channel.
Any paste sent to my GoPaste server will (usually) be automatically linked in the channel.
If you don't want this, I encouge you to use a paste site like GitHub's [gist][1].

Removal
-------

Pastes are automaticlaly removed after a fixed time.  Currently, that time is 1 hour.
I also only save a certain number of submissions, and old ones are the first ones to be deleted.
If you need longer-term pastes, I encourage you to make use of a service like [gist][1].

[1]: https://gist.github.com/ "GitHub's gist web snippet service"
