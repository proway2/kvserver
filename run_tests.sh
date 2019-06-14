#!/bin/bash
GO=`which go`
cd $GOPATH/github/proway2/kvserver

$GO test github.com/proway2/kvserver/kvstorage \
github.com/proway2/kvserver/vacuum \
github.com/proway2/kvserver/router \
-cover github.com/proway2/kvserver/kvstorage \
github.com/proway2/kvserver/vacuum \
github.com/proway2/kvserver/router 
