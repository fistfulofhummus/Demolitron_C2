#!/bin/bash
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1
cd Bushido
go build -o clientDebug.exe 
