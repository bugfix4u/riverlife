#!/bin/bash

PROJECT_NAME="riverlife"
ROOT_PROJECT=/root/rlcode
ROOT_SRC=$GOPATH/src/$PROJECT_NAME

# User passed in the microservice to "Go" build
SERVICE=$1
echo -e "Compiling ${SERVICE}"

# Build Go code for each micro-service
cd ${ROOT_SRC}/cmd/${SERVICE}

GO_FILES=$(find . -name "*.go" | head -n1)
if [ -z ${GO_FILES} ];then
   echo -e "\t...nothing to build"
   exit 0
fi

go build
if [ $? -eq 0 ];then
   echo -e "\t... Build: SUCCESS"
else
   echo -e "\t... ERROR: build failed"
   exit -1
fi

go install
if [ $? -eq 0 ];then
   echo -e "\t... Install: SUCCESS"
else
   echo -e "\t... ERROR: install failed"
   exit -1
fi

cp ${ROOT_PROJECT}/bin/${SERVICE} ../../builds/${SERVICE}

echo -e "\tCompiled ${SERVICE}: SUCCESS"
