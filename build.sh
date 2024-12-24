#!/bin/bash

OUTPUT_DIR="/opt/filespace"
MAIN_PACKAGE="./cmd/server/main.go"
ENV_PATH="./.env"
BINARY_NAME="run"
SERVICE_FILE="filespace.service"
SERVICE_PATH="/etc/systemd/system/$SERVICE_FILE"

sudo mkdir -p $OUTPUT_DIR

echo "Building the application..."
sudo go build -ldflags="-s -w" -o $OUTPUT_DIR/$BINARY_NAME $MAIN_PACKAGE

if [ $? -eq 0 ]; then
    echo "Build successful. Binary located at $OUTPUT_DIR/$BINARY_NAME"
else
    echo "Build failed. Check errors above."
	exit 1
fi

if [ -f "./$SERVICE_FILE" ]; then
	echo "Moving $SERVICE_FILE to $SERVICE_PATH..."
	sudo cp ./$SERVICE_FILE $SERVICE_PATH

	echo "Generating .env file..."
	if go run ./pkg/secret/secret.go; then
		echo ".env file generated successfully."
	else
		echo "Failed to generate .env file."
		exit 1
	fi

	echo "Moving .env file to $OUTPUT_DIR..."
	sudo cp $ENV_PATH $OUTPUT_DIR/.env

	if [ $? -eq 0 ]; then
		echo "Service file moved successfully to $SERVICE_PATH."
		sudo systemctl daemon-reload
		sudo systemctl enable $SERVICE_FILE --now 
		sudo systemctl restart caddy
		echo "Service started and enabled."
	else
		echo "Failed to move service file. Check permissions or path."
		exit 1
	fi
else
	echo "Service file $SERVICE_FILE not found in the root directory."
	exit 1
fi
