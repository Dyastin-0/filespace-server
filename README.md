# Filespace Build and Deploy Script

This repository includes a custom `build.sh` script to simplify the build and deployment process for the **Filespace** application.

## Quick Start

1. Clone the repository:

   SSH  
   ```bash
   git clone git@github.com:Dyastin-0/filespace-server.git
   ```
   HTTPS
   ```bash
   git clone https://github.com/Dyastin-0/filespace-server.git
   ```

   ```bash
   cd filespace-server
   ```


2. Create a `.env` file and set the `SECRETS_SERVICE_ACCOUNT` to the service account json (in a single line). 

   ```bash
   nano .env
   ```


3. Make the build script executable:

   ```bash
   chmod +x build.sh
   ```


4. Run the script:

   ```bash
   ./build.sh
   ```

## What the Script Does
The script is tightly integrated with `Caddy` & `systemd` and performs the following tasks:

1. **Builds the Application**:
   - Compiles the Go application into a binary named `run` and places it in `/opt/filespace`.

2. **Sets Up Secrets**:
   - Runs the `secret.go` file to generate an `.env` file and moves it along with the `secretsaccesor.json` file to `/opt/filespace`.

3. **Configures the Service**:
   - Copies the `filespace-v2.service` file to `/etc/systemd/system/` and reloads the systemd daemon to register the service.

4. **Manages the Service**:
   - Checks if the service is active and restarts it if necessary. If inactive, it starts the service.

5. **Restarts the Caddy Server**:
   - Ensures the Caddy reverse proxy is restarted to apply any updates.

6. **Verifies Deployment**:
   - Displays the status of the `filespace-v2` service for confirmation.

## What are Things for

1. **`secretsaccesor.json`**
   - Used by Google Cloud services like Storage and Secret Manager. The `.env` file contains the variable `GOOGLE_APPLICATION_CREDENTIALS=./secretsaccesor.json`, which points to this file for authentication.

2. **`.env`**
   - Stores secret keys retrieved from Secret Manager, including credentials for email services and JWT token keys. It's critical for securely configuring the application.

3. **`filespace-v2.service`**
   - A systemd service file that defines how the Filespace application runs as a background service.

## File Locations

- **Binary**: `/opt/filespace/run`
- **Environment File**: `/opt/filespace/.env`
- **Secrets File**: `/opt/filespace/secretsaccesor.json`
- **Service File**: `/etc/systemd/system/filespace-v2.service`

## Useful Commands

- Check service logs:

  ```bash
  sudo journalctl -u filespace-v2.service
  ```

- Reload systemd manually:

  ```bash
  sudo systemctl daemon-reload
  ```

- Restart the service:

  ```bash
  sudo systemctl restart filespace-v2.service
  ```

- Check service status:

  ```bash
  sudo systemctl status filespace-v2.service
  ```

---

Keep secrets safe and ensure permissions are correctly configured for production use.

