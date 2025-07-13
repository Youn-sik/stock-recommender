import re
import logging
from typing import List, Dict
import nltk
from textblob import TextBlob

# Try to import Korean language processing
try:
    from konlpy.tag import Okt
    KONLPY_AVAILABLE = True
except ImportError:
    KONLPY_AVAILABLE = False
    logging.warning("KoNLPy not available. Korean sentiment analysis will be limited.")

logger = logging.getLogger(__name__)

class SentimentAnalyzer:
    """뉴스 감성 분석기"""
    
    def __init__(self):
        # 영어 감성 분석용
        self.positive_words = [
            '상승', '증가', '호재', '급등', '강세', '돌파', '신고가', '성장',
            '긍정', '좋은', '우수', '개선', '성공', '이익', '수익', '투자',
            'rise', 'increase', 'positive', 'growth', 'profit', 'gain', 'bullish'
        ]
        
        self.negative_words = [
            '하락', '감소', '악재', '급락', '약세', '붕괴', '신저가', '손실',
            '부정', '나쁜', '악화', '실패', '위험', '우려', '불안', '매도',
            'fall', 'decrease', 'negative', 'loss', 'decline', 'bearish', 'sell'
        ]
        
        # 주식 관련 종목 패턴
        self.stock_patterns = [
            # 한국 종목 (6자리 숫자)
            r'\b\d{6}\b',
            # 미국 종목 (대문자 2-5자리)
            r'\b[A-Z]{2,5}\b',
            # 종목명 패턴
            r'(삼성전자|SK하이닉스|NAVER|LG화학|POSCO|현대차|기아|셀트리온)',
            r'(Apple|Microsoft|Google|Amazon|Tesla|NVIDIA|Meta|Netflix)'
        ]
        
        # KoNLPy 초기화
        if KONLPY_AVAILABLE:
            try:
                self.okt = Okt()
                logger.info("KoNLPy Okt initialized successfully")
            except Exception as e:
                logger.warning(f"Failed to initialize KoNLPy: {str(e)}")
                self.okt = None
        else:
            self.okt = None
    
    def analyze(self, text: str) -> float:
        """
        텍스트 감성 분석
        Returns: -1.0 (매우 부정) ~ 1.0 (매우 긍정)
        """
        if not text:
            return 0.0
        
        try:
            # 영어 텍스트에 대한 TextBlob 분석
            blob = TextBlob(text)
            english_sentiment = blob.sentiment.polarity
            
            # 한국어 텍스트에 대한 키워드 기반 분석
            korean_sentiment = self._analyze_korean_keywords(text)
            
            # 두 결과를 결합 (영어: 70%, 한국어: 30%)
            combined_sentiment = english_sentiment * 0.7 + korean_sentiment * 0.3
            
            # -1.0 ~ 1.0 범위로 정규화
            return max(-1.0, min(1.0, combined_sentiment))
            
        except Exception as e:
            logger.error(f"Error in sentiment analysis: {str(e)}")
            return 0.0
    
    def _analyze_korean_keywords(self, text: str) -> float:
        """한국어 키워드 기반 감성 분석"""
        text_lower = text.lower()
        
        positive_count = sum(1 for word in self.positive_words if word in text_lower)
        negative_count = sum(1 for word in self.negative_words if word in text_lower)
        
        total_count = positive_count + negative_count
        if total_count == 0:
            return 0.0
        
        # 긍정/부정 비율로 점수 계산
        sentiment_score = (positive_count - negative_count) / total_count
        return sentiment_score
    
    def extract_symbols(self, text: str) -> List[str]:
        """텍스트에서 관련 종목 추출"""
        symbols = []
        
        try:
            for pattern in self.stock_patterns:
                matches = re.findall(pattern, text)
                symbols.extend(matches)
            
            # 중복 제거 및 정리
            unique_symbols = list(set(symbols))
            
            # 유효한 종목만 필터링
            valid_symbols = []
            for symbol in unique_symbols:
                if self._is_valid_symbol(symbol):
                    valid_symbols.append(symbol)
            
            return valid_symbols[:10]  # 최대 10개
            
        except Exception as e:
            logger.error(f"Error extracting symbols: {str(e)}")
            return []
    
    def _is_valid_symbol(self, symbol: str) -> bool:
        """유효한 종목 코드인지 확인"""
        # 6자리 숫자 (한국 종목)
        if re.match(r'^\d{6}$', symbol):
            return True
        
        # 2-5자리 대문자 (미국 종목)
        if re.match(r'^[A-Z]{2,5}$', symbol):
            # 일반적인 단어 제외
            common_words = {'THE', 'AND', 'FOR', 'ARE', 'BUT', 'NOT', 'YOU', 'ALL', 'CAN', 'HAD', 'HER', 'WAS', 'ONE', 'OUR', 'OUT', 'DAY', 'GET', 'HAS', 'HIM', 'HIS', 'HOW', 'ITS', 'MAY', 'NEW', 'NOW', 'OLD', 'SEE', 'TWO', 'WHO', 'BOY', 'DID', 'ITS', 'LET', 'OWN', 'SAY', 'SHE', 'TOO', 'USE'}
            return symbol not in common_words
        
        # 알려진 종목명
        known_symbols = {'삼성전자', 'SK하이닉스', 'NAVER', 'LG화학', 'POSCO', 'Apple', 'Microsoft', 'Google', 'Amazon', 'Tesla', 'NVIDIA', 'Meta', 'Netflix'}
        return symbol in known_symbols
    
    def extract_keywords(self, text: str, max_keywords: int = 10) -> List[str]:
        """텍스트에서 키워드 추출"""
        keywords = []
        
        try:
            # 한국어 키워드 추출
            if self.okt and KONLPY_AVAILABLE:
                korean_keywords = self._extract_korean_keywords(text)
                keywords.extend(korean_keywords)
            
            # 영어 키워드 추출
            english_keywords = self._extract_english_keywords(text)
            keywords.extend(english_keywords)
            
            # 중복 제거 및 정리
            unique_keywords = []
            seen = set()
            for keyword in keywords:
                if keyword.lower() not in seen and len(keyword) > 1:
                    unique_keywords.append(keyword)
                    seen.add(keyword.lower())
            
            return unique_keywords[:max_keywords]
            
        except Exception as e:
            logger.error(f"Error extracting keywords: {str(e)}")
            return []
    
    def _extract_korean_keywords(self, text: str) -> List[str]:
        """한국어 키워드 추출"""
        if not self.okt:
            return []
        
        try:
            # 명사와 형용사 추출
            tokens = self.okt.pos(text)
            keywords = []
            
            for word, pos in tokens:
                # 명사(Noun) 또는 형용사(Adjective)
                if pos in ['Noun', 'Adjective'] and len(word) > 1:
                    keywords.append(word)
            
            return keywords
            
        except Exception as e:
            logger.error(f"Error in Korean keyword extraction: {str(e)}")
            return []
    
    def _extract_english_keywords(self, text: str) -> List[str]:
        """영어 키워드 추출"""
        try:
            # 간단한 정규표현식으로 영어 단어 추출
            words = re.findall(r'\b[A-Za-z]{3,}\b', text)
            
            # 불용어 제거
            stopwords = {'the', 'and', 'for', 'are', 'but', 'not', 'you', 'all', 'can', 'had', 'her', 'was', 'one', 'our', 'out', 'day', 'get', 'has', 'him', 'his', 'how', 'its', 'may', 'new', 'now', 'old', 'see', 'two', 'who', 'boy', 'did', 'let', 'own', 'say', 'she', 'too', 'use', 'with', 'from', 'they', 'know', 'want', 'been', 'good', 'much', 'some', 'time', 'very', 'when', 'come', 'here', 'just', 'like', 'long', 'make', 'many', 'over', 'such', 'take', 'than', 'them', 'well', 'were'}
            
            keywords = [word for word in words if word.lower() not in stopwords]
            return keywords
            
        except Exception as e:
            logger.error(f"Error in English keyword extraction: {str(e)}")
            return []
    
    def get_sentiment_label(self, score: float) -> str:
        """감성 점수를 라벨로 변환"""
        if score > 0.3:
            return "positive"
        elif score < -0.3:
            return "negative"
        else:
            return "neutral"