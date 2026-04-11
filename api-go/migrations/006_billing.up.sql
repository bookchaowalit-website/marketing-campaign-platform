-- Billing & Subscription Migration - UP
-- Creates tables for subscription plans, billing, invoices, and payment methods

-- =====================================================
-- 1. SUBSCRIPTION PLANS
-- =====================================================

CREATE TABLE IF NOT EXISTS subscription_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100),
    description TEXT,
    price_monthly DECIMAL(10,2) NOT NULL,
    price_yearly DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    features JSONB DEFAULT '[]',
    limits JSONB DEFAULT '{}',
    trial_days INTEGER DEFAULT 0,
    is_public BOOLEAN DEFAULT true,
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    stripe_product_id VARCHAR(255),
    stripe_price_id_monthly VARCHAR(255),
    stripe_price_id_yearly VARCHAR(255),
    paypal_plan_id_monthly VARCHAR(255),
    paypal_plan_id_yearly VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 2. SUBSCRIPTIONS
-- =====================================================

CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES subscription_plans(id),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'trialing', 'past_due', 'canceled', 'incomplete', 'incomplete_expired', 'unpaid', 'paused')),
    billing_cycle VARCHAR(20) NOT NULL DEFAULT 'monthly' CHECK (billing_cycle IN ('monthly', 'yearly')),
    current_period_start TIMESTAMP WITH TIME ZONE,
    current_period_end TIMESTAMP WITH TIME ZONE,
    trial_start TIMESTAMP WITH TIME ZONE,
    trial_end TIMESTAMP WITH TIME ZONE,
    canceled_at TIMESTAMP WITH TIME ZONE,
    cancel_at_period_end BOOLEAN DEFAULT false,
    cancellation_reason TEXT,
    paused_at TIMESTAMP WITH TIME ZONE,
    stripe_subscription_id VARCHAR(255) UNIQUE,
    stripe_customer_id VARCHAR(255),
    paypal_subscription_id VARCHAR(255) UNIQUE,
    paypal_agreement_id VARCHAR(255),
    payment_provider VARCHAR(20) CHECK (payment_provider IN ('stripe', 'paypal')),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id)
);

-- =====================================================
-- 3. PAYMENT METHODS
-- =====================================================

CREATE TABLE IF NOT EXISTS payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('card', 'paypal', 'bank_account')),
    provider VARCHAR(20) NOT NULL CHECK (provider IN ('stripe', 'paypal')),
    provider_payment_method_id VARCHAR(255),
    is_default BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    -- Card details (tokenized)
    card_last4 VARCHAR(4),
    card_brand VARCHAR(20) CHECK (card_brand IN ('visa', 'mastercard', 'amex', 'discover', 'jcb', 'diners', 'unionpay', 'other')),
    card_exp_month INTEGER,
    card_exp_year INTEGER,
    card_funding VARCHAR(20) CHECK (card_funding IN ('credit', 'debit', 'prepaid', 'unknown')),
    card_country VARCHAR(2),
    -- PayPal details
    paypal_email VARCHAR(255),
    -- Bank account details
    bank_name VARCHAR(100),
    bank_last4 VARCHAR(4),
    bank_country VARCHAR(2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 4. INVOICES
-- =====================================================

CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    subscription_id UUID REFERENCES subscriptions(id) ON DELETE SET NULL,
    stripe_invoice_id VARCHAR(255) UNIQUE,
    paypal_invoice_id VARCHAR(255) UNIQUE,
    invoice_number VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'open', 'paid', 'void', 'uncollectible', 'pending')),
    amount_due DECIMAL(12,2) NOT NULL,
    amount_paid DECIMAL(12,2) DEFAULT 0,
    amount_remaining DECIMAL(12,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'USD',
    tax_amount DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,2) DEFAULT 0,
    discount_amount DECIMAL(12,2) DEFAULT 0,
    discount_coupon VARCHAR(100),
    subtotal DECIMAL(12,2) NOT NULL,
    total DECIMAL(12,2) NOT NULL,
    invoice_pdf TEXT,
    invoice_url TEXT,
    due_at TIMESTAMP WITH TIME ZONE,
    paid_at TIMESTAMP WITH TIME ZONE,
    voided_at TIMESTAMP WITH TIME ZONE,
    attempted_count INTEGER DEFAULT 0,
    next_payment_attempt TIMESTAMP WITH TIME ZONE,
    billing_reason VARCHAR(50) CHECK (billing_reason IN ('manual', 'subscription_cycle', 'subscription_update', 'subscription_threshold', 'upcoming', 'subscription_create', 'subscription', 'one_off')),
    description TEXT,
    statement_descriptor VARCHAR(22),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 5. INVOICE ITEMS
-- =====================================================

CREATE TABLE IF NOT EXISTS invoice_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    subscription_id UUID REFERENCES subscriptions(id) ON DELETE SET NULL,
    stripe_invoice_item_id VARCHAR(255),
    type VARCHAR(20) NOT NULL CHECK (type IN ('subscription', 'one_time', 'usage', 'discount', 'tax')),
    description VARCHAR(255) NOT NULL,
    quantity INTEGER DEFAULT 1,
    unit_amount DECIMAL(12,2) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    period_start TIMESTAMP WITH TIME ZONE,
    period_end TIMESTAMP WITH TIME ZONE,
    discount_amounts JSONB DEFAULT '[]',
    tax_amounts JSONB DEFAULT '[]',
    proration BOOLEAN DEFAULT false,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 6. USAGE RECORDS (Metered Billing)
-- =====================================================

CREATE TABLE IF NOT EXISTS usage_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    subscription_id UUID REFERENCES subscriptions(id) ON DELETE SET NULL,
    metric_type VARCHAR(50) NOT NULL CHECK (metric_type IN ('ai_credits', 'sms_sent', 'emails_sent', 'contacts', 'storage_mb', 'api_calls', 'chat_sessions', 'webhook_calls')),
    quantity DECIMAL(10,2) NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,4) DEFAULT 0,
    amount DECIMAL(12,2) DEFAULT 0,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    billing_period_start TIMESTAMP WITH TIME ZONE,
    billing_period_end TIMESTAMP WITH TIME ZONE,
    invoice_item_id UUID REFERENCES invoice_items(id),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 7. USAGE SUMMARY (Aggregated)
-- =====================================================

CREATE TABLE IF NOT EXISTS usage_summaries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    subscription_id UUID REFERENCES subscriptions(id) ON DELETE SET NULL,
    metric_type VARCHAR(50) NOT NULL,
    period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    total_quantity DECIMAL(10,2) DEFAULT 0,
    total_amount DECIMAL(12,2) DEFAULT 0,
    limit_amount DECIMAL(10,2),
    overage_quantity DECIMAL(10,2) DEFAULT 0,
    overage_amount DECIMAL(12,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, metric_type, period_start, period_end)
);

-- =====================================================
-- 8. COUPONS & DISCOUNTS
-- =====================================================

CREATE TABLE IF NOT EXISTS coupons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    discount_type VARCHAR(20) NOT NULL CHECK (discount_type IN ('percentage', 'fixed_amount', 'trial_extension')),
    discount_percent DECIMAL(5,2),
    discount_amount DECIMAL(12,2),
    trial_extension_days INTEGER,
    currency VARCHAR(3) DEFAULT 'USD',
    duration VARCHAR(20) NOT NULL CHECK (duration IN ('once', 'repeating', 'forever')),
    duration_months INTEGER,
    max_redemptions INTEGER,
    redemptions_count INTEGER DEFAULT 0,
    min_amount DECIMAL(12,2),
    applies_to_plans UUID[] DEFAULT '{}',
    first_time_transaction BOOLEAN DEFAULT false, -- Only for first purchase
    valid_from TIMESTAMP WITH TIME ZONE,
    valid_until TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT true,
    stripe_coupon_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 9. PROMOTION CODES
-- =====================================================

CREATE TABLE IF NOT EXISTS promotion_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    coupon_id UUID NOT NULL REFERENCES coupons(id) ON DELETE CASCADE,
    code VARCHAR(50) UNIQUE NOT NULL,
    active BOOLEAN DEFAULT true,
    code_type VARCHAR(20) CHECK (code_type IN ('public', 'unique', 'bulk')),
    max_uses INTEGER,
    uses_count INTEGER DEFAULT 0,
    expires_at TIMESTAMP WITH TIME ZONE,
    restrictions JSONB DEFAULT '{}', -- { minimum_amount, first_time_transaction, etc. }
    stripe_promotion_code_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 10. CUSTOMER BALANCE
-- =====================================================

CREATE TABLE IF NOT EXISTS customer_balances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('credit', 'debit')),
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    description TEXT,
    reference_id UUID,
    reference_type VARCHAR(50),
    stripe_balance_transaction_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 11. TAX SETTINGS
-- =====================================================

CREATE TABLE IF NOT EXISTS tax_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    tax_enabled BOOLEAN DEFAULT false,
    tax_type VARCHAR(20) CHECK (tax_type IN ('vat', 'sales_tax', 'gst', 'igst', 'qst', 'pst', 'hst')),
    default_tax_rate DECIMAL(5,2) DEFAULT 0,
    tax_id VARCHAR(100),
    tax_id_type VARCHAR(50),
    tax_exempt BOOLEAN DEFAULT false,
    automatic_tax BOOLEAN DEFAULT false,
    stripe_tax_settings_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id)
);

-- =====================================================
-- 12. BILLING ADDRESS
-- =====================================================

CREATE TABLE IF NOT EXISTS billing_addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255),
    line1 VARCHAR(255) NOT NULL,
    line2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(2) NOT NULL,
    phone VARCHAR(50),
    tax_id VARCHAR(100),
    is_default BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 13. WEBHOOK EVENTS
-- =====================================================

CREATE TABLE IF NOT EXISTS billing_webhook_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider VARCHAR(20) NOT NULL CHECK (provider IN ('stripe', 'paypal')),
    event_id VARCHAR(255) UNIQUE NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
    subscription_id UUID REFERENCES subscriptions(id) ON DELETE SET NULL,
    invoice_id UUID REFERENCES invoices(id) ON DELETE SET NULL,
    payload JSONB NOT NULL,
    processed BOOLEAN DEFAULT false,
    processing_attempts INTEGER DEFAULT 0,
    last_processing_attempt TIMESTAMP WITH TIME ZONE,
    processing_error TEXT,
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 14. INDEXES
-- =====================================================

CREATE INDEX idx_subscription_plans_slug ON subscription_plans(slug);
CREATE INDEX idx_subscription_plans_active ON subscription_plans(is_active);

CREATE INDEX idx_subscriptions_org_id ON subscriptions(organization_id);
CREATE INDEX idx_subscriptions_plan_id ON subscriptions(plan_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_stripe_sub_id ON subscriptions(stripe_subscription_id);
CREATE INDEX idx_subscriptions_stripe_cust_id ON subscriptions(stripe_customer_id);
CREATE INDEX idx_subscriptions_paypal_sub_id ON subscriptions(paypal_subscription_id);
CREATE INDEX idx_subscriptions_current_period ON subscriptions(current_period_start, current_period_end);

CREATE INDEX idx_payment_methods_org_id ON payment_methods(organization_id);
CREATE INDEX idx_payment_methods_default ON payment_methods(is_default);

CREATE INDEX idx_invoices_org_id ON invoices(organization_id);
CREATE INDEX idx_invoices_subscription_id ON invoices(subscription_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_stripe_id ON invoices(stripe_invoice_id);
CREATE INDEX idx_invoices_due_at ON invoices(due_at);

CREATE INDEX idx_invoice_items_invoice_id ON invoice_items(invoice_id);
CREATE INDEX idx_invoice_items_subscription_id ON invoice_items(subscription_id);

CREATE INDEX idx_usage_records_org_id ON usage_records(organization_id);
CREATE INDEX idx_usage_records_subscription_id ON usage_records(subscription_id);
CREATE INDEX idx_usage_records_metric_type ON usage_records(metric_type);
CREATE INDEX idx_usage_records_recorded_at ON usage_records(recorded_at);

CREATE INDEX idx_usage_summaries_org_id ON usage_summaries(organization_id);
CREATE INDEX idx_usage_summaries_metric_type ON usage_summaries(metric_type);
CREATE INDEX idx_usage_summaries_period ON usage_summaries(period_start, period_end);

CREATE INDEX idx_coupons_code ON coupons(code);
CREATE INDEX idx_coupons_active ON coupons(is_active);

CREATE INDEX idx_promotion_codes_code ON promotion_codes(code);
CREATE INDEX idx_promotion_codes_coupon_id ON promotion_codes(coupon_id);

CREATE INDEX idx_customer_balances_org_id ON customer_balances(organization_id);

CREATE INDEX idx_tax_settings_org_id ON tax_settings(organization_id);

CREATE INDEX idx_billing_addresses_org_id ON billing_addresses(organization_id);

CREATE INDEX idx_billing_webhook_events_event_id ON billing_webhook_events(event_id);
CREATE INDEX idx_billing_webhook_events_type ON billing_webhook_events(event_type);
CREATE INDEX idx_billing_webhook_events_processed ON billing_webhook_events(processed);
CREATE INDEX idx_billing_webhook_events_org_id ON billing_webhook_events(organization_id);

-- =====================================================
-- 15. TRIGGERS
-- =====================================================

CREATE TRIGGER update_subscription_plans_updated_at BEFORE UPDATE ON subscription_plans
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_subscriptions_updated_at BEFORE UPDATE ON subscriptions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payment_methods_updated_at BEFORE UPDATE ON payment_methods
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_invoices_updated_at BEFORE UPDATE ON invoices
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_coupons_updated_at BEFORE UPDATE ON coupons
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_usage_summaries_updated_at BEFORE UPDATE ON usage_summaries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tax_settings_updated_at BEFORE UPDATE ON tax_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_billing_addresses_updated_at BEFORE UPDATE ON billing_addresses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- 16. DEFAULT SUBSCRIPTION PLANS
-- =====================================================

INSERT INTO subscription_plans (name, slug, display_name, description, price_monthly, price_yearly, currency, features, limits, trial_days, is_public, is_active, sort_order) VALUES
('Starter', 'starter', 'Starter Plan', 'Perfect for small businesses and individuals starting their marketing journey', 29.00, 290.00, 'USD', 
 '["Campaign management", "Basic analytics", "Email marketing (1000 subscribers)", "Social media scheduling", "Landing page builder (3 pages)"]'::jsonb,
 '{"contacts": 500, "emails_monthly": 5000, "sms_monthly": 100, "ai_credits_monthly": 50, "landing_pages": 3, "social_accounts": 3, "automation_workflows": 2}'::jsonb,
 14, true, true, 1),
('Professional', 'professional', 'Professional Plan', 'Ideal for growing businesses with advanced marketing needs', 99.00, 990.00, 'USD',
 '["All Starter features", "Advanced analytics & reporting", "Email marketing (5000 subscribers)", "CRM & lead scoring", "AI content generator", "SMS marketing", "Chatbot & live chat", "Unlimited landing pages", "Marketing automation"]'::jsonb,
 '{"contacts": 5000, "emails_monthly": 25000, "sms_monthly": 500, "ai_credits_monthly": 200, "landing_pages": -1, "social_accounts": 10, "automation_workflows": 10, "chatbots": 3}'::jsonb,
 14, true, true, 2),
('Enterprise', 'enterprise', 'Enterprise Plan', 'For large organizations with custom requirements and unlimited access', 299.00, 2990.00, 'USD',
 '["All Professional features", "Unlimited contacts & emails", "White-label branding", "Custom domain", "Priority support", "API access", "Advanced permissions", "Custom integrations"]'::jsonb,
 '{"contacts": -1, "emails_monthly": -1, "sms_monthly": 1000, "ai_credits_monthly": 500, "landing_pages": -1, "social_accounts": -1, "automation_workflows": -1, "chatbots": -1}'::jsonb,
 30, true, true, 3);

-- =====================================================
-- 17. ADD FOREIGN KEY TO ORGANIZATIONS
-- =====================================================

ALTER TABLE organizations ADD CONSTRAINT fk_organizations_plan_id 
    FOREIGN KEY (plan_id) REFERENCES subscription_plans(id) ON DELETE SET NULL;