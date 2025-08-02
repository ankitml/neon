#!/usr/bin/env python3
"""
BM25 Search Validator for ParadeDB
Tests various search queries to validate the BM25 index functionality.
"""
import os
import sys
from typing import List, Dict, Any
import logging
from sqlalchemy import create_engine, text
from sqlalchemy.orm import sessionmaker
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Setup logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)


class BM25Validator:
    def __init__(self, database_url: str = None):
        """Initialize the validator with database connection."""
        self.database_url = database_url or os.getenv("DATABASE_URL")
        if not self.database_url:
            raise ValueError("DATABASE_URL must be provided or set in environment")
        
        self.engine = create_engine(self.database_url)
        self.Session = sessionmaker(bind=self.engine)
        self.test_results = []
    
    def run_test_query(self, test_name: str, query: str, expected_min_results: int = 1) -> Dict[str, Any]:
        """Run a test query and return results with validation."""
        logger.info(f"Running test: {test_name}")
        
        try:
            with self.Session() as session:
                result = session.execute(text(query)).fetchall()
                
                test_result = {
                    'test_name': test_name,
                    'query': query,
                    'result_count': len(result),
                    'expected_min': expected_min_results,
                    'passed': len(result) >= expected_min_results,
                    'sample_results': result[:3] if result else [],
                    'error': None
                }
                
                if test_result['passed']:
                    logger.info(f"✅ {test_name}: Found {len(result)} results")
                    if result and len(result[0]) >= 3:
                        logger.info(f"   Top result: \"{result[0][1][:60]}...\" - {result[0][2]}")
                    elif result:
                        logger.info(f"   Result: {result[0]}")
                else:
                    logger.warning(f"❌ {test_name}: Expected at least {expected_min_results}, got {len(result)}")
                
                return test_result
                
        except Exception as e:
            logger.error(f"❌ {test_name}: Error - {e}")
            return {
                'test_name': test_name,
                'query': query,
                'result_count': 0,
                'expected_min': expected_min_results,
                'passed': False,
                'sample_results': [],
                'error': str(e)
            }
    
    def test_basic_search(self):
        """Test basic keyword search."""
        return self.run_test_query(
            "Basic Search - 'courage'",
            """
            SELECT id, quote, author, paradedb.score(id) as relevance
            FROM quotes 
            WHERE quote @@@ 'courage'
            ORDER BY paradedb.score(id) DESC
            LIMIT 5
            """,
            expected_min_results=3
        )
    
    def test_author_search(self):
        """Test author-specific search."""
        return self.run_test_query(
            "Author Search - 'Einstein'",
            """
            SELECT id, quote, author, paradedb.score(id) as relevance
            FROM quotes 
            WHERE author @@@ 'Einstein'
            ORDER BY paradedb.score(id) DESC
            LIMIT 5
            """,
            expected_min_results=1
        )
    
    def test_category_search(self):
        """Test category-based search."""
        return self.run_test_query(
            "Category Search - 'inspiration'",
            """
            SELECT id, quote, author, category, paradedb.score(id) as relevance
            FROM quotes 
            WHERE category @@@ 'inspiration'
            ORDER BY paradedb.score(id) DESC
            LIMIT 5
            """,
            expected_min_results=1
        )
    
    def test_tags_search(self):
        """Test tag-based search."""
        return self.run_test_query(
            "Tags Search - 'motivational'",
            """
            SELECT id, quote, author, tags, paradedb.score(id) as relevance
            FROM quotes 
            WHERE tags @@@ 'motivational'
            ORDER BY paradedb.score(id) DESC
            LIMIT 5
            """,
            expected_min_results=1
        )
    
    def test_phrase_search(self):
        """Test exact phrase search."""
        return self.run_test_query(
            "Phrase Search - \"never give up\"",
            """
            SELECT id, quote, author, paradedb.score(id) as relevance
            FROM quotes 
            WHERE quote @@@ '"never give up"'
            ORDER BY paradedb.score(id) DESC
            LIMIT 5
            """,
            expected_min_results=1
        )
    
    def test_boolean_search(self):
        """Test boolean AND/OR search."""
        return self.run_test_query(
            "Boolean Search - 'success AND failure'",
            """
            SELECT id, quote, author, paradedb.score(id) as relevance
            FROM quotes 
            WHERE quote @@@ 'success' AND quote @@@ 'failure'
            ORDER BY paradedb.score(id) DESC
            LIMIT 5
            """,
            expected_min_results=1
        )
    
    def test_fuzzy_search(self):
        """Test fuzzy matching."""
        return self.run_test_query(
            "Fuzzy Search - 'perseverence' (misspelled)",
            """
            SELECT id, quote, author, paradedb.score(id) as relevance
            FROM quotes 
            WHERE id @@@ paradedb.fuzzy_term('quote', 'perseverence', 1)
            ORDER BY paradedb.score(id) DESC
            LIMIT 5
            """,
            expected_min_results=1
        )
    
    def test_highlighting(self):
        """Test result highlighting."""
        return self.run_test_query(
            "Highlighting - 'innovation' with snippets",
            """
            SELECT id, 
                   paradedb.snippet(quote) as highlighted_quote,
                   author,
                   paradedb.score(id) as relevance
            FROM quotes 
            WHERE quote @@@ 'innovation'
            ORDER BY paradedb.score(id) DESC
            LIMIT 3
            """,
            expected_min_results=1
        )
    
    def test_multi_field_search(self):
        """Test searching across multiple fields."""
        return self.run_test_query(
            "Multi-field Search - 'wisdom' in quote OR author",
            """
            SELECT id, quote, author, category, paradedb.score(id) as relevance
            FROM quotes 
            WHERE quote @@@ 'wisdom' OR author @@@ 'wisdom'
            ORDER BY paradedb.score(id) DESC
            LIMIT 5
            """,
            expected_min_results=2
        )
    
    def test_relevance_scoring(self):
        """Test that relevance scores are working properly."""
        return self.run_test_query(
            "Relevance Scoring - Multiple 'love' matches",
            """
            SELECT id, quote, author, paradedb.score(id) as relevance
            FROM quotes 
            WHERE quote @@@ 'love'
            ORDER BY paradedb.score(id) DESC
            LIMIT 5
            """,
            expected_min_results=5
        )
    
    def test_index_exists(self):
        """Verify the BM25 index exists."""
        return self.run_test_query(
            "Index Verification - Check quotes_search_idx exists",
            """
            SELECT indexname, indexdef
            FROM pg_indexes 
            WHERE tablename = 'quotes' 
            AND indexname = 'quotes_search_idx'
            """,
            expected_min_results=1
        )
    
    def run_all_tests(self) -> Dict[str, Any]:
        """Run all validation tests."""
        logger.info("🚀 Starting BM25 Search Validation Tests")
        logger.info("=" * 60)
        
        # List of all test methods
        test_methods = [
            self.test_index_exists,
            self.test_basic_search,
            self.test_author_search,
            self.test_category_search,
            self.test_tags_search,
            self.test_phrase_search,
            self.test_boolean_search,
            self.test_fuzzy_search,
            self.test_highlighting,
            self.test_multi_field_search,
            self.test_relevance_scoring,
        ]
        
        # Run all tests
        for test_method in test_methods:
            result = test_method()
            self.test_results.append(result)
        
        # Calculate summary
        passed_tests = [r for r in self.test_results if r['passed']]
        failed_tests = [r for r in self.test_results if not r['passed']]
        
        summary = {
            'total_tests': len(self.test_results),
            'passed': len(passed_tests),
            'failed': len(failed_tests),
            'success_rate': (len(passed_tests) / len(self.test_results)) * 100,
            'test_results': self.test_results
        }
        
        # Print summary
        logger.info("=" * 60)
        logger.info("📊 BM25 VALIDATION SUMMARY")
        logger.info(f"Total Tests: {summary['total_tests']}")
        logger.info(f"Passed: {summary['passed']} ✅")
        logger.info(f"Failed: {summary['failed']} ❌")
        logger.info(f"Success Rate: {summary['success_rate']:.1f}%")
        
        if failed_tests:
            logger.info("\n❌ Failed Tests:")
            for test in failed_tests:
                logger.info(f"  - {test['test_name']}: {test.get('error', 'No results found')}")
        
        if summary['success_rate'] >= 80:
            logger.info("\n🎉 BM25 search functionality is working well!")
        else:
            logger.warning("\n⚠️  Some search features may need attention.")
        
        return summary


def main():
    """Main function for CLI usage."""
    try:
        validator = BM25Validator()
        summary = validator.run_all_tests()
        
        # Exit with error code if tests failed
        if summary['failed'] > 0:
            sys.exit(1)
        else:
            logger.info("\n✅ All BM25 validation tests passed!")
            
    except Exception as e:
        logger.error(f"❌ Fatal error during validation: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()