#!/bin/bash

GO=go
SRC=./src/cctvbnews.go
TASK=cctvbnews

${GO} build ${SRC} 
mv ${TASK} ./bin

