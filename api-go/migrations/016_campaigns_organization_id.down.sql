-- Add organization_id to campaigns - DOWN

DROP INDEX IF EXISTS idx_campaigns_organization_id;
ALTER TABLE campaigns DROP COLUMN IF EXISTS organization_id;
