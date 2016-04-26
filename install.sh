#!/bin/bash

CUR=$(cd "$(dirname "$0")"; pwd)"/"
GO=go
SRC=./src/cctvbnews.go
TASK=cctvbnews

${GO} build ${SRC} 
mv ${TASK} ./bin

zCUR=`python fileName.py ${CUR}`
sed -i "s/\${DIR}/${zCUR}/g" ./conf/config.ini
sed -i "s/\${DIR}/${zCUR}/g" ./conf/log.xml
