#!/bin/sh
cd src
go build server.go boardactions.go
./server
