from sqlalchemy.orm import Session
from typing import List, Dict, Any
from datetime import datetime, timedelta
import random


class AnalyticsService:
    """Service for marketing analytics"""
    
    def __init__(self, db: Session):
        self.db = db
    
    def get_summary(self, days: int = 30) -> Dict[str, Any]:
        """Get analytics summary for the specified period"""
        # Mock data - in production, this would query actual analytics data
        return {
            "period_days": days,
            "total_visitors": random.randint(50000, 200000),
            "total_sessions": random.randint(80000, 300000),
            "page_views": random.randint(200000, 800000),
            "bounce_rate": round(random.uniform(30, 60), 2),
            "avg_session_duration": random.randint(120, 400),
            "conversions": random.randint(1000, 10000),
            "conversion_rate": round(random.uniform(1, 5), 2),
            "revenue": round(random.uniform(50000, 300000), 2),
        }
    
    def get_trends(self, metric: str, period: str) -> Dict[str, Any]:
        """Get trend analysis for a specific metric"""
        # Generate mock trend data
        data_points = []
        end_date = datetime.utcnow()
        
        for i in range(30):
            date = end_date - timedelta(days=i)
            data_points.append({
                "date": date.strftime("%Y-%m-%d"),
                "value": random.randint(1000, 10000),
            })
        
        data_points.reverse()
        
        return {
            "metric": metric,
            "period": period,
            "data": data_points,
            "trend": random.choice(["increasing", "decreasing", "stable"]),
            "change_percent": round(random.uniform(-20, 50), 2),
        }
    
    def get_channel_performance(self) -> List[Dict[str, Any]]:
        """Get performance breakdown by marketing channel"""
        channels = [
            {"name": "Organic Search", "visitors": random.randint(20000, 80000), "conversions": random.randint(500, 2000)},
            {"name": "Paid Search", "visitors": random.randint(10000, 50000), "conversions": random.randint(300, 1500)},
            {"name": "Social Media", "visitors": random.randint(5000, 30000), "conversions": random.randint(100, 800)},
            {"name": "Email", "visitors": random.randint(3000, 15000), "conversions": random.randint(200, 1000)},
            {"name": "Direct", "visitors": random.randint(10000, 40000), "conversions": random.randint(400, 1600)},
            {"name": "Referral", "visitors": random.randint(2000, 10000), "conversions": random.randint(50, 400)},
        ]
        
        for channel in channels:
            channel["conversion_rate"] = round(
                (channel["conversions"] / channel["visitors"]) * 100, 2
            )
        
        return channels
    
    def get_conversion_funnel(self) -> Dict[str, Any]:
        """Get conversion funnel analysis"""
        funnel = {
            "stages": [
                {"name": "Visitors", "count": random.randint(100000, 200000)},
                {"name": "Landing Page Views", "count": random.randint(80000, 160000)},
                {"name": "Product Page Views", "count": random.randint(40000, 80000)},
                {"name": "Add to Cart", "count": random.randint(10000, 30000)},
                {"name": "Checkout Started", "count": random.randint(5000, 15000)},
                {"name": "Purchase Completed", "count": random.randint(3000, 10000)},
            ],
        }
        
        # Calculate conversion rates between stages
        for i in range(1, len(funnel["stages"])):
            prev_count = funnel["stages"][i-1]["count"]
            curr_count = funnel["stages"][i]["count"]
            funnel["stages"][i]["conversion_rate"] = round((curr_count / prev_count) * 100, 2)
        
        return funnel
    
    def get_attribution(self, model: str) -> Dict[str, Any]:
        """Get marketing attribution analysis"""
        channels = ["Organic Search", "Paid Search", "Social Media", "Email", "Direct", "Referral"]
        
        # Different attribution models distribute credit differently
        if model == "linear":
            weights = {channel: round(100 / len(channels), 2) for channel in channels}
        elif model == "time_decay":
            weights = {
                "Direct": 30,
                "Email": 25,
                "Paid Search": 20,
                "Organic Search": 15,
                "Social Media": 7,
                "Referral": 3,
            }
        elif model == "position_based":
            weights = {
                "Organic Search": 40,
                "Direct": 30,
                "Paid Search": 15,
                "Email": 10,
                "Social Media": 3,
                "Referral": 2,
            }
        else:
            weights = {channel: round(100 / len(channels), 2) for channel in channels}
        
        return {
            "model": model,
            "attribution": weights,
            "total_conversions": random.randint(5000, 15000),
        }
    
    def calculate_roi(self, channel: str = None) -> Dict[str, Any]:
        """Calculate ROI for marketing channels"""
        channels_data = [
            {"name": "Organic Search", "spend": 5000, "revenue": 50000},
            {"name": "Paid Search", "spend": 20000, "revenue": 80000},
            {"name": "Social Media", "spend": 8000, "revenue": 35000},
            {"name": "Email", "spend": 2000, "revenue": 25000},
            {"name": "Content Marketing", "spend": 10000, "revenue": 45000},
        ]
        
        if channel:
            channels_data = [c for c in channels_data if c["name"] == channel]
        
        for c in channels_data:
            c["roi"] = round(((c["revenue"] - c["spend"]) / c["spend"]) * 100, 2)
            c["profit"] = c["revenue"] - c["spend"]
        
        return {
            "channels": channels_data,
            "total_spend": sum(c["spend"] for c in channels_data),
            "total_revenue": sum(c["revenue"] for c in channels_data),
            "total_roi": round(
                ((sum(c["revenue"] for c in channels_data) - sum(c["spend"] for c in channels_data)) 
                 / sum(c["spend"] for c in channels_data)) * 100, 2
            ),
        }
