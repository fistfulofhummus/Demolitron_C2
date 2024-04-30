#!/bin/bash
export GOOS=windows
export GOARCH=amd64
cd Bushido
go build -ldflags -H=windowsgui
