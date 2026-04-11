-- Billing & Subscription Migration - DOWN

ALTER TABLE organizations DROP CONSTRAINT IF EXISTS fk_organizations_plan_id;

DROP TRIGGER IF EXISTS update_billing_addresses_updated_at ON billing_addresses;
DROP TRIGGER IF EXISTS update_tax_settings_updated_at ON tax_settings;
DROP TRIGGER IF EXISTS update_usage_summaries_updated_at ON usage_summaries;
DROP TRIGGER IF EXISTS update_coupons_updated_at ON coupons;
DROP TRIGGER IF EXISTS update_invoices_updated_at ON invoices;
DROP TRIGGER IF EXISTS update_payment_methods_updated_at ON payment_methods;
DROP TRIGGER IF EXISTS update_subscriptions_updated_at ON subscriptions;
DROP TRIGGER IF EXISTS update_subscription_plans_updated_at ON subscription_plans;

DROP INDEX IF EXISTS idx_billing_webhook_events_org_id;
DROP INDEX IF EXISTS idx_billing_webhook_events_processed;
DROP INDEX IF EXISTS idx_billing_webhook_events_type;
DROP INDEX IF EXISTS idx_billing_webhook_events_event_id;
DROP INDEX IF EXISTS idx_billing_addresses_org_id;
DROP INDEX IF EXISTS idx_tax_settings_org_id;
DROP INDEX IF EXISTS idx_customer_balances_org_id;
DROP INDEX IF EXISTS idx_promotion_codes_coupon_id;
DROP INDEX IF EXISTS idx_promotion_codes_code;
DROP INDEX IF EXISTS idx_coupons_active;
DROP INDEX IF EXISTS idx_coupons_code;
DROP INDEX IF EXISTS idx_usage_summaries_period;
DROP INDEX IF EXISTS idx_usage_summaries_metric_type;
DROP INDEX IF EXISTS idx_usage_summaries_org_id;
DROP INDEX IF EXISTS idx_usage_records_recorded_at;
DROP INDEX IF EXISTS idx_usage_records_metric_type;
DROP INDEX IF EXISTS idx_usage_records_subscription_id;
DROP INDEX IF EXISTS idx_usage_records_org_id;
DROP INDEX IF EXISTS idx_invoice_items_subscription_id;
DROP INDEX IF EXISTS idx_invoice_items_invoice_id;
DROP INDEX IF EXISTS idx_invoices_due_at;
DROP INDEX IF EXISTS idx_invoices_stripe_id;
DROP INDEX IF EXISTS idx_invoices_status;
DROP INDEX IF EXISTS idx_invoices_subscription_id;
DROP INDEX IF EXISTS idx_invoices_org_id;
DROP INDEX IF EXISTS idx_payment_methods_default;
DROP INDEX IF EXISTS idx_payment_methods_org_id;
DROP INDEX IF EXISTS idx_subscriptions_current_period;
DROP INDEX IF EXISTS idx_subscriptions_paypal_sub_id;
DROP INDEX IF EXISTS idx_subscriptions_stripe_cust_id;
DROP INDEX IF EXISTS idx_subscriptions_stripe_sub_id;
DROP INDEX IF EXISTS idx_subscriptions_status;
DROP INDEX IF EXISTS idx_subscriptions_plan_id;
DROP INDEX IF EXISTS idx_subscriptions_org_id;
DROP INDEX IF EXISTS idx_subscription_plans_active;
DROP INDEX IF EXISTS idx_subscription_plans_slug;

DELETE FROM subscription_plans WHERE slug IN ('starter', 'professional', 'enterprise');

DROP TABLE IF EXISTS billing_webhook_events;
DROP TABLE IF EXISTS billing_addresses;
DROP TABLE IF EXISTS tax_settings;
DROP TABLE IF EXISTS customer_balances;
DROP TABLE IF EXISTS promotion_codes;
DROP TABLE IF EXISTS coupons;
DROP TABLE IF EXISTS usage_summaries;
DROP TABLE IF EXISTS usage_records;
DROP TABLE IF EXISTS invoice_items;
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS payment_methods;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS subscription_plans;