#!/bin/sh
  oldgo=$GOPATH
  GOPATH=$(cd `dirname $0`; pwd)
  GOBIN=$GOPATH/bin

  scp email-server root@121.42.147.36:/tmp
  ssh -o StrictHostKeyChecking=no root@121.42.147.36 \"/tmp/email-server\"

  GOPATH=$oldgo
  GOBIN=$GOPATH/bin
