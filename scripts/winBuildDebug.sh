#!/bin/zsh
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1
cd cmd/Bushido/
go build -o Compiled/clientDebug.exe
