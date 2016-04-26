#!/bin/bash

CUR=$(cd "$(dirname "$0")"; pwd)"/"
GO=go
SRC=./src/cctvbnews.go
TASK=$1

${GO} build ${SRC} 
mv cctvbnews ./bin/${TASK}

zCUR=`python fileName.py ${CUR}`
sed -i "s/\${DIR}/${zCUR}/g" ./conf/config.ini
sed -i "s/\${DIR}/${zCUR}/g" ./conf/log.xml
sed -i "s/\${DIR}/${zCUR}/g" ./bin/control.sh
sed -i "s/cctvbnews/${TASK}/g" ./bin/control.sh
