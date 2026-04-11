package repository

import (
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

type refreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *sql.DB) domain.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at, revoked)
		VALUES ($1, $2, $3, $4, $5, false)
	`
	now := time.Now()
	token.CreatedAt = now

	_, err := r.db.Exec(query,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.CreatedAt,
	)
	return err
}

func (r *refreshTokenRepository) GetByTokenHash(tokenHash string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked
		FROM refresh_tokens
		WHERE token_hash = $1 AND revoked = false AND expires_at > NOW()
	`
	token := &domain.RefreshToken{}
	err := r.db.QueryRow(query, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.Revoked,
	)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *refreshTokenRepository) Revoke(tokenHash string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true, revoked_at = NOW()
		WHERE token_hash = $1
	`
	_, err := r.db.Exec(query, tokenHash)
	return err
}

func (r *refreshTokenRepository) IsBlacklisted(tokenHash string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM token_blacklist WHERE token_hash = $1 AND expires_at > NOW()`
	err := r.db.QueryRow(query, tokenHash).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *refreshTokenRepository) Blacklist(token *domain.TokenBlacklist) error {
	query := `
		INSERT INTO token_blacklist (id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (token_hash) DO NOTHING
	`
	now := time.Now()
	token.CreatedAt = now

	_, err := r.db.Exec(query,
		token.ID,
		token.TokenHash,
		token.ExpiresAt,
		token.CreatedAt,
	)
	return err
}

func (r *refreshTokenRepository) DeleteExpired() error {
	query := `DELETE FROM refresh_tokens WHERE revoked = true OR expires_at < NOW()`
	_, err := r.db.Exec(query)
	return err
}
