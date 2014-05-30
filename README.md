cavemode
========

Cavemode is a stand-alone implementation of the Orchestrate.io API for use in caves (or anywhere else without access to wifi).

The Orchestrate.io API is fully supported.  Lucene search capability is minimally supported, but is probably sufficient for use in a test/development environment.

Written in go, cavemode leverages the gorest and go-leveldb golang packages. Requires go v1.2.

Runs on MacOSX and Linux.

#### Getting started

##### Downloading the binary for Linux or MacOSX

Download the binary from https://github.com/jimcar/cavemode/releases, and install it so that can be found in $PATH.

##### Building cavemode from source

$ cd $WORK_DIR

$ git clone https://github.com/jimcar/cavemode.git cavemode

$ export GOPATH=$WORK_DIR/cavemode

$ export PATH=$GOPATH/bin:$PATH

$ go get code.google.com/p/gorest

$ go get code.google.com/p/go-leveldb

$ cd $WORK_DIR/src/github.com/jimcar/cavemode

$ go install


##### Environment variables

CAVEMODE_DB_DIR specifies the directory where the leveldb files reside.  The default is $HOME/.cavemode/leveldb-files.

CAVEMODE_JSON_INDENT, when set to "true", prettifies the json response body. Useful when exercising the API with curl or a browser.

CAVEMODE_EXTEND_GRAPH_DEPTH, when set to "true", doubles the maximum graph depth from the default of three hops to six.

##### Running cavemode

$ cavemode

Cavemode listens on port 8787, so point your client to localhost:8787 and start making requests.

