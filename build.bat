@echo off
setlocal enabledelayedexpansion

REM Get version from git or use default
for /f "tokens=*" %%i in ('git describe --tags --always 2^>nul') do set VERSION=%%i
if "!VERSION!"=="" set VERSION=1.0.1

REM Get current date in YYYY-MM-DD format
for /f "tokens=2,3,4 delims=/ " %%a in ('echo %date%') do set BUILD_DATE=%%c-%%a-%%b

REM Get commit hash
for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set COMMIT_HASH=%%i
if "!COMMIT_HASH!"=="" set COMMIT_HASH=unknown

REM Build with version information
go mod tidy
go build -ldflags "-X github.com/acheddir/git-contrib/cmd.Version=!VERSION! -X github.com/acheddir/git-contrib/cmd.BuildDate=!BUILD_DATE! -X github.com/acheddir/git-contrib/cmd.CommitHash=!COMMIT_HASH!"

echo Build completed successfully.