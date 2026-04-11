-- External Affiliate Integrations Migration - UP
-- Creates tables for external affiliate network integrations (Shopee, TikTok, etc.)

-- External Affiliate Accounts - Stores API credentials for each platform
CREATE TABLE IF NOT EXISTS external_affiliate_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL,
    platform VARCHAR(50) NOT NULL CHECK (platform IN ('shopee', 'tiktok', 'lazada', 'amazon', 'custom')),
    platform_account_id VARCHAR(255) NOT NULL,
    platform_account_name VARCHAR(255),
    platform_username VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'disconnected', 'expired', 'error')),
    
    -- API Credentials (encrypted in production)
    api_key TEXT,
    api_secret TEXT,
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP WITH TIME ZONE,
    
    -- Sync settings
    auto_sync BOOLEAN DEFAULT true,
    sync_interval_minutes INTEGER DEFAULT 60,
    last_sync_at TIMESTAMP WITH TIME ZONE,
    last_sync_status VARCHAR(50) DEFAULT 'pending' CHECK (last_sync_status IN ('pending', 'success', 'partial', 'failed')),
    last_sync_error TEXT,
    
    -- Stats (cached from last sync)
    total_earnings DECIMAL(15,2) DEFAULT 0,
    total_orders INTEGER DEFAULT 0,
    total_clicks BIGINT DEFAULT 0,
    total_products INTEGER DEFAULT 0,
    pending_commission DECIMAL(15,2) DEFAULT 0,
    
    -- Metadata
    connection_data JSONB, -- Additional platform-specific data
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, platform, platform_account_id)
);

-- External Affiliate Products - Products available for promotion
CREATE TABLE IF NOT EXISTS external_affiliate_products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID REFERENCES external_affiliate_accounts(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL,
    product_id VARCHAR(255) NOT NULL,
    product_name VARCHAR(500) NOT NULL,
    product_url VARCHAR(1000),
    product_image_url VARCHAR(1000),
    category VARCHAR(255),
    brand VARCHAR(255),
    
    -- Pricing
    original_price DECIMAL(15,2),
    sale_price DECIMAL(15,2),
    currency VARCHAR(10) DEFAULT 'USD',
    
    -- Commission info
    commission_rate DECIMAL(5,4) DEFAULT 0,
    commission_type VARCHAR(50) DEFAULT 'percentage' CHECK (commission_type IN ('percentage', 'fixed')),
    estimated_earnings DECIMAL(15,2) DEFAULT 0,
    
    -- Stats
    total_clicks BIGINT DEFAULT 0,
    total_orders INTEGER DEFAULT 0,
    total_earnings DECIMAL(15,2) DEFAULT 0,
    
    -- Product status
    product_status VARCHAR(50) DEFAULT 'active' CHECK (product_status IN ('active', 'out_of_stock', 'discontinued', 'unknown')),
    
    -- Tracking
    affiliate_link VARCHAR(1000),
    short_link VARCHAR(100),
    
    -- Metadata
    product_data JSONB, -- Full product data from platform
    last_synced_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(account_id, product_id)
);

-- External Affiliate Orders - Orders tracked through affiliate links
CREATE TABLE IF NOT EXISTS external_affiliate_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID REFERENCES external_affiliate_accounts(id) ON DELETE CASCADE,
    product_id UUID REFERENCES external_affiliate_products(id) ON DELETE SET NULL,
    
    -- Order identification
    platform VARCHAR(50) NOT NULL,
    order_id VARCHAR(255) NOT NULL,
    parent_order_id VARCHAR(255), -- For grouped orders
    
    -- Order details
    order_status VARCHAR(50) NOT NULL CHECK (order_status IN ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled', 'returned', 'refunded')),
    order_date TIMESTAMP WITH TIME ZONE,
    order_items INTEGER DEFAULT 1,
    
    -- Amounts
    order_amount DECIMAL(15,2) NOT NULL,
    product_amount DECIMAL(15,2),
    shipping_amount DECIMAL(15,2),
    discount_amount DECIMAL(15,2),
    currency VARCHAR(10) DEFAULT 'USD',
    
    -- Commission
    commission_rate DECIMAL(5,4),
    commission_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    commission_status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (commission_status IN ('pending', 'confirmed', 'paid', 'cancelled', 'adjusted')),
    commission_paid_at TIMESTAMP WITH TIME ZONE,
    
    -- Customer info (anonymized)
    customer_id VARCHAR(255),
    
    -- Attribution
    traffic_source VARCHAR(100),
    utm_source VARCHAR(100),
    utm_medium VARCHAR(100),
    utm_campaign VARCHAR(255),
    
    -- Metadata
    order_data JSONB, -- Full order data from platform
    synced_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(account_id, order_id)
);

-- External Affiliate Earnings - Daily/period earnings summary
CREATE TABLE IF NOT EXISTS external_affiliate_earnings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID REFERENCES external_affiliate_accounts(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL,
    date DATE NOT NULL,
    
    -- Earnings breakdown
    gross_sales DECIMAL(15,2) DEFAULT 0,
    net_sales DECIMAL(15,2) DEFAULT 0,
    orders_count INTEGER DEFAULT 0,
    items_sold INTEGER DEFAULT 0,
    
    -- Commission
    commission_earned DECIMAL(15,2) DEFAULT 0,
    commission_pending DECIMAL(15,2) DEFAULT 0,
    commission_confirmed DECIMAL(15,2) DEFAULT 0,
    
    -- Traffic
    clicks BIGINT DEFAULT 0,
    impressions BIGINT DEFAULT 0,
    
    -- Performance
    conversion_rate DECIMAL(5,2) DEFAULT 0,
    epc DECIMAL(15,4) DEFAULT 0, -- Earnings Per Click
    
    currency VARCHAR(10) DEFAULT 'USD',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(account_id, date)
);

-- External Affiliate Sync Logs - Track sync operations
CREATE TABLE IF NOT EXISTS external_affiliate_sync_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID REFERENCES external_affiliate_accounts(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL,
    
    sync_type VARCHAR(50) NOT NULL CHECK (sync_type IN ('full', 'incremental', 'products', 'orders', 'earnings')),
    status VARCHAR(50) NOT NULL CHECK (status IN ('started', 'success', 'partial', 'failed')),
    
    -- Timing
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_ms BIGINT,
    
    -- Results
    records_processed INTEGER DEFAULT 0,
    records_created INTEGER DEFAULT 0,
    records_updated INTEGER DEFAULT 0,
    records_failed INTEGER DEFAULT 0,
    
    -- Error info
    error_message TEXT,
    error_details JSONB,
    
    -- API info
    api_calls_count INTEGER DEFAULT 0,
    rate_limit_hit BOOLEAN DEFAULT false,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- External Affiliate Links - Custom tracking links created by user
CREATE TABLE IF NOT EXISTS external_affiliate_links (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID REFERENCES external_affiliate_accounts(id) ON DELETE CASCADE,
    product_id UUID REFERENCES external_affiliate_products(id) ON DELETE SET NULL,
    
    link_code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Links
    original_url VARCHAR(1000) NOT NULL,
    affiliate_url VARCHAR(1000) NOT NULL,
    short_url VARCHAR(100),
    
    -- Tracking parameters
    utm_source VARCHAR(100),
    utm_medium VARCHAR(100),
    utm_campaign VARCHAR(255),
    utm_term VARCHAR(255),
    utm_content VARCHAR(255),
    sub_id VARCHAR(255),
    
    -- Stats
    total_clicks BIGINT DEFAULT 0,
    unique_clicks BIGINT DEFAULT 0,
    total_orders INTEGER DEFAULT 0,
    total_earnings DECIMAL(15,2) DEFAULT 0,
    
    -- Settings
    is_active BOOLEAN DEFAULT true,
    expires_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_ext_aff_accounts_user ON external_affiliate_accounts(user_id);
CREATE INDEX idx_ext_aff_accounts_org ON external_affiliate_accounts(organization_id);
CREATE INDEX idx_ext_aff_accounts_platform ON external_affiliate_accounts(platform);
CREATE INDEX idx_ext_aff_accounts_status ON external_affiliate_accounts(status);

CREATE INDEX idx_ext_aff_products_account ON external_affiliate_products(account_id);
CREATE INDEX idx_ext_aff_products_platform ON external_affiliate_products(platform);
CREATE INDEX idx_ext_aff_products_status ON external_affiliate_products(product_status);

CREATE INDEX idx_ext_aff_orders_account ON external_affiliate_orders(account_id);
CREATE INDEX idx_ext_aff_orders_platform ON external_affiliate_orders(platform);
CREATE INDEX idx_ext_aff_orders_status ON external_affiliate_orders(order_status);
CREATE INDEX idx_ext_aff_orders_date ON external_affiliate_orders(order_date);
CREATE INDEX idx_ext_aff_orders_commission_status ON external_affiliate_orders(commission_status);

CREATE INDEX idx_ext_aff_earnings_account ON external_affiliate_earnings(account_id);
CREATE INDEX idx_ext_aff_earnings_date ON external_affiliate_earnings(date);
CREATE INDEX idx_ext_aff_earnings_platform ON external_affiliate_earnings(platform);

CREATE INDEX idx_ext_aff_sync_logs_account ON external_affiliate_sync_logs(account_id);
CREATE INDEX idx_ext_aff_sync_logs_status ON external_affiliate_sync_logs(status);
CREATE INDEX idx_ext_aff_sync_logs_started ON external_affiliate_sync_logs(started_at);

CREATE INDEX idx_ext_aff_links_code ON external_affiliate_links(link_code);
CREATE INDEX idx_ext_aff_links_account ON external_affiliate_links(account_id);

-- Triggers for updated_at
CREATE TRIGGER update_ext_aff_accounts_updated_at BEFORE UPDATE ON external_affiliate_accounts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ext_aff_products_updated_at BEFORE UPDATE ON external_affiliate_products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ext_aff_orders_updated_at BEFORE UPDATE ON external_affiliate_orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ext_aff_earnings_updated_at BEFORE UPDATE ON external_affiliate_earnings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ext_aff_links_updated_at BEFORE UPDATE ON external_affiliate_links FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();