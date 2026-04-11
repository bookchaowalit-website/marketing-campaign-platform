package service

import (
	"context"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

// ABTestService handles business logic for A/B tests
type ABTestService struct {
	repo domain.ABTestRepository
}

// NewABTestService creates a new A/B test service
func NewABTestService(repo domain.ABTestRepository) *ABTestService {
	return &ABTestService{repo: repo}
}

// CreateABTest creates a new A/B test
func (s *ABTestService) CreateABTest(ctx context.Context, test *domain.ABTest) error {
	test.CreatedAt = time.Now()
	test.UpdatedAt = time.Now()
	return s.repo.Create(test)
}

// GetABTest retrieves an A/B test by ID
func (s *ABTestService) GetABTest(ctx context.Context, id string) (*domain.ABTest, error) {
	return s.repo.GetByID(id)
}

// GetAllABTests retrieves all A/B tests with pagination for a user
func (s *ABTestService) GetAllABTests(ctx context.Context, userID string, limit, offset int) ([]*domain.ABTest, error) {
	return s.repo.GetAll(userID, limit, offset)
}

// UpdateABTest updates an existing A/B test
func (s *ABTestService) UpdateABTest(ctx context.Context, test *domain.ABTest) error {
	test.UpdatedAt = time.Now()
	return s.repo.Update(test)
}

// DeleteABTest deletes an A/B test
func (s *ABTestService) DeleteABTest(ctx context.Context, id string) error {
	return s.repo.Delete(id)
}

// GetByStatus retrieves A/B tests by status for a user
func (s *ABTestService) GetByStatus(ctx context.Context, userID, status string) ([]*domain.ABTest, error) {
	return s.repo.GetByStatus(userID, status)
}

// CompleteABTest completes an A/B test
func (s *ABTestService) CompleteABTest(ctx context.Context, id string) error {
	test, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	now := time.Now()
	test.Status = "completed"
	test.EndDate = &now
	test.UpdatedAt = now
	return s.repo.Update(test)
}

// StartABTest starts a draft A/B test
func (s *ABTestService) StartABTest(ctx context.Context, id string) error {
	test, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	test.Status = "running"
	test.StartDate = time.Now()
	test.UpdatedAt = time.Now()
	return s.repo.Update(test)
}
