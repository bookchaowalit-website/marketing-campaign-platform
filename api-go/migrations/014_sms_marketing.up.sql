-- Phase 2: SMS Marketing

-- ===========================
-- SMS LISTS
-- ===========================

CREATE TABLE IF NOT EXISTS sms_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    subscriber_count INT DEFAULT 0,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'archived')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ===========================
-- SMS SUBSCRIBERS
-- ===========================

CREATE TABLE IF NOT EXISTS sms_subscribers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL,
    list_id UUID REFERENCES sms_lists(id) ON DELETE SET NULL,
    phone VARCHAR(20) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    custom_fields JSONB,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'opted_out', 'bounced')),
    subscribed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    opted_out_at TIMESTAMP WITH TIME ZONE,
    opt_out_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(org_id, phone)
);

-- ===========================
-- SMS CAMPAIGNS
-- ===========================

CREATE TABLE IF NOT EXISTS sms_campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL,
    list_id UUID REFERENCES sms_lists(id),
    name VARCHAR(255) NOT NULL,
    message_template TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'scheduled', 'sending', 'sent', 'failed', 'cancelled')),
    scheduled_at TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,
    sent_count INT DEFAULT 0,
    delivered_count INT DEFAULT 0,
    failed_count INT DEFAULT 0,
    opt_out_count INT DEFAULT 0,
    cost DECIMAL(10,4) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ===========================
-- SMS LOGS
-- ===========================

CREATE TABLE IF NOT EXISTS sms_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES sms_campaigns(id) ON DELETE CASCADE,
    subscriber_id UUID REFERENCES sms_subscribers(id) ON DELETE SET NULL,
    phone VARCHAR(20) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'delivered', 'failed', 'undelivered')),
    message_sid VARCHAR(255),
    sent_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    cost DECIMAL(10,4),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ===========================
-- INDEXES
-- ===========================

CREATE INDEX idx_sms_lists_org ON sms_lists(org_id);
CREATE INDEX idx_sms_lists_status ON sms_lists(status);

CREATE INDEX idx_sms_subscribers_org ON sms_subscribers(org_id);
CREATE INDEX idx_sms_subscribers_list ON sms_subscribers(list_id);
CREATE INDEX idx_sms_subscribers_phone ON sms_subscribers(phone);
CREATE INDEX idx_sms_subscribers_status ON sms_subscribers(status);

CREATE INDEX idx_sms_campaigns_org ON sms_campaigns(org_id);
CREATE INDEX idx_sms_campaigns_list ON sms_campaigns(list_id);
CREATE INDEX idx_sms_campaigns_status ON sms_campaigns(status);
CREATE INDEX idx_sms_campaigns_scheduled ON sms_campaigns(scheduled_at);

CREATE INDEX idx_sms_logs_campaign ON sms_logs(campaign_id);
CREATE INDEX idx_sms_logs_subscriber ON sms_logs(subscriber_id);
CREATE INDEX idx_sms_logs_status ON sms_logs(status);
CREATE INDEX idx_sms_logs_phone ON sms_logs(phone);

-- ===========================
-- TRIGGERS
-- ===========================

CREATE TRIGGER update_sms_lists_updated_at BEFORE UPDATE ON sms_lists FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_sms_subscribers_updated_at BEFORE UPDATE ON sms_subscribers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_sms_campaigns_updated_at BEFORE UPDATE ON sms_campaigns FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
