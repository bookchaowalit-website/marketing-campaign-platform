-- Rollback initial migration

DROP TRIGGER IF EXISTS update_keyword_rankings_updated_at ON keyword_rankings;
DROP TRIGGER IF EXISTS update_seo_metrics_updated_at ON seo_metrics;
DROP TRIGGER IF EXISTS update_competitors_updated_at ON competitors;
DROP TRIGGER IF EXISTS update_ab_tests_updated_at ON ab_tests;
DROP TRIGGER IF EXISTS update_campaigns_updated_at ON campaigns;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_analytics_events_created_at;
DROP INDEX IF EXISTS idx_analytics_events_user_id;
DROP INDEX IF EXISTS idx_seo_metrics_user_id;
DROP INDEX IF EXISTS idx_competitors_user_id;
DROP INDEX IF EXISTS idx_ab_tests_status;
DROP INDEX IF EXISTS idx_ab_tests_user_id;
DROP INDEX IF EXISTS idx_campaigns_type;
DROP INDEX IF EXISTS idx_campaigns_status;
DROP INDEX IF EXISTS idx_campaigns_user_id;

DROP TABLE IF EXISTS analytics_events;
DROP TABLE IF EXISTS keyword_rankings;
DROP TABLE IF EXISTS seo_metrics;
DROP TABLE IF EXISTS competitors;
DROP TABLE IF EXISTS ab_tests;
DROP TABLE IF EXISTS campaigns;
DROP TABLE IF EXISTS users;