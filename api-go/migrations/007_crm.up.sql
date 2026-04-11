-- CRM Module Migration - UP
-- Creates tables for contacts, deals, pipeline stages, lead scoring, and activities

-- =====================================================
-- 1. CONTACTS
-- =====================================================

CREATE TABLE IF NOT EXISTS contacts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(50),
    phone_country VARCHAR(5),
    company VARCHAR(255),
    title VARCHAR(255),
    website VARCHAR(255),
    linkedin_url VARCHAR(500),
    twitter_url VARCHAR(500),
    facebook_url VARCHAR(500),
    instagram_url VARCHAR(500),
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    address_city VARCHAR(100),
    address_state VARCHAR(100),
    address_postal_code VARCHAR(20),
    address_country VARCHAR(2),
    avatar_url VARCHAR(500),
    custom_fields JSONB DEFAULT '{}',
    lead_score INTEGER DEFAULT 0,
    lead_status VARCHAR(20) DEFAULT 'new' CHECK (lead_status IN ('new', 'contacted', 'qualified', 'proposal', 'negotiation', 'converted', 'churned')),
    lead_type VARCHAR(20) CHECK (lead_type IN ('cold', 'warm', 'hot')),
    lifecycle_stage VARCHAR(20) CHECK (lifecycle_stage IN ('subscriber', 'lead', 'marketing_qualified', 'sales_qualified', 'opportunity', 'customer', 'evangelist', 'other')),
    source VARCHAR(100) CHECK (source IN ('website', 'referral', 'social', 'email', 'ads', 'event', 'phone', 'chat', 'affiliate', 'organic', 'import', 'api', 'other')),
    source_details JSONB DEFAULT '{}',
    tags TEXT[] DEFAULT '{}',
    owner_id UUID REFERENCES users(id),
    created_by UUID REFERENCES users(id),
    last_contacted_at TIMESTAMP WITH TIME ZONE,
    last_engaged_at TIMESTAMP WITH TIME ZONE,
    next_follow_up_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    is_primary BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    merged_to UUID REFERENCES contacts(id), -- For merging duplicate contacts
    merged_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 2. CONTACT GROUPS / SEGMENTS
-- =====================================================

CREATE TABLE IF NOT EXISTS contact_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(20) DEFAULT 'static' CHECK (type IN ('static', 'dynamic')),
    filter_rules JSONB DEFAULT '{}', -- For dynamic segments
    contact_count INTEGER DEFAULT 0,
    parent_id UUID REFERENCES contact_groups(id) ON DELETE CASCADE,
    color VARCHAR(7) DEFAULT '#1976D2',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Contact-Group Relationship
CREATE TABLE IF NOT EXISTS contact_group_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    contact_id UUID NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES contact_groups(id) ON DELETE CASCADE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(contact_id, group_id)
);

-- =====================================================
-- 3. PIPELINE STAGES
-- =====================================================

CREATE TABLE IF NOT EXISTS pipeline_stages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(7) DEFAULT '#1976D2',
    position INTEGER DEFAULT 0,
    probability INTEGER DEFAULT 0 CHECK (probability >= 0 AND probability <= 100),
    is_won_stage BOOLEAN DEFAULT false,
    is_lost_stage BOOLEAN DEFAULT false,
    is_system_stage BOOLEAN DEFAULT false,
    requirements TEXT, -- Requirements to enter this stage
    automation_rules JSONB DEFAULT '{}',
    deal_count INTEGER DEFAULT 0,
    deal_value DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 4. PIPELINES
-- =====================================================

CREATE TABLE IF NOT EXISTS pipelines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_default BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    currency VARCHAR(3) DEFAULT 'USD',
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Pipeline-Stage Relationship
CREATE TABLE IF NOT EXISTS pipeline_stage_mappings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pipeline_id UUID NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
    stage_id UUID NOT NULL REFERENCES pipeline_stages(id) ON DELETE CASCADE,
    position INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(pipeline_id, stage_id)
);

-- =====================================================
-- 5. DEALS
-- =====================================================

CREATE TABLE IF NOT EXISTS deals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    pipeline_id UUID NOT NULL REFERENCES pipelines(id),
    stage_id UUID NOT NULL REFERENCES pipeline_stages(id),
    contact_id UUID REFERENCES contacts(id) ON DELETE SET NULL,
    primary_contact_id UUID REFERENCES contacts(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    value DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(20) DEFAULT 'open' CHECK (status IN ('open', 'won', 'lost', 'abandoned')),
    probability INTEGER DEFAULT 0 CHECK (probability >= 0 AND probability <= 100),
    expected_close_date DATE,
    actual_close_date DATE,
    closed_by UUID REFERENCES users(id),
    lost_reason VARCHAR(100),
    lost_reason_details TEXT,
    won_reason VARCHAR(100),
    won_reason_details TEXT,
    owner_id UUID REFERENCES users(id),
    created_by UUID REFERENCES users(id),
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    type VARCHAR(50) CHECK (type IN ('new_business', 'existing_business', 'renewal', 'upsell', 'cross_sell')),
    products JSONB DEFAULT '[]', -- [{ name, quantity, unit_price, total }]
    competitors JSONB DEFAULT '[]',
    custom_fields JSONB DEFAULT '{}',
    tags TEXT[] DEFAULT '{}',
    last_activity_at TIMESTAMP WITH TIME ZONE,
    next_activity_at TIMESTAMP WITH TIME ZONE,
    last_stage_change_at TIMESTAMP WITH TIME ZONE,
    time_in_current_stage INTEGER DEFAULT 0, -- Days in current stage
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 6. DEAL CONTACTS (Multiple contacts per deal)
-- =====================================================

CREATE TABLE IF NOT EXISTS deal_contacts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    deal_id UUID NOT NULL REFERENCES deals(id) ON DELETE CASCADE,
    contact_id UUID NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'participant', -- primary, participant, influencer, decision_maker
    is_primary BOOLEAN DEFAULT false,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(deal_id, contact_id)
);

-- =====================================================
-- 7. LEAD SCORING RULES
-- =====================================================

CREATE TABLE IF NOT EXISTS lead_scoring_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    field VARCHAR(100) NOT NULL, -- email_domain, company, title, source, country, engagement_score, etc.
    operator VARCHAR(20) NOT NULL CHECK (operator IN ('equals', 'not_equals', 'contains', 'not_contains', 'starts_with', 'ends_with', 'exists', 'not_exists', 'greater_than', 'less_than', 'between', 'in', 'not_in')),
    value TEXT,
    points INTEGER NOT NULL CHECK (points >= -100 AND points <= 100),
    weight INTEGER DEFAULT 1,
    is_active BOOLEAN DEFAULT true,
    category VARCHAR(50) DEFAULT 'demographic' CHECK (category IN ('demographic', 'behavioral', 'engagement', 'custom')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 8. LEAD SCORING PROFILES
-- =====================================================

CREATE TABLE IF NOT EXISTS lead_scoring_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_default BOOLEAN DEFAULT false,
    min_score INTEGER DEFAULT 0,
    max_score INTEGER DEFAULT 100,
    score_ranges JSONB DEFAULT '{}', -- { "hot": { min: 80, max: 100 }, "warm": { min: 50, max: 79 }, "cold": { min: 0, max: 49 } }
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 9. CONTACT ACTIVITIES (Timeline)
-- =====================================================

CREATE TABLE IF NOT EXISTS contact_activities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    contact_id UUID NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    deal_id UUID REFERENCES deals(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN (
        'email_sent', 'email_received', 'email_opened', 'email_clicked', 'email_bounced',
        'call_outbound', 'call_inbound', 'call_missed',
        'meeting_scheduled', 'meeting_completed', 'meeting_cancelled',
        'note_added', 'task_created', 'task_completed', 'task_skipped',
        'page_view', 'form_submit', 'download', 'purchase',
        'chat_started', 'chat_message', 'chat_closed',
        'social_follow', 'social_engagement', 'social_mention',
        'deal_created', 'deal_stage_changed', 'deal_won', 'deal_lost',
        'lead_score_changed', 'status_changed', 'owner_changed', 'tag_added', 'tag_removed',
        'sms_sent', 'sms_received',
        'other'
    )),
    title VARCHAR(255),
    description TEXT,
    details JSONB DEFAULT '{}',
    performed_by UUID REFERENCES users(id),
    performed_by_type VARCHAR(20) DEFAULT 'user' CHECK (performed_by_type IN ('user', 'system', 'automation', 'contact')),
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    source VARCHAR(50), -- manual, email, phone, chat, website, automation, api
    source_url VARCHAR(500),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 10. TASKS
-- =====================================================

CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    contact_id UUID REFERENCES contacts(id) ON DELETE SET NULL,
    deal_id UUID REFERENCES deals(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) DEFAULT 'general' CHECK (type IN ('call', 'email', 'meeting', 'follow_up', 'general', 'reminder')),
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled', 'overdue')),
    due_date TIMESTAMP WITH TIME ZONE,
    due_time VARCHAR(10), -- HH:MM format
    completed_at TIMESTAMP WITH TIME ZONE,
    completed_by UUID REFERENCES users(id),
    owner_id UUID NOT NULL REFERENCES users(id),
    assigned_by UUID REFERENCES users(id),
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    reminder_at TIMESTAMP WITH TIME ZONE,
    reminder_sent BOOLEAN DEFAULT false,
    is_recurring BOOLEAN DEFAULT false,
    recurring_pattern VARCHAR(50) CHECK (recurring_pattern IN ('daily', 'weekly', 'monthly', 'yearly', 'custom')),
    recurring_config JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 11. NOTES
-- =====================================================

CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    contact_id UUID REFERENCES contacts(id) ON DELETE SET NULL,
    deal_id UUID REFERENCES deals(id) ON DELETE SET NULL,
    task_id UUID REFERENCES tasks(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    type VARCHAR(20) DEFAULT 'general' CHECK (type IN ('general', 'call', 'meeting', 'email', 'important', 'private')),
    is_private BOOLEAN DEFAULT false,
    mentions TEXT[] DEFAULT '{}', -- Users mentioned in the note
    created_by UUID NOT NULL REFERENCES users(id),
    edited_by UUID REFERENCES users(id),
    edited_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 12. EMAIL TRACKING
-- =====================================================

CREATE TABLE IF NOT EXISTS crm_email_tracking (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    contact_id UUID NOT NULL REFERENCES contacts(id),
    campaign_id UUID REFERENCES email_campaigns(id),
    event_type VARCHAR(20) NOT NULL CHECK (event_type IN ('sent', 'delivered', 'opened', 'clicked', 'bounced', 'unsubscribed', 'complained', 'replied')),
    message_id VARCHAR(255),
    email_subject VARCHAR(500),
    link_url VARCHAR(1000),
    link_text VARCHAR(255),
    user_agent TEXT,
    ip_address VARCHAR(45),
    location JSONB DEFAULT '{}',
    device VARCHAR(50),
    browser VARCHAR(50),
    os VARCHAR(50),
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 13. CUSTOM FIELDS DEFINITIONS
-- =====================================================

CREATE TABLE IF NOT EXISTS custom_field_definitions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN ('contact', 'deal', 'company', 'task')),
    name VARCHAR(100) NOT NULL,
    label VARCHAR(255) NOT NULL,
    field_type VARCHAR(20) NOT NULL CHECK (field_type IN ('text', 'number', 'date', 'datetime', 'select', 'multiselect', 'checkbox', 'radio', 'url', 'email', 'phone', 'textarea', 'currency', 'percentage', 'file')),
    options JSONB DEFAULT '[]', -- For select, multiselect, radio
    default_value TEXT,
    placeholder VARCHAR(255),
    is_required BOOLEAN DEFAULT false,
    is_unique BOOLEAN DEFAULT false,
    is_searchable BOOLEAN DEFAULT true,
    is_visible BOOLEAN DEFAULT true,
    position INTEGER DEFAULT 0,
    validation_rules JSONB DEFAULT '{}',
    group_name VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, entity_type, name)
);

-- =====================================================
-- 14. INDEXES
-- =====================================================

CREATE INDEX idx_contacts_org_id ON contacts(organization_id);
CREATE INDEX idx_contacts_email ON contacts(email);
CREATE INDEX idx_contacts_phone ON contacts(phone);
CREATE INDEX idx_contacts_company ON contacts(company);
CREATE INDEX idx_contacts_lead_status ON contacts(lead_status);
CREATE INDEX idx_contacts_lead_score ON contacts(lead_score);
CREATE INDEX idx_contacts_owner_id ON contacts(owner_id);
CREATE INDEX idx_contacts_source ON contacts(source);
CREATE INDEX idx_contacts_created_at ON contacts(created_at);
CREATE INDEX idx_contacts_is_active ON contacts(is_active);

CREATE INDEX idx_contact_groups_org_id ON contact_groups(organization_id);
CREATE INDEX idx_contact_group_members_contact ON contact_group_members(contact_id);
CREATE INDEX idx_contact_group_members_group ON contact_group_members(group_id);

CREATE INDEX idx_pipeline_stages_org_id ON pipeline_stages(organization_id);
CREATE INDEX idx_pipelines_org_id ON pipelines(organization_id);
CREATE INDEX idx_pipeline_stage_mappings_pipeline ON pipeline_stage_mappings(pipeline_id);
CREATE INDEX idx_pipeline_stage_mappings_stage ON pipeline_stage_mappings(stage_id);

CREATE INDEX idx_deals_org_id ON deals(organization_id);
CREATE INDEX idx_deals_pipeline_id ON deals(pipeline_id);
CREATE INDEX idx_deals_stage_id ON deals(stage_id);
CREATE INDEX idx_deals_contact_id ON deals(contact_id);
CREATE INDEX idx_deals_status ON deals(status);
CREATE INDEX idx_deals_owner_id ON deals(owner_id);
CREATE INDEX idx_deals_expected_close ON deals(expected_close_date);
CREATE INDEX idx_deals_created_at ON deals(created_at);
CREATE INDEX idx_deals_value ON deals(value);

CREATE INDEX idx_deal_contacts_deal ON deal_contacts(deal_id);
CREATE INDEX idx_deal_contacts_contact ON deal_contacts(contact_id);

CREATE INDEX idx_lead_scoring_rules_org_id ON lead_scoring_rules(organization_id);
CREATE INDEX idx_lead_scoring_rules_field ON lead_scoring_rules(field);
CREATE INDEX idx_lead_scoring_rules_active ON lead_scoring_rules(is_active);

CREATE INDEX idx_lead_scoring_profiles_org_id ON lead_scoring_profiles(organization_id);

CREATE INDEX idx_contact_activities_org_id ON contact_activities(organization_id);
CREATE INDEX idx_contact_activities_contact_id ON contact_activities(contact_id);
CREATE INDEX idx_contact_activities_deal_id ON contact_activities(deal_id);
CREATE INDEX idx_contact_activities_type ON contact_activities(type);
CREATE INDEX idx_contact_activities_occurred_at ON contact_activities(occurred_at);
CREATE INDEX idx_contact_activities_performed_by ON contact_activities(performed_by);

CREATE INDEX idx_tasks_org_id ON tasks(organization_id);
CREATE INDEX idx_tasks_contact_id ON tasks(contact_id);
CREATE INDEX idx_tasks_deal_id ON tasks(deal_id);
CREATE INDEX idx_tasks_owner_id ON tasks(owner_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
CREATE INDEX idx_tasks_type ON tasks(type);

CREATE INDEX idx_notes_org_id ON notes(organization_id);
CREATE INDEX idx_notes_contact_id ON notes(contact_id);
CREATE INDEX idx_notes_deal_id ON notes(deal_id);
CREATE INDEX idx_notes_created_by ON notes(created_by);

CREATE INDEX idx_crm_email_tracking_org_id ON crm_email_tracking(organization_id);
CREATE INDEX idx_crm_email_tracking_contact_id ON crm_email_tracking(contact_id);
CREATE INDEX idx_crm_email_tracking_event_type ON crm_email_tracking(event_type);
CREATE INDEX idx_crm_email_tracking_occurred_at ON crm_email_tracking(occurred_at);

CREATE INDEX idx_custom_field_definitions_org_id ON custom_field_definitions(organization_id);
CREATE INDEX idx_custom_field_definitions_entity ON custom_field_definitions(entity_type);

-- =====================================================
-- 15. TRIGGERS
-- =====================================================

CREATE TRIGGER update_contacts_updated_at BEFORE UPDATE ON contacts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_contact_groups_updated_at BEFORE UPDATE ON contact_groups
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pipeline_stages_updated_at BEFORE UPDATE ON pipeline_stages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pipelines_updated_at BEFORE UPDATE ON pipelines
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_deals_updated_at BEFORE UPDATE ON deals
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_lead_scoring_rules_updated_at BEFORE UPDATE ON lead_scoring_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_lead_scoring_profiles_updated_at BEFORE UPDATE ON lead_scoring_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE ON tasks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notes_updated_at BEFORE UPDATE ON notes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_custom_field_definitions_updated_at BEFORE UPDATE ON custom_field_definitions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- 16. DEFAULT PIPELINE STAGES
-- =====================================================

INSERT INTO pipeline_stages (name, color, position, probability, is_system_stage, requirements) VALUES
('Lead In', '#9E9E9E', 0, 10, true, 'Initial contact identified'),
('Qualified', '#FF9800', 1, 20, true, 'Contact meets qualification criteria'),
('Proposal', '#2196F3', 2, 50, true, 'Proposal sent to contact'),
('Negotiation', '#FFC107', 3, 75, true, 'In negotiation phase'),
('Won', '#4CAF50', 4, 100, true, 'Deal successfully closed'),
('Lost', '#F44336', 5, 0, true, 'Deal was lost');

-- =====================================================
-- 17. DEFAULT LEAD SCORING PROFILE
-- =====================================================

INSERT INTO lead_scoring_profiles (name, description, is_default, score_ranges) VALUES
('Default Scoring Profile', 'Standard lead scoring model for B2B sales', true, '{"hot": {"min": 80, "max": 100, "label": "Hot Lead"}, "warm": {"min": 50, "max": 79, "label": "Warm Lead"}, "cold": {"min": 0, "max": 49, "label": "Cold Lead"}}'::jsonb);

-- =====================================================
-- 18. DEFAULT LEAD SCORING RULES
-- =====================================================

INSERT INTO lead_scoring_rules (organization_id, name, field, operator, value, points, category) VALUES
('00000000-0000-0000-0000-000000000001', 'Has valid email', 'email', 'exists', NULL, 5, 'demographic'),
('00000000-0000-0000-0000-000000000001', 'Has phone number', 'phone', 'exists', NULL, 5, 'demographic'),
('00000000-0000-0000-0000-000000000001', 'Has company name', 'company', 'exists', NULL, 10, 'demographic'),
('00000000-0000-0000-0000-000000000001', 'Has title/position', 'title', 'exists', NULL, 5, 'demographic'),
('00000000-0000-0000-0000-000000000001', 'Decision maker title', 'title', 'contains', 'manager,director,ceo,cfo,vp,vice president', 15, 'demographic'),
('00000000-0000-0000-0000-000000000001', 'Company website exists', 'website', 'exists', NULL, 5, 'demographic'),
('00000000-0000-0000-0000-000000000001', 'From referral', 'source', 'equals', 'referral', 20, 'behavioral'),
('00000000-0000-0000-0000-000000000001', 'From organic search', 'source', 'equals', 'organic', 10, 'behavioral'),
('00000000-0000-0000-0000-000000000001', 'From paid ads', 'source', 'equals', 'ads', 5, 'behavioral'),
('00000000-0000-0000-0000-000000000001', 'Requested demo', 'lead_status', 'equals', 'qualified', 25, 'behavioral');

-- =====================================================
-- 19. FUNCTIONS FOR LEAD SCORING
-- =====================================================

CREATE OR REPLACE FUNCTION calculate_contact_lead_score(contact_id UUID)
RETURNS INTEGER AS $$
DECLARE
    total_score INTEGER := 0;
    rule_record RECORD;
    contact_record contacts;
BEGIN
    SELECT * FROM contacts WHERE id = contact_id INTO contact_record;
    
    IF contact_record IS NULL THEN
        RETURN 0;
    END IF;
    
    FOR rule_record IN 
        SELECT field, operator, value, points 
        FROM lead_scoring_rules 
        WHERE organization_id = contact_record.organization_id AND is_active = true
    LOOP
        CASE rule_record.operator
            WHEN 'exists' THEN
                CASE rule_record.field
                    WHEN 'email' THEN IF contact_record.email IS NOT NULL AND contact_record.email != '' THEN total_score := total_score + rule_record.points; END IF;
                    WHEN 'phone' THEN IF contact_record.phone IS NOT NULL AND contact_record.phone != '' THEN total_score := total_score + rule_record.points; END IF;
                    WHEN 'company' THEN IF contact_record.company IS NOT NULL AND contact_record.company != '' THEN total_score := total_score + rule_record.points; END IF;
                    WHEN 'title' THEN IF contact_record.title IS NOT NULL AND contact_record.title != '' THEN total_score := total_score + rule_record.points; END IF;
                    WHEN 'website' THEN IF contact_record.website IS NOT NULL AND contact_record.website != '' THEN total_score := total_score + rule_record.points; END IF;
                END CASE;
            WHEN 'equals' THEN
                CASE rule_record.field
                    WHEN 'source' THEN IF contact_record.source = rule_record.value THEN total_score := total_score + rule_record.points; END IF;
                    WHEN 'lead_status' THEN IF contact_record.lead_status = rule_record.value THEN total_score := total_score + rule_record.points; END IF;
                END CASE;
            WHEN 'contains' THEN
                CASE rule_record.field
                    WHEN 'title' THEN IF contact_record.title ILIKE '%' || rule_record.value || '%' THEN total_score := total_score + rule_record.points; END IF;
                END CASE;
        END CASE;
    END LOOP;
    
    UPDATE contacts SET lead_score = total_score WHERE id = contact_id;
    
    RETURN total_score;
END;
$$ LANGUAGE plpgsql;