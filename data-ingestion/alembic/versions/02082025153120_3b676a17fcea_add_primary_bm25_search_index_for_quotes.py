"""add comprehensive BM25 search index for quotes

Revision ID: 3b676a17fcea
Revises: 2af765446024
Create Date: 2025-08-02 15:31:20.384805

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '3b676a17fcea'
down_revision = '2af765446024'
branch_labels = None
depends_on = None


def upgrade() -> None:
    # Create comprehensive BM25 search index for quotes (includes all searchable fields)
    op.execute("""
        CREATE INDEX quotes_search_idx ON quotes
        USING bm25 (id, quote, author, category, tags)
        WITH (key_field='id')
    """)


def downgrade() -> None:
    # Drop primary BM25 search index
    op.execute("DROP INDEX IF EXISTS quotes_search_idx")