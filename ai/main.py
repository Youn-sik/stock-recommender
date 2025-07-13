from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
import logging
from datetime import datetime

from interface.decision_engine import MockDecisionEngine, TechnicalIndicators, AIDecision
from interface.models import (
    DecisionRequest, 
    DecisionResponse, 
    HealthResponse,
    ModelUpdateRequest,
    ModelUpdateResponse
)

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Create FastAPI app
app = FastAPI(
    title="Stock AI Decision Service",
    description="AI-powered stock trading decision service",
    version="1.0.0"
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Initialize decision engine
decision_engine = MockDecisionEngine()

@app.get("/", response_model=dict)
async def root():
    """Root endpoint"""
    return {
        "service": "Stock AI Decision Service",
        "version": "1.0.0",
        "status": "running",
        "timestamp": datetime.now().isoformat()
    }

@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    return HealthResponse(
        status="healthy",
        timestamp=datetime.now().isoformat(),
        service="ai-decision-service",
        version="1.0.0"
    )

@app.post("/api/v1/decision", response_model=DecisionResponse)
async def make_decision(request: DecisionRequest):
    """
    Make trading decision based on technical indicators
    """
    try:
        logger.info(f"Making decision for symbol: {request.symbol}")
        
        # Convert request to internal format
        indicators = TechnicalIndicators(
            symbol=request.symbol,
            timestamp=request.timestamp or datetime.now().isoformat(),
            price=request.price.dict() if request.price else {},
            indicators=request.indicators
        )
        
        # Get decision from engine
        decision = decision_engine.make_decision(indicators)
        
        # Convert to response format
        response = DecisionResponse(
            symbol=decision.symbol,
            decision=decision.decision,
            confidence=decision.confidence,
            reasoning=decision.reasoning,
            timestamp=decision.timestamp,
            metadata={
                "engine": "mock",
                "version": "1.0.0",
                "processing_time": "fast"
            }
        )
        
        logger.info(f"Decision made for {request.symbol}: {decision.decision} (confidence: {decision.confidence})")
        return response
        
    except Exception as e:
        logger.error(f"Error making decision: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Failed to make decision: {str(e)}")

@app.post("/api/v1/model/update", response_model=ModelUpdateResponse)
async def update_model(request: ModelUpdateRequest):
    """
    Update ML model with new training data (TODO: Implement actual ML training)
    """
    try:
        logger.info("Model update requested")
        
        # For now, just return success (actual implementation will train ML models)
        success = decision_engine.update_model(request.training_data)
        
        return ModelUpdateResponse(
            success=success,
            message="Model update scheduled for future implementation",
            timestamp=datetime.now().isoformat(),
            metadata={
                "data_points": len(request.training_data) if request.training_data else 0,
                "status": "scheduled"
            }
        )
        
    except Exception as e:
        logger.error(f"Error updating model: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Failed to update model: {str(e)}")

@app.get("/api/v1/models/status")
async def get_model_status():
    """
    Get current model status and information
    """
    return {
        "current_model": "mock_decision_engine",
        "version": "1.0.0",
        "status": "active",
        "last_updated": "2024-07-13T00:00:00Z",
        "performance_metrics": {
            "accuracy": "N/A (mock model)",
            "precision": "N/A (mock model)",
            "recall": "N/A (mock model)"
        },
        "todo": [
            "Implement LSTM price prediction model",
            "Implement Random Forest classification model",
            "Implement ensemble decision making",
            "Add real-time model updating",
            "Add backtesting capabilities",
            "Add performance monitoring"
        ]
    }

if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=8001,
        reload=True,
        log_level="info"
    )