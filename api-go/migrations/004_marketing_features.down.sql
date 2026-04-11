-- Marketing Features Migration - DOWN
-- Removes all tables for Email Marketing, Social Media, Landing Pages, Automation, Influencer Marketing, and Reporting

-- Drop triggers
DROP TRIGGER IF EXISTS update_email_lists_updated_at ON email_lists;
DROP TRIGGER IF EXISTS update_email_subscribers_updated_at ON email_subscribers;
DROP TRIGGER IF EXISTS update_email_templates_updated_at ON email_templates;
DROP TRIGGER IF EXISTS update_email_campaigns_updated_at ON email_campaigns;
DROP TRIGGER IF EXISTS update_social_accounts_updated_at ON social_accounts;
DROP TRIGGER IF EXISTS update_social_posts_updated_at ON social_posts;
DROP TRIGGER IF EXISTS update_landing_pages_updated_at ON landing_pages;
DROP TRIGGER IF EXISTS update_automation_workflows_updated_at ON automation_workflows;
DROP TRIGGER IF EXISTS update_influencers_updated_at ON influencers;
DROP TRIGGER IF EXISTS update_custom_reports_updated_at ON custom_reports;

-- Drop indexes
DROP INDEX IF EXISTS idx_email_lists_org;
DROP INDEX IF EXISTS idx_email_subscribers_org;
DROP INDEX IF EXISTS idx_email_subscribers_email;
DROP INDEX IF EXISTS idx_email_subscribers_status;
DROP INDEX IF EXISTS idx_email_campaigns_org;
DROP INDEX IF EXISTS idx_email_campaigns_status;
DROP INDEX IF EXISTS idx_email_events_campaign;
DROP INDEX IF EXISTS idx_email_events_subscriber;
DROP INDEX IF EXISTS idx_email_events_type;
DROP INDEX IF EXISTS idx_email_events_created;
DROP INDEX IF EXISTS idx_social_accounts_org;
DROP INDEX IF EXISTS idx_social_accounts_platform;
DROP INDEX IF EXISTS idx_social_posts_org;
DROP INDEX IF EXISTS idx_social_posts_status;
DROP INDEX IF EXISTS idx_social_post_schedules_time;
DROP INDEX IF EXISTS idx_social_analytics_account;
DROP INDEX IF EXISTS idx_landing_pages_org;
DROP INDEX IF EXISTS idx_landing_pages_slug;
DROP INDEX IF EXISTS idx_landing_pages_status;
DROP INDEX IF EXISTS idx_landing_page_leads_org;
DROP INDEX IF EXISTS idx_landing_page_leads_email;
DROP INDEX IF EXISTS idx_automation_workflows_org;
DROP INDEX IF EXISTS idx_automation_workflows_status;
DROP INDEX IF EXISTS idx_automation_contacts_workflow;
DROP INDEX IF EXISTS idx_automation_contacts_status;
DROP INDEX IF EXISTS idx_contact_tags_org;
DROP INDEX IF EXISTS idx_influencers_org;
DROP INDEX IF EXISTS idx_influencers_status;
DROP INDEX IF EXISTS idx_influencer_campaigns_org;
DROP INDEX IF EXISTS idx_influencer_collaborations_campaign;
DROP INDEX IF EXISTS idx_influencer_content_collab;
DROP INDEX IF EXISTS idx_custom_reports_org;
DROP INDEX IF EXISTS idx_report_snapshots_report;
DROP INDEX IF EXISTS idx_custom_dashboards_org;
DROP INDEX IF EXISTS idx_analytics_events_unified_org;
DROP INDEX IF EXISTS idx_analytics_events_unified_type;
DROP INDEX IF EXISTS idx_analytics_events_unified_created;
DROP INDEX IF EXISTS idx_marketing_attribution_org;

-- Drop tables in reverse order of dependencies

-- Reporting
DROP TABLE IF EXISTS marketing_attribution;
DROP TABLE IF EXISTS analytics_events_unified;
DROP TABLE IF EXISTS dashboard_widgets;
DROP TABLE IF EXISTS custom_dashboards;
DROP TABLE IF EXISTS report_snapshots;
DROP TABLE IF EXISTS report_widgets;
DROP TABLE IF EXISTS custom_reports;
DROP TABLE IF EXISTS report_templates;

-- Influencer Marketing
DROP TABLE IF EXISTS influencer_payments;
DROP TABLE IF EXISTS influencer_content;
DROP TABLE IF EXISTS influencer_collaborations;
DROP TABLE IF EXISTS influencer_campaigns;
DROP TABLE IF EXISTS influencer_profiles;
DROP TABLE IF EXISTS influencers;

-- Automation
DROP TABLE IF EXISTS contact_tag_relations;
DROP TABLE IF EXISTS contact_tags;
DROP TABLE IF EXISTS automation_events;
DROP TABLE IF EXISTS automation_contacts;
DROP TABLE IF EXISTS automation_goals;
DROP TABLE IF EXISTS automation_steps;
DROP TABLE IF EXISTS automation_workflows;

-- Landing Pages
DROP TABLE IF EXISTS landing_page_ab_tests;
DROP TABLE IF EXISTS landing_page_leads;
DROP TABLE IF EXISTS landing_page_forms;
DROP TABLE IF EXISTS landing_page_components;
DROP TABLE IF EXISTS landing_pages;

-- Social Media
DROP TABLE IF EXISTS content_calendar_events;
DROP TABLE IF EXISTS social_analytics;
DROP TABLE IF EXISTS social_post_schedules;
DROP TABLE IF EXISTS social_posts;
DROP TABLE IF EXISTS social_accounts;

-- Email Marketing
DROP TABLE IF EXISTS email_events;
DROP TABLE IF EXISTS email_campaign_lists;
DROP TABLE IF EXISTS email_campaigns;
DROP TABLE IF EXISTS email_templates;
DROP TABLE IF EXISTS email_list_subscribers;
DROP TABLE IF EXISTS email_subscribers;
DROP TABLE IF EXISTS email_lists;