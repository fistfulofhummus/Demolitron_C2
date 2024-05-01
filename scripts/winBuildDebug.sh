#!/bin/bash
export GOOS=windows
export GOARCH=amd64
cd Bushido
go build -o clientDebug.exe
