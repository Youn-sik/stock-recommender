import os
import asyncio
import logging
from typing import List, Dict, Optional
from datetime import datetime
import asyncpg
import json

logger = logging.getLogger(__name__)

class DatabaseManager:
    """데이터베이스 연결 및 뉴스 데이터 관리"""
    
    def __init__(self):
        self.connection_string = self._build_connection_string()
        self.pool = None
    
    def _build_connection_string(self) -> str:
        """환경 변수에서 데이터베이스 연결 문자열 생성"""
        host = os.getenv('DB_HOST', 'localhost')
        port = os.getenv('DB_PORT', '5432')
        user = os.getenv('DB_USER', 'stockuser')
        password = os.getenv('DB_PASSWORD', 'stockpass')
        database = os.getenv('DB_NAME', 'stockdb')
        
        return f"postgresql://{user}:{password}@{host}:{port}/{database}"
    
    async def initialize(self):
        """데이터베이스 연결 풀 초기화"""
        try:
            self.pool = await asyncpg.create_pool(
                self.connection_string,
                min_size=2,
                max_size=10,
                command_timeout=60
            )
            logger.info("Database connection pool initialized")
        except Exception as e:
            logger.error(f"Failed to initialize database pool: {str(e)}")
            raise
    
    async def close(self):
        """연결 풀 종료"""
        if self.pool:
            await self.pool.close()
            logger.info("Database connection pool closed")
    
    async def save_articles(self, articles: List[Dict]) -> int:
        """뉴스 기사들을 데이터베이스에 저장"""
        if not self.pool:
            await self.initialize()
        
        saved_count = 0
        
        try:
            async with self.pool.acquire() as conn:
                for article in articles:
                    try:
                        # 중복 확인 (제목과 URL 기반)
                        existing = await conn.fetchrow(
                            "SELECT id FROM news_articles WHERE title = $1 OR url = $2",
                            article.get('title'),
                            article.get('url')
                        )
                        
                        if existing:
                            continue  # 이미 존재하는 기사
                        
                        # 새 기사 저장
                        await conn.execute("""
                            INSERT INTO news_articles (
                                title, content, url, source, sentiment_score,
                                keywords, related_symbols, published_at, created_at
                            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
                        """,
                            article.get('title'),
                            article.get('content'),
                            article.get('url'),
                            article.get('source'),
                            article.get('sentiment_score'),
                            json.dumps(article.get('keywords', [])),
                            json.dumps(article.get('related_symbols', [])),
                            datetime.fromisoformat(article.get('published_at', datetime.now().isoformat())),
                            datetime.now()
                        )
                        
                        saved_count += 1
                        
                    except Exception as e:
                        logger.error(f"Error saving article '{article.get('title', 'Unknown')}': {str(e)}")
                        continue
        
        except Exception as e:
            logger.error(f"Error in save_articles: {str(e)}")
        
        logger.info(f"Saved {saved_count} new articles to database")
        return saved_count
    
    async def get_unprocessed_articles(self, limit: int = 100) -> List[Dict]:
        """감성 분석이 처리되지 않은 기사들 조회"""
        if not self.pool:
            await self.initialize()
        
        try:
            async with self.pool.acquire() as conn:
                rows = await conn.fetch("""
                    SELECT id, title, content, url, source, published_at
                    FROM news_articles 
                    WHERE sentiment_score IS NULL 
                    ORDER BY created_at DESC 
                    LIMIT $1
                """, limit)
                
                articles = []
                for row in rows:
                    articles.append({
                        'id': row['id'],
                        'title': row['title'],
                        'content': row['content'],
                        'url': row['url'],
                        'source': row['source'],
                        'published_at': row['published_at'].isoformat() if row['published_at'] else None
                    })
                
                return articles
                
        except Exception as e:
            logger.error(f"Error getting unprocessed articles: {str(e)}")
            return []
    
    async def update_article_analysis(
        self, 
        article_id: int, 
        sentiment_score: float, 
        related_symbols: List[str], 
        keywords: List[str]
    ) -> bool:
        """기사의 감성 분석 결과 업데이트"""
        if not self.pool:
            await self.initialize()
        
        try:
            async with self.pool.acquire() as conn:
                await conn.execute("""
                    UPDATE news_articles 
                    SET sentiment_score = $1, 
                        related_symbols = $2, 
                        keywords = $3
                    WHERE id = $4
                """,
                    sentiment_score,
                    json.dumps(related_symbols),
                    json.dumps(keywords),
                    article_id
                )
                
                return True
                
        except Exception as e:
            logger.error(f"Error updating article analysis for ID {article_id}: {str(e)}")
            return False
    
    async def get_recent_news(
        self, 
        symbol: Optional[str] = None, 
        hours: int = 24, 
        limit: int = 50
    ) -> List[Dict]:
        """최근 뉴스 조회 (선택적으로 종목별)"""
        if not self.pool:
            await self.initialize()
        
        try:
            async with self.pool.acquire() as conn:
                if symbol:
                    # 특정 종목 관련 뉴스
                    rows = await conn.fetch("""
                        SELECT id, title, content, url, source, sentiment_score,
                               keywords, related_symbols, published_at, created_at
                        FROM news_articles 
                        WHERE related_symbols @> $1 
                        AND created_at >= NOW() - INTERVAL '%s hours'
                        ORDER BY published_at DESC 
                        LIMIT $2
                    """ % hours, json.dumps([symbol]), limit)
                else:
                    # 전체 뉴스
                    rows = await conn.fetch("""
                        SELECT id, title, content, url, source, sentiment_score,
                               keywords, related_symbols, published_at, created_at
                        FROM news_articles 
                        WHERE created_at >= NOW() - INTERVAL '%s hours'
                        ORDER BY published_at DESC 
                        LIMIT $1
                    """ % hours, limit)
                
                articles = []
                for row in rows:
                    articles.append({
                        'id': row['id'],
                        'title': row['title'],
                        'content': row['content'],
                        'url': row['url'],
                        'source': row['source'],
                        'sentiment_score': row['sentiment_score'],
                        'keywords': json.loads(row['keywords']) if row['keywords'] else [],
                        'related_symbols': json.loads(row['related_symbols']) if row['related_symbols'] else [],
                        'published_at': row['published_at'].isoformat() if row['published_at'] else None,
                        'created_at': row['created_at'].isoformat() if row['created_at'] else None
                    })
                
                return articles
                
        except Exception as e:
            logger.error(f"Error getting recent news: {str(e)}")
            return []
    
    async def get_sentiment_summary(self, symbol: Optional[str] = None, hours: int = 24) -> Dict:
        """감성 분석 요약 통계"""
        if not self.pool:
            await self.initialize()
        
        try:
            async with self.pool.acquire() as conn:
                if symbol:
                    row = await conn.fetchrow("""
                        SELECT 
                            COUNT(*) as total_articles,
                            AVG(sentiment_score) as avg_sentiment,
                            COUNT(CASE WHEN sentiment_score > 0.3 THEN 1 END) as positive_count,
                            COUNT(CASE WHEN sentiment_score < -0.3 THEN 1 END) as negative_count,
                            COUNT(CASE WHEN sentiment_score BETWEEN -0.3 AND 0.3 THEN 1 END) as neutral_count
                        FROM news_articles 
                        WHERE related_symbols @> $1 
                        AND created_at >= NOW() - INTERVAL '%s hours'
                        AND sentiment_score IS NOT NULL
                    """ % hours, json.dumps([symbol]))
                else:
                    row = await conn.fetchrow("""
                        SELECT 
                            COUNT(*) as total_articles,
                            AVG(sentiment_score) as avg_sentiment,
                            COUNT(CASE WHEN sentiment_score > 0.3 THEN 1 END) as positive_count,
                            COUNT(CASE WHEN sentiment_score < -0.3 THEN 1 END) as negative_count,
                            COUNT(CASE WHEN sentiment_score BETWEEN -0.3 AND 0.3 THEN 1 END) as neutral_count
                        FROM news_articles 
                        WHERE created_at >= NOW() - INTERVAL '%s hours'
                        AND sentiment_score IS NOT NULL
                    """ % hours)
                
                return {
                    'total_articles': row['total_articles'] or 0,
                    'avg_sentiment': float(row['avg_sentiment']) if row['avg_sentiment'] else 0.0,
                    'positive_count': row['positive_count'] or 0,
                    'negative_count': row['negative_count'] or 0,
                    'neutral_count': row['neutral_count'] or 0,
                    'symbol': symbol,
                    'time_range_hours': hours
                }
                
        except Exception as e:
            logger.error(f"Error getting sentiment summary: {str(e)}")
            return {
                'total_articles': 0,
                'avg_sentiment': 0.0,
                'positive_count': 0,
                'negative_count': 0,
                'neutral_count': 0,
                'symbol': symbol,
                'time_range_hours': hours
            }