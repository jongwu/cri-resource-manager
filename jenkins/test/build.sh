#!/bin/bash

set -ex

#prepare golang binary
apt update; apt install -y wget make
dir=$(mktemp -d)

echo "==> Download golang package"
export proxy="http://10.61.40.54:7890/"
export http_proxy=$proxy
export https_proxy=$proxy

pushd $dir
wget --no-check-certificate https://mirrors.aliyun.com/golang/go1.23.0.linux-amd64.tar.gz
tar xvf go*
[ -d /usr/local/go ] && [ -d /usr/local/go.bk ] && rm -rf /usr/local/go.bk; mv /usr/local/go /usr/local/go.bk
mv go/ /usr/local/
popd

rm -rf $dir

current_path=$(dirname  $(readlink -f "$0"))
jenkins_root_path=`realpath $current_path/../../..`

if ! [ -d $current_path/../output/modules ];then
    mkdir -p $current_path/../output/modules
fi

if ! [ -d $current_path/../output/scripts ];then
    mkdir -p $current_path/../output/scripts
fi

if ! [ -d $current_path/../output/logs ];then
    mkdir -p $current_path/../output/logs
fi

cp $current_path/run.sh $current_path/../output/scripts
cp $current_path/prepare.sh $current_path/../output/scripts
cp $current_path/config.yaml $current_path/../output/scripts

cd $current_path/../..

crirm_dir="$(pwd)/run_tests.out/run_cri-rm"

if [ -d ${crirm_dir} ];then
  rm -rf ${crirm_dir}
fi

mkdir -p ${crirm_dir}
cd ${crirm_dir}

crirm_dir=$(find $jenkins_root_path/ -name cri-resource-manager)
pushd $crirm_dir

export PATH=$PATH:/usr/local/go/bin
[ -f go.mod ] || go mod init
go env -w GOPROXY='https://goproxy.cn'

echo "==> build cri-rm"
make >> crirm_build.log 2>&1

lsdate=`date "+%Y-%m-%d"`
pushd ../
tar -czf ${jenkins_root_path}/cri-rm-${lsdate}.tar.gz ./cri-resource-manager
popd
cp ${jenkins_root_path}/cri-rm-${lsdate}.tar.gz ${jenkins_root_path}/cri-resource-manager/jenkins/output/modules
popd

