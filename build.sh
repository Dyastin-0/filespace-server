#!/bin/bash

APP="filespace"
OUTPUT_DIR="/opt/filespace"
MAIN_PACKAGE="./cmd/server/main.go"
BINARY_NAME="run"
ENV_FILE="./.env"
ENV_PATH=$OUTPUT_DIR/.env
SECRET_FILE="./secretsaccesor.json"
SECRET_PATH=$OUTPUT_DIR/secretsaccesor.json
SERVICE_FILE="filespace-v2.service"
SERVICE_PATH="/etc/systemd/system/$SERVICE_FILE"

copy_file() {
    local source_file=$1
    local dest_file=$2
    echo "$APP: Copying $source_file to $dest_file..."
    sudo cp "$source_file" "$dest_file"
    if [ $? -eq 0 ]; then
        echo "$APP: $source_file successfully copied to $dest_file"
    else
        echo "$APP: Failed to move $source_file. Check permissions or path."
        exit 1
    fi
}

copy_file ./$SERVICE_FILE $SERVICE_PATH

sudo mkdir -p $OUTPUT_DIR

echo "$APP: Building the binary..."
sudo go build -ldflags="-s -w" -o $OUTPUT_DIR/$BINARY_NAME $MAIN_PACKAGE
if [ $? -eq 0 ]; then
    echo "$APP: Build successful. Binary located at $OUTPUT_DIR/$BINARY_NAME"
else
    echo "$APP: Build failed. Check errors above."
    exit 1
fi

echo "$APP: Setting up secrets..."
if ! go run ./pkg/secret/secret.go; then
    echo "Failed to generate .env file"
    exit 1
else
    echo "$APP: secrets set up complete"
fi

copy_file $ENV_FILE $ENV_PATH
copy_file $SECRET_FILE $SECRET_PATH

echo "$APP: Reloading systemd daemon..."
sudo systemctl daemon-reload
echo "$APP: Daemon reloaded"

if systemctl is-active --quiet $SERVICE_FILE; then
	echo "$APP: Restarting the service..."
    sudo systemctl restart $SERVICE_FILE
	echo "$APP: Service restarted"
else
	echo "$APP: Starting the service..."
    sudo systemctl start $SERVICE_FILE
	echo "$APP: Service started"
fi

echo "$APP: Restarting Caddy server..."
sudo systemctl restart caddy
echo "$APP: Caddy restarted"

echo "$APP: Service build and deployment complete"

sudo systemctl status $SERVICE_FILE
