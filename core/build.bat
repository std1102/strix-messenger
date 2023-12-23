@ECHO OFF
set GOOS=js
set GOARCH=wasm
go build -o .\static\main.wasm .\main.go
ECHO build completed, wasm file is located in folder 'static'
PAUSE