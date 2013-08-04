The contents of this repository were originally incorporated into Apache Thrift
as of version 0.7 or commit revision 1072478.*  However, unfortunately upstream
development ceased and necessitated continued development here; thusly, these
files have been rebased against Thrift revision 1362341 at
http://svn.apache.org/repos/asf/thrift/tags/thrift-0.8.0.

This project will continued to be maintained for any further developments for
Thrift for Go as needed (e.g., Go Version 1 compatibility).

Add this Git repository to your Thrift checkout.  It will overwrite a few files
to add options for Go.

Currently generates code for the following protocols:

1. Binary Protocol (with test cases)
2. Fast Binary Protocol (with test cases)
3. Standard Thrift JSON Protocol (with test cases)
4. A (custom) simple JSON Protocol (with test cases)
5. Services (compiles and runs against Java, assumed to work elsewhere)

Tested on Mac OS X 10.6 (Snow Leopard) and Linux ca. 2010 derivatives.

# Manual installation

Install build-essential, maven, libtool and autoconf. For a debian (derived) linux distribution run:
``sudo apt-get install build-essential maven libtool autoconf``

Get the master branch for thrift4go:

``git clone https://github.com/pomack/thrift4go.git``


Get the thrift-0.9.0.tar.gz and extract it:

``wget -c https://dist.apache.org/repos/dist/release/thrift/0.9.0/thrift-0.9.0.tar.gz``

``tar xzvf thrift-0.9.0.tar.gz``


Set environment variables (used by the script) and execute the merge script:

``export THRIFT=$(pwd)/thrift-0.9.0``

``export THRIFT4GO=$(pwd)/thrift4go``

``bash ./thrift4go/scripts/merge_and_build.sh -b``


Configure and make thrift:

``cd thrift-0.9.0``

``./configure --with-go``

``make``


The compiler is located at thrift-0.9.0/compiler/cpp/thrift.
Install the Go library.

``ln -s $(pwd)/lib/go/src/thrift $GOPATH/src/thrift``

``go install thrift``

<!-- Get the master branch for thrift4go.
``git clone https://github.com/pomack/thrift4go.git``

Build the thrift with thrift4go compiler:
``cd thrift4go``
``make``

Build the thrift4go Go library (package thrift).
``ln -s $(pwd)/lib/go/src/thrift $GOPATH/src/thrift``
``go install thrift``

Your compiler is located at
thrift4go/compiler/cpp/ -->





<!-- To install locally, perform the following:
  ``go get github.com/pomack/thrift4go/lib/go/thrift``

Or, to build manually:

    cp -R <thrift_src_dir>/lib/go/src/thrift <go_package_src_directory>/
    go install thrift

Four files for thrift compiler (last tested on August 12, 2012):

1. ``configure.ac``
2. ``lib/Makefile.am``
3. ``lib/go/Makefile.am``
4. ``compiler/cpp/src/generate/t_go_generator.cc`` -->

# Tutorial files

A tutorial has been created in the thrift4go/tutorial/go/src directory.

To use this tutorial, run the following:
    ``thrift -r --gen go <thrift_src_dir>/tutorial/tutorial.thrift``
Build the files in the ``gen-go`` directory using ``go install`` as
appropriate.
Build the files in the ``<thrift_src_dir>/tutorial/go`` directory with
``go install`` as appropriate.
Run the server from ``<thrift_src_dir>/tutorial/go/TutorialServerClient`` and
run the client from either
``<thrift_src_dir>/tutorial/go/TutorialServerClient`` or
``gen-go/tutorial/Calculator/Calculator-remote``.

Make sure you specify the same protocol for both the server and client.

# Basic Walkthrough

``thrift --gen go <thrift_src_dir>/test/ThriftTest.thrift``

This will create ``gen-go/thrift/test/*.go`` and associated files/directories.

- ``gen-go/thrift/test/ThriftTest.go`` shows a service and client base
implementation with the associated interfaces and the ability to send/receive
or serialize/deserialize as necessary.

- ThriftTestClient is a client library designed to access the ThriftTest
service.  No changes would need to be made here.

- A ``ThriftTest/ThriftTest-remote.go`` and associated Makefile is also made
available so you can access a remote service implementing the ThriftTest
interface and see how the client side works under the covers.  The command-line
arguments use the custom JSON parser, so you can just pass in JSON strings as
arguments when you need to populate a struct, which I find better than any
other alternative.

You can use thrift annotation to generate go tags on struct fields, like
struct S {
        1: optional i64 X,
        2: optional i64 Y,
	3: optional i64 Id (go.tag = "json:\"structid,string\"" ), 
}
will generate
type S struct {
        X            int64    `json:"X"`                // 1
        Y            int64    `json:"Y"`                // 2
        Id           int64    `json:"structid,string"`  // 3
}

- ThriftTestProcessor implements the server side and you would want to implement
the server handlers using ``NewThriftTestProcessor()``.

- You just pass in your handler that implements the ``IThriftTest`` interface
and make sure you import the appropriate package.  Package directories/names
are shown in the relevant Makefile.

- One unique thing about Go is that to have a publicly available
function/variable, the first letter has to be capitalized, so all exportable
functions/variables have the first letter capitalized, but since the Thrift
files normally don't, they assume any serialization uses the capitalization
found in the Thrift interface definition file itself.

# Patching into Mainline Thrift
This package is targeted to Thrift stable, which at the time of writing this,
is 0.8.0.  Please give the ``merge_and_build.sh`` script a run for more
information.  It will perform the grungework of merging a local checkout of
Thrift stable and thrift4go.

Also, ``GOROOT`` must be set to the directory that contains the ``go`` binary when configuring thrift.

# Areas for Future Assistance

- Providing qualification tests that automatically build against Thrift stable
as well as HEAD.

- Improving idiomaticness of generated Thrift.

- Generating Go interface code that would comply with ``gofmt`` tool.

# Continuous Integration

[![Build Status](https://secure.travis-ci.org/pomack/thrift4go.png?branch=master)](http://travis-ci.org/pomack/thrift4go)


[*] - See https://issues.apache.org/jira/browse/THRIFT-625.
