-- Seed data for development

-- Create default organization
INSERT INTO organizations (id, name, slug, plan) VALUES 
    ('00000000-0000-0000-0000-000000000001', 'Demo Organization', 'demo', 'pro');

-- Create demo user (password: demo123)
INSERT INTO users (id, email, password_hash, full_name, role) VALUES 
    ('00000000-0000-0000-0000-000000000001', 
     'demo@example.com', 
     '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzS3MebAJu',
     'Demo User',
     'admin');

-- Link user to organization
INSERT INTO organization_members (organization_id, user_id, role) VALUES 
    ('00000000-0000-0000-0000-000000000001', 
     '00000000-0000-0000-0000-000000000001', 
     'owner');

-- Create demo campaigns
INSERT INTO campaigns (id, organization_id, name, type, status, budget, spent, impressions, clicks, conversions, start_date, end_date) VALUES 
    ('10000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'Summer Sale 2024', 'SEM', 'active', 10000.00, 6500.00, 125000, 4500, 180, '2024-06-01', '2024-08-31'),
    ('10000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'SEO Organic Growth', 'SEO', 'active', 5000.00, 3200.00, 85000, 3200, 95, '2024-05-01', NULL),
    ('10000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', 'Facebook Brand Awareness', 'Social', 'paused', 8000.00, 4100.00, 250000, 8900, 220, '2024-04-15', NULL),
    ('10000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000001', 'Email Retargeting Q2', 'Email', 'completed', 3000.00, 2850.00, 45000, 2100, 310, '2024-04-01', '2024-06-30');

-- Create demo A/B tests
INSERT INTO ab_tests (id, organization_id, name, variant, status, visitors, conversions, conversion_rate, confidence, start_date) VALUES 
    ('20000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'Homepage CTA Button Color', 'A', 'running', 12500, 875, 0.07, 0.85, '2024-06-01'),
    ('20000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'Checkout Page Layout', 'B', 'running', 8900, 534, 0.06, 0.72, '2024-06-05'),
    ('20000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', 'Email Subject Line Test', 'A', 'completed', 25000, 2250, 0.09, 0.98, '2024-05-15');

-- Create demo competitors
INSERT INTO competitors (id, organization_id, competitor_name, domain, type, market_share, seo_score, sem_spend, keywords, backlinks, social_followers) VALUES 
    ('30000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'Competitor A', 'competitor-a.com', 'direct', 28.5, 85, 45000.00, 1250, 8500, 125000),
    ('30000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'Competitor B', 'competitor-b.com', 'direct', 22.3, 78, 38000.00, 980, 6200, 98000),
    ('30000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', 'Competitor C', 'competitor-c.com', 'indirect', 15.8, 72, 25000.00, 750, 4100, 67000);

-- Create demo SEO keywords
INSERT INTO seo_keywords (organization_id, keyword, domain, current_position, previous_position, search_volume, difficulty, cpc, url) VALUES 
    ('00000000-0000-0000-0000-000000000001', 'marketing tools', 'yoursite.com', 3, 5, 12500, 65, 8.50, 'https://yoursite.com/marketing-tools'),
    ('00000000-0000-0000-0000-000000000001', 'seo software', 'yoursite.com', 5, 4, 8900, 72, 12.30, 'https://yoursite.com/seo'),
    ('00000000-0000-0000-0000-000000000001', 'marketing automation', 'yoursite.com', 8, 10, 6700, 58, 15.00, 'https://yoursite.com/automation'),
    ('00000000-0000-0000-0000-000000000001', 'digital marketing platform', 'yoursite.com', 12, 15, 5400, 70, 10.20, 'https://yoursite.com/platform'),
    ('00000000-0000-0000-0000-000000000001', 'competitor analysis tool', 'yoursite.com', 7, 8, 4200, 55, 9.80, 'https://yoursite.com/competitor-analysis');

-- Create demo analytics data
INSERT INTO analytics_daily (organization_id, date, channel, visitors, sessions, page_views, conversions, revenue, spend) VALUES 
    ('00000000-0000-0000-0000-000000000001', CURRENT_DATE - 1, 'Organic Search', 2500, 3200, 8500, 85, 8500.00, 0),
    ('00000000-0000-0000-0000-000000000001', CURRENT_DATE - 1, 'Paid Search', 1800, 2100, 5200, 72, 7200.00, 1500.00),
    ('00000000-0000-0000-0000-000000000001', CURRENT_DATE - 1, 'Social Media', 1200, 1500, 3800, 35, 3500.00, 500.00),
    ('00000000-0000-0000-0000-000000000001', CURRENT_DATE - 1, 'Email', 800, 950, 2200, 48, 4800.00, 100.00);
