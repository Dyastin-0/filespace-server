package html

import (
	"bytes"
	"errors"
	"html/template"
	"regexp"
)

func isHTML(input string) bool {
	re := regexp.MustCompile(`<[a-z][\s\S]*>.*?</[a-z]+>`)
	return re.MatchString(input)
}

func Check(input string) (string, error) {
	var buf bytes.Buffer

	if !isHTML(input) {
		return "", errors.New("invalid input. expected HTML content")
	}

	tmpl, err := template.New("safeHTML").Parse(input)
	if err != nil {
		return "", err
	}

	if err := tmpl.Execute(&buf, nil); err != nil {
		return "", err
	}
	return buf.String(), nil
}
