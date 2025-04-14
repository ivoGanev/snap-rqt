@echo off
set GOOS=windows
set GOARCH=amd64

:: Ensure the build directory exists
if not exist build mkdir build

:: Build the Go executable
go build -o build\snap-rq.exe

:: Check if build succeeded
if errorlevel 1 (
    echo Build failed. Exiting.
    exit /b %errorlevel%
)

:: Run the executable
build\snap-rq.exe
