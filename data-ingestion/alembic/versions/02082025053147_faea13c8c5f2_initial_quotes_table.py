"""initial quotes table

Revision ID: faea13c8c5f2
Revises: 
Create Date: 2025-08-02 05:31:47.625460

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = 'faea13c8c5f2'
down_revision = None
branch_labels = None
depends_on = None


def upgrade() -> None:
    # Execute the SQL from migration.sql
    op.execute("""
        -- Main quotes table with array column for tags
        CREATE TABLE quotes (
            id SERIAL PRIMARY KEY,
            quote TEXT NOT NULL,
            author VARCHAR(255) NOT NULL,
            tags TEXT[] NOT NULL DEFAULT '{}',
            popularity DECIMAL(10, 8),
            category VARCHAR(100),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );

        -- Indexes for better query performance
        CREATE INDEX idx_quotes_author ON quotes(author);
        CREATE INDEX idx_quotes_category ON quotes(category);
        CREATE INDEX idx_quotes_popularity ON quotes(popularity DESC);

        -- GIN index for efficient array operations on tags
        CREATE INDEX idx_quotes_tags ON quotes USING GIN (tags);
    """)


def downgrade() -> None:
    # Drop the table and all indexes
    op.execute("DROP TABLE IF EXISTS quotes CASCADE;")