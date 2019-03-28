#!/bin/bash

export GOPATH=$(pwd)
mv vendor src
go get $1
mv src vendor

echo "Done."