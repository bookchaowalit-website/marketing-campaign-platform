-- Add social media tables for publishing

-- Social accounts (connected Facebook/Instagram/etc.)
CREATE TABLE IF NOT EXISTS social_publish_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    platform VARCHAR(50) NOT NULL CHECK (platform IN ('facebook', 'instagram', 'twitter', 'linkedin')),
    platform_id VARCHAR(255) NOT NULL,
    account_name VARCHAR(255),
    page_id VARCHAR(255),
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    token_expiry TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) DEFAULT 'connected' CHECK (status IN ('connected', 'disconnected', 'error')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, platform, platform_id)
);

-- Social posts
CREATE TABLE IF NOT EXISTS social_publish_posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    account_id UUID REFERENCES social_publish_accounts(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    media_urls TEXT[],
    hashtags TEXT[],
    platform VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'scheduled', 'published', 'failed')),
    scheduled_at TIMESTAMP WITH TIME ZONE,
    published_at TIMESTAMP WITH TIME ZONE,
    platform_post_id VARCHAR(255),
    error_msg TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Social post analytics
CREATE TABLE IF NOT EXISTS social_publish_post_analytics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID REFERENCES social_publish_posts(id) ON DELETE CASCADE UNIQUE,
    platform VARCHAR(50) NOT NULL,
    impressions INTEGER DEFAULT 0,
    reach INTEGER DEFAULT 0,
    likes INTEGER DEFAULT 0,
    comments INTEGER DEFAULT 0,
    shares INTEGER DEFAULT 0,
    saves INTEGER DEFAULT 0,
    clicks INTEGER DEFAULT 0,
    engagement_rate DECIMAL(5,2) DEFAULT 0,
    fetched_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_social_publish_accounts_org_id ON social_publish_accounts(organization_id);
CREATE INDEX idx_social_publish_accounts_platform ON social_publish_accounts(platform);
CREATE INDEX idx_social_publish_posts_org_id ON social_publish_posts(organization_id);
CREATE INDEX idx_social_publish_posts_account_id ON social_publish_posts(account_id);
CREATE INDEX idx_social_publish_posts_status ON social_publish_posts(status);
CREATE INDEX idx_social_publish_posts_scheduled_at ON social_publish_posts(scheduled_at);
CREATE INDEX idx_social_publish_analytics_post_id ON social_publish_post_analytics(post_id);

-- Trigger for updated_at
CREATE TRIGGER update_social_publish_accounts_updated_at BEFORE UPDATE ON social_publish_accounts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_social_publish_posts_updated_at BEFORE UPDATE ON social_publish_posts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
