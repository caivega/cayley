#!/bin/bash

SRC="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

rsync -avP --delete $GOPATH/src/github.com/caivega/cayley/{static,docs,templates} $SRC/
