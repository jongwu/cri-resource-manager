#!/bin/bash

set -ex
cur_date="`date +%Y_%m_%d`"

#unzip cri-rm package
current_path=$(dirname  $(readlink -f "$0"))
modules_path=$current_path"/../modules"
cd $modules_path
if ! [ -d $modules_path/cri-rm ];then
  mkdir $modules_path/cri-rm
fi
tar -xf cri-rm*.tar.gz -C $modules_path/cri-rm;

pushd $modules_path/cri-rm/cri-resource-manager
export PATH=$PATH:/usr/local/go/bin
[ -f go.mod ] || go mod init
go test ./... >> crirm_test.log 2>&1
popd

log_path=$current_path"/../logs"
mkdir -p $log_path
cp $modules_path/cri-rm/cri-resource-manager/*.log $log_path
cp $log_path/crirm_test.log $log_path/run_result.txt
sed -i 's/^ok/case ~ PASSED/g' $log_path/run_result.txt
sed -i 's/FAIL/case ~ FAILED/g' $log_path/run_result.txt
