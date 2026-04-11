"""Add billing and subscription tables

Revision ID: 004
Revises: 003
Create Date: 2024-01-18

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.postgresql import UUID, JSONB


revision: str = '004'
down_revision: Union[str, None] = '003'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # Subscription plans
    op.create_table(
        'subscription_plans',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('name', sa.String(255), nullable=False),
        sa.Column('slug', sa.String(255), nullable=False, unique=True),
        sa.Column('price_monthly', sa.Numeric(10, 2), server_default='0'),
        sa.Column('price_yearly', sa.Numeric(10, 2), server_default='0'),
        sa.Column('features', JSONB),
        sa.Column('limits', JSONB),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Subscriptions
    op.create_table(
        'subscriptions',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('org_id', UUID(as_uuid=True), sa.ForeignKey('organizations.id', ondelete='CASCADE'), unique=True),
        sa.Column('plan_id', UUID(as_uuid=True), sa.ForeignKey('subscription_plans.id')),
        sa.Column('status', sa.String(50), server_default='trial'),
        sa.Column('billing_cycle', sa.String(20), server_default='monthly'),
        sa.Column('current_period_start', sa.DateTime(timezone=True)),
        sa.Column('current_period_end', sa.DateTime(timezone=True)),
        sa.Column('stripe_customer_id', sa.String(255)),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.Column('updated_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Invoices
    op.create_table(
        'invoices',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('org_id', UUID(as_uuid=True), sa.ForeignKey('organizations.id', ondelete='CASCADE')),
        sa.Column('subscription_id', UUID(as_uuid=True), sa.ForeignKey('subscriptions.id')),
        sa.Column('stripe_invoice_id', sa.String(255), unique=True),
        sa.Column('status', sa.String(50), server_default='draft'),
        sa.Column('amount', sa.Numeric(10, 2), server_default='0'),
        sa.Column('tax', sa.Numeric(10, 2), server_default='0'),
        sa.Column('total', sa.Numeric(10, 2), server_default='0'),
        sa.Column('billing_reason', sa.String(100)),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Usage records
    op.create_table(
        'usage_records',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('org_id', UUID(as_uuid=True), sa.ForeignKey('organizations.id', ondelete='CASCADE')),
        sa.Column('subscription_id', UUID(as_uuid=True), sa.ForeignKey('subscriptions.id')),
        sa.Column('metric_type', sa.String(100), nullable=False),
        sa.Column('quantity', sa.Integer(), server_default='0'),
        sa.Column('unit_price', sa.Numeric(10, 4), server_default='0'),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Indexes
    op.create_index('idx_subscriptions_org_id', 'subscriptions', ['org_id'])
    op.create_index('idx_invoices_org_id', 'invoices', ['org_id'])
    op.create_index('idx_usage_records_org_id', 'usage_records', ['org_id'])


def downgrade() -> None:
    op.drop_index('idx_usage_records_org_id', 'usage_records')
    op.drop_index('idx_invoices_org_id', 'invoices')
    op.drop_index('idx_subscriptions_org_id', 'subscriptions')

    op.drop_table('usage_records')
    op.drop_table('invoices')
    op.drop_table('subscriptions')
    op.drop_table('subscription_plans')
