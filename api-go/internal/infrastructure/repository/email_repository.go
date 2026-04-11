package repository

import (
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

type passwordResetTokenRepository struct {
	db *sql.DB
}

// NewPasswordResetTokenRepository creates a new password reset token repository
func NewPasswordResetTokenRepository(db *sql.DB) domain.PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(token *domain.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (id, user_id, token, expires_at, created_at, used)
		VALUES ($1, $2, $3, $4, $5, false)
	`
	now := time.Now()
	token.CreatedAt = now

	_, err := r.db.Exec(query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	)
	return err
}

func (r *passwordResetTokenRepository) GetByToken(token string) (*domain.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at
		FROM password_reset_tokens
		WHERE token = $1 AND used = false AND expires_at > NOW()
	`
	t := &domain.PasswordResetToken{}
	err := r.db.QueryRow(query, token).Scan(
		&t.ID, &t.UserID, &t.Token, &t.ExpiresAt, &t.Used, &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *passwordResetTokenRepository) MarkAsUsed(token string) error {
	query := `UPDATE password_reset_tokens SET used = true WHERE token = $1`
	_, err := r.db.Exec(query, token)
	return err
}

func (r *passwordResetTokenRepository) DeleteExpired() error {
	query := `DELETE FROM password_reset_tokens WHERE expires_at < NOW() OR used = true`
	_, err := r.db.Exec(query)
	return err
}

// Email verification token repository
type emailVerificationTokenRepository struct {
	db *sql.DB
}

// NewEmailVerificationTokenRepository creates a new email verification token repository
func NewEmailVerificationTokenRepository(db *sql.DB) domain.EmailVerificationTokenRepository {
	return &emailVerificationTokenRepository{db: db}
}

func (r *emailVerificationTokenRepository) Create(token *domain.EmailVerificationToken) error {
	query := `
		INSERT INTO email_verification_tokens (id, user_id, token, expires_at, created_at, used)
		VALUES ($1, $2, $3, $4, $5, false)
	`
	now := time.Now()
	token.CreatedAt = now

	_, err := r.db.Exec(query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	)
	return err
}

func (r *emailVerificationTokenRepository) GetByToken(token string) (*domain.EmailVerificationToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at
		FROM email_verification_tokens
		WHERE token = $1 AND used = false AND expires_at > NOW()
	`
	t := &domain.EmailVerificationToken{}
	err := r.db.QueryRow(query, token).Scan(
		&t.ID, &t.UserID, &t.Token, &t.ExpiresAt, &t.Used, &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *emailVerificationTokenRepository) MarkAsUsed(token string) error {
	query := `UPDATE email_verification_tokens SET used = true WHERE token = $1`
	_, err := r.db.Exec(query, token)
	return err
}

func (r *emailVerificationTokenRepository) DeleteByUserID(userID string) error {
	query := `DELETE FROM email_verification_tokens WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *emailVerificationTokenRepository) DeleteExpired() error {
	query := `DELETE FROM email_verification_tokens WHERE expires_at < NOW() OR used = true`
	_, err := r.db.Exec(query)
	return err
}
