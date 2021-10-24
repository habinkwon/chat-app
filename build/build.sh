#!/bin/bash
APP=chat-server
REPO=$GOSRC/chat-app

cd $REPO
go build -o build/$APP.out
cd $REPO/build
docker build -t $APP .
