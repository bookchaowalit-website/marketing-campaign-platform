from typing import List, Dict, Any
from datetime import datetime
from app.schemas.seo import SEOMetrics, SEOMetricsResponse, KeywordAnalysis, SEOAuditResponse, SEOAuditIssue
import random


class SEOAnalyzer:
    """Service for SEO analysis and optimization"""
    
    def get_domain_metrics(self, domain: str) -> SEOMetricsResponse:
        """Get SEO metrics for a domain"""
        # Mock data - in production, this would integrate with SEO APIs
        metrics = SEOMetrics(
            domain=domain,
            organic_traffic=random.randint(10000, 100000),
            keywords=random.randint(500, 5000),
            avg_position=round(random.uniform(10, 30), 2),
            backlinks=random.randint(1000, 20000),
            domain_authority=random.randint(30, 80),
            page_speed=random.randint(60, 95),
            mobile_friendly=random.choice([True, True, True, False]),
            ssl_enabled=True,
            last_updated=datetime.utcnow(),
        )
        
        trends = {
            "traffic_change": round(random.uniform(-20, 50), 2),
            "keywords_change": round(random.uniform(-10, 30), 2),
            "position_change": round(random.uniform(-5, 10), 2),
        }
        
        recommendations = self._generate_recommendations(metrics)
        
        return SEOMetricsResponse(
            metrics=metrics,
            trends=trends,
            recommendations=recommendations
        )
    
    def get_keywords(self, domain: str, limit: int = 50) -> List[KeywordAnalysis]:
        """Get keyword analysis for a domain"""
        # Mock keyword data
        keywords_pool = [
            "marketing tools", "seo software", "marketing automation",
            "digital marketing", "content marketing", "social media marketing",
            "email marketing", "ppc management", "marketing analytics",
            "competitor analysis", "keyword research", "link building"
        ]
        
        keywords = []
        for i, keyword in enumerate(keywords_pool[:limit]):
            keywords.append(KeywordAnalysis(
                keyword=keyword,
                position=random.randint(1, 50),
                search_volume=random.randint(1000, 50000),
                difficulty=random.randint(20, 90),
                cpc=round(random.uniform(0.5, 15.0), 2),
                trend=random.choice(["up", "down", "stable"]),
                url=f"https://{domain}/page/{i+1}"
            ))
        
        return keywords
    
    def run_audit(self, domain: str) -> SEOAuditResponse:
        """Run comprehensive SEO audit"""
        metrics_response = self.get_domain_metrics(domain)
        
        # Generate audit issues
        issues = self._generate_audit_issues(metrics_response.metrics)
        
        # Calculate overall score
        overall_score = self._calculate_audit_score(metrics_response.metrics, issues)
        
        return SEOAuditResponse(
            domain=domain,
            overall_score=overall_score,
            audit_date=datetime.utcnow(),
            issues=issues,
            metrics=metrics_response.metrics,
            recommendations=metrics_response.recommendations
        )
    
    def get_backlinks(self, domain: str, limit: int = 100) -> List[Dict[str, Any]]:
        """Get backlink analysis"""
        # Mock backlink data
        backlinks = []
        for i in range(limit):
            backlinks.append({
                "source_domain": f"source{i}.example.com",
                "source_url": f"https://source{i}.example.com/page",
                "target_url": f"https://{domain}/page",
                "anchor_text": f"anchor text {i+1}",
                "domain_authority": random.randint(20, 90),
                "page_authority": random.randint(15, 80),
                "follow": random.choice([True, False]),
                "first_seen": datetime.utcnow(),
            })
        
        return backlinks
    
    def compare_domains(self, domains: List[str]) -> List[Dict[str, Any]]:
        """Compare SEO performance across multiple domains"""
        comparison = []
        
        for domain in domains:
            metrics = self.get_domain_metrics(domain)
            comparison.append({
                "domain": domain,
                "organic_traffic": metrics.metrics.organic_traffic,
                "keywords": metrics.metrics.keywords,
                "avg_position": metrics.metrics.avg_position,
                "domain_authority": metrics.metrics.domain_authority,
                "page_speed": metrics.metrics.page_speed,
            })
        
        # Sort by domain authority
        comparison.sort(key=lambda x: x["domain_authority"], reverse=True)
        
        return comparison
    
    def track_keyword(self, keyword: str, domain: str) -> Dict[str, Any]:
        """Track keyword ranking"""
        return {
            "keyword": keyword,
            "domain": domain,
            "current_position": random.randint(1, 50),
            "previous_position": random.randint(1, 50),
            "change": random.randint(-10, 10),
            "search_volume": random.randint(1000, 50000),
            "difficulty": random.randint(20, 90),
            "tracked_at": datetime.utcnow(),
        }
    
    def _generate_recommendations(self, metrics: SEOMetrics) -> List[str]:
        """Generate SEO recommendations based on metrics"""
        recommendations = []
        
        if metrics.page_speed < 70:
            recommendations.append("Improve page speed: optimize images and enable caching")
        
        if metrics.domain_authority < 40:
            recommendations.append("Build domain authority through quality content and backlinks")
        
        if metrics.keywords < 500:
            recommendations.append("Expand keyword coverage with targeted content")
        
        if metrics.avg_position > 20:
            recommendations.append("Focus on improving rankings for top keywords")
        
        if not metrics.mobile_friendly:
            recommendations.append("Make website mobile-friendly")
        
        if not metrics.ssl_enabled:
            recommendations.append("Enable HTTPS/SSL certificate")
        
        return recommendations
    
    def _generate_audit_issues(self, metrics: SEOMetrics) -> List[SEOAuditIssue]:
        """Generate audit issues based on metrics"""
        issues = []
        
        if metrics.page_speed < 50:
            issues.append(SEOAuditIssue(
                severity="critical",
                category="performance",
                issue="Very slow page speed",
                recommendation="Optimize images, enable caching, and minimize code",
                impact="High bounce rate and poor user experience"
            ))
        elif metrics.page_speed < 70:
            issues.append(SEOAuditIssue(
                severity="warning",
                category="performance",
                issue="Slow page speed",
                recommendation="Consider optimizing images and enabling compression",
                impact="May affect rankings and user experience"
            ))
        
        if metrics.domain_authority < 30:
            issues.append(SEOAuditIssue(
                severity="warning",
                category="authority",
                issue="Low domain authority",
                recommendation="Build quality backlinks and create valuable content",
                impact="Difficulty ranking for competitive keywords"
            ))
        
        if not metrics.mobile_friendly:
            issues.append(SEOAuditIssue(
                severity="critical",
                category="mobile",
                issue="Website not mobile-friendly",
                recommendation="Implement responsive design",
                impact="Poor mobile user experience and lower mobile rankings"
            ))
        
        return issues
    
    def _calculate_audit_score(self, metrics: SEOMetrics, issues: List[SEOAuditIssue]) -> int:
        """Calculate overall audit score"""
        base_score = 100
        
        # Deduct points for issues
        for issue in issues:
            if issue.severity == "critical":
                base_score -= 20
            elif issue.severity == "warning":
                base_score -= 10
            else:
                base_score -= 5
        
        return max(0, min(100, base_score))
