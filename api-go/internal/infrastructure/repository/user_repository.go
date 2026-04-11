package repository

import (
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into the database
func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.Role,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, is_active, email_verified, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, is_active, email_verified, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update updates a user
func (r *userRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, full_name = $2, role = $3, is_active = $4, updated_at = $5
		WHERE id = $6
	`
	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		user.Email,
		user.FullName,
		user.Role,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)
	return err
}

// Delete deletes a user
func (r *userRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// MarkEmailVerified marks a user's email as verified
func (r *userRepository) MarkEmailVerified(id string) error {
	query := `UPDATE users SET email_verified = true WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// ResetPassword updates a user's password
func (r *userRepository) ResetPassword(id, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`
	_, err := r.db.Exec(query, passwordHash, id)
	return err
}
