"""Add CRM module tables

Revision ID: 005
Revises: 004
Create Date: 2024-01-19

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.postgresql import UUID, JSONB


revision: str = '005'
down_revision: Union[str, None] = '004'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # Contacts
    op.create_table(
        'contacts',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('org_id', UUID(as_uuid=True), sa.ForeignKey('organizations.id', ondelete='CASCADE')),
        sa.Column('name', sa.String(255), nullable=False),
        sa.Column('email', sa.String(255)),
        sa.Column('phone', sa.String(50)),
        sa.Column('company', sa.String(255)),
        sa.Column('title', sa.String(255)),
        sa.Column('lead_score', sa.Integer(), server_default='0'),
        sa.Column('lead_status', sa.String(50), server_default='new'),
        sa.Column('lifecycle_stage', sa.String(50), server_default='lead'),
        sa.Column('source', sa.String(100)),
        sa.Column('custom_data', JSONB),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.Column('updated_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Pipelines
    op.create_table(
        'pipelines',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('org_id', UUID(as_uuid=True), sa.ForeignKey('organizations.id', ondelete='CASCADE')),
        sa.Column('name', sa.String(255), nullable=False),
        sa.Column('description', sa.Text()),
        sa.Column('is_default', sa.Boolean(), server_default='false'),
        sa.Column('currency', sa.String(3), server_default='USD'),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Pipeline stages
    op.create_table(
        'pipeline_stages',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('pipeline_id', UUID(as_uuid=True), sa.ForeignKey('pipelines.id', ondelete='CASCADE')),
        sa.Column('name', sa.String(255), nullable=False),
        sa.Column('color', sa.String(7)),
        sa.Column('position', sa.Integer(), server_default='0'),
        sa.Column('probability', sa.Integer(), server_default='0'),
        sa.Column('is_won_stage', sa.Boolean(), server_default='false'),
        sa.Column('is_lost_stage', sa.Boolean(), server_default='false'),
    )

    # Deals
    op.create_table(
        'deals',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('org_id', UUID(as_uuid=True), sa.ForeignKey('organizations.id', ondelete='CASCADE')),
        sa.Column('pipeline_id', UUID(as_uuid=True), sa.ForeignKey('pipelines.id')),
        sa.Column('stage_id', UUID(as_uuid=True), sa.ForeignKey('pipeline_stages.id')),
        sa.Column('contact_id', UUID(as_uuid=True), sa.ForeignKey('contacts.id')),
        sa.Column('name', sa.String(255), nullable=False),
        sa.Column('value', sa.Numeric(15, 2), server_default='0'),
        sa.Column('status', sa.String(50), server_default='open'),
        sa.Column('probability', sa.Integer(), server_default='0'),
        sa.Column('priority', sa.String(20), server_default='medium'),
        sa.Column('expected_close_date', sa.DateTime(timezone=True)),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.Column('updated_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Tasks
    op.create_table(
        'tasks',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('org_id', UUID(as_uuid=True), sa.ForeignKey('organizations.id', ondelete='CASCADE')),
        sa.Column('contact_id', UUID(as_uuid=True), sa.ForeignKey('contacts.id')),
        sa.Column('deal_id', UUID(as_uuid=True), sa.ForeignKey('deals.id')),
        sa.Column('title', sa.String(255), nullable=False),
        sa.Column('type', sa.String(50), server_default='call'),
        sa.Column('priority', sa.String(20), server_default='medium'),
        sa.Column('status', sa.String(50), server_default='pending'),
        sa.Column('due_date', sa.DateTime(timezone=True)),
        sa.Column('completed_at', sa.DateTime(timezone=True)),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Notes
    op.create_table(
        'notes',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('org_id', UUID(as_uuid=True), sa.ForeignKey('organizations.id', ondelete='CASCADE')),
        sa.Column('contact_id', UUID(as_uuid=True), sa.ForeignKey('contacts.id')),
        sa.Column('deal_id', UUID(as_uuid=True), sa.ForeignKey('deals.id')),
        sa.Column('content', sa.Text(), nullable=False),
        sa.Column('is_private', sa.Boolean(), server_default='false'),
        sa.Column('created_by', UUID(as_uuid=True), sa.ForeignKey('users.id')),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Indexes
    op.create_index('idx_contacts_org_id', 'contacts', ['org_id'])
    op.create_index('idx_deals_org_id', 'deals', ['org_id'])
    op.create_index('idx_tasks_org_id', 'tasks', ['org_id'])
    op.create_index('idx_notes_org_id', 'notes', ['org_id'])


def downgrade() -> None:
    op.drop_index('idx_notes_org_id', 'notes')
    op.drop_index('idx_tasks_org_id', 'tasks')
    op.drop_index('idx_deals_org_id', 'deals')
    op.drop_index('idx_contacts_org_id', 'contacts')

    op.drop_table('notes')
    op.drop_table('tasks')
    op.drop_table('deals')
    op.drop_table('pipeline_stages')
    op.drop_table('pipelines')
    op.drop_table('contacts')
