package mail

import (
	"log"
	"os"
	"testing"

	mailTemplate "filespace/pkg/mail/template"

	godotenv "github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestSendPlainTextEmail(t *testing.T) {
	tests := []struct {
		options *Message
		wantErr bool
	}{
		{&Message{
			To:      os.Getenv("TEST_EMAIL"),
			Subject: "", Body: "",
			ContentType: PlainTextEmail}, true},
		{&Message{
			To:          os.Getenv("TEST_EMAIL"),
			Subject:     "Test",
			Body:        "Test",
			ContentType: PlainTextEmail}, false},
		{&Message{
			To:      os.Getenv("TEST_EMAIL"),
			Subject: "Test", Body: "",
			ContentType: PlainTextEmail}, true},
	}

	for _, tt := range tests {
		if err := SendPlainTextEmail(tt.options); (err != nil) != tt.wantErr {
			t.Errorf("SendPlainTextEmail() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func TestSendHTMLEmail(t *testing.T) {
	tests := []struct {
		options *Message
		wantErr bool
	}{
		{&Message{
			To:      os.Getenv("TEST_EMAIL"),
			Subject: "", Body: "",
			ContentType: PlainTextEmail}, true},
		{&Message{
			To:      os.Getenv("TEST_EMAIL"),
			Subject: "Test", Body: "Test",
			ContentType: HTMLTextEmail}, true},
		{&Message{
			To:      os.Getenv("TEST_EMAIL"),
			Subject: "Test",
			Body: mailTemplate.Default("Test",
				"A test email.",
				"test.com",
				"Test Link",
			),
			ContentType: HTMLTextEmail}, false},
	}

	for _, tt := range tests {
		if err := SendHTMLEmail(tt.options); (err != nil) != tt.wantErr {
			t.Errorf("SendHTMLEmail() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
