package templates

import (
	"fmt"
	"strings"
)

func escapeEmail(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(text, ".", "&#46;"), "@", "&#64;")
}

func Default(title string, message string, link string, linkTitle string) string {
	escapedTitle := escapeEmail(title)
	escapedMessage := escapeEmail(message)
	escapedLinkTitle := escapeEmail(linkTitle)

	linkSection := ""
	if link != "" {
		linkSection = fmt.Sprintf(`
            <a href="%s" style="display: inline-block; font-size: 14px; color: #ffffff; background-color: #007bff; padding: 10px 20px; text-decoration: none; border-radius: 5px;">
                %s
            </a>`, link, escapedLinkTitle,
		)
	}

	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>%s</title>
        <style>
            body {
                margin: 0;
                padding: 0;
                font-family: Arial, sans-serif;
                color: #333333;
                background-color: #f8f9fa;
            }
            .container {
                max-width: 400px;
                margin: 14px auto;
                padding: 14px;
                background-color: #ffffff;
                border: 1px solid #dddddd;
                border-radius: 10px;
                box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            }
            .header {
                text-align: center;
                padding: 10px 0;
            }
            .header h1 {
                display: inline;
                font-size: 24px;
                font-weight: bold;
                color: #007bff;
                margin: 0;
            }
            .header h1 span {
                color: #333333;
            }
            .content {
                text-align: center;
                padding: 10px 0;
            }
            .content h2 {
                font-size: 18px;
                color: #333333;
                margin: 0;
            }
            .content p {
                font-size: 14px;
                color: #666666;
                line-height: 1.5;
                margin: 10px 0;
            }
            .footer {
                text-align: center;
                padding: 10px 0;
                font-size: 12px;
                color: #666666;
            }
            .footer a {
                color: #007bff;
                text-decoration: none;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1><span style="color: #007bff;">File</span><span style="color: #333333;">space</span></h1>
            </div>
            <div class="content">
                <h2>%s</h2>
                <p>%s</p>
                %s
            </div>
            <div class="footer">
                <p>If you have any questions, contact us at
                    <a href="mailto:support@filespace.dyastin.tech">support@filespace.dyastin.tech</a>.
                </p>
            </div>
        </div>
    </body>
    </html>
    `, escapedTitle, escapedTitle, escapedMessage, linkSection)
}
