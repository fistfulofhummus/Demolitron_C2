#!/bin/bash
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1
cd cmd/Bushido/
go build --ldflags -H=windowsgui -o Compiled/client.exe
