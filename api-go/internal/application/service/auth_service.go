package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrTokenRevoked       = errors.New("token revoked")
	ErrEmailAlreadyVerified = errors.New("email already verified")
	ErrInvalidResetToken  = errors.New("invalid reset token")
)

type AuthService struct {
	userRepo                domain.UserRepository
	refreshRepo             domain.RefreshTokenRepository
	passwordResetRepo       domain.PasswordResetTokenRepository
	emailVerificationRepo   domain.EmailVerificationTokenRepository
	emailService            domain.EmailService
	jwtSecret               string
	frontendURL             string
}

func NewAuthService(
	userRepo domain.UserRepository,
	refreshRepo domain.RefreshTokenRepository,
	passwordResetRepo domain.PasswordResetTokenRepository,
	emailVerificationRepo domain.EmailVerificationTokenRepository,
	emailService domain.EmailService,
	jwtSecret string,
	frontendURL string,
) *AuthService {
	return &AuthService{
		userRepo:              userRepo,
		refreshRepo:           refreshRepo,
		passwordResetRepo:     passwordResetRepo,
		emailVerificationRepo: emailVerificationRepo,
		emailService:          emailService,
		jwtSecret:             jwtSecret,
		frontendURL:           frontendURL,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type AuthResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	User         *domain.User `json:"user"`
}

// RefreshTokenResponse contains new access token and refresh token
type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Register creates a new user
func (s *AuthService) Register(ctx context.Context, email, password, fullName string) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user with proper UUID
	user := &domain.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		Role:         "user",
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens
	token, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    86400, // 24 hours
		User:         user,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	token, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    86400, // 24 hours
		User:         user,
	}, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *AuthService) generateToken(user *domain.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// generateTokens creates both access and refresh tokens
func (s *AuthService) generateTokens(user *domain.User) (string, string, error) {
	// Generate access token
	accessToken, err := s.generateToken(user)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken := uuid.New().String()
	tokenHash := sha256.Sum256([]byte(refreshToken))
	tokenHashStr := hex.EncodeToString(tokenHash[:])

	refresh := &domain.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: tokenHashStr,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.refreshRepo.Create(refresh); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken generates a new access token using a valid refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	// Hash the refresh token
	tokenHash := sha256.Sum256([]byte(refreshToken))
	tokenHashStr := hex.EncodeToString(tokenHash[:])

	// Get refresh token from database
	storedToken, err := s.refreshRepo.GetByTokenHash(tokenHashStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Check if token is revoked
	if storedToken.Revoked {
		return nil, ErrTokenRevoked
	}

	// Get user
	user, err := s.userRepo.GetByID(storedToken.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Revoke old refresh token
	if err := s.refreshRepo.Revoke(tokenHashStr); err != nil {
		return nil, err
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    86400, // 24 hours
	}, nil
}

// Logout blacklists the access token and revokes the refresh token
func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	// Blacklist access token
	if accessToken != "" {
		accessTokenHash := sha256.Sum256([]byte(accessToken))
		accessTokenHashStr := hex.EncodeToString(accessTokenHash[:])
		
		// Get token expiration from claims
		claims, err := s.ValidateToken(accessToken)
		if err == nil {
			blacklist := &domain.TokenBlacklist{
				ID:        uuid.New().String(),
				TokenHash: accessTokenHashStr,
				ExpiresAt: time.Unix(claims.ExpiresAt.Unix(), 0),
			}
			s.refreshRepo.Blacklist(blacklist)
		}
	}

	// Revoke refresh token
	if refreshToken != "" {
		refreshTokenHash := sha256.Sum256([]byte(refreshToken))
		refreshTokenHashStr := hex.EncodeToString(refreshTokenHash[:])
		s.refreshRepo.Revoke(refreshTokenHashStr)
	}

	return nil
}

// IsTokenBlacklisted checks if a token is in the blacklist
func (s *AuthService) IsTokenBlacklisted(tokenString string) bool {
	tokenHash := sha256.Sum256([]byte(tokenString))
	tokenHashStr := hex.EncodeToString(tokenHash[:])
	
	blacklisted, _ := s.refreshRepo.IsBlacklisted(tokenHashStr)
	return blacklisted
}

// RequestPasswordReset generates a password reset token and sends email
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		// Don't reveal if email exists or not
		return nil
	}

	// Generate reset token
	token := s.generateSecureToken()
	resetToken := &domain.PasswordResetToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
	}

	if err := s.passwordResetRepo.Create(resetToken); err != nil {
		return err
	}

	// Send password reset email
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, token)
	return s.emailService.SendPasswordResetEmail(user.Email, resetURL)
}

// ResetPassword resets user password using a valid token
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Get reset token
	resetToken, err := s.passwordResetRepo.GetByToken(token)
	if err != nil {
		return ErrInvalidResetToken
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update user password
	if err := s.userRepo.ResetPassword(resetToken.UserID, string(hashedPassword)); err != nil {
		return err
	}

	// Mark reset token as used
	return s.passwordResetRepo.MarkAsUsed(token)
}

// SendVerificationEmail generates and sends email verification token
func (s *AuthService) SendVerificationEmail(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if user.EmailVerified {
		return ErrEmailAlreadyVerified
	}

	// Delete any existing verification tokens
	s.emailVerificationRepo.DeleteByUserID(userID)

	// Generate verification token
	token := s.generateSecureToken()
	verificationToken := &domain.EmailVerificationToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour expiry
	}

	if err := s.emailVerificationRepo.Create(verificationToken); err != nil {
		return err
	}

	// Send verification email
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", s.frontendURL, token)
	return s.emailService.SendVerificationEmail(user.Email, verificationURL)
}

// VerifyEmail verifies user email using a valid token
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	// Get verification token
	verificationToken, err := s.emailVerificationRepo.GetByToken(token)
	if err != nil {
		return ErrInvalidToken
	}

	// Update user email_verified status
	if err := s.userRepo.MarkEmailVerified(verificationToken.UserID); err != nil {
		return err
	}

	// Mark verification token as used
	return s.emailVerificationRepo.MarkAsUsed(token)
}

// generateSecureToken generates a cryptographically secure random token
func (s *AuthService) generateSecureToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, 48)
	for i := range token {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		token[i] = charset[num.Int64()]
	}
	return string(token)
}
