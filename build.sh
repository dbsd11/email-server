#!/bin/sh
  oldgo=$GOPATH
  GOPATH=$(cd `dirname $0`; pwd)
  GOBIN=$GOPATH/bin

  docker run --rm -ti -v "$PWD":/build -v "$PWD"/src:/go/src -e TARGETS="linux/amd64" karalabe/xgo-latest github.com/dbsd11/email-server/src

  if [ -f src-linux-amd64 ];then
    mv src-linux-amd64 email-server
  fi

  GOPATH=$oldgo
  GOBIN=$GOPATH/bin
