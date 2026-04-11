"""Add full_name column and migrate data from first_name/last_name

Revision ID: 002
Revises: 001
Create Date: 2024-01-16

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = '002'
down_revision: Union[str, None] = '001'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # Add full_name column
    op.add_column('users', sa.Column('full_name', sa.String(255)))
    
    # Migrate existing data: combine first_name and last_name
    op.execute("""
        UPDATE users 
        SET full_name = TRIM(COALESCE(first_name, '') || ' ' || COALESCE(last_name, ''))
        WHERE full_name IS NULL
    """)
    
    # Update the default admin user's full_name
    op.execute("""
        UPDATE users 
        SET full_name = 'Admin User' 
        WHERE email = 'admin@marketing.com' AND full_name IS NULL
    """)


def downgrade() -> None:
    op.drop_column('users', 'full_name')
