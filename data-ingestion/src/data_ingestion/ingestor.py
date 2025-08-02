#!/usr/bin/env python3
"""
Quotes ingestor that reads JSON file and upserts data into PostgreSQL using MERGE.
"""
import json
import os
import sys
import time
from pathlib import Path
from typing import List, Dict, Any
import logging
from sqlalchemy import create_engine, text
from sqlalchemy.orm import sessionmaker
from sqlalchemy.exc import OperationalError, DisconnectionError
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Setup logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)


class QuotesIngestor:
    def __init__(self, database_url: str = None):
        """Initialize the ingestor with database connection."""
        self.database_url = database_url or os.getenv("DATABASE_URL")
        if not self.database_url:
            raise ValueError("DATABASE_URL must be provided or set in environment")
        
        # Create engine with connection pool settings
        self.engine = create_engine(
            self.database_url,
            pool_size=5,
            max_overflow=10,
            pool_pre_ping=True,
            pool_recycle=3600
        )
        self.Session = sessionmaker(bind=self.engine)
    
    def load_quotes_json(self, file_path: str) -> List[Dict[str, Any]]:
        """Load quotes from JSON file."""
        logger.info(f"Loading quotes from {file_path}")
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                quotes = json.load(f)
            
            logger.info(f"Loaded {len(quotes)} quotes from JSON")
            return quotes
        except Exception as e:
            logger.error(f"Error loading JSON file: {e}")
            raise
    
    def normalize_quote_data(self, quote_data: Dict[str, Any]) -> Dict[str, Any]:
        """Normalize a single quote record for database insertion."""
        return {
            'quote': quote_data.get('Quote', '').strip(),
            'author': quote_data.get('Author', '').strip(),
            'tags': quote_data.get('Tags', []),
            'popularity': float(quote_data.get('Popularity', 0.0)) if quote_data.get('Popularity') else None,
            'category': quote_data.get('Category', '').strip() if quote_data.get('Category') else None
        }
    
    def upsert_batch_with_retry(self, quotes_batch: List[Dict[str, Any]], batch_num: int, max_retries: int = 3) -> None:
        """Upsert a batch of quotes with retry logic."""
        upsert_sql = text("""
            INSERT INTO quotes (quote, author, tags, popularity, category, created_at, updated_at)
            VALUES (:quote, :author, :tags, :popularity, :category, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
            ON CONFLICT (quote, author) DO UPDATE SET
                tags = EXCLUDED.tags,
                popularity = EXCLUDED.popularity,
                category = EXCLUDED.category,
                updated_at = CURRENT_TIMESTAMP
        """)
        
        normalized_batch = [self.normalize_quote_data(quote) for quote in quotes_batch]
        
        for attempt in range(max_retries):
            try:
                with self.Session() as session:
                    session.execute(upsert_sql, normalized_batch)
                    session.commit()
                    logger.info(f"Processed batch {batch_num} of {len(quotes_batch)} quotes")
                    return
            except (OperationalError, DisconnectionError) as e:
                logger.warning(f"Connection error on batch {batch_num}, attempt {attempt + 1}/{max_retries}: {e}")
                if attempt < max_retries - 1:
                    sleep_time = 2 ** attempt  # Exponential backoff
                    logger.info(f"Retrying in {sleep_time} seconds...")
                    time.sleep(sleep_time)
                else:
                    logger.error(f"Failed to process batch {batch_num} after {max_retries} attempts")
                    raise
    
    def get_table_stats(self, session) -> Dict[str, int]:
        """Get statistics about the quotes table."""
        stats_sql = text("""
            SELECT 
                COUNT(*) as total_quotes,
                COUNT(DISTINCT author) as unique_authors,
                COUNT(DISTINCT category) as unique_categories
            FROM quotes
        """)
        
        result = session.execute(stats_sql).fetchone()
        return {
            'total_quotes': result.total_quotes,
            'unique_authors': result.unique_authors,
            'unique_categories': result.unique_categories
        }
    
    def ingest_quotes(self, json_file_path: str) -> Dict[str, Any]:
        """Main method to ingest quotes from JSON file."""
        logger.info("Starting quotes ingestion process")
        
        try:
            # Load quotes from JSON
            quotes_data = self.load_quotes_json(json_file_path)
            
            if not quotes_data:
                logger.warning("No quotes found in JSON file")
                return {'success': False, 'message': 'No quotes found'}
            
            # Process in batches of 1000 with separate transactions
            batch_size = 1000
            total_quotes = len(quotes_data)
            total_batches = (total_quotes + batch_size - 1) // batch_size
            
            logger.info(f"Processing {total_quotes} quotes in {total_batches} batches of {batch_size}")
            
            for i in range(0, total_quotes, batch_size):
                batch = quotes_data[i:i + batch_size]
                batch_num = i // batch_size + 1
                
                self.upsert_batch_with_retry(batch, batch_num)
                logger.info(f"Completed batch {batch_num}/{total_batches}")
            
            # Get final stats
            with self.Session() as session:
                stats = self.get_table_stats(session)
                
                logger.info(f"Ingestion completed successfully. Stats: {stats}")
                
                return {
                    'success': True,
                    'processed_quotes': total_quotes,
                    'table_stats': stats
                }
                
        except Exception as e:
            logger.error(f"Error during ingestion: {e}")
            return {'success': False, 'error': str(e)}


def main():
    """Main function for CLI usage."""
    if len(sys.argv) != 2:
        print("Usage: python ingestor.py <path_to_quotes.json>")
        sys.exit(1)
    
    json_file_path = sys.argv[1]
    
    if not Path(json_file_path).exists():
        print(f"Error: File {json_file_path} not found")
        sys.exit(1)
    
    try:
        ingestor = QuotesIngestor()
        result = ingestor.ingest_quotes(json_file_path)
        
        if result['success']:
            print(f"✅ Successfully ingested quotes!")
            print(f"📊 Processed: {result['processed_quotes']} quotes")
            print(f"📈 Table stats: {result['table_stats']}")
        else:
            print(f"❌ Ingestion failed: {result.get('error', result.get('message'))}")
            sys.exit(1)
            
    except Exception as e:
        print(f"❌ Fatal error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()