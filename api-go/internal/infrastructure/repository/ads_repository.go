package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/google/uuid"
)

// ===========================
// AD PLATFORM REPOSITORY
// ===========================

type adPlatformRepository struct {
	db *sql.DB
}

func NewAdPlatformRepository(db *sql.DB) domain.AdPlatformRepository {
	return &adPlatformRepository{db: db}
}

func (r *adPlatformRepository) Create(p *domain.AdPlatform) error {
	query := `
		INSERT INTO ad_platforms (id, org_id, platform, account_id, api_credentials, status, connected_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now()
	p.ID = uuid.New().String()
	p.ConnectedAt = now
	_, err := r.db.Exec(query, p.ID, p.OrganizationID, p.Platform, p.AccountID, "", p.Status, p.ConnectedAt)
	return err
}

func (r *adPlatformRepository) GetByID(id string) (*domain.AdPlatform, error) {
	query := `
		SELECT id, org_id, platform, account_id, status, connected_at, last_sync_at, error_message, created_at, updated_at
		FROM ad_platforms WHERE id = $1
	`
	p := &domain.AdPlatform{}
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.OrganizationID, &p.Platform, &p.AccountID, &p.Status,
		&p.ConnectedAt, &p.LastSyncAt, &p.ErrorMessage, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *adPlatformRepository) ListByOrg(orgID string) ([]*domain.AdPlatform, error) {
	query := `
		SELECT id, org_id, platform, account_id, status, connected_at, last_sync_at, error_message, created_at, updated_at
		FROM ad_platforms WHERE org_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var platforms []*domain.AdPlatform
	for rows.Next() {
		p := &domain.AdPlatform{}
		err := rows.Scan(
			&p.ID, &p.OrganizationID, &p.Platform, &p.AccountID, &p.Status,
			&p.ConnectedAt, &p.LastSyncAt, &p.ErrorMessage, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		platforms = append(platforms, p)
	}
	return platforms, nil
}

func (r *adPlatformRepository) Update(p *domain.AdPlatform) error {
	query := `
		UPDATE ad_platforms SET platform=$1, account_id=$2, status=$3, updated_at=$4 WHERE id=$5
	`
	p.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, p.Platform, p.AccountID, p.Status, p.UpdatedAt, p.ID)
	return err
}

func (r *adPlatformRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM ad_platforms WHERE id = $1`, id)
	return err
}

func (r *adPlatformRepository) UpdateSyncStatus(id string, status string, errMsg string) error {
	now := time.Now()
	query := `
		UPDATE ad_platforms SET status=$1, error_message=$2, last_sync_at=$3, updated_at=$4 WHERE id=$5
	`
	_, err := r.db.Exec(query, status, errMsg, now, now, id)
	return err
}

// ===========================
// AD CAMPAIGN REPOSITORY
// ===========================

type adCampaignRepository struct {
	db *sql.DB
}

func NewAdCampaignRepository(db *sql.DB) domain.AdCampaignRepository {
	return &adCampaignRepository{db: db}
}

func (r *adCampaignRepository) Create(c *domain.AdCampaign) error {
	query := `
		INSERT INTO ad_campaigns (id, org_id, unified_campaign_id, platform_id, platform_campaign_id, name, status,
			budget, spent, impressions, clicks, conversions, cpc, cpm, ctr, start_date, end_date, targeting)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
	`
	now := time.Now()
	c.ID = uuid.New().String()
	c.CreatedAt = now
	c.UpdatedAt = now

	var targetingJSON []byte
	if c.Targeting != nil {
		targetingJSON, _ = json.Marshal(c.Targeting)
	}

	_, err := r.db.Exec(query, c.ID, c.OrganizationID, c.UnifiedCampaignID, c.PlatformID,
		c.PlatformCampaignID, c.Name, c.Status, c.Budget, c.Spent, c.Impressions, c.Clicks,
		c.Conversions, c.CPC, c.CPM, c.CTR, c.StartDate, c.EndDate, targetingJSON)
	return err
}

func (r *adCampaignRepository) GetByID(id string) (*domain.AdCampaign, error) {
	query := `
		SELECT id, org_id, unified_campaign_id, platform_id, platform_campaign_id, name, status,
			budget, spent, impressions, clicks, conversions, cpc, cpm, ctr, start_date, end_date,
			targeting, created_at, updated_at
		FROM ad_campaigns WHERE id = $1
	`
	c := &domain.AdCampaign{}
	var targetingRaw []byte
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.OrganizationID, &c.UnifiedCampaignID, &c.PlatformID, &c.PlatformCampaignID,
		&c.Name, &c.Status, &c.Budget, &c.Spent, &c.Impressions, &c.Clicks, &c.Conversions,
		&c.CPC, &c.CPM, &c.CTR, &c.StartDate, &c.EndDate, &targetingRaw, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(targetingRaw) > 0 {
		json.Unmarshal(targetingRaw, &c.Targeting)
	}
	return c, nil
}

func (r *adCampaignRepository) ListByOrg(orgID, status string, limit, offset int) ([]*domain.AdCampaign, int, error) {
	countQuery := `SELECT COUNT(*) FROM ad_campaigns WHERE org_id = $1`
	args := []interface{}{orgID}
	argIdx := 2

	if status != "" {
		countQuery += ` AND status = $` + itoa(argIdx)
		args = append(args, status)
		argIdx++
	}

	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, org_id, unified_campaign_id, platform_id, platform_campaign_id, name, status,
			budget, spent, impressions, clicks, conversions, cpc, cpm, ctr, start_date, end_date,
			created_at, updated_at
		FROM ad_campaigns WHERE org_id = $1`
	if status != "" {
		query += ` AND status = $` + itoa(argIdx)
		args = append(args, status)
		argIdx++
	}
	query += ` ORDER BY created_at DESC LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var campaigns []*domain.AdCampaign
	for rows.Next() {
		c := &domain.AdCampaign{}
		err := rows.Scan(
			&c.ID, &c.OrganizationID, &c.UnifiedCampaignID, &c.PlatformID, &c.PlatformCampaignID,
			&c.Name, &c.Status, &c.Budget, &c.Spent, &c.Impressions, &c.Clicks, &c.Conversions,
			&c.CPC, &c.CPM, &c.CTR, &c.StartDate, &c.EndDate, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		campaigns = append(campaigns, c)
	}
	return campaigns, total, nil
}

func (r *adCampaignRepository) ListByPlatform(platformID string) ([]*domain.AdCampaign, error) {
	query := `
		SELECT id, org_id, unified_campaign_id, platform_id, platform_campaign_id, name, status,
			budget, spent, impressions, clicks, conversions, cpc, cpm, ctr, start_date, end_date,
			created_at, updated_at
		FROM ad_campaigns WHERE platform_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, platformID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []*domain.AdCampaign
	for rows.Next() {
		c := &domain.AdCampaign{}
		err := rows.Scan(
			&c.ID, &c.OrganizationID, &c.UnifiedCampaignID, &c.PlatformID, &c.PlatformCampaignID,
			&c.Name, &c.Status, &c.Budget, &c.Spent, &c.Impressions, &c.Clicks, &c.Conversions,
			&c.CPC, &c.CPM, &c.CTR, &c.StartDate, &c.EndDate, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, c)
	}
	return campaigns, nil
}

func (r *adCampaignRepository) Update(c *domain.AdCampaign) error {
	query := `
		UPDATE ad_campaigns SET name=$1, status=$2, budget=$3, spent=$4,
			impressions=$5, clicks=$6, conversions=$7, cpc=$8, cpm=$9, ctr=$10,
			start_date=$11, end_date=$12, updated_at=$13 WHERE id=$14
	`
	c.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, c.Name, c.Status, c.Budget, c.Spent,
		c.Impressions, c.Clicks, c.Conversions, c.CPC, c.CPM, c.CTR,
		c.StartDate, c.EndDate, c.UpdatedAt, c.ID)
	return err
}

func (r *adCampaignRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM ad_campaigns WHERE id = $1`, id)
	return err
}

func (r *adCampaignRepository) UpdateStatus(id, status string) error {
	_, err := r.db.Exec(`UPDATE ad_campaigns SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return err
}

// ===========================
// AD GROUP REPOSITORY
// ===========================

type adGroupRepository struct {
	db *sql.DB
}

func NewAdGroupRepository(db *sql.DB) domain.AdGroupRepository {
	return &adGroupRepository{db: db}
}

func (r *adGroupRepository) Create(g *domain.AdGroup) error {
	query := `
		INSERT INTO ad_groups (id, campaign_id, platform_ad_group_id, name, status, budget, spent, impressions, clicks, conversions)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`
	now := time.Now()
	g.ID = uuid.New().String()
	g.CreatedAt = now
	g.UpdatedAt = now

	_, err := r.db.Exec(query, g.ID, g.CampaignID, g.PlatformAdGroupID, g.Name,
		g.Status, g.Budget, g.Spent, g.Impressions, g.Clicks, g.Conversions)
	return err
}

func (r *adGroupRepository) GetByID(id string) (*domain.AdGroup, error) {
	query := `
		SELECT id, campaign_id, platform_ad_group_id, name, status, budget, spent,
			impressions, clicks, conversions, created_at, updated_at
		FROM ad_groups WHERE id = $1
	`
	g := &domain.AdGroup{}
	err := r.db.QueryRow(query, id).Scan(
		&g.ID, &g.CampaignID, &g.PlatformAdGroupID, &g.Name, &g.Status,
		&g.Budget, &g.Spent, &g.Impressions, &g.Clicks, &g.Conversions,
		&g.CreatedAt, &g.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (r *adGroupRepository) ListByCampaign(campaignID string) ([]*domain.AdGroup, error) {
	query := `
		SELECT id, campaign_id, platform_ad_group_id, name, status, budget, spent,
			impressions, clicks, conversions, created_at, updated_at
		FROM ad_groups WHERE campaign_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*domain.AdGroup
	for rows.Next() {
		g := &domain.AdGroup{}
		err := rows.Scan(
			&g.ID, &g.CampaignID, &g.PlatformAdGroupID, &g.Name, &g.Status,
			&g.Budget, &g.Spent, &g.Impressions, &g.Clicks, &g.Conversions,
			&g.CreatedAt, &g.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (r *adGroupRepository) Update(g *domain.AdGroup) error {
	query := `
		UPDATE ad_groups SET name=$1, status=$2, budget=$3, spent=$4,
			impressions=$5, clicks=$6, conversions=$7, updated_at=$8 WHERE id=$9
	`
	g.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, g.Name, g.Status, g.Budget, g.Spent,
		g.Impressions, g.Clicks, g.Conversions, g.UpdatedAt, g.ID)
	return err
}

func (r *adGroupRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM ad_groups WHERE id = $1`, id)
	return err
}

// ===========================
// AD CREATIVE REPOSITORY
// ===========================

type adCreativeRepository struct {
	db *sql.DB
}

func NewAdCreativeRepository(db *sql.DB) domain.AdCreativeRepository {
	return &adCreativeRepository{db: db}
}

func (r *adCreativeRepository) Create(c *domain.AdCreative) error {
	query := `
		INSERT INTO ad_creatives (id, ad_group_id, type, headline, description, image_url, video_url,
			landing_page_url, status, approval_status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`
	now := time.Now()
	c.ID = uuid.New().String()
	c.CreatedAt = now
	c.UpdatedAt = now

	_, err := r.db.Exec(query, c.ID, c.AdGroupID, c.Type, c.Headline, c.Description,
		c.ImageURL, c.VideoURL, c.LandingPageURL, c.Status, c.ApprovalStatus)
	return err
}

func (r *adCreativeRepository) GetByID(id string) (*domain.AdCreative, error) {
	query := `
		SELECT id, ad_group_id, type, headline, description, image_url, video_url,
			landing_page_url, status, approval_status, platform_creative_id, created_at, updated_at
		FROM ad_creatives WHERE id = $1
	`
	c := &domain.AdCreative{}
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.AdGroupID, &c.Type, &c.Headline, &c.Description,
		&c.ImageURL, &c.VideoURL, &c.LandingPageURL, &c.Status, &c.ApprovalStatus,
		&c.PlatformCreativeID, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *adCreativeRepository) ListByAdGroup(adGroupID string) ([]*domain.AdCreative, error) {
	query := `
		SELECT id, ad_group_id, type, headline, description, image_url, video_url,
			landing_page_url, status, approval_status, platform_creative_id, created_at, updated_at
		FROM ad_creatives WHERE ad_group_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, adGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creatives []*domain.AdCreative
	for rows.Next() {
		c := &domain.AdCreative{}
		err := rows.Scan(
			&c.ID, &c.AdGroupID, &c.Type, &c.Headline, &c.Description,
			&c.ImageURL, &c.VideoURL, &c.LandingPageURL, &c.Status, &c.ApprovalStatus,
			&c.PlatformCreativeID, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		creatives = append(creatives, c)
	}
	return creatives, nil
}

func (r *adCreativeRepository) Update(c *domain.AdCreative) error {
	query := `
		UPDATE ad_creatives SET type=$1, headline=$2, description=$3, image_url=$4,
			video_url=$5, landing_page_url=$6, status=$7, approval_status=$8, updated_at=$9 WHERE id=$10
	`
	c.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, c.Type, c.Headline, c.Description, c.ImageURL,
		c.VideoURL, c.LandingPageURL, c.Status, c.ApprovalStatus, c.UpdatedAt, c.ID)
	return err
}

func (r *adCreativeRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM ad_creatives WHERE id = $1`, id)
	return err
}

// ===========================
// AD KEYWORD REPOSITORY
// ===========================

type adKeywordRepository struct {
	db *sql.DB
}

func NewAdKeywordRepository(db *sql.DB) domain.AdKeywordRepository {
	return &adKeywordRepository{db: db}
}

func (r *adKeywordRepository) Create(k *domain.AdKeyword) error {
	query := `
		INSERT INTO ad_keywords (id, ad_group_id, keyword, match_type, bid_amount, max_cpc,
			impressions, clicks, conversions, cost, avg_position, quality_score, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`
	now := time.Now()
	k.ID = uuid.New().String()
	k.CreatedAt = now
	k.UpdatedAt = now

	_, err := r.db.Exec(query, k.ID, k.AdGroupID, k.Keyword, k.MatchType, k.BidAmount,
		k.MaxCPC, k.Impressions, k.Clicks, k.Conversions, k.Cost, k.AvgPosition,
		k.QualityScore, k.Status)
	return err
}

func (r *adKeywordRepository) GetByID(id string) (*domain.AdKeyword, error) {
	query := `
		SELECT id, ad_group_id, keyword, match_type, bid_amount, max_cpc, impressions, clicks,
			conversions, cost, avg_position, quality_score, status, created_at, updated_at
		FROM ad_keywords WHERE id = $1
	`
	k := &domain.AdKeyword{}
	err := r.db.QueryRow(query, id).Scan(
		&k.ID, &k.AdGroupID, &k.Keyword, &k.MatchType, &k.BidAmount, &k.MaxCPC,
		&k.Impressions, &k.Clicks, &k.Conversions, &k.Cost, &k.AvgPosition,
		&k.QualityScore, &k.Status, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return k, nil
}

func (r *adKeywordRepository) ListByAdGroup(adGroupID string) ([]*domain.AdKeyword, error) {
	query := `
		SELECT id, ad_group_id, keyword, match_type, bid_amount, max_cpc, impressions, clicks,
			conversions, cost, avg_position, quality_score, status, created_at, updated_at
		FROM ad_keywords WHERE ad_group_id = $1 ORDER BY keyword
	`
	rows, err := r.db.Query(query, adGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []*domain.AdKeyword
	for rows.Next() {
		k := &domain.AdKeyword{}
		err := rows.Scan(
			&k.ID, &k.AdGroupID, &k.Keyword, &k.MatchType, &k.BidAmount, &k.MaxCPC,
			&k.Impressions, &k.Clicks, &k.Conversions, &k.Cost, &k.AvgPosition,
			&k.QualityScore, &k.Status, &k.CreatedAt, &k.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		keywords = append(keywords, k)
	}
	return keywords, nil
}

func (r *adKeywordRepository) Update(k *domain.AdKeyword) error {
	query := `
		UPDATE ad_keywords SET keyword=$1, match_type=$2, bid_amount=$3, max_cpc=$4,
			impressions=$5, clicks=$6, conversions=$7, cost=$8, avg_position=$9,
			quality_score=$10, status=$11, updated_at=$12 WHERE id=$13
	`
	k.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, k.Keyword, k.MatchType, k.BidAmount, k.MaxCPC,
		k.Impressions, k.Clicks, k.Conversions, k.Cost, k.AvgPosition,
		k.QualityScore, k.Status, k.UpdatedAt, k.ID)
	return err
}

func (r *adKeywordRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM ad_keywords WHERE id = $1`, id)
	return err
}

// ===========================
// AD PERFORMANCE REPOSITORY
// ===========================

type adPerformanceRepository struct {
	db *sql.DB
}

func NewAdPerformanceRepository(db *sql.DB) domain.AdPerformanceRepository {
	return &adPerformanceRepository{db: db}
}

func (r *adPerformanceRepository) Create(p *domain.AdPerformanceDaily) error {
	query := `
		INSERT INTO ad_performance_daily (id, ad_id, date, impressions, clicks, conversions, cost, ctr, cpc, cpa, roas)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()
	_, err := r.db.Exec(query, p.ID, p.AdID, p.Date, p.Impressions, p.Clicks,
		p.Conversions, p.Cost, p.CTR, p.CPC, p.CPA, p.ROAS)
	return err
}

func (r *adPerformanceRepository) GetByAdID(adID string, startDate, endDate time.Time) ([]*domain.AdPerformanceDaily, error) {
	query := `
		SELECT id, ad_id, date, impressions, clicks, conversions, cost, ctr, cpc, cpa, roas, created_at
		FROM ad_performance_daily WHERE ad_id = $1 AND date >= $2 AND date <= $3 ORDER BY date
	`
	rows, err := r.db.Query(query, adID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []*domain.AdPerformanceDaily
	for rows.Next() {
		p := &domain.AdPerformanceDaily{}
		err := rows.Scan(
			&p.ID, &p.AdID, &p.Date, &p.Impressions, &p.Clicks, &p.Conversions,
			&p.Cost, &p.CTR, &p.CPC, &p.CPA, &p.ROAS, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		performances = append(performances, p)
	}
	return performances, nil
}

func (r *adPerformanceRepository) Upsert(p *domain.AdPerformanceDaily) error {
	query := `
		INSERT INTO ad_performance_daily (id, ad_id, date, impressions, clicks, conversions, cost, ctr, cpc, cpa, roas)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (ad_id, date) DO UPDATE SET
			impressions = $4, clicks = $5, conversions = $6, cost = $7,
			ctr = $8, cpc = $9, cpa = $10, roas = $11
	`
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()
	_, err := r.db.Exec(query, p.ID, p.AdID, p.Date, p.Impressions, p.Clicks,
		p.Conversions, p.Cost, p.CTR, p.CPC, p.CPA, p.ROAS)
	return err
}
