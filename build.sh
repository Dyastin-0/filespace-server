#!/bin/bash

OUTPUT_DIR="/opt/filespace"
MAIN_PACKAGE="./cmd/server/main.go"
BINARY_NAME="filespace"

mkdir -p $OUTPUT_DIR

echo "Building the application..."
go build -ldflags="-s -w" -o $OUTPUT_DIR/$BINARY_NAME $MAIN_PACKAGE

if [ $? -eq 0 ]; then
  echo "Build successful. Binary located at $OUTPUT_DIR/$BINARY_NAME"
else
  echo "Build failed. Check errors above."
  exit 1
fi