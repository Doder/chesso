package services

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func SendPasswordResetEmailWithResend(to, token string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY environment variable not set")
	}

	fromEmail := os.Getenv("RESEND_FROM_EMAIL")
	if fromEmail == "" {
		return fmt.Errorf("RESEND_FROM_EMAIL environment variable not set")
	}

	// Validate email domain for development
	if to == "" {
		return fmt.Errorf("recipient email is empty")
	}

	// Debug logging
	fmt.Printf("Sending email with Resend:\n")
	fmt.Printf("- API Key: %s...\n", apiKey[:10])
	fmt.Printf("- From: %s\n", fromEmail)
	fmt.Printf("- To: %s\n", to)

	client := resend.NewClient(apiKey)

	resetURL := fmt.Sprintf("https://chesso.org/app/reset-password/%s", token)

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Password Reset - Chesso</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            line-height: 1.6; 
            color: #333; 
            margin: 0; 
            padding: 0; 
            background-color: #f4f4f4; 
        }
        .container { 
            max-width: 600px; 
            margin: 20px auto; 
            background-color: white; 
            border-radius: 8px; 
            box-shadow: 0 2px 10px rgba(0,0,0,0.1); 
            overflow: hidden;
        }
        .header { 
            background-color: #171717; 
            color: white; 
            text-align: center; 
            padding: 30px 20px; 
        }
        .content { 
            padding: 30px; 
        }
        .button {

            display: inline-block; 
            background-color: #171717; 
            color: white !important; 
            padding: 14px 28px; 
            text-decoration: none; 
            border-radius: 6px; 
            font-weight: bold;
            margin: 20px 0;
        }
        .button:hover { 
            background-color:#171717; 
        }
        .footer { 
            background-color: #f8f9fa; 
            padding: 20px; 
            text-align: center; 
            font-size: 12px; 
            color: #6b7280; 
        }
        .warning {
            background-color: #fef3c7;
            border-left: 4px solid #f59e0b;
            padding: 15px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
            <p>Chesso Chess Application</p>
        </div>
        <div class="content">
            <h2>Hello!</h2>
            <p>You requested a password reset for your Chesso account. Click the button below to reset your password:</p>
            
            <div style="text-align: center;">
                <a href="%s" class="button">Reset My Password</a>
            </div>
            
            <p>Or copy and paste this link in your browser:</p>
            <p style="word-break: break-all; background-color: #f8f9fa; padding: 10px; border-radius: 4px;"><a href="%s">%s</a></p>
            
            <div class="warning">
                <strong>Important:</strong> This link will expire in 1 hour for security reasons.
            </div>
            
            <p>If you didn't request this password reset, please ignore this email. Your account remains secure.</p>
        </div>
        <div class="footer">
            <p>This is an automated message from Chesso.</p>
            <p>© 2024 Chesso Chess Application</p>
        </div>
    </div>
</body>
</html>`, resetURL, resetURL, resetURL)

	textContent := fmt.Sprintf(`Password Reset - Chesso

Hello!

You requested a password reset for your Chesso account.

Reset your password by visiting this link:
%s

This link will expire in 1 hour.

If you didn't request this password reset, please ignore this email.

---
Chesso Chess Application
This is an automated message.`, resetURL)

	params := &resend.SendEmailRequest{
		From:    fromEmail,
		To:      []string{to},
		Subject: "Password Reset - Chesso",
		Html:    htmlContent,
		Text:    textContent,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Printf("Resend API Error: %v\n", err)
		return fmt.Errorf("failed to send email via Resend: %v", err)
	}

	fmt.Printf("Email sent successfully! ID: %s\n", sent.Id)
	return nil
}
