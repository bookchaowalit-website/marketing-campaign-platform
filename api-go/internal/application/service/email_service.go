package service

import (
	"fmt"
	"log"
)

// StubEmailService is a stub implementation for development
// Replace with real implementation (SendGrid, AWS SES, etc.) in production
type StubEmailService struct {
	frontendURL string
}

// NewStubEmailService creates a stub email service
func NewStubEmailService(frontendURL string) *StubEmailService {
	return &StubEmailService{frontendURL: frontendURL}
}

func (s *StubEmailService) SendPasswordResetEmail(to, resetURL string) error {
	// In production, send actual email via SendGrid/AWS SES
	log.Printf("[STUB] Password reset email to %s: %s", to, resetURL)
	fmt.Printf("\n🔑 [DEV MODE] Password Reset URL: %s\n\n", resetURL)
	return nil
}

func (s *StubEmailService) SendVerificationEmail(to, verificationURL string) error {
	// In production, send actual email via SendGrid/AWS SES
	log.Printf("[STUB] Verification email to %s: %s", to, verificationURL)
	fmt.Printf("\n✅ [DEV MODE] Verification URL: %s\n\n", verificationURL)
	return nil
}

func (s *StubEmailService) SendWelcomeEmail(to, name string) error {
	log.Printf("[STUB] Welcome email to %s: Welcome %s!", to, name)
	return nil
}

func (s *StubEmailService) SendEmail(to, name, subject, htmlContent, textContent string) error {
	log.Printf("[STUB] Marketing email to %s (%s): %s", to, name, subject)
	return nil
}
