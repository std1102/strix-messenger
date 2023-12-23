#!/bin/bash

GOOS=js GOARCH=wasm go build -o ./static/main.wasm ./main.go

echo "build completed, wasm file is located in folder 'static'"