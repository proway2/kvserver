#!/bin/bash
GO=`which go`
cd $GOPATH

$GO test kvserver/element \
kvserver/kvstorage \
kvserver/vacuum \
kvserver \
-cover kvserver/element \
kvserver/kvstorage \
kvserver/vacuum \
kvserver
