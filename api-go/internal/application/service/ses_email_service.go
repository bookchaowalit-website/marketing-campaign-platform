package service

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// SESEmailService implements EmailService using AWS SES
type SESEmailService struct {
	client      *ses.SES
	fromEmail   string
	fromName    string
	frontendURL string
}

// NewSESEmailService creates a new AWS SES email service
func NewSESEmailService(region, fromEmail, fromName, frontendURL string) (*SESEmailService, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &SESEmailService{
		client:      ses.New(sess),
		fromEmail:   fromEmail,
		fromName:    fromName,
		frontendURL: frontendURL,
	}, nil
}

func (s *SESEmailService) sendEmail(to, subject, htmlContent string) error {
	from := fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)

	input := &ses.SendEmailInput{
		Source: aws.String(from),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(to)},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Data:    aws.String(subject),
				Charset: aws.String("UTF-8"),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Data:    aws.String(htmlContent),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	}

	_, err := s.client.SendEmail(input)
	if err != nil {
		return fmt.Errorf("failed to send email via SES: %w", err)
	}

	log.Printf("Email sent to %s via SES: %s", to, subject)
	return nil
}

func (s *SESEmailService) SendPasswordResetEmail(to, resetURL string) error {
	subject := "Reset Your Password"
	htmlContent := s.passwordResetTemplate(resetURL)
	return s.sendEmail(to, subject, htmlContent)
}

func (s *SESEmailService) SendVerificationEmail(to, verificationURL string) error {
	subject := "Verify Your Email Address"
	htmlContent := s.verificationTemplate(verificationURL)
	return s.sendEmail(to, subject, htmlContent)
}

func (s *SESEmailService) SendWelcomeEmail(to, name string) error {
	subject := "Welcome to Marketing SaaS!"
	htmlContent := s.welcomeTemplate(name)
	return s.sendEmail(to, subject, htmlContent)
}

// SendEmail sends a marketing email with custom subject and content
func (s *SESEmailService) SendEmail(to, name, subject, htmlContent, textContent string) error {
	from := fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)

	input := &ses.SendEmailInput{
		Source: aws.String(from),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(to)},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Data:    aws.String(subject),
				Charset: aws.String("UTF-8"),
			},
			Body: &ses.Body{
				Html: &ses.Content{
					Data:    aws.String(htmlContent),
					Charset: aws.String("UTF-8"),
				},
				Text: &ses.Content{
					Data:    aws.String(textContent),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	}

	_, err := s.client.SendEmail(input)
	if err != nil {
		return fmt.Errorf("failed to send marketing email via SES: %w", err)
	}

	log.Printf("Marketing email sent to %s via SES: %s", to, subject)
	return nil
}

// HTML Templates (same as SendGrid)
func (s *SESEmailService) passwordResetTemplate(resetURL string) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; margin: 0; padding: 20px; }
		.container { max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
		.header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 40px 20px; text-align: center; }
		.header h1 { margin: 0; font-size: 24px; }
		.content { padding: 40px 30px; }
		.content p { color: #333; line-height: 1.6; }
		.button { display: inline-block; background: #667eea; color: white; text-decoration: none; padding: 14px 32px; border-radius: 6px; font-weight: 600; margin: 20px 0; }
		.url-fallback { background: #f8f9fa; border: 1px solid #dee2e6; padding: 10px; border-radius: 4px; word-break: break-all; font-size: 12px; color: #666; margin-top: 15px; }
		.footer { background: #f8f9fa; padding: 20px; text-align: center; color: #666; font-size: 12px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header"><h1>🔑 Reset Your Password</h1></div>
		<div class="content">
			<p>Hi there,</p>
			<p>We received a request to reset your password. Click the button below to create a new password:</p>
			<p style="text-align: center;"><a href="%s" class="button">Reset Password</a></p>
			<p>This link will expire in <strong>1 hour</strong>.</p>
			<p>If you didn't request this, you can safely ignore this email.</p>
			<div class="url-fallback"><strong>Button not working?</strong> Copy and paste this URL:<br>%s</div>
		</div>
		<div class="footer"><p>&copy; 2024 Marketing SaaS. All rights reserved.</p></div>
	</div>
</body>
</html>
`
	return fmt.Sprintf(tmpl, resetURL, resetURL)
}

func (s *SESEmailService) verificationTemplate(verificationURL string) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; margin: 0; padding: 20px; }
		.container { max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
		.header { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); color: white; padding: 40px 20px; text-align: center; }
		.header h1 { margin: 0; font-size: 24px; }
		.content { padding: 40px 30px; }
		.content p { color: #333; line-height: 1.6; }
		.button { display: inline-block; background: #11998e; color: white; text-decoration: none; padding: 14px 32px; border-radius: 6px; font-weight: 600; margin: 20px 0; }
		.url-fallback { background: #f8f9fa; border: 1px solid #dee2e6; padding: 10px; border-radius: 4px; word-break: break-all; font-size: 12px; color: #666; margin-top: 15px; }
		.footer { background: #f8f9fa; padding: 20px; text-align: center; color: #666; font-size: 12px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header"><h1>✅ Verify Your Email</h1></div>
		<div class="content">
			<p>Welcome to Marketing SaaS!</p>
			<p>Please verify your email address by clicking the button below:</p>
			<p style="text-align: center;"><a href="%s" class="button">Verify Email Address</a></p>
			<p>This link will expire in <strong>24 hours</strong>.</p>
			<div class="url-fallback"><strong>Button not working?</strong> Copy and paste this URL:<br>%s</div>
		</div>
		<div class="footer"><p>&copy; 2024 Marketing SaaS. All rights reserved.</p></div>
	</div>
</body>
</html>
`
	return fmt.Sprintf(tmpl, verificationURL, verificationURL)
}

func (s *SESEmailService) welcomeTemplate(name string) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; margin: 0; padding: 20px; }
		.container { max-width: 600px; margin: 0 auto; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
		.header { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); color: white; padding: 40px 20px; text-align: center; }
		.header h1 { margin: 0; font-size: 24px; }
		.content { padding: 40px 30px; }
		.content p, .content ul { color: #333; line-height: 1.8; }
		.button { display: inline-block; background: #f5576c; color: white; text-decoration: none; padding: 14px 32px; border-radius: 6px; font-weight: 600; margin: 20px 0; }
		.footer { background: #f8f9fa; padding: 20px; text-align: center; color: #666; font-size: 12px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header"><h1>🎉 Welcome, %s!</h1></div>
		<div class="content">
			<p>Thank you for joining Marketing SaaS Platform!</p>
			<p>Here's what you can do:</p>
			<ul>
				<li>🚀 Create and manage marketing campaigns</li>
				<li>📊 Track performance with real-time analytics</li>
				<li>🧪 Run A/B tests to optimize results</li>
				<li>🔍 Analyze competitors and SEO metrics</li>
				<li>📧 Manage email marketing campaigns</li>
			</ul>
			<p style="text-align: center;"><a href="%s/dashboard" class="button">Go to Dashboard</a></p>
		</div>
		<div class="footer"><p>&copy; 2024 Marketing SaaS. All rights reserved.</p></div>
	</div>
</body>
</html>
`
	return fmt.Sprintf(tmpl, name, s.frontendURL)
}
