#!/usr/bin/env bash
export CGO_ENABLED=0
export VERSION=$(eval "git describe --tags --abbrev=0")
export BRANCH=$(eval "git branch | grep \* | cut -d ' ' -f2")
export REV=$(eval "git rev-parse HEAD")
go build -o defiarb -ldflags "-X main.Branch=`echo $BRANCH` -X main.Revision=`echo $REV` -X main.Version=`echo "$VERSION"` " main.go
echo "$VERSION ($BRANCH - $REV)"
