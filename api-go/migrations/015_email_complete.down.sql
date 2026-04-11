-- Email Marketing Complete Migration - DOWN

DROP TRIGGER IF EXISTS update_email_sends_updated_at ON email_sends;
DROP TRIGGER IF EXISTS update_email_ab_tests_updated_at ON email_campaign_ab_tests;

DROP INDEX IF EXISTS idx_email_sends_campaign_status;
DROP INDEX IF EXISTS idx_email_sends_opened_at;
DROP INDEX IF EXISTS idx_email_sends_sent_at;
DROP INDEX IF EXISTS idx_email_sends_status;
DROP INDEX IF EXISTS idx_email_sends_subscriber;
DROP INDEX IF EXISTS idx_email_sends_campaign;

DROP INDEX IF EXISTS idx_email_ab_tests_campaign_b;
DROP INDEX IF EXISTS idx_email_ab_tests_campaign_a;
DROP INDEX IF EXISTS idx_email_ab_tests_status;

ALTER TABLE email_templates
    DROP COLUMN IF EXISTS subject_line,
    DROP COLUMN IF EXISTS preheader_text,
    DROP COLUMN IF EXISTS plain_text_content;

DROP TABLE IF EXISTS email_sends;
DROP TABLE IF EXISTS email_campaign_ab_tests;
