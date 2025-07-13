import asyncio
import aiohttp
import logging
from bs4 import BeautifulSoup
from datetime import datetime, timedelta
from typing import List, Dict
import re

logger = logging.getLogger(__name__)

class NaverFinanceCrawler:
    """네이버 금융 뉴스 크롤러"""
    
    def __init__(self):
        self.base_url = "https://finance.naver.com"
        self.news_url = "https://finance.naver.com/news/mainnews.naver"
        self.headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        }
        
    async def crawl_articles(self, keywords: List[str] = None, max_articles: int = 50) -> List[Dict]:
        """네이버 금융 뉴스 크롤링"""
        logger.info(f"Starting Naver Finance crawling (max: {max_articles} articles)")
        
        articles = []
        
        try:
            async with aiohttp.ClientSession(headers=self.headers) as session:
                # 메인 뉴스 페이지 크롤링
                main_articles = await self._crawl_main_news(session, max_articles // 2)
                articles.extend(main_articles)
                
                # 키워드 기반 검색 크롤링
                if keywords:
                    for keyword in keywords[:3]:  # 최대 3개 키워드
                        search_articles = await self._crawl_search_news(session, keyword, max_articles // 6)
                        articles.extend(search_articles)
                
        except Exception as e:
            logger.error(f"Error crawling Naver Finance: {str(e)}")
        
        # 중복 제거
        unique_articles = self._remove_duplicates(articles)
        logger.info(f"Collected {len(unique_articles)} unique articles from Naver Finance")
        
        return unique_articles[:max_articles]
    
    async def _crawl_main_news(self, session: aiohttp.ClientSession, max_articles: int) -> List[Dict]:
        """메인 뉴스 페이지 크롤링"""
        articles = []
        
        try:
            async with session.get(self.news_url) as response:
                if response.status != 200:
                    logger.warning(f"Failed to fetch main news page: {response.status}")
                    return articles
                
                html = await response.text()
                soup = BeautifulSoup(html, 'html.parser')
                
                # 뉴스 링크 추출
                news_links = soup.select('.newsList .news_item .link_news')
                
                # 각 뉴스 기사 상세 정보 수집
                tasks = []
                for link in news_links[:max_articles]:
                    article_url = link.get('href')
                    if article_url and not article_url.startswith('http'):
                        article_url = self.base_url + article_url
                    
                    if article_url:
                        tasks.append(self._crawl_article_detail(session, article_url))
                
                if tasks:
                    article_results = await asyncio.gather(*tasks, return_exceptions=True)
                    articles = [result for result in article_results if isinstance(result, dict)]
                
        except Exception as e:
            logger.error(f"Error crawling main news: {str(e)}")
        
        return articles
    
    async def _crawl_search_news(self, session: aiohttp.ClientSession, keyword: str, max_articles: int) -> List[Dict]:
        """키워드 검색 뉴스 크롤링"""
        articles = []
        
        try:
            search_url = f"https://finance.naver.com/news/news_search.naver?q={keyword}"
            
            async with session.get(search_url) as response:
                if response.status != 200:
                    logger.warning(f"Failed to search for keyword '{keyword}': {response.status}")
                    return articles
                
                html = await response.text()
                soup = BeautifulSoup(html, 'html.parser')
                
                # 검색 결과에서 뉴스 링크 추출
                news_links = soup.select('.newsList .news_item .link_news')
                
                # 각 뉴스 기사 상세 정보 수집
                tasks = []
                for link in news_links[:max_articles]:
                    article_url = link.get('href')
                    if article_url and not article_url.startswith('http'):
                        article_url = self.base_url + article_url
                    
                    if article_url:
                        tasks.append(self._crawl_article_detail(session, article_url))
                
                if tasks:
                    article_results = await asyncio.gather(*tasks, return_exceptions=True)
                    articles = [result for result in article_results if isinstance(result, dict)]
                
        except Exception as e:
            logger.error(f"Error searching for keyword '{keyword}': {str(e)}")
        
        return articles
    
    async def _crawl_article_detail(self, session: aiohttp.ClientSession, url: str) -> Dict:
        """개별 기사 상세 내용 크롤링"""
        try:
            async with session.get(url) as response:
                if response.status != 200:
                    return None
                
                html = await response.text()
                soup = BeautifulSoup(html, 'html.parser')
                
                # 제목 추출
                title_elem = soup.select_one('.articleSubject')
                title = title_elem.get_text().strip() if title_elem else ""
                
                # 내용 추출
                content_elem = soup.select_one('.articleCont')
                content = ""
                if content_elem:
                    # 불필요한 태그 제거
                    for script in content_elem(["script", "style"]):
                        script.decompose()
                    content = content_elem.get_text().strip()
                
                # 발행일 추출
                date_elem = soup.select_one('.article_info .dates')
                published_at = datetime.now()
                if date_elem:
                    date_text = date_elem.get_text().strip()
                    published_at = self._parse_date(date_text)
                
                # 유효한 기사인지 확인
                if not title or len(content) < 50:
                    return None
                
                return {
                    'title': title,
                    'content': content,
                    'url': url,
                    'published_at': published_at.isoformat(),
                    'source': 'naver_finance'
                }
                
        except Exception as e:
            logger.error(f"Error crawling article detail {url}: {str(e)}")
            return None
    
    def _parse_date(self, date_text: str) -> datetime:
        """날짜 텍스트를 datetime 객체로 변환"""
        try:
            # "2024.07.13 15:30" 형식 파싱
            date_pattern = r'(\d{4})\.(\d{2})\.(\d{2})\s+(\d{2}):(\d{2})'
            match = re.search(date_pattern, date_text)
            
            if match:
                year, month, day, hour, minute = map(int, match.groups())
                return datetime(year, month, day, hour, minute)
            
            # 상대 시간 파싱 ("1시간 전", "2일 전" 등)
            if "분 전" in date_text:
                minutes = int(re.search(r'(\d+)분', date_text).group(1))
                return datetime.now() - timedelta(minutes=minutes)
            elif "시간 전" in date_text:
                hours = int(re.search(r'(\d+)시간', date_text).group(1))
                return datetime.now() - timedelta(hours=hours)
            elif "일 전" in date_text:
                days = int(re.search(r'(\d+)일', date_text).group(1))
                return datetime.now() - timedelta(days=days)
            
        except Exception as e:
            logger.warning(f"Failed to parse date '{date_text}': {str(e)}")
        
        return datetime.now()
    
    def _remove_duplicates(self, articles: List[Dict]) -> List[Dict]:
        """중복 기사 제거"""
        seen_titles = set()
        unique_articles = []
        
        for article in articles:
            title = article.get('title', '')
            if title and title not in seen_titles:
                seen_titles.add(title)
                unique_articles.append(article)
        
        return unique_articles