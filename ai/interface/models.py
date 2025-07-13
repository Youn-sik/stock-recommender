from pydantic import BaseModel, Field
from typing import Dict, List, Optional, Any
from datetime import datetime

class StockPrice(BaseModel):
    """Stock price data model"""
    open_price: Optional[float] = None
    high_price: Optional[float] = None
    low_price: Optional[float] = None
    close_price: Optional[float] = None
    volume: Optional[int] = None

class DecisionRequest(BaseModel):
    """Request model for AI decision endpoint"""
    symbol: str = Field(..., description="Stock symbol")
    market: str = Field(..., description="Market (KR/US)")
    price: Optional[StockPrice] = Field(None, description="Current price data")
    indicators: Dict[str, float] = Field(..., description="Technical indicators")
    news_score: Optional[float] = Field(None, description="News sentiment score")
    timestamp: Optional[str] = Field(None, description="Request timestamp")
    metadata: Optional[Dict[str, Any]] = Field(None, description="Additional metadata")

    class Config:
        schema_extra = {
            "example": {
                "symbol": "005930",
                "market": "KR",
                "price": {
                    "open_price": 71000.0,
                    "high_price": 72500.0,
                    "low_price": 70500.0,
                    "close_price": 72000.0,
                    "volume": 12000000
                },
                "indicators": {
                    "rsi": 65.4,
                    "macd": 1250.5,
                    "macd_signal": 980.2,
                    "sma_20": 70800.0,
                    "sma_50": 69500.0,
                    "ema_12": 71200.0,
                    "ema_26": 70100.0,
                    "bollinger_upper": 73000.0,
                    "bollinger_lower": 68000.0,
                    "stochastic_k": 75.2,
                    "williams_r": -25.5,
                    "atr": 1500.0,
                    "obv": 150000000.0
                },
                "news_score": 0.3,
                "timestamp": "2024-07-13T15:30:00Z"
            }
        }

class DecisionResponse(BaseModel):
    """Response model for AI decision endpoint"""
    symbol: str = Field(..., description="Stock symbol")
    decision: str = Field(..., description="Trading decision (BUY/SELL/HOLD)")
    confidence: float = Field(..., ge=0.0, le=1.0, description="Confidence level (0.0-1.0)")
    reasoning: List[str] = Field(..., description="List of reasons for the decision")
    timestamp: str = Field(..., description="Decision timestamp")
    metadata: Optional[Dict[str, Any]] = Field(None, description="Additional response metadata")

    class Config:
        schema_extra = {
            "example": {
                "symbol": "005930",
                "decision": "BUY",
                "confidence": 0.75,
                "reasoning": [
                    "RSI oversold condition",
                    "MACD bullish crossover",
                    "Price above SMA20",
                    "High volume supports signal"
                ],
                "timestamp": "2024-07-13T15:30:00Z",
                "metadata": {
                    "engine": "mock",
                    "version": "1.0.0",
                    "processing_time": "fast"
                }
            }
        }

class ModelUpdateRequest(BaseModel):
    """Request model for model update endpoint"""
    training_data: List[Dict[str, Any]] = Field(..., description="Training data")
    model_type: Optional[str] = Field("auto", description="Model type to update")
    validation_split: Optional[float] = Field(0.2, ge=0.0, le=0.5, description="Validation split ratio")
    metadata: Optional[Dict[str, Any]] = Field(None, description="Update metadata")

    class Config:
        schema_extra = {
            "example": {
                "training_data": [
                    {
                        "symbol": "005930",
                        "indicators": {"rsi": 65.4, "macd": 1250.5},
                        "target": "BUY",
                        "outcome": 1.05,  # 5% gain
                        "timestamp": "2024-07-10T09:00:00Z"
                    }
                ],
                "model_type": "ensemble",
                "validation_split": 0.2
            }
        }

class ModelUpdateResponse(BaseModel):
    """Response model for model update endpoint"""
    success: bool = Field(..., description="Update success status")
    message: str = Field(..., description="Update result message")
    timestamp: str = Field(..., description="Update timestamp")
    metadata: Optional[Dict[str, Any]] = Field(None, description="Update result metadata")

    class Config:
        schema_extra = {
            "example": {
                "success": True,
                "message": "Model update completed successfully",
                "timestamp": "2024-07-13T15:30:00Z",
                "metadata": {
                    "data_points": 1000,
                    "validation_accuracy": 0.85,
                    "training_time": "5m 30s"
                }
            }
        }

class HealthResponse(BaseModel):
    """Health check response model"""
    status: str = Field(..., description="Service status")
    timestamp: str = Field(..., description="Health check timestamp")
    service: str = Field(..., description="Service name")
    version: str = Field(..., description="Service version")

    class Config:
        schema_extra = {
            "example": {
                "status": "healthy",
                "timestamp": "2024-07-13T15:30:00Z",
                "service": "ai-decision-service",
                "version": "1.0.0"
            }
        }

class ModelStatus(BaseModel):
    """Model status information"""
    model_name: str = Field(..., description="Model name")
    version: str = Field(..., description="Model version")
    status: str = Field(..., description="Model status")
    last_updated: str = Field(..., description="Last update timestamp")
    performance_metrics: Dict[str, Any] = Field(..., description="Performance metrics")

class BatchDecisionRequest(BaseModel):
    """Batch decision request for multiple symbols"""
    requests: List[DecisionRequest] = Field(..., description="List of decision requests")
    parallel: Optional[bool] = Field(True, description="Process in parallel")

class BatchDecisionResponse(BaseModel):
    """Batch decision response"""
    responses: List[DecisionResponse] = Field(..., description="List of decision responses")
    success_count: int = Field(..., description="Number of successful decisions")
    error_count: int = Field(..., description="Number of failed decisions")
    processing_time: float = Field(..., description="Total processing time in seconds")

# Error response models
class ErrorResponse(BaseModel):
    """Error response model"""
    error: str = Field(..., description="Error type")
    message: str = Field(..., description="Error message")
    timestamp: str = Field(..., description="Error timestamp")
    details: Optional[Dict[str, Any]] = Field(None, description="Additional error details")