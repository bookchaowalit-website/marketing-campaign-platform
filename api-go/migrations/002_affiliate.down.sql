-- Affiliate Program Migration - DOWN
-- Removes all affiliate tables

DROP TRIGGER IF EXISTS update_affiliates_updated_at ON affiliates;
DROP TRIGGER IF EXISTS update_affiliate_links_updated_at ON affiliate_links;
DROP TRIGGER IF EXISTS update_referrals_updated_at ON referrals;
DROP TRIGGER IF EXISTS update_commissions_updated_at ON commissions;
DROP TRIGGER IF EXISTS update_affiliate_payouts_updated_at ON affiliate_payouts;
DROP TRIGGER IF EXISTS update_affiliate_settings_updated_at ON affiliate_settings;

DROP INDEX IF EXISTS idx_affiliates_user_id;
DROP INDEX IF EXISTS idx_affiliates_organization_id;
DROP INDEX IF EXISTS idx_affiliates_code;
DROP INDEX IF EXISTS idx_affiliates_status;
DROP INDEX IF EXISTS idx_affiliate_links_affiliate_id;
DROP INDEX IF EXISTS idx_affiliate_links_code;
DROP INDEX IF EXISTS idx_affiliate_links_campaign;
DROP INDEX IF EXISTS idx_referrals_affiliate_id;
DROP INDEX IF EXISTS idx_referrals_link_id;
DROP INDEX IF EXISTS idx_referrals_status;
DROP INDEX IF EXISTS idx_referrals_referred_user;
DROP INDEX IF EXISTS idx_commissions_affiliate_id;
DROP INDEX IF EXISTS idx_commissions_status;
DROP INDEX IF EXISTS idx_commissions_referral_id;
DROP INDEX IF EXISTS idx_affiliate_payouts_affiliate_id;
DROP INDEX IF EXISTS idx_affiliate_payouts_status;
DROP INDEX IF EXISTS idx_affiliate_settings_organization;

DROP TABLE IF EXISTS affiliate_payouts;
DROP TABLE IF EXISTS commissions;
DROP TABLE IF EXISTS referrals;
DROP TABLE IF EXISTS affiliate_links;
DROP TABLE IF EXISTS affiliate_settings;
DROP TABLE IF EXISTS affiliates;