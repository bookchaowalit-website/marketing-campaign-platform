package domain

import "time"

// RefreshToken represents a refresh token entity
type RefreshToken struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	TokenHash   string    `json:"-"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	Revoked     bool      `json:"revoked"`
	RevokedAt   time.Time `json:"-"`
}

// TokenBlacklist represents a blacklisted token
type TokenBlacklist struct {
	ID         string    `json:"id"`
	TokenHash  string    `json:"-"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// RefreshTokenRepository defines the interface for refresh token operations
type RefreshTokenRepository interface {
	Create(token *RefreshToken) error
	GetByTokenHash(tokenHash string) (*RefreshToken, error)
	Revoke(tokenHash string) error
	IsBlacklisted(tokenHash string) (bool, error)
	Blacklist(token *TokenBlacklist) error
	DeleteExpired() error
}
