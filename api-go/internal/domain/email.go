package domain

import "time"

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// EmailVerificationToken represents an email verification token
type EmailVerificationToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// PasswordResetTokenRepository defines the interface for password reset operations
type PasswordResetTokenRepository interface {
	Create(token *PasswordResetToken) error
	GetByToken(token string) (*PasswordResetToken, error)
	MarkAsUsed(token string) error
	DeleteExpired() error
}

// EmailVerificationTokenRepository defines the interface for email verification
type EmailVerificationTokenRepository interface {
	Create(token *EmailVerificationToken) error
	GetByToken(token string) (*EmailVerificationToken, error)
	MarkAsUsed(token string) error
	DeleteByUserID(userID string) error
	DeleteExpired() error
}

// EmailService defines the interface for sending emails
type EmailService interface {
	SendPasswordResetEmail(to, resetURL string) error
	SendVerificationEmail(to, verificationURL string) error
	SendWelcomeEmail(to, name string) error
	// SendEmail sends a marketing email (added for email campaigns)
	SendEmail(to, name, subject, htmlContent, textContent string) error
}
