import asyncio
import logging
import os
from datetime import datetime
import schedule
import time

from crawlers.naver_finance import NaverFinanceCrawler
from crawlers.daum_finance import DaumFinanceCrawler
from sentiment.analyzer import SentimentAnalyzer
from database.connection import DatabaseManager
from queue_client.rabbitmq import QueueClient

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('logs/crawler.log'),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)

class NewscrawlerService:
    """뉴스 크롤링 서비스"""
    
    def __init__(self):
        self.db_manager = DatabaseManager()
        self.queue_client = QueueClient()
        self.sentiment_analyzer = SentimentAnalyzer()
        
        # 크롤러 초기화
        self.crawlers = {
            'naver': NaverFinanceCrawler(),
            'daum': DaumFinanceCrawler(),
        }
        
        logger.info("News Crawler Service initialized")
    
    async def crawl_news(self, keywords=None):
        """뉴스 크롤링 및 처리"""
        logger.info("Starting news crawling cycle")
        
        if not keywords:
            keywords = ['주식', '증권', '코스피', '나스닥', '삼성전자', 'Apple', 'NVIDIA']
        
        all_articles = []
        
        # 각 크롤러로부터 뉴스 수집
        for name, crawler in self.crawlers.items():
            try:
                logger.info(f"Crawling from {name}")
                articles = await crawler.crawl_articles(keywords)
                
                for article in articles:
                    article['source'] = name
                    # 감성 분석 수행
                    article['sentiment_score'] = self.sentiment_analyzer.analyze(
                        article.get('title', '') + ' ' + article.get('content', '')
                    )
                
                all_articles.extend(articles)
                logger.info(f"Collected {len(articles)} articles from {name}")
                
            except Exception as e:
                logger.error(f"Error crawling from {name}: {str(e)}")
        
        # 데이터베이스에 저장
        if all_articles:
            saved_count = await self.db_manager.save_articles(all_articles)
            logger.info(f"Saved {saved_count} articles to database")
            
            # 메시지 큐에 알림 발송
            await self.queue_client.publish_news_update(len(all_articles))
        
        return all_articles
    
    async def process_sentiment_analysis(self):
        """저장된 뉴스의 감성 분석 처리"""
        logger.info("Processing sentiment analysis for unprocessed articles")
        
        # 감성 분석이 안 된 기사들 조회
        unprocessed_articles = await self.db_manager.get_unprocessed_articles()
        
        processed_count = 0
        for article in unprocessed_articles:
            try:
                text = article.get('title', '') + ' ' + article.get('content', '')
                sentiment_score = self.sentiment_analyzer.analyze(text)
                
                # 관련 종목 추출
                related_symbols = self.sentiment_analyzer.extract_symbols(text)
                
                # 키워드 추출
                keywords = self.sentiment_analyzer.extract_keywords(text)
                
                # 데이터베이스 업데이트
                await self.db_manager.update_article_analysis(
                    article['id'],
                    sentiment_score,
                    related_symbols,
                    keywords
                )
                
                processed_count += 1
                
            except Exception as e:
                logger.error(f"Error processing article {article.get('id')}: {str(e)}")
        
        logger.info(f"Processed sentiment analysis for {processed_count} articles")
        return processed_count
    
    async def run_scheduled_crawling(self):
        """스케줄된 크롤링 실행"""
        try:
            await self.crawl_news()
            await self.process_sentiment_analysis()
        except Exception as e:
            logger.error(f"Error in scheduled crawling: {str(e)}")

def run_crawler_service():
    """크롤러 서비스 실행"""
    crawler_service = NewsrawlerService()
    
    # 스케줄 설정
    schedule.every(30).minutes.do(
        lambda: asyncio.run(crawler_service.run_scheduled_crawling())
    )
    
    # 시작시 한 번 실행
    asyncio.run(crawler_service.run_scheduled_crawling())
    
    # 스케줄 실행
    while True:
        schedule.run_pending()
        time.sleep(60)  # 1분마다 스케줄 체크

if __name__ == "__main__":
    logger.info("Starting News Crawler Service")
    
    try:
        run_crawler_service()
    except KeyboardInterrupt:
        logger.info("News Crawler Service stopped by user")
    except Exception as e:
        logger.error(f"Fatal error: {str(e)}")
        raise