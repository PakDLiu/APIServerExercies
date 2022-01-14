#!/bin/bash

if [[ $(pwd) =~ .*"scripts".* ]]; then
  cd ..
fi

go generate ./...
go build
go test ./...

docker build -f docker/Dockerfile -t apiserverexercise.azurecr.io/server:latest .
