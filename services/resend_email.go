package services

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

type TrainingReminderData struct {
	Username       string
	Openings       []OpeningTrainingData
	TotalPositions int
}

type OpeningTrainingData struct {
	Name          string
	Side          string
	PositionCount int
}

func SendTrainingReminderEmail(to string, data TrainingReminderData) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY environment variable not set")
	}

	fromEmail := os.Getenv("RESEND_FROM_EMAIL")
	if fromEmail == "" {
		return fmt.Errorf("RESEND_FROM_EMAIL environment variable not set")
	}

	if to == "" {
		return fmt.Errorf("recipient email is empty")
	}

	fmt.Printf("Sending training reminder with Resend:\n")
	fmt.Printf("- API Key: %s...\n", apiKey[:10])
	fmt.Printf("- From: %s\n", fromEmail)
	fmt.Printf("- To: %s\n", to)

	client := resend.NewClient(apiKey)

	// Generate opening list HTML
	openingListHTML := ""
	for _, opening := range data.Openings {
		sideName := "White"
		if opening.Side == "b" {
			sideName = "Black"
		}

		openingListHTML += fmt.Sprintf(`
			<tr>
				<td style="padding: 12px; border-bottom: 1px solid #e5e7eb;">
					<div style="display: flex; align-items: center; gap: 8px;">
						<span style="font-weight: 500;">%s</span>
					</div>
				</td>
				<td style="padding: 12px; border-bottom: 1px solid #e5e7eb;">
					<span style="padding: 4px 8px; background-color: %s; color: %s; border-radius: 4px; font-size: 12px; font-weight: 500;">
						%s
					</span>
				</td>
				<td style="padding: 12px; border-bottom: 1px solid #e5e7eb; text-align: right;">
					<span style="font-weight: 600; color: #059669;">%d positions</span>
				</td>
			</tr>`,
			opening.Name,
			map[string]string{"b": "#1f2937", "w": "#f3f4f6"}[opening.Side],
			map[string]string{"b": "white", "w": "#1f2937"}[opening.Side],
			sideName, opening.PositionCount)
	}

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Time for Chess Training! - Chesso</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6; 
            color: #374151; 
            margin: 0; 
            padding: 0; 
            background-color: #f9fafb; 
        }
        .container { 
            max-width: 600px; 
            margin: 20px auto; 
            background-color: white; 
            border-radius: 12px; 
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1); 
            overflow: hidden;
        }
        .header { 
            background: linear-gradient(135deg, #1f2937 0%%, #111827 100%%);
            color: white; 
            text-align: center; 
            padding: 40px 20px; 
        }
        .chess-icon {
            display: flex;
            justify-content: center;
            align-items: center;
            font-size: 48px;
            margin-bottom: 16px;
        }
        .chess-icon img {
            height: 48px;
            width: 48px;
            margin-right: 12px;
        }
        .content { 
            padding: 32px; 
        }
        .greeting {
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 16px;
        }
        .stats-card {
            background: linear-gradient(135deg, #059669 0%%, #047857 100%%);
            color: white;
            padding: 24px;
            border-radius: 12px;
            text-align: center;
            margin: 24px 0;
        }
        .stats-number {
            font-size: 36px;
            font-weight: 700;
            margin-bottom: 8px;
        }
        .openings-table {
            width: 100%%;
            border-collapse: collapse;
            margin: 24px 0;
            border-radius: 8px;
            overflow: hidden;
            border: 1px solid #e5e7eb;
        }
        .table-header {
            background-color: #f9fafb;
            font-weight: 600;
            color: #374151;
        }
        .table-header th {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #e5e7eb;
        }
        .button {
            display: inline-block; 
            background: linear-gradient(135deg, #1f2937 0%%, #111827 100%%);
            color: white !important; 
            padding: 16px 32px; 
            text-decoration: none; 
            border-radius: 8px; 
            font-weight: 600;
            font-size: 16px;
            margin: 24px 0;
            transition: transform 0.2s ease;
        }
        .button:hover { 
            transform: translateY(-2px);
        }
        .footer { 
            background-color: #f9fafb; 
            padding: 24px; 
            text-align: center; 
            font-size: 14px; 
            color: #6b7280; 
        }
        .tip {
            background-color: #fef3c7;
            border-left: 4px solid #f59e0b;
            padding: 16px;
            margin: 24px 0;
            border-radius: 4px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1 class="logo">Chesso</h1>
            <h1>Time for Chess Training!</h1>
            <p>Your positions are ready for review</p>
        </div>
        <div class="content">
            <div class="greeting">Hello %s!</div>
            <p>Your chess positions are waiting for some practice. It's time to sharpen your tactical skills!</p>
            
            <div class="stats-card">
                <div class="stats-number">%d</div>
                <div>Positions Ready to Train</div>
            </div>
            
            <h3 style="margin-top: 32px; margin-bottom: 16px; color: #1f2937;">Your Training Opportunities:</h3>
            
            <table class="openings-table">
                <thead class="table-header">
                    <tr>
                        <th>Opening</th>
                        <th>Side</th>
                        <th>Positions</th>
                    </tr>
                </thead>
                <tbody>
                    %s
                </tbody>
            </table>
            
            <div style="text-align: center; margin: 32px 0;">
                <a href="https://chesso.org/app/train" class="button">Start Training Now</a>
            </div>
            
            <div class="tip">
                <strong>💡 Training Tip:</strong> Regular practice with spaced repetition helps improve your pattern recognition and tactical awareness. Even 10-15 minutes of focused training can make a significant difference!
            </div>
        </div>
        <div class="footer">
            <p>This is an automated training reminder from Chesso.</p>
            <p>© 2024 Chesso Chess Application</p>
            <p style="margin-top: 16px; font-size: 12px;">
                <a href="https://chesso.org/unsubscribe" style="color: #6b7280;">Unsubscribe from training reminders</a>
            </p>
        </div>
    </div>
</body>
</html>`, data.Username, data.TotalPositions, openingListHTML)

	// Generate opening list for text version
	openingListText := ""
	for _, opening := range data.Openings {
		sideName := "White"
		if opening.Side == "b" {
			sideName = "Black"
		}
		openingListText += fmt.Sprintf("• %s (%s): %d positions\n", opening.Name, sideName, opening.PositionCount)
	}

	textContent := fmt.Sprintf(`Time for Chess Training! - Chesso

Hello %s!

Your chess positions are waiting for some practice. It's time to sharpen your tactical skills!

POSITIONS READY TO TRAIN: %d

Your Training Opportunities:
%s

Start training now: https://chesso.org/app/train

💡 Training Tip: Regular practice with spaced repetition helps improve your pattern recognition and tactical awareness. Even 10-15 minutes of focused training can make a significant difference!

---
This is an automated training reminder from Chesso.
© 2024 Chesso Chess Application

Unsubscribe from training reminders: https://chesso.org/unsubscribe`, data.Username, data.TotalPositions, openingListText)

	params := &resend.SendEmailRequest{
		From:    fromEmail,
		To:      []string{to},
		Subject: fmt.Sprintf("🏋️ %d Chess Positions Ready for Training!", data.TotalPositions),
		Html:    htmlContent,
		Text:    textContent,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Printf("Resend API Error: %v\n", err)
		return fmt.Errorf("failed to send training reminder email via Resend: %v", err)
	}

	fmt.Printf("Training reminder email sent successfully! ID: %s\n", sent.Id)
	return nil
}

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
