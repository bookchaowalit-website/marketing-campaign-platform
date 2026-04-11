from typing import Dict, Any, List
from datetime import datetime, timedelta
import random


class PredictiveAnalytics:
    """Service for predictive analytics using ML models"""
    
    def forecast_performance(self, days: int = 30) -> Dict[str, Any]:
        """Forecast future performance using ML models"""
        # Mock forecasting - in production, this would use actual ML models
        forecast_data = []
        base_date = datetime.utcnow()
        
        # Base metrics
        base_visitors = 5000
        base_conversions = 150
        base_revenue = 15000
        
        for i in range(days):
            date = base_date + timedelta(days=i)
            
            # Add some trend and seasonality
            trend_factor = 1 + (i * 0.01)  # Slight upward trend
            seasonal_factor = 1 + 0.1 * random.sin(i * 0.5)  # Weekly seasonality
            
            visitors = int(base_visitors * trend_factor * seasonal_factor + random.randint(-200, 200))
            conversions = int(base_conversions * trend_factor * seasonal_factor + random.randint(-10, 10))
            revenue = round(base_revenue * trend_factor * seasonal_factor + random.uniform(-1000, 1000), 2)
            
            forecast_data.append({
                "date": date.strftime("%Y-%m-%d"),
                "visitors": visitors,
                "conversions": conversions,
                "revenue": revenue,
                "conversion_rate": round((conversions / visitors) * 100, 2) if visitors > 0 else 0,
            })
        
        # Calculate summary statistics
        avg_visitors = sum(d["visitors"] for d in forecast_data) / len(forecast_data)
        avg_conversions = sum(d["conversions"] for d in forecast_data) / len(forecast_data)
        total_revenue = sum(d["revenue"] for d in forecast_data)
        
        return {
            "forecast_period_days": days,
            "generated_at": datetime.utcnow().isoformat(),
            "summary": {
                "avg_daily_visitors": int(avg_visitors),
                "avg_daily_conversions": int(avg_conversions),
                "total_forecasted_revenue": round(total_revenue, 2),
                "growth_rate": round(((forecast_data[-1]["visitors"] - forecast_data[0]["visitors"]) / forecast_data[0]["visitors"]) * 100, 2),
            },
            "forecast": forecast_data,
            "confidence_interval": {
                "visitors_lower": int(avg_visitors * 0.8),
                "visitors_upper": int(avg_visitors * 1.2),
                "revenue_lower": round(total_revenue * 0.8, 2),
                "revenue_upper": round(total_revenue * 1.2, 2),
            },
        }
    
    def predict_churn_risk(self, customer_data: Dict[str, Any]) -> Dict[str, Any]:
        """Predict customer churn risk"""
        # Mock churn prediction
        risk_score = random.uniform(0, 100)
        
        if risk_score > 70:
            risk_level = "high"
            recommendations = [
                "Reach out with personalized offer",
                "Schedule check-in call",
                "Review account usage patterns"
            ]
        elif risk_score > 40:
            risk_level = "medium"
            recommendations = [
                "Send engagement email",
                "Offer additional training resources"
            ]
        else:
            risk_level = "low"
            recommendations = [
                "Continue regular engagement",
                "Ask for referral or testimonial"
            ]
        
        return {
            "risk_score": round(risk_score, 2),
            "risk_level": risk_level,
            "recommendations": recommendations,
            "predicted_churn_date": (datetime.utcnow() + timedelta(days=random.randint(30, 90))).strftime("%Y-%m-%d"),
        }
    
    def optimize_budget_allocation(self, total_budget: float, channels: List[str]) -> Dict[str, Any]:
        """Optimize marketing budget allocation across channels"""
        # Mock optimization - in production, this would use historical performance data
        channel_weights = {
            "Organic Search": 0.20,
            "Paid Search": 0.30,
            "Social Media": 0.15,
            "Email": 0.10,
            "Content Marketing": 0.15,
            "Display Ads": 0.10,
        }
        
        allocation = {}
        for channel in channels:
            weight = channel_weights.get(channel, 1 / len(channels))
            allocated_budget = round(total_budget * weight, 2)
            
            # Mock expected ROI for each channel
            expected_roi = random.uniform(150, 400)
            
            allocation[channel] = {
                "budget": allocated_budget,
                "percentage": round(weight * 100, 2),
                "expected_roi": round(expected_roi, 2),
                "expected_revenue": round(allocated_budget * (expected_roi / 100), 2),
            }
        
        total_expected_revenue = sum(a["expected_revenue"] for a in allocation.values())
        
        return {
            "total_budget": total_budget,
            "allocation": allocation,
            "total_expected_revenue": round(total_expected_revenue, 2),
            "overall_expected_roi": round(((total_expected_revenue - total_budget) / total_budget) * 100, 2),
        }
