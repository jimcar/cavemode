cavemode
========

Cavemode is a stand-alone implementation of the [Orchestrate.io API](https://orchestrate.io/docs/api/) for use in caves (or anywhere else without access to wifi).

The Orchestrate.io API is fully supported.  Search capability is minimally supported, but is probably sufficient for use in a test/development environment.

Written in go, cavemode leverages the gorest and go-leveldb golang packages. Requires go v1.2.

Runs on MacOSX and Linux.

#### Getting started

##### Downloading the binary for Linux or MacOSX

Download the binary from https://github.com/jimcar/cavemode/releases, and install it so that can be found in $PATH.

##### Building from source

To build from source, you must have a Go 1.2 development environment.

###### MacOSX: you may need to install mercurial using homebrew: "brew install mercurial"

$ cd $WORK_DIR

$ git clone https://github.com/jimcar/cavemode.git cavemode

$ cd cavemode

$ mkdir bin pkg

$ export GOPATH=$WORK_DIR/cavemode

$ export PATH=$GOPATH/bin:$PATH

$ go get code.google.com/p/gorest

$ go get code.google.com/p/go-leveldb

###### MacOSX: Temporary workaround for go-leveldb compilation issue (not detecting OS)
Add the following line to src/code.google.com/p/go-leveldb/port/port_posix.h, before line 11:

 #define OS_MACOSX

$ cd $WORK_DIR/src/github.com/jimcar/cavemode

$ go install


##### Environment variables

CAVEMODE_DB_DIR specifies the directory where the leveldb files reside.  The default is $HOME/.cavemode/leveldb-files.

CAVEMODE_JSON_INDENT, when set to "true", prettifies the json response body. Useful when exercising the API with curl or a browser.

CAVEMODE_EXTEND_GRAPH_DEPTH, when set to "true", doubles the maximum graph depth from the default of three hops to six.

##### Running cavemode

$ cavemode

Cavemode listens on port 8787, so point your client to localhost:8787 and start making requests.

#### Status

Cavemode is fully-functional, but remains rough around the edges. There are a number of known issues, and certainly many more yet to be discovered.

Community participation is highly encouraged, so make a pull request and improve something!
