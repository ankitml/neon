#!/usr/bin/env python3
"""
Quick script to check quotes table data.
"""
import os
from sqlalchemy import create_engine, text
from dotenv import load_dotenv

load_dotenv()

database_url = os.getenv("DATABASE_URL")
engine = create_engine(database_url)

with engine.connect() as conn:
    # Basic stats
    result = conn.execute(text("""
        SELECT 
            COUNT(*) as total_quotes,
            COUNT(DISTINCT author) as unique_authors,
            COUNT(DISTINCT category) as unique_categories,
            MIN(created_at) as first_created,
            MAX(updated_at) as last_updated
        FROM quotes
    """)).fetchone()
    
    print(f"📊 Quotes Table Stats:")
    print(f"  Total quotes: {result.total_quotes:,}")
    print(f"  Unique authors: {result.unique_authors:,}")
    print(f"  Unique categories: {result.unique_categories:,}")
    print(f"  First created: {result.first_created}")
    print(f"  Last updated: {result.last_updated}")
    
    if result.total_quotes > 0:
        # Sample data
        sample = conn.execute(text("""
            SELECT quote, author, array_length(tags, 1) as tag_count, category
            FROM quotes 
            ORDER BY popularity DESC NULLS LAST 
            LIMIT 3
        """)).fetchall()
        
        print(f"\n📝 Sample quotes (top by popularity):")
        for i, row in enumerate(sample, 1):
            print(f"  {i}. \"{row.quote[:60]}...\" - {row.author}")
            print(f"     Tags: {row.tag_count}, Category: {row.category}")