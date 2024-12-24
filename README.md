# Table of Contents

- [Introduction](#introduction)
- [Endpoints](#filespace-endpoints)
- [Build and Deploy](#filespace-build-and-deploy-script)

# Introduction

***Filespace*** is a simple cloud storage utilizing the ***Google Cloud Storage API***. It includes features such as uploading, deleting, moving, sharing and viewing files, and creating folders.

This service is coupled with the [Filespace Frontend](https://github.com/Dyastin-0/filespace) built with Vite; the [repo](https://github.com/Dyastin-0/filespace) also includes the v1 API written in ***Node.js***.

# Filespace Endpoints

1. Files `/api/v2/files`

   - `GET /api/v2/files`

      Files metadata are fetched based on the present `Bearer <t>` on the `Authorization` header of the request automatically handled by the `JWT` middleware.
   
      Returns an array of:

      ```go
         type Metadata struct {
         	Name        string
         	Link        string
         	Owner       string
         	Size        int64
         	Updated     time.Time
         	ContentType string
         	Created     time.Time
         	Type        string
         }
      ```

      on success

   - `POST /api/v2/files`
  
      Uploads the files inside the `formdata`; automatically handles the distinction between a file and a folder. Similar to `GET`, files are uploaded based on the present `Bearer <t>` on the `Authorization` header of the request.

      Expects:

      ```go
         files := r.MultipartForm.File["files"]
         path := r.FormValue("path")
         folder := r.FormValue("folder")
      ```
      Returns: `status 201` on success

   - `DELETE /api/v2/files`

      Deletes the specified files in the request body of the authenticated user; automatically handles the distinction between a file and a folder.
      
      Expects:

      ```go
         type DeleteBody struct {
            Files []string `json:"files"`
         }
      ```

      Returns: `status 200` on success

   - `POST /api/v2/files/move`

      Moves the specified file into another location, this process includes copying the specified file into the `TargetPath`, or files if the specified path is a folder, and deletes the file located on the previous path.

      Expects:

      ```go
         type file struct {
            Name string `json:"name"`
            Path string `json:"path"`
            Type string `json:"type"`
         }
         
         type MoveBody struct {
            File       file   `json:"file"`
            TargetPath string `json:"targetPath"`
         }
      ```

      Returns: `status 200` on success

   - `POST /api/v2/files/share`
  
      Sends the specified file to the specified email on the request body.

      Expects:

      ```go
         type expiration struct {
            Value int64  `json:"value"`
            Str   string `json:"text"`
         }
         
         type ShareBody struct {
            Email string     `json:"email"`
            File  string     `json:"file"`
            Exp   expiration `json:"expiration"`
         }
      ```

      Returns: `status 200` on success

2. Authentication `/api/v2/auth`

   - `POST /api/v2/auth`

      Expects:

      ```go
         type Body struct {
            Email: string `json:"email"`
            Password: string `json:"password"`
         }
      ```
      Returns:
      
      ```go
         type Response struct {
         AccessToken string `json:"accessToken"`
         User        User
         }
         
         type User struct {
            ID          string               `bson:"_id,omitempty"`
            Username    string               `bson:"username"`
            Email       string               `bson:"email"`
            Roles       []string             `bson:"roles"`
            ImageURL    string               `bson:"profileImageURL"`
            UsedStorage primitive.Decimal128 `bson:"usedStorage"`
         }
      ```
      on success

   - `POST /api/v2/auth/sign-up`

      Expects:

      ```go
         type SignupBody struct {
            Email    string `json:"email"`
            Password string `json:"password"`
            Username string `json:"username"`
         }
      ```

      returns `status 201` on success

   - `POST /api/v2/auth/refresh`

      Expects:

      ```go
         r.Cookie("jwt")
      ```

      returns the same data as `POST /api/v2/auth` on success

   - `/api/v2/auth/verify`

      Expects:

      ```go
         query.Get("t")
      ```
      Returns the same data as `POST /api/v2/auth` on success

   - `/api/v2/auth/send-verification`

      Expects:

      ```go
         type VerificationBody struct {
            Email string `json:"email"`
         }
      ```

      Returns: `status 200` on success

   - `/api/v2/auth/recover`

      Expects:

      ```go
      type RecoverBody struct { //WHAT BODY?
         Token       string `json:"t"`
         NewPassword string `json:"newPassword"`
      }
      ```

      Returns: `status 200` on success, does not automatically log ins user as `POST /api/v2/auth/verify` do

   - `/api/v2/auth//send-recovery`

      Expects:

      ```go
         type SendRecoveryBody struct {
         Email string `json:"email"`
      }
      ```

      Returns: `status 200` on success

   - `/api/v2/auth/log-out`

      Simply clears the jwt in the cookie, if present remove it from the database

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

2. Create a `.env` file and set the `SECRETS_SERVICE_ACCOUNT` to the service account json (in a single line):
   
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
