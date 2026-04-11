-- Marketing Features Migration - UP
-- Creates tables for Email Marketing, Social Media, Landing Pages, Automation, Influencer Marketing, and Reporting

-- =====================================================
-- 1. EMAIL MARKETING MODULE
-- =====================================================

-- Email Lists
CREATE TABLE IF NOT EXISTS email_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'archived', 'deleted')),
    subscriber_count INTEGER DEFAULT 0,
    bounce_count INTEGER DEFAULT 0,
    unsubscribe_count INTEGER DEFAULT 0,
    double_optin BOOLEAN DEFAULT false,
    welcome_email_enabled BOOLEAN DEFAULT false,
    welcome_email_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Email Subscribers
CREATE TABLE IF NOT EXISTS email_subscribers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(50),
    company VARCHAR(255),
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'unsubscribed', 'bounced', 'complained', 'pending')),
    source VARCHAR(100),
    ip_address VARCHAR(45),
    user_agent TEXT,
    custom_fields JSONB,
    tags TEXT[],
    subscribed_at TIMESTAMP WITH TIME ZONE,
    unsubscribed_at TIMESTAMP WITH TIME ZONE,
    bounced_at TIMESTAMP WITH TIME ZONE,
    bounce_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, email)
);

-- List-Subscriber Relationship
CREATE TABLE IF NOT EXISTS email_list_subscribers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    list_id UUID REFERENCES email_lists(id) ON DELETE CASCADE,
    subscriber_id UUID REFERENCES email_subscribers(id) ON DELETE CASCADE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(list_id, subscriber_id)
);

-- Email Templates
CREATE TABLE IF NOT EXISTS email_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    subject VARCHAR(500),
    preview_text VARCHAR(500),
    from_name VARCHAR(255),
    from_email VARCHAR(255),
    reply_to VARCHAR(255),
    html_content TEXT,
    text_content TEXT,
    design_json JSONB,
    thumbnail_url VARCHAR(500),
    category VARCHAR(100),
    template_type VARCHAR(50) DEFAULT 'custom' CHECK (template_type IN ('custom', 'welcome', 'newsletter', 'promotional', 'transactional', 'automation')),
    is_public BOOLEAN DEFAULT false,
    usage_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Email Campaigns
CREATE TABLE IF NOT EXISTS email_campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    subject VARCHAR(500) NOT NULL,
    preview_text VARCHAR(500),
    from_name VARCHAR(255),
    from_email VARCHAR(255),
    reply_to VARCHAR(255),
    template_id UUID REFERENCES email_templates(id),
    html_content TEXT,
    text_content TEXT,
    campaign_type VARCHAR(50) DEFAULT 'regular' CHECK (campaign_type IN ('regular', 'automated', 'ab_test', 'rss')),
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'scheduled', 'sending', 'sent', 'paused', 'cancelled')),
    recipient_count INTEGER DEFAULT 0,
    sent_count INTEGER DEFAULT 0,
    delivered_count INTEGER DEFAULT 0,
    opened_count INTEGER DEFAULT 0,
    clicked_count INTEGER DEFAULT 0,
    bounced_count INTEGER DEFAULT 0,
    unsubscribed_count INTEGER DEFAULT 0,
    complained_count INTEGER DEFAULT 0,
    open_rate DECIMAL(5,2) DEFAULT 0,
    click_rate DECIMAL(5,2) DEFAULT 0,
    bounce_rate DECIMAL(5,2) DEFAULT 0,
    unsubscribe_rate DECIMAL(5,2) DEFAULT 0,
    schedule_time TIMESTAMP WITH TIME ZONE,
    sent_at TIMESTAMP WITH TIME ZONE,
    tracking_enabled BOOLEAN DEFAULT true,
    google_analytics BOOLEAN DEFAULT false,
    utm_source VARCHAR(100),
    utm_medium VARCHAR(100),
    utm_campaign VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Campaign-List Relationship
CREATE TABLE IF NOT EXISTS email_campaign_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID REFERENCES email_campaigns(id) ON DELETE CASCADE,
    list_id UUID REFERENCES email_lists(id) ON DELETE CASCADE,
    UNIQUE(campaign_id, list_id)
);

-- Email Events (opens, clicks, etc.)
CREATE TABLE IF NOT EXISTS email_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    campaign_id UUID REFERENCES email_campaigns(id) ON DELETE CASCADE,
    subscriber_id UUID REFERENCES email_subscribers(id) ON DELETE SET NULL,
    event_type VARCHAR(50) NOT NULL CHECK (event_type IN ('sent', 'delivered', 'opened', 'clicked', 'bounced', 'unsubscribed', 'complained')),
    email VARCHAR(255),
    ip_address VARCHAR(45),
    user_agent TEXT,
    url VARCHAR(1000),
    link_text VARCHAR(500),
    bounce_type VARCHAR(50),
    bounce_reason TEXT,
    user_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 2. SOCIAL MEDIA SCHEDULER
-- =====================================================

-- Social Media Accounts
CREATE TABLE IF NOT EXISTS social_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    platform VARCHAR(50) NOT NULL CHECK (platform IN ('facebook', 'instagram', 'twitter', 'linkedin', 'tiktok', 'youtube', 'pinterest', 'threads')),
    account_id VARCHAR(255) NOT NULL,
    account_name VARCHAR(255),
    account_username VARCHAR(255),
    account_avatar VARCHAR(500),
    account_type VARCHAR(50) CHECK (account_type IN ('personal', 'business', 'creator')),
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'expired', 'disconnected')),
    auto_publish BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, platform, account_id)
);

-- Social Media Posts
CREATE TABLE IF NOT EXISTS social_posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    title VARCHAR(255),
    content TEXT NOT NULL,
    media_urls TEXT[],
    media_type VARCHAR(50) CHECK (media_type IN ('image', 'video', 'carousel', 'story', 'reel')),
    link_url VARCHAR(500),
    hashtags TEXT[],
    mentions TEXT[],
    platform_specific JSONB,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'scheduled', 'publishing', 'published', 'failed', 'deleted')),
    is_recurring BOOLEAN DEFAULT false,
    recurring_pattern VARCHAR(50) CHECK (recurring_pattern IN ('daily', 'weekly', 'monthly', 'custom')),
    recurring_end_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP WITH TIME ZONE
);

-- Post Schedules
CREATE TABLE IF NOT EXISTS social_post_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID REFERENCES social_posts(id) ON DELETE CASCADE,
    account_id UUID REFERENCES social_accounts(id) ON DELETE CASCADE,
    scheduled_time TIMESTAMP WITH TIME ZONE NOT NULL,
    platform_post_id VARCHAR(255),
    platform_permalink VARCHAR(500),
    status VARCHAR(50) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'publishing', 'published', 'failed')),
    error_message TEXT,
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Social Media Analytics
CREATE TABLE IF NOT EXISTS social_analytics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_schedule_id UUID REFERENCES social_post_schedules(id) ON DELETE CASCADE,
    account_id UUID REFERENCES social_accounts(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL,
    impressions BIGINT DEFAULT 0,
    reach BIGINT DEFAULT 0,
    likes INTEGER DEFAULT 0,
    comments INTEGER DEFAULT 0,
    shares INTEGER DEFAULT 0,
    saves INTEGER DEFAULT 0,
    clicks INTEGER DEFAULT 0,
    video_views INTEGER DEFAULT 0,
    video_view_rate DECIMAL(5,2) DEFAULT 0,
    engagement_rate DECIMAL(5,2) DEFAULT 0,
    engagement_count INTEGER DEFAULT 0,
    analytics_data JSONB,
    fetched_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Content Calendar Events
CREATE TABLE IF NOT EXISTS content_calendar_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_type VARCHAR(50) DEFAULT 'post' CHECK (event_type IN ('post', 'campaign', 'holiday', 'reminder', 'meeting')),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    all_day BOOLEAN DEFAULT false,
    color VARCHAR(20),
    tags TEXT[],
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 3. LANDING PAGE BUILDER
-- =====================================================

-- Landing Pages
CREATE TABLE IF NOT EXISTS landing_pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    title VARCHAR(500),
    meta_description VARCHAR(1000),
    favicon_url VARCHAR(500),
    custom_domain VARCHAR(255),
    html_content TEXT,
    design_json JSONB,
    css_styles TEXT,
    head_scripts TEXT,
    body_scripts TEXT,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'unpublished', 'archived')),
    is_template BOOLEAN DEFAULT false,
    template_category VARCHAR(100),
    thumbnail_url VARCHAR(500),
    view_count BIGINT DEFAULT 0,
    unique_visitors BIGINT DEFAULT 0,
    conversion_count INTEGER DEFAULT 0,
    conversion_rate DECIMAL(5,2) DEFAULT 0,
    published_at TIMESTAMP WITH TIME ZONE,
    seo_indexable BOOLEAN DEFAULT true,
    password_protected BOOLEAN DEFAULT false,
    password VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, slug)
);

-- Landing Page Components
CREATE TABLE IF NOT EXISTS landing_page_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id UUID REFERENCES landing_pages(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES landing_page_components(id) ON DELETE CASCADE,
    component_type VARCHAR(100) NOT NULL,
    name VARCHAR(255),
    props JSONB,
    styles JSONB,
    content TEXT,
    order_index INTEGER DEFAULT 0,
    is_visible BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Landing Page Forms
CREATE TABLE IF NOT EXISTS landing_page_forms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    page_id UUID REFERENCES landing_pages(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    form_type VARCHAR(50) DEFAULT 'lead' CHECK (form_type IN ('lead', 'contact', 'signup', 'newsletter', 'custom')),
    fields JSONB,
    submit_button_text VARCHAR(100) DEFAULT 'Submit',
    success_message TEXT,
    redirect_url VARCHAR(500),
    email_notification BOOLEAN DEFAULT true,
    notification_emails TEXT[],
    webhook_url VARCHAR(500),
    double_optin BOOLEAN DEFAULT false,
    list_id UUID REFERENCES email_lists(id),
    conversion_goal VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Landing Page Leads
CREATE TABLE IF NOT EXISTS landing_page_leads (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    page_id UUID REFERENCES landing_pages(id) ON DELETE SET NULL,
    form_id UUID REFERENCES landing_page_forms(id) ON DELETE SET NULL,
    email VARCHAR(255),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(50),
    company VARCHAR(255),
    job_title VARCHAR(255),
    custom_data JSONB,
    utm_source VARCHAR(100),
    utm_medium VARCHAR(100),
    utm_campaign VARCHAR(255),
    utm_term VARCHAR(255),
    utm_content VARCHAR(255),
    ip_address VARCHAR(45),
    user_agent TEXT,
    referrer_url VARCHAR(500),
    status VARCHAR(50) DEFAULT 'new' CHECK (status IN ('new', 'contacted', 'qualified', 'converted', 'lost')),
    value DECIMAL(15,2),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Landing Page A/B Tests
CREATE TABLE IF NOT EXISTS landing_page_ab_tests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    page_id_a UUID REFERENCES landing_pages(id),
    page_id_b UUID REFERENCES landing_pages(id),
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'running', 'paused', 'completed')),
    traffic_split INTEGER DEFAULT 50,
    visits_a INTEGER DEFAULT 0,
    conversions_a INTEGER DEFAULT 0,
    visits_b INTEGER DEFAULT 0,
    conversions_b INTEGER DEFAULT 0,
    conversion_rate_a DECIMAL(5,2) DEFAULT 0,
    conversion_rate_b DECIMAL(5,2) DEFAULT 0,
    winner VARCHAR(10),
    confidence DECIMAL(5,2) DEFAULT 0,
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 4. MARKETING AUTOMATION
-- =====================================================

-- Automation Workflows
CREATE TABLE IF NOT EXISTS automation_workflows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    trigger_type VARCHAR(100) NOT NULL CHECK (trigger_type IN ('list_subscription', 'tag_added', 'form_submission', 'page_visit', 'purchase', 'date', 'api', 'manual')),
    trigger_config JSONB,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'paused', 'archived')),
    total_contacts INTEGER DEFAULT 0,
    active_contacts INTEGER DEFAULT 0,
    completed_contacts INTEGER DEFAULT 0,
    goals_met INTEGER DEFAULT 0,
    revenue DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Automation Steps
CREATE TABLE IF NOT EXISTS automation_steps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID REFERENCES automation_workflows(id) ON DELETE CASCADE,
    parent_step_id UUID REFERENCES automation_steps(id) ON DELETE CASCADE,
    step_type VARCHAR(100) NOT NULL CHECK (step_type IN ('email', 'wait', 'condition', 'action', 'goal', 'split')),
    name VARCHAR(255),
    order_index INTEGER DEFAULT 0,
    config JSONB,
    delay_minutes INTEGER DEFAULT 0,
    delay_type VARCHAR(50) CHECK (delay_type IN ('minutes', 'hours', 'days', 'specific_time')),
    conditions JSONB,
    actions JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Automation Contacts (Contact journey through workflows)
CREATE TABLE IF NOT EXISTS automation_contacts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID REFERENCES automation_workflows(id) ON DELETE CASCADE,
    subscriber_id UUID REFERENCES email_subscribers(id) ON DELETE CASCADE,
    current_step_id UUID REFERENCES automation_steps(id) ON DELETE SET NULL,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'paused', 'exited', 'failed')),
    entered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    exited_at TIMESTAMP WITH TIME ZONE,
    goal_met BOOLEAN DEFAULT false,
    goal_met_at TIMESTAMP WITH TIME ZONE,
    revenue DECIMAL(15,2) DEFAULT 0,
    step_history JSONB,
    custom_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Automation Events
CREATE TABLE IF NOT EXISTS automation_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    workflow_id UUID REFERENCES automation_workflows(id) ON DELETE CASCADE,
    automation_contact_id UUID REFERENCES automation_contacts(id) ON DELETE CASCADE,
    subscriber_id UUID REFERENCES email_subscribers(id) ON DELETE CASCADE,
    step_id UUID REFERENCES automation_steps(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL,
    event_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Automation Goals
CREATE TABLE IF NOT EXISTS automation_goals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID REFERENCES automation_workflows(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    goal_type VARCHAR(100) NOT NULL CHECK (goal_type IN ('purchase', 'page_visit', 'link_click', 'form_submit', 'tag_added', 'custom')),
    goal_config JSONB,
    conversion_value DECIMAL(15,2),
    conversions INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Contact Tags
CREATE TABLE IF NOT EXISTS contact_tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(20),
    description TEXT,
    contact_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, name)
);

-- Contact-Tag Relationship
CREATE TABLE IF NOT EXISTS contact_tag_relations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    contact_id UUID REFERENCES email_subscribers(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES contact_tags(id) ON DELETE CASCADE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(contact_id, tag_id)
);

-- =====================================================
-- 5. INFLUENCER MARKETING
-- =====================================================

-- Influencers
CREATE TABLE IF NOT EXISTS influencers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    bio TEXT,
    profile_image VARCHAR(500),
    category VARCHAR(100),
    niches TEXT[],
    location VARCHAR(255),
    country VARCHAR(100),
    language VARCHAR(50),
    gender VARCHAR(20),
    age_range VARCHAR(50),
    status VARCHAR(50) DEFAULT 'prospect' CHECK (status IN ('prospect', 'outreached', 'negotiating', 'active', 'completed', 'declined', 'blacklisted')),
    rating DECIMAL(3,2) DEFAULT 0,
    notes TEXT,
    custom_fields JSONB,
    total_collaborations INTEGER DEFAULT 0,
    total_reach BIGINT DEFAULT 0,
    total_engagement BIGINT DEFAULT 0,
    total_earnings DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Influencer Social Profiles
CREATE TABLE IF NOT EXISTS influencer_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    influencer_id UUID REFERENCES influencers(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL,
    username VARCHAR(255),
    profile_url VARCHAR(500),
    profile_image VARCHAR(500),
    bio TEXT,
    followers BIGINT DEFAULT 0,
    following BIGINT DEFAULT 0,
    posts INTEGER DEFAULT 0,
    avg_likes INTEGER DEFAULT 0,
    avg_comments INTEGER DEFAULT 0,
    avg_shares INTEGER DEFAULT 0,
    engagement_rate DECIMAL(5,2) DEFAULT 0,
    verified BOOLEAN DEFAULT false,
    audience_gender_male DECIMAL(5,2) DEFAULT 0,
    audience_gender_female DECIMAL(5,2) DEFAULT 0,
    audience_age_distribution JSONB,
    audience_locations JSONB,
    last_synced_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(influencer_id, platform)
);

-- Influencer Campaigns
CREATE TABLE IF NOT EXISTS influencer_campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    campaign_type VARCHAR(50) CHECK (campaign_type IN ('brand_awareness', 'product_launch', 'content_creation', 'event', 'giveaway', 'affiliate')),
    budget DECIMAL(15,2) DEFAULT 0,
    spent DECIMAL(15,2) DEFAULT 0,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'paused', 'completed', 'cancelled')),
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    goals JSONB,
    requirements TEXT,
    deliverables TEXT,
    hashtags TEXT[],
    mentions TEXT[],
    tracking_links TEXT[],
    total_reach BIGINT DEFAULT 0,
    total_impressions BIGINT DEFAULT 0,
    total_engagement BIGINT DEFAULT 0,
    total_clicks BIGINT DEFAULT 0,
    total_conversions INTEGER DEFAULT 0,
    roi DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Influencer Collaborations
CREATE TABLE IF NOT EXISTS influencer_collaborations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID REFERENCES influencer_campaigns(id) ON DELETE CASCADE,
    influencer_id UUID REFERENCES influencers(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'outreached' CHECK (status IN ('outreached', 'negotiating', 'accepted', 'in_progress', 'content_review', 'published', 'completed', 'declined', 'cancelled')),
    negotiation_amount DECIMAL(15,2),
    agreed_amount DECIMAL(15,2),
    currency VARCHAR(10) DEFAULT 'USD',
    payment_status VARCHAR(50) DEFAULT 'pending' CHECK (payment_status IN ('pending', 'processing', 'paid', 'failed', 'refunded')),
    paid_amount DECIMAL(15,2) DEFAULT 0,
    deliverables JSONB,
    content_deadline TIMESTAMP WITH TIME ZONE,
    publish_date TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    outreach_date TIMESTAMP WITH TIME ZONE,
    response_date TIMESTAMP WITH TIME ZONE,
    accepted_date TIMESTAMP WITH TIME ZONE,
    completed_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Influencer Content
CREATE TABLE IF NOT EXISTS influencer_content (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    collaboration_id UUID REFERENCES influencer_collaborations(id) ON DELETE CASCADE,
    content_type VARCHAR(50) CHECK (content_type IN ('post', 'story', 'reel', 'video', 'blog', 'review')),
    platform VARCHAR(50),
    content_url VARCHAR(500),
    media_urls TEXT[],
    caption TEXT,
    hashtags TEXT[],
    mentions TEXT[],
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'submitted', 'approved', 'revision', 'published', 'rejected')),
    review_notes TEXT,
    impressions BIGINT DEFAULT 0,
    reach BIGINT DEFAULT 0,
    likes INTEGER DEFAULT 0,
    comments INTEGER DEFAULT 0,
    shares INTEGER DEFAULT 0,
    saves INTEGER DEFAULT 0,
    clicks INTEGER DEFAULT 0,
    engagement_rate DECIMAL(5,2) DEFAULT 0,
    published_at TIMESTAMP WITH TIME ZONE,
    analytics_fetched_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Influencer Payments
CREATE TABLE IF NOT EXISTS influencer_payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    collaboration_id UUID REFERENCES influencer_collaborations(id) ON DELETE CASCADE,
    influencer_id UUID REFERENCES influencers(id) ON DELETE CASCADE,
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    payment_method VARCHAR(50) CHECK (payment_method IN ('bank_transfer', 'paypal', 'wise', 'check', 'other')),
    payment_reference VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'refunded')),
    notes TEXT,
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- 6. REPORTING & ANALYTICS
-- =====================================================

-- Custom Reports
CREATE TABLE IF NOT EXISTS custom_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    report_type VARCHAR(50) CHECK (report_type IN ('marketing', 'sales', 'social', 'email', 'web', 'custom')),
    data_sources JSONB,
    metrics JSONB,
    dimensions JSONB,
    filters JSONB,
    date_range_type VARCHAR(50) DEFAULT 'last_30_days',
    custom_date_start DATE,
    custom_date_end DATE,
    comparison_enabled BOOLEAN DEFAULT false,
    comparison_type VARCHAR(50),
    chart_config JSONB,
    table_config JSONB,
    is_scheduled BOOLEAN DEFAULT false,
    schedule_frequency VARCHAR(50) CHECK (schedule_frequency IN ('daily', 'weekly', 'monthly', 'quarterly')),
    schedule_recipients TEXT[],
    last_run_at TIMESTAMP WITH TIME ZONE,
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Report Widgets
CREATE TABLE IF NOT EXISTS report_widgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID REFERENCES custom_reports(id) ON DELETE CASCADE,
    widget_type VARCHAR(100) NOT NULL CHECK (widget_type IN ('kpi', 'chart', 'table', 'gauge', 'funnel', 'map', 'list')),
    title VARCHAR(255),
    data_source VARCHAR(100),
    metric VARCHAR(100),
    dimension VARCHAR(100),
    filters JSONB,
    chart_type VARCHAR(50),
    chart_config JSONB,
    order_index INTEGER DEFAULT 0,
    width INTEGER DEFAULT 6,
    height INTEGER DEFAULT 4,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Report Snapshots (for historical data)
CREATE TABLE IF NOT EXISTS report_snapshots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID REFERENCES custom_reports(id) ON DELETE CASCADE,
    snapshot_date DATE NOT NULL,
    data JSONB,
    summary JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Custom Dashboards
CREATE TABLE IF NOT EXISTS custom_dashboards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    layout JSONB,
    is_default BOOLEAN DEFAULT false,
    is_public BOOLEAN DEFAULT false,
    shared_with JSONB,
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Dashboard Widgets
CREATE TABLE IF NOT EXISTS dashboard_widgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dashboard_id UUID REFERENCES custom_dashboards(id) ON DELETE CASCADE,
    widget_type VARCHAR(100) NOT NULL,
    title VARCHAR(255),
    config JSONB,
    data_source VARCHAR(100),
    refresh_interval INTEGER DEFAULT 300,
    position_x INTEGER DEFAULT 0,
    position_y INTEGER DEFAULT 0,
    width INTEGER DEFAULT 4,
    height INTEGER DEFAULT 3,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Analytics Events (unified tracking)
CREATE TABLE IF NOT EXISTS analytics_events_unified (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_category VARCHAR(100),
    event_action VARCHAR(100),
    event_label VARCHAR(255),
    event_value DECIMAL(15,2),
    user_id UUID,
    session_id VARCHAR(255),
    device_type VARCHAR(50),
    browser VARCHAR(100),
    os VARCHAR(100),
    country VARCHAR(100),
    city VARCHAR(100),
    ip_address VARCHAR(45),
    page_url VARCHAR(500),
    referrer_url VARCHAR(500),
    utm_source VARCHAR(100),
    utm_medium VARCHAR(100),
    utm_campaign VARCHAR(255),
    utm_term VARCHAR(255),
    utm_content VARCHAR(255),
    custom_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Marketing Attribution
CREATE TABLE IF NOT EXISTS marketing_attribution (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    conversion_id VARCHAR(255),
    conversion_type VARCHAR(100),
    conversion_value DECIMAL(15,2),
    conversion_date TIMESTAMP WITH TIME ZONE,
    touchpoints JSONB,
    attribution_model VARCHAR(50) DEFAULT 'last_touch' CHECK (attribution_model IN ('first_touch', 'last_touch', 'linear', 'time_decay', 'position_based', 'custom')),
    attributed_channels JSONB,
    attributed_campaigns JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Report Templates
CREATE TABLE IF NOT EXISTS report_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    description TEXT,
    preview_image VARCHAR(500),
    config JSONB,
    is_public BOOLEAN DEFAULT false,
    usage_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- INDEXES
-- =====================================================

-- Email Marketing Indexes
CREATE INDEX idx_email_lists_org ON email_lists(organization_id);
CREATE INDEX idx_email_subscribers_org ON email_subscribers(organization_id);
CREATE INDEX idx_email_subscribers_email ON email_subscribers(email);
CREATE INDEX idx_email_subscribers_status ON email_subscribers(status);
CREATE INDEX idx_email_campaigns_org ON email_campaigns(organization_id);
CREATE INDEX idx_email_campaigns_status ON email_campaigns(status);
CREATE INDEX idx_email_events_campaign ON email_events(campaign_id);
CREATE INDEX idx_email_events_subscriber ON email_events(subscriber_id);
CREATE INDEX idx_email_events_type ON email_events(event_type);
CREATE INDEX idx_email_events_created ON email_events(created_at);

-- Social Media Indexes
CREATE INDEX idx_social_accounts_org ON social_accounts(organization_id);
CREATE INDEX idx_social_accounts_platform ON social_accounts(platform);
CREATE INDEX idx_social_posts_org ON social_posts(organization_id);
CREATE INDEX idx_social_posts_status ON social_posts(status);
CREATE INDEX idx_social_post_schedules_time ON social_post_schedules(scheduled_time);
CREATE INDEX idx_social_analytics_account ON social_analytics(account_id);

-- Landing Page Indexes
CREATE INDEX idx_landing_pages_org ON landing_pages(organization_id);
CREATE INDEX idx_landing_pages_slug ON landing_pages(slug);
CREATE INDEX idx_landing_pages_status ON landing_pages(status);
CREATE INDEX idx_landing_page_leads_org ON landing_page_leads(organization_id);
CREATE INDEX idx_landing_page_leads_email ON landing_page_leads(email);

-- Automation Indexes
CREATE INDEX idx_automation_workflows_org ON automation_workflows(organization_id);
CREATE INDEX idx_automation_workflows_status ON automation_workflows(status);
CREATE INDEX idx_automation_contacts_workflow ON automation_contacts(workflow_id);
CREATE INDEX idx_automation_contacts_status ON automation_contacts(status);
CREATE INDEX idx_contact_tags_org ON contact_tags(organization_id);

-- Influencer Indexes
CREATE INDEX idx_influencers_org ON influencers(organization_id);
CREATE INDEX idx_influencers_status ON influencers(status);
CREATE INDEX idx_influencer_campaigns_org ON influencer_campaigns(organization_id);
CREATE INDEX idx_influencer_collaborations_campaign ON influencer_collaborations(campaign_id);
CREATE INDEX idx_influencer_content_collab ON influencer_content(collaboration_id);

-- Reporting Indexes
CREATE INDEX idx_custom_reports_org ON custom_reports(organization_id);
CREATE INDEX idx_report_snapshots_report ON report_snapshots(report_id);
CREATE INDEX idx_custom_dashboards_org ON custom_dashboards(organization_id);
CREATE INDEX idx_analytics_events_unified_org ON analytics_events_unified(organization_id);
CREATE INDEX idx_analytics_events_unified_type ON analytics_events_unified(event_type);
CREATE INDEX idx_analytics_events_unified_created ON analytics_events_unified(created_at);
CREATE INDEX idx_marketing_attribution_org ON marketing_attribution(organization_id);

-- =====================================================
-- TRIGGERS
-- =====================================================

CREATE TRIGGER update_email_lists_updated_at BEFORE UPDATE ON email_lists FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_email_subscribers_updated_at BEFORE UPDATE ON email_subscribers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_email_templates_updated_at BEFORE UPDATE ON email_templates FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_email_campaigns_updated_at BEFORE UPDATE ON email_campaigns FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_social_accounts_updated_at BEFORE UPDATE ON social_accounts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_social_posts_updated_at BEFORE UPDATE ON social_posts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_landing_pages_updated_at BEFORE UPDATE ON landing_pages FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_automation_workflows_updated_at BEFORE UPDATE ON automation_workflows FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_influencers_updated_at BEFORE UPDATE ON influencers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_custom_reports_updated_at BEFORE UPDATE ON custom_reports FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();