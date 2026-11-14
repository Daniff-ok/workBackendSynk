@echo off
setlocal
cd /d "%~dp0"
if errorlevel 1 goto :error

where git >nul 2>&1
if errorlevel 1 (
    echo ERROR: Install Git for Windows and add it to PATH.
    goto :error
)

where go >nul 2>&1
if errorlevel 1 (
    echo ERROR: Install Go and add it to PATH.
    goto :error
)

echo Updating files from Git...
git pull --ff-only
if errorlevel 1 (
    echo ERROR: Git update failed. Check the connection and local changes.
    goto :error
)

echo Starting server...
go run .
if errorlevel 1 goto :error
exit /b 0

:error
echo.
echo Server stopped or could not be started.
pause
exit /b 1
