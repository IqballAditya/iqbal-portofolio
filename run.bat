@echo off
title Iqbal Netflix Portfolio
echo ==================================================
echo   Starting Iqbal Netflix Portfolio Server...
echo ==================================================
echo Opening http://localhost:8080 in your default browser...
start http://localhost:8080
echo Running server on port 8080. Press Ctrl+C to stop.
echo.
"%~dp0bin\iqbal-portfolio.exe"
pause
