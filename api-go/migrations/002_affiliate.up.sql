-- Affiliate Program Migration - UP
-- Creates tables for affiliate/referral program functionality

-- Affiliates table - Partner accounts
CREATE TABLE IF NOT EXISTS affiliates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL,
    affiliate_code VARCHAR(50) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'suspended', 'approved')),
    commission_rate DECIMAL(5,4) NOT NULL DEFAULT 0.1000, -- 10% default commission
    commission_type VARCHAR(50) NOT NULL DEFAULT 'percentage' CHECK (commission_type IN ('percentage', 'fixed')),
    fixed_commission_amount DECIMAL(15,2) DEFAULT 0,
    total_clicks BIGINT NOT NULL DEFAULT 0,
    total_signups BIGINT NOT NULL DEFAULT 0,
    total_conversions BIGINT NOT NULL DEFAULT 0,
    total_earnings DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_paid DECIMAL(15,2) NOT NULL DEFAULT 0,
    balance DECIMAL(15,2) NOT NULL DEFAULT 0,
    minimum_payout DECIMAL(15,2) NOT NULL DEFAULT 50.00,
    paypal_email VARCHAR(255),
    bank_name VARCHAR(100),
    bank_account VARCHAR(100),
    bank_routing VARCHAR(50),
    tax_id VARCHAR(50),
    website VARCHAR(255),
    description TEXT,
    promo_materials TEXT, -- JSON array of promotional materials
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    approved_at TIMESTAMP WITH TIME ZONE,
    approved_by UUID REFERENCES users(id)
);

-- Affiliate Links table - Tracking links for campaigns
CREATE TABLE IF NOT EXISTS affiliate_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    affiliate_id UUID REFERENCES affiliates(id) ON DELETE CASCADE,
    campaign_id UUID REFERENCES campaigns(id) ON DELETE SET NULL,
    link_code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    destination_url VARCHAR(500) NOT NULL,
    click_count BIGINT NOT NULL DEFAULT 0,
    signup_count BIGINT NOT NULL DEFAULT 0,
    conversion_count BIGINT NOT NULL DEFAULT 0,
    earnings DECIMAL(15,2) NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Referrals table - Tracks individual referrals
CREATE TABLE IF NOT EXISTS referrals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    affiliate_id UUID REFERENCES affiliates(id) ON DELETE CASCADE,
    affiliate_link_id UUID REFERENCES affiliate_links(id) ON DELETE SET NULL,
    referred_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    referral_code VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'signup', 'qualified', 'converted', 'cancelled')),
    click_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    signup_timestamp TIMESTAMP WITH TIME ZONE,
    conversion_timestamp TIMESTAMP WITH TIME ZONE,
    order_id VARCHAR(100),
    order_amount DECIMAL(15,2) DEFAULT 0,
    commission_amount DECIMAL(15,2) DEFAULT 0,
    commission_status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (commission_status IN ('pending', 'approved', 'paid', 'cancelled', 'refunded')),
    ip_address VARCHAR(45),
    user_agent TEXT,
    referrer_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Affiliate Payouts table - Payment records to affiliates
-- (created before Commissions: Commissions.payout_id references this table)
CREATE TABLE IF NOT EXISTS affiliate_payouts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    affiliate_id UUID REFERENCES affiliates(id) ON DELETE CASCADE,
    amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
    payment_method VARCHAR(50) NOT NULL CHECK (payment_method IN ('paypal', 'bank_transfer', 'stripe', 'manual')),
    payment_reference VARCHAR(100), -- Transaction ID or reference
    payment_notes TEXT,
    commissions JSONB, -- Array of commission IDs included in this payout
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP WITH TIME ZONE,
    processed_by UUID REFERENCES users(id),
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Commissions table - Commission records
CREATE TABLE IF NOT EXISTS commissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    affiliate_id UUID REFERENCES affiliates(id) ON DELETE CASCADE,
    referral_id UUID REFERENCES referrals(id) ON DELETE SET NULL,
    amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'paid', 'cancelled', 'refunded')),
    type VARCHAR(50) NOT NULL DEFAULT 'referral' CHECK (type IN ('referral', 'bonus', 'adjustment', 'referral_bonus')),
    description TEXT,
    reference_id VARCHAR(100), -- External reference (order ID, transaction ID)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    approved_at TIMESTAMP WITH TIME ZONE,
    approved_by UUID REFERENCES users(id),
    paid_at TIMESTAMP WITH TIME ZONE,
    payout_id UUID REFERENCES affiliate_payouts(id)
);

-- Affiliate Settings table - Program configuration per organization
CREATE TABLE IF NOT EXISTS affiliate_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL UNIQUE,
    program_name VARCHAR(255) DEFAULT 'Affiliate Program',
    program_description TEXT,
    default_commission_rate DECIMAL(5,4) DEFAULT 0.1000,
    default_commission_type VARCHAR(50) DEFAULT 'percentage' CHECK (default_commission_type IN ('percentage', 'fixed')),
    default_fixed_commission DECIMAL(15,2) DEFAULT 0,
    cookie_duration_days INTEGER DEFAULT 30,
    minimum_payout_amount DECIMAL(15,2) DEFAULT 50.00,
    payout_schedule VARCHAR(50) DEFAULT 'monthly' CHECK (payout_schedule IN ('weekly', 'biweekly', 'monthly', 'quarterly')),
    require_approval BOOLEAN DEFAULT true,
    allow_custom_commission_rates BOOLEAN DEFAULT true,
    allow_affiliate_links BOOLEAN DEFAULT true,
    terms_and_conditions TEXT,
    welcome_message TEXT,
    email_notifications BOOLEAN DEFAULT true,
    notify_on_signup BOOLEAN DEFAULT true,
    notify_on_conversion BOOLEAN DEFAULT true,
    notify_on_payout BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for affiliate tables
CREATE INDEX idx_affiliates_user_id ON affiliates(user_id);
CREATE INDEX idx_affiliates_organization_id ON affiliates(organization_id);
CREATE INDEX idx_affiliates_code ON affiliates(affiliate_code);
CREATE INDEX idx_affiliates_status ON affiliates(status);
CREATE INDEX idx_affiliate_links_affiliate_id ON affiliate_links(affiliate_id);
CREATE INDEX idx_affiliate_links_code ON affiliate_links(link_code);
CREATE INDEX idx_affiliate_links_campaign ON affiliate_links(campaign_id);
CREATE INDEX idx_referrals_affiliate_id ON referrals(affiliate_id);
CREATE INDEX idx_referrals_link_id ON referrals(affiliate_link_id);
CREATE INDEX idx_referrals_status ON referrals(status);
CREATE INDEX idx_referrals_referred_user ON referrals(referred_user_id);
CREATE INDEX idx_commissions_affiliate_id ON commissions(affiliate_id);
CREATE INDEX idx_commissions_status ON commissions(status);
CREATE INDEX idx_commissions_referral_id ON commissions(referral_id);
CREATE INDEX idx_affiliate_payouts_affiliate_id ON affiliate_payouts(affiliate_id);
CREATE INDEX idx_affiliate_payouts_status ON affiliate_payouts(status);
CREATE INDEX idx_affiliate_settings_organization ON affiliate_settings(organization_id);

-- Triggers for updated_at on affiliate tables
CREATE TRIGGER update_affiliates_updated_at BEFORE UPDATE ON affiliates FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_affiliate_links_updated_at BEFORE UPDATE ON affiliate_links FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_referrals_updated_at BEFORE UPDATE ON referrals FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_commissions_updated_at BEFORE UPDATE ON commissions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_affiliate_payouts_updated_at BEFORE UPDATE ON affiliate_payouts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_affiliate_settings_updated_at BEFORE UPDATE ON affiliate_settings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default affiliate settings for the demo organization
INSERT INTO affiliate_settings (organization_id, program_name, program_description)
VALUES ('00000000-0000-0000-0000-000000000001', 'BookMarketing Affiliate Program', 'Join our affiliate program and earn commissions by referring new customers to our marketing SaaS platform.')
ON CONFLICT (organization_id) DO NOTHING;