from pydantic import BaseModel, Field
from typing import Optional, List
from datetime import datetime
from enum import Enum


class KeywordTrend(str, Enum):
    UP = "up"
    DOWN = "down"
    STABLE = "stable"


class KeywordAnalysis(BaseModel):
    """Schema for keyword analysis"""
    keyword: str
    position: int
    search_volume: int
    difficulty: int
    cpc: float
    trend: KeywordTrend
    url: Optional[str] = None


class SEOMetrics(BaseModel):
    """Schema for SEO metrics"""
    domain: str
    organic_traffic: int
    keywords: int
    avg_position: float
    backlinks: int
    domain_authority: int
    page_speed: int
    mobile_friendly: bool
    ssl_enabled: bool
    last_updated: datetime


class SEOMetricsResponse(BaseModel):
    """Schema for SEO metrics response"""
    metrics: SEOMetrics
    trends: dict
    recommendations: List[str]


class SEOAuditIssue(BaseModel):
    """Schema for SEO audit issue"""
    severity: str  # critical, warning, info
    category: str
    issue: str
    recommendation: str
    impact: str


class SEOAuditResponse(BaseModel):
    """Schema for SEO audit response"""
    domain: str
    overall_score: int
    audit_date: datetime
    issues: List[SEOAuditIssue]
    metrics: SEOMetrics
    recommendations: List[str]


class BacklinkAnalysis(BaseModel):
    """Schema for backlink analysis"""
    source_domain: str
    source_url: str
    target_url: str
    anchor_text: str
    domain_authority: int
    page_authority: int
    follow: bool
    first_seen: datetime


class DomainComparison(BaseModel):
    """Schema for domain comparison"""
    domain: str
    organic_traffic: int
    keywords: int
    avg_position: float
    domain_authority: int
    market_share: float
