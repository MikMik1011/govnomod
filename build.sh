#!/bin/bash

CGO_ENABLED=1 \
GOOS=linux \
GOARCH=386 \
go build -buildmode=c-shared -o build/govnomod.so src/main.go

if [ $? -ne 0 ]; then
    echo "Error: Build failed!"
    exit 1
fi

echo "Build succeeded! govnomod.so is created in the 'build' directory."