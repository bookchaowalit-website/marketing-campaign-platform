-- Email Marketing Complete Migration - UP
-- Adds A/B testing, individual send tracking, and enhanced templates

-- ===========================
-- EMAIL CAMPAIGN A/B TESTS
-- ===========================
CREATE TABLE IF NOT EXISTS email_campaign_ab_tests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id_a UUID NOT NULL REFERENCES email_campaigns(id) ON DELETE CASCADE,
    campaign_id_b UUID NOT NULL REFERENCES email_campaigns(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'running' CHECK (status IN ('running', 'completed', 'cancelled')),
    winner VARCHAR(1) CHECK (winner IN ('A', 'B', NULL)),
    test_metric VARCHAR(50) DEFAULT 'open_rate' CHECK (test_metric IN ('open_rate', 'click_rate', 'conversion_rate')),
    sample_size_percent INTEGER DEFAULT 20 CHECK (sample_size_percent > 0 AND sample_size_percent <= 100),
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CHECK (campaign_id_a != campaign_id_b)
);

-- ===========================
-- EMAIL SENDS (individual send tracking)
-- ===========================
CREATE TABLE IF NOT EXISTS email_sends (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES email_campaigns(id) ON DELETE CASCADE,
    subscriber_id UUID NOT NULL REFERENCES email_subscribers(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'delivered', 'opened', 'clicked', 'bounced', 'unsubscribed', 'complained')),
    send_error TEXT,
    sent_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    opened_at TIMESTAMP WITH TIME ZONE,
    clicked_at TIMESTAMP WITH TIME ZONE,
    bounced_at TIMESTAMP WITH TIME ZONE,
    bounce_type VARCHAR(50),
    bounce_reason TEXT,
    unsubscribed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(campaign_id, subscriber_id)
);

-- ===========================
-- ENHANCE EMAIL TEMPLATES
-- ===========================
ALTER TABLE email_templates
    ADD COLUMN IF NOT EXISTS subject_line VARCHAR(500),
    ADD COLUMN IF NOT EXISTS preheader_text VARCHAR(500),
    ADD COLUMN IF NOT EXISTS plain_text_content TEXT;

-- ===========================
-- INDEXES
-- ===========================
CREATE INDEX IF NOT EXISTS idx_email_ab_tests_status ON email_campaign_ab_tests(status);
CREATE INDEX IF NOT EXISTS idx_email_ab_tests_campaign_a ON email_campaign_ab_tests(campaign_id_a);
CREATE INDEX IF NOT EXISTS idx_email_ab_tests_campaign_b ON email_campaign_ab_tests(campaign_id_b);

CREATE INDEX IF NOT EXISTS idx_email_sends_campaign ON email_sends(campaign_id);
CREATE INDEX IF NOT EXISTS idx_email_sends_subscriber ON email_sends(subscriber_id);
CREATE INDEX IF NOT EXISTS idx_email_sends_status ON email_sends(status);
CREATE INDEX IF NOT EXISTS idx_email_sends_sent_at ON email_sends(sent_at);
CREATE INDEX IF NOT EXISTS idx_email_sends_opened_at ON email_sends(opened_at);
CREATE INDEX IF NOT EXISTS idx_email_sends_campaign_status ON email_sends(campaign_id, status);

-- ===========================
-- TRIGGERS
-- ===========================
DROP TRIGGER IF EXISTS update_email_ab_tests_updated_at ON email_campaign_ab_tests;
CREATE TRIGGER update_email_ab_tests_updated_at
    BEFORE UPDATE ON email_campaign_ab_tests
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_email_sends_updated_at ON email_sends;
CREATE TRIGGER update_email_sends_updated_at
    BEFORE UPDATE ON email_sends
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
