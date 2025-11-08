package helpers

import (
	"fmt"
	"net/smtp"

	"api.workzen.odoo/constants"
)

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// GetEmailConfig retrieves email configuration from environment variables
func GetEmailConfig() *EmailConfig {
	return &EmailConfig{
		SMTPHost:     constants.SMTPHost,
		SMTPPort:     constants.SMTPPort,
		SMTPUser:     constants.SMTPUsername,
		SMTPPassword: constants.SMTPPassword,
		FromEmail:    constants.SMTPFromEmail,
		FromName:     constants.SMTPFromName,
	}
}

// SendEmail sends an email using SMTP
func SendEmail(to, subject, body string) error {
	config := GetEmailConfig()

	if constants.ServerMode == "development" {
		// In development mode, just log the email instead of sending
		fmt.Printf("\n=== EMAIL (Development Mode) ===\n")
		fmt.Printf("To: %s\n", to)
		fmt.Printf("Subject: %s\n", subject)
		fmt.Printf("Body:\n%s\n", body)
		fmt.Printf("================================\n\n")
		return nil
	}

	if config.SMTPUser == "" || config.SMTPPassword == "" {
		return fmt.Errorf("SMTP credentials are not configured")
	}

	// Create the email message
	from := fmt.Sprintf("%s <%s>", config.FromName, config.FromEmail)
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", from, to, subject, body))

	// Set up authentication
	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPassword, config.SMTPHost)

	// Send the email
	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)
	err := smtp.SendMail(addr, auth, config.FromEmail, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendVerificationEmail sends an email verification link to the user
func SendVerificationEmail(to, token, companyName string) error {
	// In production, this should be your actual frontend URL
	frontendURL := constants.FrontendURL
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, token)

	subject := "Verify Your Email - WorkZen HRMS"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { background-color: #f9f9f9; padding: 30px; border-radius: 5px; margin-top: 20px; }
        .button { 
            display: inline-block; 
            padding: 12px 30px; 
            background-color: #4CAF50; 
            color: white; 
            text-decoration: none; 
            border-radius: 5px; 
            margin: 20px 0;
        }
        .footer { text-align: center; margin-top: 30px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to WorkZen HRMS</h1>
        </div>
        <div class="content">
            <h2>Verify Your Email Address</h2>
            <p>Hello,</p>
            <p>Thank you for registering your company <strong>%s</strong> with WorkZen HRMS.</p>
            <p>Please click the button below to verify your email address and complete your registration:</p>
            <div style="text-align: center;">
                <a href="%s" class="button">Verify Email</a>
            </div>
            <p>Or copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #666;">%s</p>
            <p><strong>Note:</strong> This verification link will expire in 24 hours.</p>
            <p>If you didn't create an account with WorkZen, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>© 2025 WorkZen HRMS. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, companyName, verificationLink, verificationLink)

	return SendEmail(to, subject, body)
}

// SendWelcomeEmail sends a welcome email after successful verification
func SendWelcomeEmail(to, firstName, companyName, username string) error {
	frontendURL := constants.FrontendURL
	loginLink := fmt.Sprintf("%s/login", frontendURL)

	subject := "Welcome to WorkZen HRMS - Your Account is Active!"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { background-color: #f9f9f9; padding: 30px; border-radius: 5px; margin-top: 20px; }
        .button { 
            display: inline-block; 
            padding: 12px 30px; 
            background-color: #4CAF50; 
            color: white; 
            text-decoration: none; 
            border-radius: 5px; 
            margin: 20px 0;
        }
        .info-box { 
            background-color: #e8f5e9; 
            padding: 15px; 
            border-radius: 5px; 
            margin: 20px 0; 
        }
        .footer { text-align: center; margin-top: 30px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Welcome to WorkZen HRMS!</h1>
        </div>
        <div class="content">
            <h2>Your Account is Now Active</h2>
            <p>Hello %s,</p>
            <p>Great news! Your company <strong>%s</strong> has been successfully registered and verified.</p>
            <p>Your account is now pending approval from our administrators. Once approved, you'll be able to:</p>
            <ul>
                <li>Manage your company employees</li>
                <li>Track attendance and leaves</li>
                <li>Process payroll</li>
                <li>Generate reports and more!</li>
            </ul>
            <div class="info-box">
                <strong>Your Login Credentials:</strong><br>
                Username: <strong>%s</strong><br>
                Use the password you created during registration.
            </div>
            <p>You will receive another notification once your account is approved by the administrators.</p>
            <div style="text-align: center;">
                <a href="%s" class="button">Go to Login</a>
            </div>
        </div>
        <div class="footer">
            <p>© 2025 WorkZen HRMS. All rights reserved.</p>
            <p>If you have any questions, please contact our support team.</p>
        </div>
    </div>
</body>
</html>
`, firstName, companyName, username, loginLink)

	return SendEmail(to, subject, body)
}

// SendEmployeeInvitationEmail sends invitation email to new employee with verification link and temporary password
func SendEmployeeInvitationEmail(to, firstName, companyName, username, tempPassword, token string) error {
	frontendURL := constants.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, token)

	subject := "Welcome to WorkZen HRMS - Verify Your Account"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { background-color: #f9f9f9; padding: 30px; border-radius: 5px; margin-top: 20px; }
        .button { 
            display: inline-block; 
            padding: 12px 30px; 
            background-color: #4CAF50; 
            color: white; 
            text-decoration: none; 
            border-radius: 5px; 
            margin: 20px 0;
        }
        .info-box { 
            background-color: #e3f2fd; 
            padding: 15px; 
            border-radius: 5px; 
            margin: 20px 0;
            border-left: 4px solid #2196F3;
        }
        .warning-box { 
            background-color: #fff3cd; 
            padding: 15px; 
            border-radius: 5px; 
            margin: 20px 0;
            border-left: 4px solid #ffc107;
        }
        .footer { text-align: center; margin-top: 30px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to %s!</h1>
        </div>
        <div class="content">
            <h2>You've Been Invited to Join WorkZen HRMS</h2>
            <p>Hello %s,</p>
            <p>Your account has been created on <strong>%s</strong>'s WorkZen HRMS system. To get started, you need to verify your email address.</p>
            
            <div class="info-box">
                <strong>📋 Your Login Credentials:</strong><br><br>
                <strong>Username:</strong> %s<br>
                <strong>Temporary Password:</strong> <code style="background: #fff; padding: 5px 10px; border-radius: 3px; font-size: 16px;">%s</code>
            </div>

            <div class="warning-box">
                <strong>⚠️ Important:</strong><br>
                Please verify your email address before attempting to login. After verification, you can login and change your password.
            </div>

            <p><strong>Step 1:</strong> Click the button below to verify your email address:</p>
            <div style="text-align: center;">
                <a href="%s" class="button">Verify Email Address</a>
            </div>
            
            <p>Or copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #666; font-size: 12px;">%s</p>

            <p><strong>Step 2:</strong> After verification, login with your credentials and change your password.</p>

            <p><strong>Note:</strong> This verification link will expire in 24 hours.</p>

            <h3>What You Can Do:</h3>
            <ul>
                <li>Track your attendance</li>
                <li>Apply for leaves</li>
                <li>View your payroll information</li>
                <li>Access company documents</li>
                <li>Update your profile</li>
            </ul>
        </div>
        <div class="footer">
            <p>© 2025 WorkZen HRMS. All rights reserved.</p>
            <p>If you didn't expect this email or have questions, please contact your HR department.</p>
        </div>
    </div>
</body>
</html>
`, companyName, firstName, companyName, username, tempPassword, verificationLink, verificationLink)

	return SendEmail(to, subject, body)
}
