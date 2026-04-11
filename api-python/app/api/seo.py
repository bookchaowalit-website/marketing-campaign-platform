from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session
from typing import Optional, List
from app.core.database import get_db
from app.schemas.seo import SEOMetricsResponse, KeywordAnalysis, SEOAuditResponse
from app.services.seo_service import SEOService
from app.services.seo_analyzer import SEOAnalyzer
from app.services.seo_data_sources import get_seo_service

router = APIRouter()


@router.get("/sources")
async def get_sources():
    """Get available SEO data sources"""
    seo_service = get_seo_service()
    return {
        "sources": seo_service.get_available_sources(),
        "configured": len(seo_service.get_available_sources()) > 0
    }


@router.get("/metrics")
async def get_seo_metrics(
    domain: str = Query(..., description="Domain to analyze"),
    source: Optional[str] = Query(None, description="Data source (google_search_console, ahrefs, semrush)"),
    db: Session = Depends(get_db)
):
    """Get SEO metrics for a domain from real data sources"""
    seo_service = get_seo_service()
    
    # If no real sources configured, fall back to mock
    if not seo_service.get_available_sources():
        analyzer = SEOAnalyzer()
        return analyzer.get_domain_metrics(domain)
    
    return await seo_service.get_domain_metrics(domain, source)


@router.get("/keywords")
async def get_keywords(
    domain: str = Query(..., description="Domain to analyze"),
    source: Optional[str] = Query(None, description="Data source"),
    limit: int = Query(50, ge=1, le=500),
    db: Session = Depends(get_db)
):
    """Get keywords for a domain from real data sources"""
    seo_service = get_seo_service()
    
    if not seo_service.get_available_sources():
        analyzer = SEOAnalyzer()
        return analyzer.get_keywords(domain, limit=limit)
    
    return await seo_service.get_keywords(domain, source, limit)


@router.get("/backlinks")
async def get_backlinks(
    domain: str = Query(..., description="Domain to analyze"),
    source: Optional[str] = Query(None, description="Data source"),
    limit: int = Query(100, ge=1, le=1000),
    db: Session = Depends(get_db)
):
    """Get backlink analysis from real data sources"""
    seo_service = get_seo_service()
    
    if not seo_service.get_available_sources():
        analyzer = SEOAnalyzer()
        return analyzer.get_backlinks(domain, limit=limit)
    
    return await seo_service.get_backlinks(domain, source, limit)


@router.get("/rankings")
async def get_ranking_positions(
    domain: str = Query(..., description="Domain to analyze"),
    source: Optional[str] = Query(None, description="Data source"),
    db: Session = Depends(get_db)
):
    """Get ranking positions from real data sources"""
    seo_service = get_seo_service()
    
    if not seo_service.get_available_sources():
        raise HTTPException(status_code=400, detail="No real SEO data sources configured")
    
    return await seo_service.get_ranking_positions(domain, source)


@router.get("/audit", response_model=SEOAuditResponse)
async def run_seo_audit(
    domain: str = Query(..., description="Domain to audit"),
    db: Session = Depends(get_db)
):
    """Run comprehensive SEO audit (uses mock data + real sources if available)"""
    analyzer = SEOAnalyzer()
    audit = analyzer.run_audit(domain)
    return audit


@router.get("/competitors/compare")
async def compare_seo_performance(
    domains: str = Query(..., description="Comma-separated domains to compare"),
    source: Optional[str] = Query(None, description="Data source"),
    db: Session = Depends(get_db)
):
    """Compare SEO performance across multiple domains"""
    seo_service = get_seo_service()
    domain_list = [d.strip() for d in domains.split(",")]
    
    if not seo_service.get_available_sources():
        analyzer = SEOAnalyzer()
        return analyzer.compare_domains(domain_list)
    
    results = {}
    for domain in domain_list:
        results[domain] = await seo_service.get_domain_metrics(domain, source)
    
    return results


@router.post("/track")
async def track_keyword(
    keyword: str = Query(...),
    domain: str = Query(...),
    db: Session = Depends(get_db)
):
    """Track keyword ranking for a domain"""
    analyzer = SEOAnalyzer()
    return analyzer.track_keyword(keyword, domain)
