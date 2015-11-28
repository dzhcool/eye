#!/bin/bash
# go编译
ROOT_DIR=`pwd`
export GOPATH=${ROOT_DIR}
echo "[GOROOT]${GOROOT}"

go get github.com/dzhcool/eye
go install main

