-- Phase 2: PPC/Paid Ads Management

-- ===========================
-- AD PLATFORMS
-- ===========================

CREATE TABLE IF NOT EXISTS ad_platforms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL,
    platform VARCHAR(50) NOT NULL CHECK (platform IN ('google_ads', 'facebook_ads', 'tiktok_ads', 'linkedin_ads', 'twitter_ads')),
    account_id VARCHAR(255) NOT NULL,
    api_credentials TEXT NOT NULL, -- Encrypted JSON
    status VARCHAR(50) DEFAULT 'connected' CHECK (status IN ('connected', 'disconnected', 'error', 'pending')),
    connected_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_sync_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(org_id, platform, account_id)
);

-- ===========================
-- AD CAMPAIGNS
-- ===========================

CREATE TABLE IF NOT EXISTS ad_campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL,
    unified_campaign_id UUID REFERENCES unified_campaignes(id) ON DELETE SET NULL,
    platform_id UUID REFERENCES ad_platforms(id) ON DELETE CASCADE,
    platform_campaign_id VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'paused', 'ended', 'archived')),
    budget DECIMAL(15,2) DEFAULT 0,
    spent DECIMAL(15,2) DEFAULT 0,
    impressions BIGINT DEFAULT 0,
    clicks BIGINT DEFAULT 0,
    conversions BIGINT DEFAULT 0,
    cpc DECIMAL(10,4) DEFAULT 0,
    cpm DECIMAL(10,4) DEFAULT 0,
    ctr DECIMAL(5,4) DEFAULT 0,
    start_date DATE,
    end_date DATE,
    targeting JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ===========================
-- AD GROUPS
-- ===========================

CREATE TABLE IF NOT EXISTS ad_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES ad_campaigns(id) ON DELETE CASCADE,
    platform_ad_group_id VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'paused', 'removed')),
    budget DECIMAL(15,2) DEFAULT 0,
    spent DECIMAL(15,2) DEFAULT 0,
    impressions BIGINT DEFAULT 0,
    clicks BIGINT DEFAULT 0,
    conversions BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ===========================
-- AD CREATIVES
-- ===========================

CREATE TABLE IF NOT EXISTS ad_creatives (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ad_group_id UUID NOT NULL REFERENCES ad_groups(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('text', 'image', 'video', 'carousel')),
    headline VARCHAR(255),
    description TEXT,
    image_url TEXT,
    video_url TEXT,
    landing_page_url TEXT,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'paused', 'disapproved', 'pending_review')),
    approval_status VARCHAR(50) DEFAULT 'pending' CHECK (approval_status IN ('pending', 'approved', 'disapproved', 'under_review')),
    platform_creative_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ===========================
-- AD KEYWORDS
-- ===========================

CREATE TABLE IF NOT EXISTS ad_keywords (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ad_group_id UUID NOT NULL REFERENCES ad_groups(id) ON DELETE CASCADE,
    keyword VARCHAR(500) NOT NULL,
    match_type VARCHAR(50) NOT NULL CHECK (match_type IN ('exact', 'phrase', 'broad', 'negative')),
    bid_amount DECIMAL(10,4) DEFAULT 0,
    max_cpc DECIMAL(10,4),
    impressions BIGINT DEFAULT 0,
    clicks BIGINT DEFAULT 0,
    conversions BIGINT DEFAULT 0,
    cost DECIMAL(15,4) DEFAULT 0,
    avg_position DECIMAL(5,2),
    quality_score DECIMAL(3,1),
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'paused', 'removed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ===========================
-- AD PERFORMANCE DAILY
-- ===========================

CREATE TABLE IF NOT EXISTS ad_performance_daily (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ad_id UUID NOT NULL REFERENCES ad_creatives(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    impressions BIGINT DEFAULT 0,
    clicks BIGINT DEFAULT 0,
    conversions BIGINT DEFAULT 0,
    cost DECIMAL(15,4) DEFAULT 0,
    ctr DECIMAL(5,4) DEFAULT 0,
    cpc DECIMAL(10,4) DEFAULT 0,
    cpa DECIMAL(10,4) DEFAULT 0,
    roas DECIMAL(10,4) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(ad_id, date)
);

-- ===========================
-- INDEXES
-- ===========================

CREATE INDEX idx_ad_platforms_org ON ad_platforms(org_id);
CREATE INDEX idx_ad_platforms_platform ON ad_platforms(platform);
CREATE INDEX idx_ad_platforms_status ON ad_platforms(status);

CREATE INDEX idx_ad_campaigns_org ON ad_campaigns(org_id);
CREATE INDEX idx_ad_campaigns_platform ON ad_campaigns(platform_id);
CREATE INDEX idx_ad_campaigns_status ON ad_campaigns(status);
CREATE INDEX idx_ad_campaigns_unified ON ad_campaigns(unified_campaign_id);
CREATE INDEX idx_ad_campaigns_dates ON ad_campaigns(start_date, end_date);

CREATE INDEX idx_ad_groups_campaign ON ad_groups(campaign_id);
CREATE INDEX idx_ad_groups_status ON ad_groups(status);

CREATE INDEX idx_ad_creatives_ad_group ON ad_creatives(ad_group_id);
CREATE INDEX idx_ad_creatives_type ON ad_creatives(type);
CREATE INDEX idx_ad_creatives_status ON ad_creatives(status);

CREATE INDEX idx_ad_keywords_ad_group ON ad_keywords(ad_group_id);
CREATE INDEX idx_ad_keywords_keyword ON ad_keywords(keyword);
CREATE INDEX idx_ad_keywords_match_type ON ad_keywords(match_type);
CREATE INDEX idx_ad_keywords_status ON ad_keywords(status);

CREATE INDEX idx_ad_performance_ad ON ad_performance_daily(ad_id);
CREATE INDEX idx_ad_performance_date ON ad_performance_daily(date);
CREATE INDEX idx_ad_performance_ad_date ON ad_performance_daily(ad_id, date);

-- ===========================
-- TRIGGERS
-- ===========================

CREATE TRIGGER update_ad_platforms_updated_at BEFORE UPDATE ON ad_platforms FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ad_campaigns_updated_at BEFORE UPDATE ON ad_campaigns FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ad_groups_updated_at BEFORE UPDATE ON ad_groups FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ad_creatives_updated_at BEFORE UPDATE ON ad_creatives FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ad_keywords_updated_at BEFORE UPDATE ON ad_keywords FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
