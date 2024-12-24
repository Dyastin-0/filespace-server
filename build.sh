#!/bin/bash

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
    echo "Copying $source_file to $dest_file..."
    sudo cp "$source_file" "$dest_file"
    if [ $? -eq 0 ]; then
        echo "[Filespace] $source_file successfully moved to $dest_file"
    else
        echo "Failed to move $source_file. Check permissions or path."
        exit 1
    fi
}

echo "Moving $SERVICE_FILE to $SERVICE_PATH..."
copy_file ./$SERVICE_FILE $SERVICE_PATH

sudo mkdir -p $OUTPUT_DIR

echo "[Filespace] Building the binary..."
go build -ldflags="-s -w" -o $OUTPUT_DIR/$BINARY_NAME $MAIN_PACKAGE
if [ $? -eq 0 ]; then
    echo "[Filespace] Build successful. Binary located at $OUTPUT_DIR/$BINARY_NAME"
else
    echo "[Filespace] Build failed. Check errors above."
    exit 1
fi

echo "Generating .env file..."
if ! go run ./pkg/secret/secret.go; then
    echo "Failed to generate .env file"
    exit 1
else
    echo ".env file generated successfully"
fi

copy_file $ENV_FILE $ENV_PATH
copy_file $SECRET_FILE $SECRET_PATH

echo "[Filespace] Reloading systemd daemon..."
sudo systemctl daemon-reload

if systemctl is-active --quiet $SERVICE_FILE; then
	echo "[Filespace] Restarting the service..."
    sudo systemctl restart $SERVICE_FILE
	echo "[Filespace] Service restarted"
else
	echo "[Filespace] Starting the service..."
    sudo systemctl start $SERVICE_FILE
	echo "[Filespace] Service started"
fi

sudo systemctl restart caddy
echo "[Filespace] Caddy restarted"

echo "[Filespace] Service build and deployment complete"
 
#  The script is pretty straightforward. It builds the binary, copies the necessary files to the output directory, generates the  .env  file, and moves the service file to the systemd directory. 
#  The script also restarts the service and Caddy server. 
#  To run the script, execute the following command: 
#  chmod +x build.sh