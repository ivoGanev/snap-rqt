@echo off
set GOOS=linux
set GOARCH=amd64

:: Ensure the build directory exists
if not exist build mkdir build

:: Build the Go executable
go build -o build\snap-rq

