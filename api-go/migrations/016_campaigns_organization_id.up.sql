-- Add organization_id to campaigns - UP
-- campaigns was created in 001_init.up.sql, before multi-tenancy
-- (organizations table) was introduced in 005_organizations.up.sql.
-- The repository layer (campaign_repository.go, dashboard_repository.go)
-- has always queried/written organization_id, but no migration ever
-- added the column, so every campaigns query has failed with
-- "column organization_id does not exist" since multi-tenancy was added.

ALTER TABLE campaigns
    ADD COLUMN IF NOT EXISTS organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_campaigns_organization_id ON campaigns(organization_id);
