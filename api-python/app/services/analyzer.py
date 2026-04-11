from typing import List, Dict, Any
from datetime import datetime
from app.models.competitor import Competitor
import random


class MarketAnalyzer:
    """Service for market analysis and competitive intelligence"""
    
    def analyze_competitor(self, competitor: Competitor) -> Dict[str, Any]:
        """Analyze a competitor's market position"""
        # Calculate market position score
        position_score = (
            competitor.market_share * 2 +
            competitor.seo_score * 0.5 +
            (competitor.keywords / 100) +
            (competitor.backlinks / 1000)
        )
        
        # Determine market position
        if position_score > 100:
            position = "market_leader"
        elif position_score > 50:
            position = "strong_competitor"
        elif position_score > 20:
            position = "growing_competitor"
        else:
            position = "emerging_competitor"
        
        # Calculate growth potential
        growth_potential = self._calculate_growth_potential(competitor)
        
        # Identify strengths and weaknesses
        strengths = self._identify_strengths(competitor)
        weaknesses = self._identify_weaknesses(competitor)
        
        return {
            "position_score": round(position_score, 2),
            "market_position": position,
            "growth_potential": growth_potential,
            "strengths": strengths,
            "weaknesses": weaknesses,
            "threat_level": self._calculate_threat_level(competitor),
        }
    
    def get_market_overview(self, competitors: List[Competitor]) -> Dict[str, Any]:
        """Get overall market overview"""
        if not competitors:
            return {"message": "No competitors to analyze"}
        
        total_market_share = sum(c.market_share for c in competitors)
        avg_seo_score = sum(c.seo_score for c in competitors) / len(competitors)
        total_sem_spend = sum(c.sem_spend for c in competitors)
        
        # Market concentration (Herfindahl-Hirschman Index)
        hhi = sum((c.market_share / total_market_share * 100) ** 2 for c in competitors)
        
        # Market concentration level
        if hhi > 2500:
            concentration = "highly_concentrated"
        elif hhi > 1500:
            concentration = "moderately_concentrated"
        else:
            concentration = "competitive"
        
        return {
            "total_competitors": len(competitors),
            "total_market_share": total_market_share,
            "average_seo_score": round(avg_seo_score, 2),
            "total_sem_spend": total_sem_spend,
            "market_concentration": concentration,
            "hhi_index": round(hhi, 2),
            "top_competitors": sorted(
                competitors,
                key=lambda c: c.market_share,
                reverse=True
            )[:3],
        }
    
    def get_recommendations(self, analysis: Dict[str, Any]) -> List[str]:
        """Get strategic recommendations based on analysis"""
        recommendations = []
        
        if analysis["market_position"] == "emerging_competitor":
            recommendations.append("Focus on niche markets to gain initial traction")
            recommendations.append("Invest in SEO to improve organic visibility")
        
        if analysis["growth_potential"] < 50:
            recommendations.append("Consider diversifying marketing channels")
            recommendations.append("Analyze competitor strategies for opportunities")
        
        if analysis["threat_level"] == "high":
            recommendations.append("Monitor competitor activities closely")
            recommendations.append("Develop differentiation strategy")
        
        for weakness in analysis.get("weaknesses", []):
            if weakness == "low_seo_score":
                recommendations.append("Improve SEO: optimize content and technical SEO")
            elif weakness == "few_backlinks":
                recommendations.append("Build quality backlinks through content marketing")
            elif weakness == "low_social_presence":
                recommendations.append("Increase social media engagement")
        
        return recommendations
    
    def _calculate_growth_potential(self, competitor: Competitor) -> float:
        """Calculate growth potential score"""
        # Simplified calculation
        base_score = 50
        
        if competitor.seo_score < 50:
            base_score += 20  # Room for SEO improvement
        
        if competitor.market_share < 10:
            base_score += 15  # Room for market share growth
        
        if competitor.keywords < 500:
            base_score += 15  # Room for keyword expansion
        
        return min(base_score, 100)
    
    def _identify_strengths(self, competitor: Competitor) -> List[str]:
        """Identify competitor strengths"""
        strengths = []
        
        if competitor.seo_score >= 80:
            strengths.append("strong_seo")
        
        if competitor.market_share >= 20:
            strengths.append("large_market_share")
        
        if competitor.backlinks >= 5000:
            strengths.append("strong_backlink_profile")
        
        if competitor.social_followers >= 100000:
            strengths.append("strong_social_presence")
        
        return strengths
    
    def _identify_weaknesses(self, competitor: Competitor) -> List[str]:
        """Identify competitor weaknesses"""
        weaknesses = []
        
        if competitor.seo_score < 50:
            weaknesses.append("low_seo_score")
        
        if competitor.backlinks < 1000:
            weaknesses.append("few_backlinks")
        
        if competitor.keywords < 300:
            weaknesses.append("limited_keyword_coverage")
        
        if competitor.social_followers < 10000:
            weaknesses.append("low_social_presence")
        
        return weaknesses
    
    def _calculate_threat_level(self, competitor: Competitor) -> str:
        """Calculate threat level"""
        score = (
            competitor.market_share * 2 +
            competitor.seo_score * 0.5
        )
        
        if score > 80:
            return "high"
        elif score > 40:
            return "medium"
        else:
            return "low"
