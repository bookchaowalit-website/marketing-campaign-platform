-- Multi-Tenancy & Organizations Migration - DOWN

DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;
DROP TRIGGER IF EXISTS update_organization_members_updated_at ON organization_members;
DROP TRIGGER IF EXISTS update_organizations_updated_at ON organizations;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_audit_logs_created_at;
DROP INDEX IF EXISTS idx_audit_logs_resource;
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP INDEX IF EXISTS idx_audit_logs_org_id;
DROP INDEX IF EXISTS idx_api_keys_prefix;
DROP INDEX IF EXISTS idx_api_keys_user_id;
DROP INDEX IF EXISTS idx_api_keys_org_id;
DROP INDEX IF EXISTS idx_roles_name;
DROP INDEX IF EXISTS idx_roles_org_id;
DROP INDEX IF EXISTS idx_team_invitations_status;
DROP INDEX IF EXISTS idx_team_invitations_org_id;
DROP INDEX IF EXISTS idx_team_invitations_token;
DROP INDEX IF EXISTS idx_team_invitations_email;
DROP INDEX IF EXISTS idx_organization_members_role;
DROP INDEX IF EXISTS idx_organization_members_user_id;
DROP INDEX IF EXISTS idx_organization_members_org_id;
DROP INDEX IF EXISTS idx_organizations_active;
DROP INDEX IF EXISTS idx_organizations_subdomain;
DROP INDEX IF EXISTS idx_organizations_custom_domain;
DROP INDEX IF EXISTS idx_organizations_slug;

DELETE FROM roles WHERE is_system_role = true;

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS organization_settings;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS team_invitations;
DROP TABLE IF EXISTS organization_members;
DROP TABLE IF EXISTS organizations;