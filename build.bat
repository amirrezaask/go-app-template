@echo off

echo Building ...
go build ./cmd/server

echo Linting ...
call lint.bat
