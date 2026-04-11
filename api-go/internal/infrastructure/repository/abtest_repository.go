package repository

import (
	"database/sql"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
)

type abTestRepository struct {
	db *sql.DB
}

// NewABTestRepository creates a new A/B test repository
func NewABTestRepository(db *sql.DB) domain.ABTestRepository {
	return &abTestRepository{db: db}
}

// Create inserts a new A/B test into the database
func (r *abTestRepository) Create(test *domain.ABTest) error {
	query := `
		INSERT INTO ab_tests (id, user_id, campaign_id, name, variant, status, visitors, conversions, conversion_rate, confidence, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	now := time.Now()
	test.CreatedAt = now
	test.UpdatedAt = now

	var endDate interface{}
	if test.EndDate != nil {
		endDate = *test.EndDate
	}

	_, err := r.db.Exec(query,
		test.ID,
		test.UserID,
		test.CampaignID,
		test.Name,
		test.Variant,
		test.Status,
		test.Visitors,
		test.Conversions,
		test.ConversionRate,
		test.Confidence,
		test.StartDate,
		endDate,
		test.CreatedAt,
		test.UpdatedAt,
	)
	return err
}

// GetByID retrieves an A/B test by ID
func (r *abTestRepository) GetByID(id string) (*domain.ABTest, error) {
	query := `
		SELECT id, user_id, campaign_id, name, variant, status, visitors, conversions, conversion_rate, confidence, start_date, end_date, created_at, updated_at
		FROM ab_tests WHERE id = $1
	`
	test := &domain.ABTest{}
	var endDate sql.NullTime
	var campaignID sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&test.ID,
		&test.UserID,
		&campaignID,
		&test.Name,
		&test.Variant,
		&test.Status,
		&test.Visitors,
		&test.Conversions,
		&test.ConversionRate,
		&test.Confidence,
		&test.StartDate,
		&endDate,
		&test.CreatedAt,
		&test.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if endDate.Valid {
		test.EndDate = &endDate.Time
	}
	if campaignID.Valid {
		test.CampaignID = campaignID.String
	}

	return test, nil
}

// GetAll retrieves all A/B tests with pagination
func (r *abTestRepository) GetAll(userID string, limit, offset int) ([]*domain.ABTest, error) {
	query := `
		SELECT id, user_id, campaign_id, name, variant, status, visitors, conversions, conversion_rate, confidence, start_date, end_date, created_at, updated_at
		FROM ab_tests WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []*domain.ABTest
	for rows.Next() {
		test := &domain.ABTest{}
		var endDate sql.NullTime
		var campaignID sql.NullString

		err := rows.Scan(
			&test.ID,
			&test.UserID,
			&campaignID,
			&test.Name,
			&test.Variant,
			&test.Status,
			&test.Visitors,
			&test.Conversions,
			&test.ConversionRate,
			&test.Confidence,
			&test.StartDate,
			&endDate,
			&test.CreatedAt,
			&test.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if endDate.Valid {
			test.EndDate = &endDate.Time
		}
		if campaignID.Valid {
			test.CampaignID = campaignID.String
		}

		tests = append(tests, test)
	}

	return tests, nil
}

// Update updates an A/B test
func (r *abTestRepository) Update(test *domain.ABTest) error {
	query := `
		UPDATE ab_tests
		SET name = $1, status = $2, visitors = $3, conversions = $4, conversion_rate = $5, confidence = $6, end_date = $7, updated_at = $8
		WHERE id = $9
	`
	test.UpdatedAt = time.Now()

	var endDate interface{}
	if test.EndDate != nil {
		endDate = *test.EndDate
	}

	_, err := r.db.Exec(query,
		test.Name,
		test.Status,
		test.Visitors,
		test.Conversions,
		test.ConversionRate,
		test.Confidence,
		endDate,
		test.UpdatedAt,
		test.ID,
	)
	return err
}

// Delete deletes an A/B test
func (r *abTestRepository) Delete(id string) error {
	query := `DELETE FROM ab_tests WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetByStatus retrieves A/B tests by status
func (r *abTestRepository) GetByStatus(userID, status string) ([]*domain.ABTest, error) {
	query := `
		SELECT id, user_id, campaign_id, name, variant, status, visitors, conversions, conversion_rate, confidence, start_date, end_date, created_at, updated_at
		FROM ab_tests WHERE user_id = $1 AND status = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []*domain.ABTest
	for rows.Next() {
		test := &domain.ABTest{}
		var endDate sql.NullTime
		var campaignID sql.NullString

		err := rows.Scan(
			&test.ID,
			&test.UserID,
			&campaignID,
			&test.Name,
			&test.Variant,
			&test.Status,
			&test.Visitors,
			&test.Conversions,
			&test.ConversionRate,
			&test.Confidence,
			&test.StartDate,
			&endDate,
			&test.CreatedAt,
			&test.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if endDate.Valid {
			test.EndDate = &endDate.Time
		}
		if campaignID.Valid {
			test.CampaignID = campaignID.String
		}

		tests = append(tests, test)
	}

	return tests, nil
}