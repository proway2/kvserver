#!/bin/bash
GO=`which go`
cd $GOPATH

$GO test kvserver/kvstorage \
kvserver/vacuum \
kvserver/router \
-cover kvserver/kvstorage \
kvserver/vacuum \
kvserver/router 
