"""Initial schema

Revision ID: 001
Revises: 
Create Date: 2024-01-15

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects.postgresql import UUID


# revision identifiers, used by Alembic.
revision: str = '001'
down_revision: Union[str, None] = None
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # Enable UUID extension
    op.execute('CREATE EXTENSION IF NOT EXISTS "uuid-ossp"')

    # Users table
    op.create_table(
        'users',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('email', sa.String(255), nullable=False, unique=True),
        sa.Column('password_hash', sa.String(255), nullable=False),
        sa.Column('first_name', sa.String(100)),
        sa.Column('last_name', sa.String(100)),
        sa.Column('role', sa.String(50), server_default='user'),
        sa.Column('is_active', sa.Boolean(), server_default='true'),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.Column('updated_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Campaigns table
    op.create_table(
        'campaigns',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('user_id', UUID(as_uuid=True), sa.ForeignKey('users.id', ondelete='CASCADE')),
        sa.Column('name', sa.String(255), nullable=False),
        sa.Column('type', sa.String(50), nullable=False),
        sa.Column('status', sa.String(50), server_default='draft'),
        sa.Column('budget', sa.Numeric(15, 2), server_default='0'),
        sa.Column('spent', sa.Numeric(15, 2), server_default='0'),
        sa.Column('impressions', sa.BigInteger(), server_default='0'),
        sa.Column('clicks', sa.BigInteger(), server_default='0'),
        sa.Column('conversions', sa.BigInteger(), server_default='0'),
        sa.Column('start_date', sa.DateTime(timezone=True), nullable=False),
        sa.Column('end_date', sa.DateTime(timezone=True)),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.Column('updated_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.CheckConstraint("type IN ('SEO', 'SEM', 'Social', 'Email', 'Display')", name='campaign_type_check'),
        sa.CheckConstraint("status IN ('active', 'paused', 'completed', 'draft')", name='campaign_status_check'),
    )

    # A/B Tests table
    op.create_table(
        'ab_tests',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('user_id', UUID(as_uuid=True), sa.ForeignKey('users.id', ondelete='CASCADE')),
        sa.Column('campaign_id', UUID(as_uuid=True), sa.ForeignKey('campaigns.id', ondelete='SET NULL')),
        sa.Column('name', sa.String(255), nullable=False),
        sa.Column('variant', sa.String(10), nullable=False),
        sa.Column('status', sa.String(50), server_default='draft'),
        sa.Column('visitors', sa.BigInteger(), server_default='0'),
        sa.Column('conversions', sa.BigInteger(), server_default='0'),
        sa.Column('conversion_rate', sa.Numeric(5, 2), server_default='0'),
        sa.Column('confidence', sa.Numeric(5, 2), server_default='0'),
        sa.Column('start_date', sa.DateTime(timezone=True), nullable=False),
        sa.Column('end_date', sa.DateTime(timezone=True)),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.Column('updated_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.CheckConstraint("variant IN ('A', 'B')", name='ab_test_variant_check'),
        sa.CheckConstraint("status IN ('running', 'completed', 'draft')", name='ab_test_status_check'),
    )

    # Competitors table
    op.create_table(
        'competitors',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('user_id', UUID(as_uuid=True), sa.ForeignKey('users.id', ondelete='CASCADE')),
        sa.Column('competitor_name', sa.String(200), nullable=False),
        sa.Column('domain', sa.String(200), nullable=False, unique=True),
        sa.Column('type', sa.String(50), server_default='direct'),
        sa.Column('market_share', sa.Numeric(5, 2), server_default='0'),
        sa.Column('seo_score', sa.Integer(), server_default='0'),
        sa.Column('sem_spend', sa.Numeric(15, 2), server_default='0'),
        sa.Column('keywords', sa.Integer(), server_default='0'),
        sa.Column('backlinks', sa.Integer(), server_default='0'),
        sa.Column('social_followers', sa.Integer(), server_default='0'),
        sa.Column('last_analyzed', sa.DateTime(timezone=True)),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.Column('updated_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.CheckConstraint("type IN ('direct', 'indirect', 'emerging')", name='competitor_type_check'),
    )

    # SEO Metrics table
    op.create_table(
        'seo_metrics',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('user_id', UUID(as_uuid=True), sa.ForeignKey('users.id', ondelete='CASCADE')),
        sa.Column('domain', sa.String(255), nullable=False),
        sa.Column('organic_traffic', sa.BigInteger(), server_default='0'),
        sa.Column('keywords_tracked', sa.Integer(), server_default='0'),
        sa.Column('avg_position', sa.Numeric(5, 2), server_default='0'),
        sa.Column('domain_authority', sa.Integer(), server_default='0'),
        sa.Column('page_speed_score', sa.Integer(), server_default='0'),
        sa.Column('backlinks', sa.Integer(), server_default='0'),
        sa.Column('referring_domains', sa.Integer(), server_default='0'),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.Column('updated_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Keyword Rankings table
    op.create_table(
        'keyword_rankings',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('seo_metric_id', UUID(as_uuid=True), sa.ForeignKey('seo_metrics.id', ondelete='CASCADE')),
        sa.Column('keyword', sa.String(255), nullable=False),
        sa.Column('position', sa.Integer(), server_default='0'),
        sa.Column('volume', sa.Integer(), server_default='0'),
        sa.Column('difficulty', sa.Integer(), server_default='0'),
        sa.Column('trend', sa.String(10), server_default='stable'),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.Column('updated_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
        sa.CheckConstraint("trend IN ('up', 'down', 'stable')", name='keyword_trend_check'),
    )

    # Analytics Events table
    op.create_table(
        'analytics_events',
        sa.Column('id', UUID(as_uuid=True), primary_key=True, server_default=sa.text('uuid_generate_v4()')),
        sa.Column('user_id', UUID(as_uuid=True), sa.ForeignKey('users.id', ondelete='CASCADE')),
        sa.Column('campaign_id', UUID(as_uuid=True), sa.ForeignKey('campaigns.id', ondelete='SET NULL')),
        sa.Column('event_type', sa.String(100), nullable=False),
        sa.Column('event_data', sa.dialects.postgresql.JSONB),
        sa.Column('channel', sa.String(100)),
        sa.Column('source', sa.String(255)),
        sa.Column('medium', sa.String(100)),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('CURRENT_TIMESTAMP')),
    )

    # Create indexes
    op.create_index('idx_campaigns_user_id', 'campaigns', ['user_id'])
    op.create_index('idx_campaigns_status', 'campaigns', ['status'])
    op.create_index('idx_campaigns_type', 'campaigns', ['type'])
    op.create_index('idx_ab_tests_user_id', 'ab_tests', ['user_id'])
    op.create_index('idx_ab_tests_status', 'ab_tests', ['status'])
    op.create_index('idx_competitors_user_id', 'competitors', ['user_id'])
    op.create_index('idx_seo_metrics_user_id', 'seo_metrics', ['user_id'])
    op.create_index('idx_analytics_events_user_id', 'analytics_events', ['user_id'])
    op.create_index('idx_analytics_events_created_at', 'analytics_events', ['created_at'])


def downgrade() -> None:
    op.drop_index('idx_analytics_events_created_at', 'analytics_events')
    op.drop_index('idx_analytics_events_user_id', 'analytics_events')
    op.drop_index('idx_seo_metrics_user_id', 'seo_metrics')
    op.drop_index('idx_competitors_user_id', 'competitors')
    op.drop_index('idx_ab_tests_status', 'ab_tests')
    op.drop_index('idx_ab_tests_user_id', 'ab_tests')
    op.drop_index('idx_campaigns_type', 'campaigns')
    op.drop_index('idx_campaigns_status', 'campaigns')
    op.drop_index('idx_campaigns_user_id', 'campaigns')

    op.drop_table('analytics_events')
    op.drop_table('keyword_rankings')
    op.drop_table('seo_metrics')
    op.drop_table('competitors')
    op.drop_table('ab_tests')
    op.drop_table('campaigns')
    op.drop_table('users')

    op.execute('DROP EXTENSION IF EXISTS "uuid-ossp"')