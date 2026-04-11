package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/bookmarketing/marketing-saas/api-go/internal/domain"
	"github.com/google/uuid"
)

// ===========================
// SMS LIST REPOSITORY
// ===========================

type smsListRepository struct {
	db *sql.DB
}

func NewSMSListRepository(db *sql.DB) domain.SMSListRepository {
	return &smsListRepository{db: db}
}

func (r *smsListRepository) Create(l *domain.SMSList) error {
	query := `
		INSERT INTO sms_lists (id, org_id, name, description, subscriber_count, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	now := time.Now()
	l.ID = uuid.New().String()
	l.CreatedAt = now
	l.UpdatedAt = now

	_, err := r.db.Exec(query, l.ID, l.OrganizationID, l.Name, l.Description, l.SubscriberCount, l.Status)
	return err
}

func (r *smsListRepository) GetByID(id string) (*domain.SMSList, error) {
	query := `
		SELECT id, org_id, name, description, subscriber_count, status, created_at, updated_at
		FROM sms_lists WHERE id = $1
	`
	l := &domain.SMSList{}
	err := r.db.QueryRow(query, id).Scan(
		&l.ID, &l.OrganizationID, &l.Name, &l.Description, &l.SubscriberCount,
		&l.Status, &l.CreatedAt, &l.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (r *smsListRepository) ListByOrg(orgID string) ([]*domain.SMSList, error) {
	query := `
		SELECT id, org_id, name, description, subscriber_count, status, created_at, updated_at
		FROM sms_lists WHERE org_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []*domain.SMSList
	for rows.Next() {
		l := &domain.SMSList{}
		err := rows.Scan(
			&l.ID, &l.OrganizationID, &l.Name, &l.Description, &l.SubscriberCount,
			&l.Status, &l.CreatedAt, &l.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		lists = append(lists, l)
	}
	return lists, nil
}

func (r *smsListRepository) Update(l *domain.SMSList) error {
	query := `
		UPDATE sms_lists SET name=$1, description=$2, subscriber_count=$3, status=$4, updated_at=$5 WHERE id=$6
	`
	l.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, l.Name, l.Description, l.SubscriberCount, l.Status, l.UpdatedAt, l.ID)
	return err
}

func (r *smsListRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM sms_lists WHERE id = $1`, id)
	return err
}

func (r *smsListRepository) UpdateSubscriberCount(id string, count int) error {
	_, err := r.db.Exec(`UPDATE sms_lists SET subscriber_count = $1, updated_at = NOW() WHERE id = $2`, count, id)
	return err
}

// ===========================
// SMS SUBSCRIBER REPOSITORY
// ===========================

type smsSubscriberRepository struct {
	db *sql.DB
}

func NewSMSSubscriberRepository(db *sql.DB) domain.SMSSubscriberRepository {
	return &smsSubscriberRepository{db: db}
}

func (r *smsSubscriberRepository) Create(s *domain.SMSSubscriber) error {
	query := `
		INSERT INTO sms_subscribers (id, org_id, list_id, phone, first_name, last_name, custom_fields, status, subscribed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	now := time.Now()
	s.ID = uuid.New().String()
	s.SubscribedAt = now
	s.CreatedAt = now
	s.UpdatedAt = now

	var customFieldsJSON []byte
	if s.CustomFields != nil {
		customFieldsJSON, _ = json.Marshal(s.CustomFields)
	}

	_, err := r.db.Exec(query, s.ID, s.OrganizationID, s.ListID, s.Phone, s.FirstName,
		s.LastName, customFieldsJSON, s.Status, s.SubscribedAt)
	return err
}

func (r *smsSubscriberRepository) GetByID(id string) (*domain.SMSSubscriber, error) {
	query := `
		SELECT id, org_id, list_id, phone, first_name, last_name, custom_fields, status,
			subscribed_at, opted_out_at, opt_out_reason, created_at, updated_at
		FROM sms_subscribers WHERE id = $1
	`
	s := &domain.SMSSubscriber{}
	var customFieldsRaw []byte
	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.OrganizationID, &s.ListID, &s.Phone, &s.FirstName, &s.LastName,
		&customFieldsRaw, &s.Status, &s.SubscribedAt, &s.OptedOutAt, &s.OptOutReason,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(customFieldsRaw) > 0 {
		json.Unmarshal(customFieldsRaw, &s.CustomFields)
	}
	return s, nil
}

func (r *smsSubscriberRepository) GetByPhone(orgID, phone string) (*domain.SMSSubscriber, error) {
	query := `
		SELECT id, org_id, list_id, phone, first_name, last_name, custom_fields, status,
			subscribed_at, opted_out_at, opt_out_reason, created_at, updated_at
		FROM sms_subscribers WHERE org_id = $1 AND phone = $2
	`
	s := &domain.SMSSubscriber{}
	var customFieldsRaw []byte
	err := r.db.QueryRow(query, orgID, phone).Scan(
		&s.ID, &s.OrganizationID, &s.ListID, &s.Phone, &s.FirstName, &s.LastName,
		&customFieldsRaw, &s.Status, &s.SubscribedAt, &s.OptedOutAt, &s.OptOutReason,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(customFieldsRaw) > 0 {
		json.Unmarshal(customFieldsRaw, &s.CustomFields)
	}
	return s, nil
}

func (r *smsSubscriberRepository) ListByOrg(orgID, status string, limit, offset int) ([]*domain.SMSSubscriber, int, error) {
	countQuery := `SELECT COUNT(*) FROM sms_subscribers WHERE org_id = $1`
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
		SELECT id, org_id, list_id, phone, first_name, last_name, status, subscribed_at, created_at
		FROM sms_subscribers WHERE org_id = $1`
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

	var subscribers []*domain.SMSSubscriber
	for rows.Next() {
		s := &domain.SMSSubscriber{}
		err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.ListID, &s.Phone, &s.FirstName, &s.LastName,
			&s.Status, &s.SubscribedAt, &s.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		subscribers = append(subscribers, s)
	}
	return subscribers, total, nil
}

func (r *smsSubscriberRepository) ListByList(listID string, status string, limit, offset int) ([]*domain.SMSSubscriber, int, error) {
	countQuery := `SELECT COUNT(*) FROM sms_subscribers WHERE list_id = $1`
	args := []interface{}{listID}
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
		SELECT id, org_id, list_id, phone, first_name, last_name, status, subscribed_at, created_at
		FROM sms_subscribers WHERE list_id = $1`
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

	var subscribers []*domain.SMSSubscriber
	for rows.Next() {
		s := &domain.SMSSubscriber{}
		err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.ListID, &s.Phone, &s.FirstName, &s.LastName,
			&s.Status, &s.SubscribedAt, &s.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		subscribers = append(subscribers, s)
	}
	return subscribers, total, nil
}

func (r *smsSubscriberRepository) Update(s *domain.SMSSubscriber) error {
	query := `
		UPDATE sms_subscribers SET phone=$1, first_name=$2, last_name=$3, status=$4, updated_at=$5 WHERE id=$6
	`
	s.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, s.Phone, s.FirstName, s.LastName, s.Status, s.UpdatedAt, s.ID)
	return err
}

func (r *smsSubscriberRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM sms_subscribers WHERE id = $1`, id)
	return err
}

func (r *smsSubscriberRepository) OptOut(id string, reason string) error {
	now := time.Now()
	_, err := r.db.Exec(`
		UPDATE sms_subscribers SET status='opted_out', opted_out_at=$1, opt_out_reason=$2, updated_at=$3 WHERE id=$4
	`, now, reason, now, id)
	return err
}

// ===========================
// SMS CAMPAIGN REPOSITORY
// ===========================

type smsCampaignRepository struct {
	db *sql.DB
}

func NewSMSCampaignRepository(db *sql.DB) domain.SMSCampaignRepository {
	return &smsCampaignRepository{db: db}
}

func (r *smsCampaignRepository) Create(c *domain.SMSCampaign) error {
	query := `
		INSERT INTO sms_campaigns (id, org_id, list_id, name, message_template, status, scheduled_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now()
	c.ID = uuid.New().String()
	c.CreatedAt = now
	c.UpdatedAt = now

	_, err := r.db.Exec(query, c.ID, c.OrganizationID, c.ListID, c.Name,
		c.MessageTemplate, c.Status, c.ScheduledAt)
	return err
}

func (r *smsCampaignRepository) GetByID(id string) (*domain.SMSCampaign, error) {
	query := `
		SELECT id, org_id, list_id, name, message_template, status, scheduled_at, sent_at,
			sent_count, delivered_count, failed_count, opt_out_count, cost, created_at, updated_at
		FROM sms_campaigns WHERE id = $1
	`
	c := &domain.SMSCampaign{}
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.OrganizationID, &c.ListID, &c.Name, &c.MessageTemplate,
		&c.Status, &c.ScheduledAt, &c.SentAt, &c.SentCount, &c.DeliveredCount,
		&c.FailedCount, &c.OptOutCount, &c.Cost, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *smsCampaignRepository) ListByOrg(orgID, status string, limit, offset int) ([]*domain.SMSCampaign, int, error) {
	countQuery := `SELECT COUNT(*) FROM sms_campaigns WHERE org_id = $1`
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
		SELECT id, org_id, list_id, name, message_template, status, scheduled_at, sent_at,
			sent_count, delivered_count, failed_count, cost, created_at
		FROM sms_campaigns WHERE org_id = $1`
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

	var campaigns []*domain.SMSCampaign
	for rows.Next() {
		c := &domain.SMSCampaign{}
		err := rows.Scan(
			&c.ID, &c.OrganizationID, &c.ListID, &c.Name, &c.MessageTemplate,
			&c.Status, &c.ScheduledAt, &c.SentAt, &c.SentCount, &c.DeliveredCount,
			&c.FailedCount, &c.Cost, &c.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		campaigns = append(campaigns, c)
	}
	return campaigns, total, nil
}

func (r *smsCampaignRepository) Update(c *domain.SMSCampaign) error {
	query := `
		UPDATE sms_campaigns SET name=$1, message_template=$2, status=$3, scheduled_at=$4,
			sent_at=$5, sent_count=$6, delivered_count=$7, failed_count=$8,
			opt_out_count=$9, cost=$10, updated_at=$11 WHERE id=$12
	`
	c.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, c.Name, c.MessageTemplate, c.Status, c.ScheduledAt,
		c.SentAt, c.SentCount, c.DeliveredCount, c.FailedCount, c.OptOutCount,
		c.Cost, c.UpdatedAt, c.ID)
	return err
}

func (r *smsCampaignRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM sms_campaigns WHERE id = $1`, id)
	return err
}

func (r *smsCampaignRepository) UpdateStatus(id, status string) error {
	now := time.Now()
	query := `UPDATE sms_campaigns SET status=$1, updated_at=$2 WHERE id=$3`
	if status == "sending" {
		query = `UPDATE sms_campaigns SET status=$1, sent_at=$2, updated_at=$3 WHERE id=$4`
		_, err := r.db.Exec(query, status, now, now, id)
		return err
	}
	_, err := r.db.Exec(query, status, now, id)
	return err
}

func (r *smsCampaignRepository) UpdateSendStats(id string, sentCount, deliveredCount, failedCount int, cost float64) error {
	now := time.Now()
	status := "sent"
	if failedCount > 0 && sentCount == failedCount {
		status = "failed"
	} else if failedCount > 0 {
		status = "sent"
	}
	_, err := r.db.Exec(`
		UPDATE sms_campaigns SET status=$1, sent_count=$2, delivered_count=$3,
			failed_count=$4, cost=$5, updated_at=$6 WHERE id=$7
	`, status, sentCount, deliveredCount, failedCount, cost, now, id)
	return err
}

// ===========================
// SMS LOG REPOSITORY
// ===========================

type smsLogRepository struct {
	db *sql.DB
}

func NewSMSLogRepository(db *sql.DB) domain.SMSLogRepository {
	return &smsLogRepository{db: db}
}

func (r *smsLogRepository) Create(l *domain.SMSLog) error {
	query := `
		INSERT INTO sms_logs (id, campaign_id, subscriber_id, phone, status)
		VALUES ($1, $2, $3, $4, $5)
	`
	l.ID = uuid.New().String()
	l.CreatedAt = time.Now()
	_, err := r.db.Exec(query, l.ID, l.CampaignID, l.SubscriberID, l.Phone, l.Status)
	return err
}

func (r *smsLogRepository) GetByID(id string) (*domain.SMSLog, error) {
	query := `
		SELECT id, campaign_id, subscriber_id, phone, status, message_sid, sent_at,
			delivered_at, error_message, cost, created_at
		FROM sms_logs WHERE id = $1
	`
	l := &domain.SMSLog{}
	err := r.db.QueryRow(query, id).Scan(
		&l.ID, &l.CampaignID, &l.SubscriberID, &l.Phone, &l.Status,
		&l.MessageSID, &l.SentAt, &l.DeliveredAt, &l.ErrorMessage, &l.Cost, &l.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (r *smsLogRepository) ListByCampaign(campaignID string) ([]*domain.SMSLog, error) {
	query := `
		SELECT id, campaign_id, subscriber_id, phone, status, message_sid, sent_at,
			delivered_at, error_message, cost, created_at
		FROM sms_logs WHERE campaign_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.SMSLog
	for rows.Next() {
		l := &domain.SMSLog{}
		err := rows.Scan(
			&l.ID, &l.CampaignID, &l.SubscriberID, &l.Phone, &l.Status,
			&l.MessageSID, &l.SentAt, &l.DeliveredAt, &l.ErrorMessage, &l.Cost, &l.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (r *smsLogRepository) ListBySubscriber(subscriberID string, limit, offset int) ([]*domain.SMSLog, int, error) {
	countQuery := `SELECT COUNT(*) FROM sms_logs WHERE subscriber_id = $1`
	var total int
	if err := r.db.QueryRow(countQuery, subscriberID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, campaign_id, subscriber_id, phone, status, message_sid, sent_at,
			delivered_at, error_message, cost, created_at
		FROM sms_logs WHERE subscriber_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, subscriberID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*domain.SMSLog
	for rows.Next() {
		l := &domain.SMSLog{}
		err := rows.Scan(
			&l.ID, &l.CampaignID, &l.SubscriberID, &l.Phone, &l.Status,
			&l.MessageSID, &l.SentAt, &l.DeliveredAt, &l.ErrorMessage, &l.Cost, &l.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}
	return logs, total, nil
}

func (r *smsLogRepository) UpdateStatus(id, status string, messageSID string, errMsg string) error {
	now := time.Now()
	query := `
		UPDATE sms_logs SET status=$1, message_sid=$2, error_message=$3, updated_at=$4
	`
	args := []interface{}{status, messageSID, errMsg, now}
	argIdx := 5

	if status == "sent" || status == "delivered" {
		query += `, sent_at=$` + itoa(argIdx)
		args = append(args, now)
		argIdx++
	}
	if status == "delivered" {
		query += `, delivered_at=$` + itoa(argIdx)
		args = append(args, now)
		argIdx++
	}

	query += ` WHERE id=$` + itoa(argIdx)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}
