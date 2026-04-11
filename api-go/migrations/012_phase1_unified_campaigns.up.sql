-- Phase 1: Unified Campaign Management, Budget Planning, Marketing Calendar, Approval Workflows

-- ===========================
-- UNIFIED CAMPAIGNS
-- ===========================

-- Campaign groups (for cross-channel campaigns)
CREATE TABLE IF NOT EXISTS campaign_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    objective VARCHAR(100) CHECK (objective IN ('awareness', 'acquisition', 'engagement', 'retention', 'revenue', 'loyalty')),
    status VARCHAR(50) DEFAULT 'planning' CHECK (status IN ('planning', 'active', 'paused', 'completed', 'archived')),
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    total_budget DECIMAL(15,2) DEFAULT 0,
    total_spent DECIMAL(15,2) DEFAULT 0,
    owner_id UUID REFERENCES users(id),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Unified campaigns table (extends existing campaigns)
-- Supports ALL online + offline channels
CREATE TABLE IF NOT EXISTS unified_campaignes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    group_id UUID REFERENCES campaign_groups(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Channel classification
    channel_category VARCHAR(50) NOT NULL CHECK (channel_category IN ('online', 'offline', 'hybrid')),
    channel VARCHAR(100) NOT NULL CHECK (channel IN (
        -- Online channels
        'google_ads', 'facebook_ads', 'instagram_ads', 'tiktok_ads', 'linkedin_ads', 'twitter_ads',
        'seo', 'email', 'sms', 'push_notifications', 'video_youtube', 'video_vimeo', 'webinar',
        'content_marketing', 'affiliate', 'influencer', 'display_ads', 'native_ads', 'podcast',
        -- Offline channels
        'tv', 'radio', 'print_magazine', 'print_newspaper', 'billboard', 'transit_ad',
        'direct_mail', 'event', 'trade_show', 'conference', 'sponsorship', 'telemarketing',
        'pos_materials', 'in_store_promotion', 'brochure', 'flyer', 'guerrilla',
        -- Hybrid
        'cross_channel', 'omnichannel'
    )),
    
    -- Budget & Performance
    budget DECIMAL(15,2) DEFAULT 0,
    spent DECIMAL(15,2) DEFAULT 0,
    cost_per_acquisition DECIMAL(10,2),
    return_on_ad_spend DECIMAL(10,2),
    
    -- Targeting
    target_audience TEXT,
    geography TEXT[],
    demographics JSONB,
    
    -- Scheduling
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'planning', 'scheduled', 'active', 'paused', 'completed', 'cancelled')),
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    timezone VARCHAR(50) DEFAULT 'UTC',
    
    -- Tracking
    utm_source VARCHAR(255),
    utm_medium VARCHAR(255),
    utm_campaign VARCHAR(255),
    tracking_url VARCHAR(500),
    
    -- Creative
    creative_brief TEXT,
    assets JSONB,
    
    -- Approval workflow
    approval_status VARCHAR(50) DEFAULT 'pending' CHECK (approval_status IN ('not_required', 'pending', 'approved', 'rejected', 'revisions_requested')),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP WITH TIME ZONE,
    rejection_reason TEXT,
    
    -- Attribution
    impressions BIGINT DEFAULT 0,
    clicks BIGINT DEFAULT 0,
    conversions BIGINT DEFAULT 0,
    revenue DECIMAL(15,2) DEFAULT 0,
    leads INTEGER DEFAULT 0,
    
    -- Metadata
    owner_id UUID REFERENCES users(id),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ===========================
-- BUDGET PLANNING
-- ===========================

-- Budget periods (monthly/quarterly/annual planning)
CREATE TABLE IF NOT EXISTS budget_periods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    period_type VARCHAR(50) CHECK (period_type IN ('monthly', 'quarterly', 'annual', 'custom')),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_budget DECIMAL(15,2) NOT NULL,
    total_allocated DECIMAL(15,2) DEFAULT 0,
    total_spent DECIMAL(15,2) DEFAULT 0,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'approved', 'active', 'closed')),
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Budget allocations per channel
CREATE TABLE IF NOT EXISTS budget_allocations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    period_id UUID REFERENCES budget_periods(id) ON DELETE CASCADE,
    channel VARCHAR(100) NOT NULL,
    channel_category VARCHAR(50) CHECK (channel_category IN ('online', 'offline', 'hybrid')),
    allocated_amount DECIMAL(15,2) NOT NULL,
    spent_amount DECIMAL(15,2) DEFAULT 0,
    forecast_amount DECIMAL(15,2),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(period_id, channel)
);

-- Budget transactions (actual spend tracking)
CREATE TABLE IF NOT EXISTS budget_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    campaign_id UUID REFERENCES unified_campaignes(id),
    allocation_id UUID REFERENCES budget_allocations(id),
    channel VARCHAR(100) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    invoice_number VARCHAR(100),
    vendor VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Budget alerts
CREATE TABLE IF NOT EXISTS budget_alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    period_id UUID REFERENCES budget_periods(id),
    campaign_id UUID REFERENCES unified_campaignes(id),
    alert_type VARCHAR(50) CHECK (alert_type IN ('threshold', 'overspend', 'forecast')),
    threshold_percent DECIMAL(5,2),
    triggered BOOLEAN DEFAULT false,
    triggered_at TIMESTAMP WITH TIME ZONE,
    notified_users UUID[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ===========================
-- MARKETING CALENDAR
-- ===========================

-- Unified marketing calendar events
CREATE TABLE IF NOT EXISTS marketing_calendar (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    campaign_id UUID REFERENCES unified_campaignes(id),
    event_type VARCHAR(100) NOT NULL CHECK (event_type IN (
        'campaign_launch', 'campaign_end', 'content_publish', 'social_post',
        'email_send', 'sms_send', 'event', 'trade_show', 'webinar',
        'deadline', 'review', 'approval', 'milestone', 'custom'
    )),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    all_day BOOLEAN DEFAULT false,
    color VARCHAR(7) DEFAULT '#1976d2',
    channel VARCHAR(100),
    status VARCHAR(50) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'in_progress', 'completed', 'cancelled', 'postponed')),
    location VARCHAR(255),
    location_url VARCHAR(500),
    attendees UUID[],
    tags TEXT[],
    reminders JSONB, -- [{type: 'email'|'push', minutes_before: 30}]
    recurrence JSONB, -- {frequency: 'daily'|'weekly'|'monthly', interval: 1, end_date: ...}
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Calendar views (saved views for users)
CREATE TABLE IF NOT EXISTS calendar_views (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    user_id UUID REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    view_type VARCHAR(50) CHECK (view_type IN ('day', 'week', 'month', 'quarter', 'year', 'timeline')),
    filters JSONB,
    channels TEXT[],
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ===========================
-- APPROVAL WORKFLOWS
-- ===========================

-- Approval workflow templates
CREATE TABLE IF NOT EXISTS approval_workflows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    trigger_type VARCHAR(100) CHECK (trigger_type IN (
        'campaign_create', 'campaign_publish', 'content_create', 'content_publish',
        'budget_create', 'budget_approve', 'vendor_create', 'custom'
    )),
    channels TEXT[],
    min_approvers INTEGER DEFAULT 1,
    approval_order VARCHAR(50) CHECK (approval_order IN ('sequential', 'parallel', 'any')),
    auto_approve_after_hours INTEGER,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Approval steps
CREATE TABLE IF NOT EXISTS approval_steps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID REFERENCES approval_workflows(id) ON DELETE CASCADE,
    step_order INTEGER NOT NULL,
    approver_role VARCHAR(100),
    approver_user UUID REFERENCES users(id),
    required BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Approval requests
CREATE TABLE IF NOT EXISTS approval_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    workflow_id UUID REFERENCES approval_workflows(id),
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    submitted_by UUID REFERENCES users(id),
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled', 'expired')),
    deadline TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    comments TEXT
);

-- Approval responses
CREATE TABLE IF NOT EXISTS approval_responses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id UUID REFERENCES approval_requests(id) ON DELETE CASCADE,
    step_id UUID REFERENCES approval_steps(id),
    approver_id UUID REFERENCES users(id),
    action VARCHAR(50) CHECK (action IN ('approved', 'rejected', 'commented', 'requested_revision')),
    comments TEXT,
    responded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ===========================
-- INDEXES
-- ===========================

CREATE INDEX idx_campaign_groups_org ON campaign_groups(organization_id);
CREATE INDEX idx_campaign_groups_status ON campaign_groups(status);

CREATE INDEX idx_unified_campaigns_org ON unified_campaignes(organization_id);
CREATE INDEX idx_unified_campaigns_channel ON unified_campaignes(channel);
CREATE INDEX idx_unified_campaigns_channel_cat ON unified_campaignes(channel_category);
CREATE INDEX idx_unified_campaigns_status ON unified_campaignes(status);
CREATE INDEX idx_unified_campaigns_group ON unified_campaignes(group_id);
CREATE INDEX idx_unified_campaigns_owner ON unified_campaignes(owner_id);
CREATE INDEX idx_unified_campaigns_dates ON unified_campaignes(start_date, end_date);

CREATE INDEX idx_budget_periods_org ON budget_periods(organization_id);
CREATE INDEX idx_budget_periods_dates ON budget_periods(start_date, end_date);
CREATE INDEX idx_budget_alloc_period ON budget_allocations(period_id);
CREATE INDEX idx_budget_txn_org ON budget_transactions(organization_id);
CREATE INDEX idx_budget_txn_campaign ON budget_transactions(campaign_id);
CREATE INDEX idx_budget_txn_date ON budget_transactions(transaction_date);

CREATE INDEX idx_budget_alerts_org ON budget_alerts(organization_id);
CREATE INDEX idx_budget_alerts_triggered ON budget_alerts(triggered);

CREATE INDEX idx_marketing_calendar_org ON marketing_calendar(organization_id);
CREATE INDEX idx_marketing_calendar_times ON marketing_calendar(start_time, end_time);
CREATE INDEX idx_marketing_calendar_type ON marketing_calendar(event_type);
CREATE INDEX idx_marketing_calendar_status ON marketing_calendar(status);
CREATE INDEX idx_marketing_calendar_campaign ON marketing_calendar(campaign_id);

CREATE INDEX idx_calendar_views_user ON calendar_views(user_id);

CREATE INDEX idx_approval_workflows_org ON approval_workflows(organization_id);
CREATE INDEX idx_approval_steps_workflow ON approval_steps(workflow_id);
CREATE INDEX idx_approval_requests_org ON approval_requests(organization_id);
CREATE INDEX idx_approval_requests_status ON approval_requests(status);
CREATE INDEX idx_approval_responses_request ON approval_responses(request_id);

-- ===========================
-- TRIGGERS
-- ===========================

CREATE TRIGGER update_campaign_groups_updated_at BEFORE UPDATE ON campaign_groups FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_unified_campaignes_updated_at BEFORE UPDATE ON unified_campaignes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_budget_periods_updated_at BEFORE UPDATE ON budget_periods FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_budget_allocations_updated_at BEFORE UPDATE ON budget_allocations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_marketing_calendar_updated_at BEFORE UPDATE ON marketing_calendar FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_approval_workflows_updated_at BEFORE UPDATE ON approval_workflows FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ===========================
-- SEED DATA
-- ===========================

-- Default approval workflow for campaign publishing
INSERT INTO approval_workflows (id, organization_id, name, description, trigger_type, channels, min_approvers, approval_order, status)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    'Campaign Publishing Approval',
    'All campaigns must be approved before going live',
    'campaign_publish',
    ARRAY['google_ads', 'facebook_ads', 'instagram_ads', 'tv', 'radio', 'print_magazine', 'billboard', 'event'],
    1,
    'sequential',
    'active'
) ON CONFLICT DO NOTHING;

-- Default budget period (current month)
INSERT INTO budget_periods (id, organization_id, name, period_type, start_date, end_date, total_budget, status)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    'Current Month Budget',
    'monthly',
    date_trunc('month', CURRENT_DATE),
    (date_trunc('month', CURRENT_DATE) + INTERVAL '1 month - 1 day')::date,
    50000.00,
    'active'
) ON CONFLICT DO NOTHING;
