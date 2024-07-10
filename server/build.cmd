@echo off
setlocal enabledelayedexpansion

REM Get the last tag sorted by version
for /f "delims=" %%i in ('git tag --sort=-v:refname ^| findstr /r /n "^" ^| findstr /b "1:"') do set "LAST_TAG=%%i"

REM Remove the line number prefix if it exists
set "LAST_TAG=%LAST_TAG:1:=%"

REM Build the Go application with the version from the last tag
go build -ldflags="-s -w -X 'main.Version=%LAST_TAG%'"

endlocal