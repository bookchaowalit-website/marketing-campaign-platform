package service

import (
	"fmt"
	"os"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// NewEmailService creates the appropriate email service based on configuration
func NewEmailService(frontendURL string) (domain.EmailService, error) {
	provider := os.Getenv("EMAIL_PROVIDER")
	if provider == "" {
		provider = "stub" // Default to stub for development
	}

	switch provider {
	case "sendgrid":
		apiKey := os.Getenv("SENDGRID_API_KEY")
		fromEmail := os.Getenv("EMAIL_FROM_ADDRESS")
		fromName := os.Getenv("EMAIL_FROM_NAME")
		if fromName == "" {
			fromName = "Marketing SaaS"
		}

		if apiKey == "" || fromEmail == "" {
			return nil, fmt.Errorf("SENDGRID_API_KEY and EMAIL_FROM_ADDRESS are required for SendGrid")
		}

		return NewSendGridEmailService(apiKey, fromEmail, fromName, frontendURL)

	case "ses":
		region := os.Getenv("AWS_REGION")
		if region == "" {
			region = "us-east-1"
		}
		fromEmail := os.Getenv("EMAIL_FROM_ADDRESS")
		fromName := os.Getenv("EMAIL_FROM_NAME")
		if fromName == "" {
			fromName = "Marketing SaaS"
		}

		if fromEmail == "" {
			return nil, fmt.Errorf("EMAIL_FROM_ADDRESS is required for SES")
		}

		return NewSESEmailService(region, fromEmail, fromName, frontendURL)

	case "stub":
		fallthrough
	default:
		return NewStubEmailService(frontendURL), nil
	}
}
