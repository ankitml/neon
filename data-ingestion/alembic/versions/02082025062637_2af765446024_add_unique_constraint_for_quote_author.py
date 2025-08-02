"""add unique constraint for quote author

Revision ID: 2af765446024
Revises: faea13c8c5f2
Create Date: 2025-08-02 06:26:37.379439

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '2af765446024'
down_revision = 'faea13c8c5f2'
branch_labels = None
depends_on = None


def upgrade() -> None:
    # Add unique constraint on quote and author combination
    op.execute("ALTER TABLE quotes ADD CONSTRAINT uq_quotes_quote_author UNIQUE (quote, author)")


def downgrade() -> None:
    # Remove unique constraint
    op.execute("ALTER TABLE quotes DROP CONSTRAINT uq_quotes_quote_author")