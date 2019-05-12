#!/bin/bash

go get

# cleanup out/ directory
if [ -d out/ ]; then
  rm -rf out/*
else 
  mkdir out/
fi

env GOOS=linux GOARCH=amd64 go build -o out/sri-amd64-linux
env GOOS=linux GOARCH=386 go build -o out/sri-386-linux
env GOOS=darwin GOARCH=amd64 go build -o out/sri-amd64-darwin
env GOOS=darwin GOARCH=386 go build -o out/sri-386-darwin
env GOOS=windows GOARCH=amd64 go build -o out/sri-amd64-windows.exe
env GOOS=windows GOARCH=386 go build -o out/sri-386-windows.exe

(cd out;sha512sum * > sri.sha512sum)
