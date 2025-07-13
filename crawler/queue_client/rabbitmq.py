import os
import json
import logging
import asyncio
from datetime import datetime
from typing import Dict, Any
import pika
import pika.exceptions

logger = logging.getLogger(__name__)

class QueueClient:
    """RabbitMQ 클라이언트 (크롤러용)"""
    
    def __init__(self):
        self.connection = None
        self.channel = None
        self._setup_connection()
    
    def _setup_connection(self):
        """RabbitMQ 연결 설정"""
        try:
            host = os.getenv('RABBITMQ_HOST', 'localhost')
            port = int(os.getenv('RABBITMQ_PORT', '5672'))
            user = os.getenv('RABBITMQ_USER', 'stockmq')
            password = os.getenv('RABBITMQ_PASS', 'stockmqpass')
            
            credentials = pika.PlainCredentials(user, password)
            parameters = pika.ConnectionParameters(
                host=host,
                port=port,
                credentials=credentials,
                heartbeat=600,
                blocked_connection_timeout=300
            )
            
            self.connection = pika.BlockingConnection(parameters)
            self.channel = self.connection.channel()
            
            # Exchange 선언
            self.channel.exchange_declare(
                exchange='news.analysis',
                exchange_type='topic',
                durable=True
            )
            
            logger.info("RabbitMQ connection established")
            
        except Exception as e:
            logger.error(f"Failed to connect to RabbitMQ: {str(e)}")
            self.connection = None
            self.channel = None
    
    def _ensure_connection(self):
        """연결 상태 확인 및 재연결"""
        if not self.connection or self.connection.is_closed:
            logger.info("Reconnecting to RabbitMQ...")
            self._setup_connection()
    
    async def publish_news_update(self, article_count: int) -> bool:
        """뉴스 업데이트 알림 발행"""
        try:
            self._ensure_connection()
            
            if not self.channel:
                logger.error("No RabbitMQ channel available")
                return False
            
            message = {
                'type': 'news_update',
                'article_count': article_count,
                'timestamp': datetime.now().isoformat(),
                'source': 'crawler'
            }
            
            self.channel.basic_publish(
                exchange='news.analysis',
                routing_key='news.update',
                body=json.dumps(message),
                properties=pika.BasicProperties(
                    content_type='application/json',
                    timestamp=int(datetime.now().timestamp())
                )
            )
            
            logger.info(f"Published news update: {article_count} articles")
            return True
            
        except Exception as e:
            logger.error(f"Error publishing news update: {str(e)}")
            return False
    
    async def publish_sentiment_analysis(self, symbol: str, sentiment_data: Dict[str, Any]) -> bool:
        """감성 분석 결과 발행"""
        try:
            self._ensure_connection()
            
            if not self.channel:
                logger.error("No RabbitMQ channel available")
                return False
            
            message = {
                'type': 'sentiment_analysis',
                'symbol': symbol,
                'sentiment_data': sentiment_data,
                'timestamp': datetime.now().isoformat(),
                'source': 'crawler'
            }
            
            self.channel.basic_publish(
                exchange='news.analysis',
                routing_key='sentiment.analysis',
                body=json.dumps(message),
                properties=pika.BasicProperties(
                    content_type='application/json',
                    timestamp=int(datetime.now().timestamp())
                )
            )
            
            logger.info(f"Published sentiment analysis for {symbol}")
            return True
            
        except Exception as e:
            logger.error(f"Error publishing sentiment analysis: {str(e)}")
            return False
    
    async def publish_crawling_status(self, status: str, details: Dict[str, Any] = None) -> bool:
        """크롤링 상태 발행"""
        try:
            self._ensure_connection()
            
            if not self.channel:
                logger.error("No RabbitMQ channel available")
                return False
            
            message = {
                'type': 'crawling_status',
                'status': status,  # 'started', 'completed', 'error'
                'details': details or {},
                'timestamp': datetime.now().isoformat(),
                'source': 'crawler'
            }
            
            self.channel.basic_publish(
                exchange='news.analysis',
                routing_key='crawling.status',
                body=json.dumps(message),
                properties=pika.BasicProperties(
                    content_type='application/json',
                    timestamp=int(datetime.now().timestamp())
                )
            )
            
            logger.info(f"Published crawling status: {status}")
            return True
            
        except Exception as e:
            logger.error(f"Error publishing crawling status: {str(e)}")
            return False
    
    def close(self):
        """연결 종료"""
        try:
            if self.channel and not self.channel.is_closed:
                self.channel.close()
            
            if self.connection and not self.connection.is_closed:
                self.connection.close()
                
            logger.info("RabbitMQ connection closed")
            
        except Exception as e:
            logger.error(f"Error closing RabbitMQ connection: {str(e)}")
    
    def __del__(self):
        """소멸자에서 연결 정리"""
        self.close()