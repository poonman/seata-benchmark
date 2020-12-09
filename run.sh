#!/bin/bash

basedir=$( cd $(dirname $0) && pwd)
cd $basedir > /dev/null

if [[ ${OS} == "Windows_NT" ]];then
    bin=seata-benchmark.exe
else
    bin=seata-benchmark
fi

ps x | grep ${bin} | grep -v grep | awk -F ' ' '{print $1}' | xargs kill -9

cd bin && ./${bin}
cd - > /dev/null

