-- External Affiliate Integrations Migration - DOWN
-- Removes all external affiliate integration tables

DROP TRIGGER IF EXISTS update_ext_aff_accounts_updated_at ON external_affiliate_accounts;
DROP TRIGGER IF EXISTS update_ext_aff_products_updated_at ON external_affiliate_products;
DROP TRIGGER IF EXISTS update_ext_aff_orders_updated_at ON external_affiliate_orders;
DROP TRIGGER IF EXISTS update_ext_aff_earnings_updated_at ON external_affiliate_earnings;
DROP TRIGGER IF EXISTS update_ext_aff_links_updated_at ON external_affiliate_links;

DROP INDEX IF EXISTS idx_ext_aff_accounts_user;
DROP INDEX IF EXISTS idx_ext_aff_accounts_org;
DROP INDEX IF EXISTS idx_ext_aff_accounts_platform;
DROP INDEX IF EXISTS idx_ext_aff_accounts_status;
DROP INDEX IF EXISTS idx_ext_aff_products_account;
DROP INDEX IF EXISTS idx_ext_aff_products_platform;
DROP INDEX IF EXISTS idx_ext_aff_products_status;
DROP INDEX IF EXISTS idx_ext_aff_orders_account;
DROP INDEX IF EXISTS idx_ext_aff_orders_platform;
DROP INDEX IF EXISTS idx_ext_aff_orders_status;
DROP INDEX IF EXISTS idx_ext_aff_orders_date;
DROP INDEX IF EXISTS idx_ext_aff_orders_commission_status;
DROP INDEX IF EXISTS idx_ext_aff_earnings_account;
DROP INDEX IF EXISTS idx_ext_aff_earnings_date;
DROP INDEX IF EXISTS idx_ext_aff_earnings_platform;
DROP INDEX IF EXISTS idx_ext_aff_sync_logs_account;
DROP INDEX IF EXISTS idx_ext_aff_sync_logs_status;
DROP INDEX IF EXISTS idx_ext_aff_sync_logs_started;
DROP INDEX IF EXISTS idx_ext_aff_links_code;
DROP INDEX IF EXISTS idx_ext_aff_links_account;

DROP TABLE IF EXISTS external_affiliate_links;
DROP TABLE IF EXISTS external_affiliate_sync_logs;
DROP TABLE IF EXISTS external_affiliate_earnings;
DROP TABLE IF EXISTS external_affiliate_orders;
DROP TABLE IF EXISTS external_affiliate_products;
DROP TABLE IF EXISTS external_affiliate_accounts;