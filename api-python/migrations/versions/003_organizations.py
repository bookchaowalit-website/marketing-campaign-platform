"""Add organizations and multi-tenancy tables

Revision ID: 003
Revises: 002
Create Date: 2024-01-17

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.postgresql import UUID, JSONB


revision: str = '003'
down_revision: Union[str, None] = '002'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # Organizations table
    op.create_table(
        'organizations',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('name', sa.String(255), nullable=False),
        sa.Column('slug', sa.String(255), nullable=False, unique=True),
        sa.Column('logo_url', sa.String(500)),
        sa.Column('custom_domain', sa.String(255), unique=True),
        sa.Column('subdomain', sa.String(255), unique=True),
        sa.Column('timezone', sa.String(100), server_default='UTC'),
        sa.Column('locale', sa.String(10), server_default='en'),
        sa.Column('currency', sa.String(3), server_default='USD'),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.Column('updated_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Organization members
    op.create_table(
        'organization_members',
        sa.Column('org_id', UUID(as_uuid=True), sa.ForeignKey('organizations.id', ondelete='CASCADE')),
        sa.Column('user_id', UUID(as_uuid=True), sa.ForeignKey('users.id', ondelete='CASCADE')),
        sa.Column('role', sa.String(50), server_default='member'),
        sa.Column('permissions', JSONB),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.UniqueConstraint('org_id', 'user_id'),
    )

    # API keys
    op.create_table(
        'api_keys',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('org_id', UUID(as_uuid=True), sa.ForeignKey('organizations.id', ondelete='CASCADE')),
        sa.Column('user_id', UUID(as_uuid=True), sa.ForeignKey('users.id', ondelete='CASCADE')),
        sa.Column('name', sa.String(255), nullable=False),
        sa.Column('key_hash', sa.String(255), nullable=False, unique=True),
        sa.Column('key_prefix', sa.String(20)),
        sa.Column('permissions', JSONB),
        sa.Column('rate_limit', sa.Integer(), server_default='1000'),
        sa.Column('expires_at', sa.DateTime(timezone=True)),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Audit logs
    op.create_table(
        'audit_logs',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('org_id', UUID(as_uuid=True), sa.ForeignKey('organizations.id', ondelete='CASCADE')),
        sa.Column('user_id', UUID(as_uuid=True), sa.ForeignKey('users.id', ondelete='SET NULL')),
        sa.Column('action', sa.String(255), nullable=False),
        sa.Column('resource_type', sa.String(100)),
        sa.Column('resource_id', UUID(as_uuid=True)),
        sa.Column('old_values', JSONB),
        sa.Column('new_values', JSONB),
        sa.Column('ip_address', sa.String(45)),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Indexes
    op.create_index('idx_org_members_org_id', 'organization_members', ['org_id'])
    op.create_index('idx_org_members_user_id', 'organization_members', ['user_id'])
    op.create_index('idx_api_keys_org_id', 'api_keys', ['org_id'])
    op.create_index('idx_api_keys_key_hash', 'api_keys', ['key_hash'])
    op.create_index('idx_audit_logs_org_id', 'audit_logs', ['org_id'])
    op.create_index('idx_audit_logs_created_at', 'audit_logs', ['created_at'])


def downgrade() -> None:
    op.drop_index('idx_audit_logs_created_at', 'audit_logs')
    op.drop_index('idx_audit_logs_org_id', 'audit_logs')
    op.drop_index('idx_api_keys_key_hash', 'api_keys')
    op.drop_index('idx_api_keys_org_id', 'api_keys')
    op.drop_index('idx_org_members_user_id', 'organization_members')
    op.drop_index('idx_org_members_org_id', 'organization_members')

    op.drop_table('audit_logs')
    op.drop_table('api_keys')
    op.drop_table('organization_members')
    op.drop_table('organizations')
