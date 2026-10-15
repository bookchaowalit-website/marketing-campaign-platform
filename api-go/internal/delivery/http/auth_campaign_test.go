package http

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const fixturePassword = "offline-password-123"
const registerJSON = `{"email":"reader@example.test","password":"offline-password-123","full_name":"Test Reader"}`
const campaignJSON = `{"name":"Autumn launch","type":"Email","budget":1250.5,"start_date":"2026-11-01"}`

// Embedding the interfaces makes unused operations fail loudly instead of
// inventing behavior. These fakes implement only the boundaries this flow uses.
type memoryUsers struct {
	domain.UserRepository
	mu        sync.Mutex
	users     map[string]domain.User
	createErr error
}

func (r *memoryUsers) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.createErr != nil {
		return r.createErr
	}
	r.users[user.ID] = *user
	return nil
}
func (r *memoryUsers) GetByID(id string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &user, nil
}
func (r *memoryUsers) GetByEmail(email string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, user := range r.users {
		if user.Email == email {
			copy := user
			return &copy, nil
		}
	}
	return nil, errors.New("not found")
}
func (r *memoryUsers) snapshot() []domain.User {
	r.mu.Lock()
	defer r.mu.Unlock()
	users := make([]domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users
}
func (r *memoryUsers) deactivate(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user := r.users[id]
	user.IsActive = false
	r.users[id] = user
}

type memoryRefresh struct {
	domain.RefreshTokenRepository
	tokens []domain.RefreshToken
}

func (r *memoryRefresh) Create(token *domain.RefreshToken) error {
	r.tokens = append(r.tokens, *token)
	return nil
}
func (r *memoryRefresh) IsBlacklisted(string) (bool, error) { return false, nil }

type memoryVerification struct {
	domain.EmailVerificationTokenRepository
}

func (*memoryVerification) DeleteByUserID(string) error                 { return nil }
func (*memoryVerification) Create(*domain.EmailVerificationToken) error { return nil }

type sentEmail struct{ to, url string }
type fixtureEmail struct {
	domain.EmailService
	sent chan sentEmail
}

func (e *fixtureEmail) SendVerificationEmail(to, url string) error {
	e.sent <- sentEmail{to, url}
	return nil
}

type memoryCampaigns struct {
	domain.CampaignRepository
	campaigns map[string]domain.Campaign
	createErr error
}

func (r *memoryCampaigns) Create(_ context.Context, campaign *domain.Campaign) error {
	if r.createErr != nil {
		return r.createErr
	}
	r.campaigns[campaign.ID] = *campaign
	return nil
}

type testAPI struct {
	router    *gin.Engine
	auth      *service.AuthService
	users     *memoryUsers
	refresh   *memoryRefresh
	campaigns *memoryCampaigns
	email     *fixtureEmail
}

func newTestAPI() *testAPI {
	gin.SetMode(gin.TestMode)
	users := &memoryUsers{users: make(map[string]domain.User)}
	refresh := &memoryRefresh{}
	email := &fixtureEmail{sent: make(chan sentEmail, 1)}
	auth := service.NewAuthService(users, refresh, nil, &memoryVerification{}, email,
		"offline-fixture-signing-key", "https://frontend.example.test")
	campaigns := &memoryCampaigns{campaigns: make(map[string]domain.Campaign)}
	authHandler := NewAuthHandler(auth)
	campaignHandler := NewCampaignHandler(service.NewCampaignService(campaigns))
	router := gin.New()
	router.POST("/api/v1/auth/register", authHandler.Register)
	router.POST("/api/v1/auth/login", authHandler.Login)
	router.POST("/api/v1/campaigns", AuthMiddleware(auth), campaignHandler.CreateCampaign)
	return &testAPI{router, auth, users, refresh, campaigns, email}
}
func (a *testAPI) post(path, body, token string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	result := httptest.NewRecorder()
	a.router.ServeHTTP(result, request)
	return result
}
func expectStatus(t *testing.T, result *httptest.ResponseRecorder, status int) {
	t.Helper()
	// Do not print response bodies: successful responses contain fixture tokens.
	if result.Code != status {
		t.Fatalf("HTTP status = %d, want %d", result.Code, status)
	}
}
func decodeResponse[T any](t *testing.T, result *httptest.ResponseRecorder) T {
	t.Helper()
	var response T
	if err := json.Unmarshal(result.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return response
}

func assertPrivateUserFields(t *testing.T, result *httptest.ResponseRecorder) {
	t.Helper()
	var envelope struct {
		User map[string]json.RawMessage `json:"user"`
	}
	if err := json.Unmarshal(result.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode user fields: %v", err)
	}
	for _, field := range []string{"password_hash", "PasswordHash", "password"} {
		if _, exists := envelope.User[field]; exists {
			t.Fatal("HTTP user object exposes a private password field")
		}
	}
}
func (a *testAPI) register(t *testing.T) service.AuthResponse {
	t.Helper()
	result := a.post("/api/v1/auth/register", registerJSON, "")
	expectStatus(t, result, http.StatusCreated)
	assertPrivateUserFields(t, result)
	response := decodeResponse[service.AuthResponse](t, result)
	// Finish the handler's background email before issuing another request,
	// so it never outlives this fixture or races reuse of Gin's pooled context.
	select {
	case email := <-a.email.sent:
		if email.to != "reader@example.test" || !strings.HasPrefix(email.url, "https://frontend.example.test/verify-email?token=") {
			t.Fatal("verification email has the wrong recipient or frontend path")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("verification email was not handed to the local fake")
	}
	return response
}

func TestRegisterPersistsHashedUserAndValidTokens(t *testing.T) {
	api := newTestAPI()
	response := api.register(t)
	if response.User == nil {
		t.Fatal("response has no user")
	}
	users := api.users.snapshot()
	if len(users) != 1 {
		t.Fatalf("persisted users = %d, want 1", len(users))
	}
	user := users[0]
	if user.Email != "reader@example.test" || user.FullName != "Test Reader" || user.Role != "user" || !user.IsActive {
		t.Fatal("persisted registration fields differ from the request/defaults")
	}
	if _, err := uuid.Parse(user.ID); err != nil {
		t.Fatal("user ID is not a UUID")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(fixturePassword)); err != nil {
		t.Fatal("stored password does not authenticate the registration password")
	}
	if user.PasswordHash == fixturePassword {
		t.Fatal("password was stored in plain text")
	}
	claims, err := api.auth.ValidateToken(response.Token)
	if err != nil || claims.UserID != user.ID || claims.Email != user.Email {
		t.Fatal("access token does not identify the registered user")
	}
	if response.ExpiresIn <= 0 || response.RefreshToken == "" {
		t.Fatal("missing token lifetime or refresh token")
	}
	hash := sha256.Sum256([]byte(response.RefreshToken))
	if len(api.refresh.tokens) != 1 || api.refresh.tokens[0].UserID != user.ID || api.refresh.tokens[0].TokenHash != hex.EncodeToString(hash[:]) {
		t.Fatal("refresh token was not persisted as the matching user/hash")
	}
}

func TestRegisterRejectsInvalidInputWithoutPersistence(t *testing.T) {
	for name, body := range map[string]string{
		"malformed": `{`, "missing email": `{"password":"offline-password-123"}`,
		"bad email":      `{"email":"invalid","password":"offline-password-123"}`,
		"short password": `{"email":"reader@example.test","password":"short"}`,
	} {
		t.Run(name, func(t *testing.T) {
			api := newTestAPI()
			expectStatus(t, api.post("/api/v1/auth/register", body, ""), http.StatusBadRequest)
			if len(api.users.snapshot()) != 0 || len(api.refresh.tokens) != 0 {
				t.Fatal("invalid registration persisted state")
			}
		})
	}
}
func TestRegisterDuplicateAndStorageError(t *testing.T) {
	t.Run("duplicate", func(t *testing.T) {
		api := newTestAPI()
		api.register(t)
		expectStatus(t, api.post("/api/v1/auth/register", registerJSON, ""), http.StatusConflict)
		if len(api.users.snapshot()) != 1 || len(api.refresh.tokens) != 1 {
			t.Fatal("duplicate registration changed stored state")
		}
	})
	t.Run("repository failure", func(t *testing.T) {
		api := newTestAPI()
		api.users.createErr = errors.New("fixture storage unavailable")
		expectStatus(t, api.post("/api/v1/auth/register", registerJSON, ""), http.StatusInternalServerError)
		if len(api.users.snapshot()) != 0 || len(api.refresh.tokens) != 0 {
			t.Fatal("failed registration issued tokens or persisted a user")
		}
	})
}

func TestLoginAndAuthenticatedCampaignCreate(t *testing.T) {
	api := newTestAPI()
	registered := api.register(t)
	login := api.post("/api/v1/auth/login", `{"email":"reader@example.test","password":"offline-password-123"}`, "")
	expectStatus(t, login, http.StatusOK)
	assertPrivateUserFields(t, login)
	response := decodeResponse[service.AuthResponse](t, login)
	if response.User == nil || response.User.ID != registered.User.ID || response.RefreshToken == "" {
		t.Fatal("login response does not identify the registered user")
	}
	if len(api.users.snapshot()) != 1 || len(api.refresh.tokens) != 2 {
		t.Fatal("login changed users or failed to persist a new refresh token")
	}
	result := api.post("/api/v1/campaigns", campaignJSON, response.Token)
	expectStatus(t, result, http.StatusCreated)
	campaign := decodeResponse[domain.Campaign](t, result)
	if _, err := uuid.Parse(campaign.ID); err != nil {
		t.Fatal("campaign ID is not a UUID")
	}
	if len(api.campaigns.campaigns) != 1 {
		t.Fatal("created campaign was not persisted")
	}
	saved := api.campaigns.campaigns[campaign.ID]
	if saved.Name != "Autumn launch" || saved.Type != "Email" || saved.Budget != 1250.5 || saved.Status != "draft" || saved.Spent != 0 {
		t.Fatal("persisted campaign does not match input/default state")
	}
	if saved.StartDate.Format("2006-01-02") != "2026-11-01" || saved.CreatedAt.IsZero() || saved.UpdatedAt.Before(saved.CreatedAt) {
		t.Fatal("campaign dates were not initialized")
	}
	if campaign.Name != saved.Name || campaign.Budget != saved.Budget || campaign.Status != saved.Status {
		t.Fatal("HTTP response differs from persisted campaign")
	}
}

func TestLoginRejectsInvalidInputAndCredentials(t *testing.T) {
	for name, body := range map[string]string{
		"malformed": `{`, "missing password": `{"email":"reader@example.test"}`,
		"bad email": `{"email":"invalid","password":"anything"}`,
	} {
		t.Run(name, func(t *testing.T) {
			api := newTestAPI()
			expectStatus(t, api.post("/api/v1/auth/login", body, ""), http.StatusBadRequest)
			if len(api.refresh.tokens) != 0 {
				t.Fatal("invalid login issued a refresh token")
			}
		})
	}
	api := newTestAPI()
	registered := api.register(t)
	for name, body := range map[string]string{
		"wrong password": `{"email":"reader@example.test","password":"wrong-password"}`,
		"unknown email":  `{"email":"unknown@example.test","password":"offline-password-123"}`,
	} {
		t.Run(name, func(t *testing.T) { expectStatus(t, api.post("/api/v1/auth/login", body, ""), http.StatusUnauthorized) })
	}
	api.users.deactivate(registered.User.ID)
	expectStatus(t, api.post("/api/v1/auth/login", `{"email":"reader@example.test","password":"offline-password-123"}`, ""), http.StatusUnauthorized)
	if len(api.refresh.tokens) != 1 {
		t.Fatal("rejected login issued refresh tokens")
	}
}

func TestCampaignCreateRejectsInvalidInputAndAuthentication(t *testing.T) {
	api := newTestAPI()
	registered := api.register(t)
	for name, body := range map[string]string{
		"malformed": `{`, "missing name": `{"type":"Email","budget":10,"start_date":"2026-11-01"}`,
		"negative budget": `{"name":"Launch","type":"Email","budget":-1,"start_date":"2026-11-01"}`,
		"bad date":        `{"name":"Launch","type":"Email","budget":10,"start_date":"tomorrow"}`,
	} {
		t.Run(name, func(t *testing.T) {
			expectStatus(t, api.post("/api/v1/campaigns", body, registered.Token), http.StatusBadRequest)
		})
	}
	expectStatus(t, api.post("/api/v1/campaigns", campaignJSON, ""), http.StatusUnauthorized)
	expectStatus(t, api.post("/api/v1/campaigns", campaignJSON, registered.Token+"tampered"), http.StatusUnauthorized)
	if len(api.campaigns.campaigns) != 0 {
		t.Fatal("rejected campaign create persisted data")
	}
}
func TestCampaignCreateStorageError(t *testing.T) {
	api := newTestAPI()
	registered := api.register(t)
	api.campaigns.createErr = errors.New("fixture storage unavailable")
	expectStatus(t, api.post("/api/v1/campaigns", campaignJSON, registered.Token), http.StatusInternalServerError)
	if len(api.campaigns.campaigns) != 0 {
		t.Fatal("failed create persisted a campaign")
	}
}
