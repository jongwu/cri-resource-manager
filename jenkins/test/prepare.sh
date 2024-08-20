#!/bin/bash

#prepare golang binary
export proxy="http://10.61.40.54:7890/"
export http_proxy=$proxy
export https_proxy=$proxy

apt update; apt install -y wget make
dir=$(mktemp -d)

pushd $dir
wget --no-check-certificate https://mirrors.aliyun.com/golang/go1.23.0.linux-amd64.tar.gz
tar xvf go*
[ -d /usr/local/go ] && [ -d /usr/local/go.bk ] && rm -rf /usr/local/go.bk; mv /usr/local/go /usr/local/go.bk
mv go/ /usr/local/
popd

rm -rf $dir
