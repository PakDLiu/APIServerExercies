#!/bin/bash

if [[ $(pwd) =~ .*"scripts".* ]]; then
  cd ..
fi

go build
docker build -f docker/Dockerfile -t apiserverexercies.azurecr.io/server:latest .
