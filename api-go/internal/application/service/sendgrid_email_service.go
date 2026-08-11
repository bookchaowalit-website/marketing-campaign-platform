package service

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridEmailService implements EmailService using SendGrid
type SendGridEmailService struct {
	apiKey      string
	fromEmail   string
	fromName    string
	frontendURL string
	templates   *template.Template
}

// NewSendGridEmailService creates a new SendGrid email service
func NewSendGridEmailService(apiKey, fromEmail, fromName, frontendURL string) (*SendGridEmailService, error) {
	// Parse email templates
	tmpl, err := template.New("email").Parse("")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &SendGridEmailService{
		apiKey:      apiKey,
		fromEmail:   fromEmail,
		fromName:    fromName,
		frontendURL: frontendURL,
		templates:   tmpl,
	}, nil
}

func (s *SendGridEmailService) sendEmail(to, subject, htmlContent string) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	toEmail := mail.NewEmail("", to)

	message := mail.NewSingleEmail(from, subject, toEmail, "", htmlContent)

	client := sendgrid.NewSendClient(s.apiKey)
	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		log.Printf("SendGrid error: Status=%d, Body=%s", response.StatusCode, response.Body)
		return fmt.Errorf("sendgrid returned status %d: %s", response.StatusCode, response.Body)
	}

	log.Printf("Email sent to %s: %s", to, subject)
	return nil
}

func (s *SendGridEmailService) SendPasswordResetEmail(to, resetURL string) error {
	subject := "Reset Your Password"
	htmlContent := s.passwordResetTemplate(resetURL)
	return s.sendEmail(to, subject, htmlContent)
}

func (s *SendGridEmailService) SendVerificationEmail(to, verificationURL string) error {
	subject := "Verify Your Email Address"
	htmlContent := s.verificationTemplate(verificationURL)
	return s.sendEmail(to, subject, htmlContent)
}

func (s *SendGridEmailService) SendWelcomeEmail(to, name string) error {
	subject := "Welcome to Marketing SaaS!"
	htmlContent := s.welcomeTemplate(name)
	return s.sendEmail(to, subject, htmlContent)
}

// SendEmail sends a marketing email with custom subject and content
func (s *SendGridEmailService) SendEmail(to, name, subject, htmlContent, textContent string) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	toEmail := mail.NewEmail(name, to)

	message := mail.NewSingleEmail(from, subject, toEmail, textContent, htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)
	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send marketing email: %w", err)
	}

	if response.StatusCode >= 400 {
		log.Printf("SendGrid marketing email error: Status=%d, Body=%s", response.StatusCode, response.Body)
		return fmt.Errorf("sendgrid returned status %d: %s", response.StatusCode, response.Body)
	}

	log.Printf("Marketing email sent to %s: %s", to, subject)
	return nil
}

// HTML Templates
func (s *SendGridEmailService) passwordResetTemplate(resetURL string) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; margin: 0; padding: 20px; }
		.container { max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
		.header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 40px 20px; text-align: center; }
		.header h1 { margin: 0; font-size: 24px; }
		.content { padding: 40px 30px; }
		.content p { color: #333; line-height: 1.6; }
		.button { display: inline-block; background: #667eea; color: white; text-decoration: none; padding: 14px 32px; border-radius: 6px; font-weight: 600; margin: 20px 0; }
		.button:hover { background: #5568d3; }
		.url-fallback { background: #f8f9fa; border: 1px solid #dee2e6; padding: 10px; border-radius: 4px; word-break: break-all; font-size: 12px; color: #666; margin-top: 15px; }
		.footer { background: #f8f9fa; padding: 20px; text-align: center; color: #666; font-size: 12px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🔑 Reset Your Password</h1>
		</div>
		<div class="content">
			<p>Hi there,</p>
			<p>We received a request to reset your password. Click the button below to create a new password:</p>
			<p style="text-align: center;">
				<a href="%s" class="button">Reset Password</a>
			</p>
			<p>This link will expire in <strong>1 hour</strong>.</p>
			<p>If you didn't request this, you can safely ignore this email.</p>
			<div class="url-fallback">
				<strong>Button not working?</strong> Copy and paste this URL into your browser:<br>
				%s
			</div>
		</div>
		<div class="footer">
			<p>This is an automated message from Marketing SaaS Platform.</p>
			<p>&copy; 2024 Marketing SaaS. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
`
	return fmt.Sprintf(tmpl, resetURL, resetURL)
}

func (s *SendGridEmailService) verificationTemplate(verificationURL string) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; margin: 0; padding: 20px; }
		.container { max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
		.header { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); color: white; padding: 40px 20px; text-align: center; }
		.header h1 { margin: 0; font-size: 24px; }
		.content { padding: 40px 30px; }
		.content p { color: #333; line-height: 1.6; }
		.button { display: inline-block; background: #11998e; color: white; text-decoration: none; padding: 14px 32px; border-radius: 6px; font-weight: 600; margin: 20px 0; }
		.button:hover { background: #0d7a6e; }
		.url-fallback { background: #f8f9fa; border: 1px solid #dee2e6; padding: 10px; border-radius: 4px; word-break: break-all; font-size: 12px; color: #666; margin-top: 15px; }
		.footer { background: #f8f9fa; padding: 20px; text-align: center; color: #666; font-size: 12px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>✅ Verify Your Email</h1>
		</div>
		<div class="content">
			<p>Welcome to Marketing SaaS!</p>
			<p>Please verify your email address by clicking the button below:</p>
			<p style="text-align: center;">
				<a href="%s" class="button">Verify Email Address</a>
			</p>
			<p>This link will expire in <strong>24 hours</strong>.</p>
			<div class="url-fallback">
				<strong>Button not working?</strong> Copy and paste this URL into your browser:<br>
				%s
			</div>
		</div>
		<div class="footer">
			<p>This is an automated message from Marketing SaaS Platform.</p>
			<p>&copy; 2024 Marketing SaaS. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
`
	return fmt.Sprintf(tmpl, verificationURL, verificationURL)
}

func (s *SendGridEmailService) welcomeTemplate(name string) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; margin: 0; padding: 20px; }
		.container { max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
		.header { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); color: white; padding: 40px 20px; text-align: center; }
		.header h1 { margin: 0; font-size: 24px; }
		.content { padding: 40px 30px; }
		.content p { color: #333; line-height: 1.6; }
		.content ul { color: #333; line-height: 1.8; }
		.button { display: inline-block; background: #f5576c; color: white; text-decoration: none; padding: 14px 32px; border-radius: 6px; font-weight: 600; margin: 20px 0; }
		.footer { background: #f8f9fa; padding: 20px; text-align: center; color: #666; font-size: 12px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🎉 Welcome, %s!</h1>
		</div>
		<div class="content">
			<p>Thank you for joining Marketing SaaS Platform!</p>
			<p>We're excited to help you grow your marketing efforts. Here's what you can do:</p>
			<ul>
				<li>🚀 Create and manage marketing campaigns</li>
				<li>📊 Track performance with real-time analytics</li>
				<li>🧪 Run A/B tests to optimize results</li>
				<li>🔍 Analyze competitors and SEO metrics</li>
				<li>📧 Manage email marketing campaigns</li>
			</ul>
			<p style="text-align: center;">
				<a href="%s/dashboard" class="button">Go to Dashboard</a>
			</p>
			<p>If you have any questions, feel free to reach out to our support team.</p>
		</div>
		<div class="footer">
			<p>This is an automated message from Marketing SaaS Platform.</p>
			<p>&copy; 2024 Marketing SaaS. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
`
	return fmt.Sprintf(tmpl, name, s.frontendURL)
}

// HTTPClient interface for mocking
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
