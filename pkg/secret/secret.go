package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

func getSecret(key string) (string, error) {
	ctx := context.Background()

	client, err := secretmanager.NewClient(ctx, option.WithCredentialsFile("./secretsaccesor.json"))
	if err != nil {
		return "", fmt.Errorf("filespace: Failed to create secret manager client: %v", err)
	}
	defer client.Close()

	secretName := fmt.Sprintf("projects/filespace-442811/secrets/%s/versions/1", key)

	req := &secretmanagerpb.AccessSecretVersionRequest{Name: secretName}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("filespace: Failed to access secret version: %v", err)
	}

	return string(result.Payload.Data), nil
}

func generateKeys() {
	envFilePath := filepath.Join("./.env")
	secrets := []string{
		"ACCESS_TOKEN_KEY",
		"REFRESH_TOKEN_KEY",
		"EMAIL_TOKEN_KEY",
		"MONGODB_URI",
		"SERVER_EMAIL",
		"SERVER_EMAIL_PASSWORD",
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
	}

	fmt.Println("filespace: Inserting secrets in .env file...")

	newVariables := make([]string, 0, len(secrets)+7)
	for _, secret := range secrets {
		value, err := getSecret(secret)
		if err != nil {
			fmt.Printf("filespace: Error retrieving secret %s: %v\n", secret, err)
			return
		}
		newVariables = append(newVariables, fmt.Sprintf("%s=%s", secret, value))
	}

	newVariables = append(newVariables,
		"BASE_CLIENT_URL=https://filespace.dyastin.tech",
		"PORT=3004",
		"VERSION=v2",
		"NODE_ENV=production",
		"GCLOUD_PROJECT_ID=filespace-442811",
		"GOOGLE_APPLICATION_CREDENTIALS=./secretsaccesor.json",
		"GOOGLE_REDIRECT_URL=https://filespace.dyastin.tech/api/v2/auth/google/callback",
		"DOMAIN=.filespace.dyastin.tech",
	)

	file, err := os.Open(envFilePath)
	if err != nil {
		fmt.Printf("filespace: Error opening .env file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("filespace: Error reading .env file: %v\n", err)
		return
	}

	firstLine := ""
	if len(lines) > 0 {
		firstLine = lines[0]
	}

	newContent := strings.Join(append([]string{firstLine}, newVariables...), "\n") + "\n"

	err = os.WriteFile(envFilePath, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("filespace: Error writing to .env file: %v\n", err)
		return
	}

	fmt.Printf("filespace: Successfully inserted secrets in %s\n", envFilePath)
}

func createSecretsAccessor() error {
	envFilePath := filepath.Join("./.env")
	tempFilePath := filepath.Join("./secretsaccesor.json")

	fmt.Println("filespace: Creating secretsaccesor.json file...")

	file, err := os.Create(tempFilePath)
	if err != nil {
		fmt.Printf("filespace: Failed to create file: %v\n", err)
		return err
	}
	defer file.Close()

	serviceAccount := os.Getenv("SECRETS_SERVICE_ACCOUNT")

	if serviceAccount == "" {
		return fmt.Errorf("filespace: SECRETS_SERVICE_ACCOUNT environment variable is not set")
	}

	_, err = file.WriteString(serviceAccount)
	if err != nil {
		return fmt.Errorf("filespace: Failed to write to secretsaccesor.json file: %v", err)
	}

	envFile, err := os.OpenFile(envFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("filespace: Failed to open .env file: %v", err)
	}
	defer envFile.Close()

	_, err = envFile.WriteString("\nGOOGLE_APPLICATION_CREDENTIALS=./secretsaccesor.json\n")
	if err != nil {
		return fmt.Errorf("filespace: Failed to write to .env file: %v", err)
	}

	fmt.Println("filespace: Successfully created secretsaccesor.json")
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("filespace: Error loading .env file")
	}
	createSecretsAccessor()
	generateKeys()
}
