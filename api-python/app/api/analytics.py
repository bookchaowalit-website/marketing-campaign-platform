from fastapi import APIRouter, Depends, Query
from sqlalchemy.orm import Session
from typing import Dict, Any
from app.core.database import get_db
from app.services.analytics_service import AnalyticsService
from app.services.predictive_analytics import PredictiveAnalytics

router = APIRouter()


@router.get("/summary")
async def get_analytics_summary(
    days: int = Query(30, ge=1, le=365),
    db: Session = Depends(get_db)
):
    """Get analytics summary for the specified period"""
    service = AnalyticsService(db)
    return service.get_summary(days)


@router.get("/trends")
async def get_trends(
    metric: str = Query("traffic", description="Metric to analyze"),
    period: str = Query("daily", description="Period: daily, weekly, monthly"),
    db: Session = Depends(get_db)
):
    """Get trend analysis for a specific metric"""
    service = AnalyticsService(db)
    return service.get_trends(metric, period)


@router.get("/channels")
async def get_channel_performance(
    db: Session = Depends(get_db)
):
    """Get performance breakdown by marketing channel"""
    service = AnalyticsService(db)
    return service.get_channel_performance()


@router.get("/conversion/funnel")
async def get_conversion_funnel(
    db: Session = Depends(get_db)
):
    """Get conversion funnel analysis"""
    service = AnalyticsService(db)
    return service.get_conversion_funnel()


@router.post("/predict")
async def predict_performance(
    days: int = Query(30, ge=1, le=90),
    db: Session = Depends(get_db)
):
    """Predict future performance using ML models"""
    predictor = PredictiveAnalytics()
    return predictor.forecast_performance(days)


@router.get("/attribution")
async def get_attribution_model(
    model: str = Query("linear", description="Attribution model: linear, time_decay, position_based"),
    db: Session = Depends(get_db)
):
    """Get marketing attribution analysis"""
    service = AnalyticsService(db)
    return service.get_attribution(model)


@router.get("/roi")
async def calculate_roi(
    channel: str = Query(None, description="Specific channel to analyze"),
    db: Session = Depends(get_db)
):
    """Calculate ROI for marketing channels"""
    service = AnalyticsService(db)
    return service.calculate_roi(channel)
