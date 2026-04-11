from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session
from typing import List, Optional
from app.core.database import get_db
from app.models.competitor import Competitor
from app.schemas.competitor import CompetitorCreate, CompetitorUpdate, CompetitorResponse
from app.services.competitor_service import CompetitorService
from app.services.analyzer import MarketAnalyzer

router = APIRouter()


@router.get("", response_model=List[CompetitorResponse])
async def get_competitors(
    skip: int = Query(0, ge=0),
    limit: int = Query(20, ge=1, le=100),
    db: Session = Depends(get_db)
):
    """Get all competitors"""
    service = CompetitorService(db)
    return service.get_competitors(skip=skip, limit=limit)


@router.get("/{competitor_id}", response_model=CompetitorResponse)
async def get_competitor(competitor_id: str, db: Session = Depends(get_db)):
    """Get a specific competitor by ID"""
    service = CompetitorService(db)
    competitor = service.get_competitor(competitor_id)
    if not competitor:
        raise HTTPException(status_code=404, detail="Competitor not found")
    return competitor


@router.post("", response_model=CompetitorResponse)
async def create_competitor(
    competitor: CompetitorCreate,
    db: Session = Depends(get_db)
):
    """Create a new competitor"""
    service = CompetitorService(db)
    return service.create_competitor(competitor)


@router.put("/{competitor_id}", response_model=CompetitorResponse)
async def update_competitor(
    competitor_id: str,
    competitor: CompetitorUpdate,
    db: Session = Depends(get_db)
):
    """Update an existing competitor"""
    service = CompetitorService(db)
    updated = service.update_competitor(competitor_id, competitor)
    if not updated:
        raise HTTPException(status_code=404, detail="Competitor not found")
    return updated


@router.delete("/{competitor_id}")
async def delete_competitor(competitor_id: str, db: Session = Depends(get_db)):
    """Delete a competitor"""
    service = CompetitorService(db)
    if not service.delete_competitor(competitor_id):
        raise HTTPException(status_code=404, detail="Competitor not found")
    return {"message": "Competitor deleted successfully"}


@router.get("/{competitor_id}/analyze")
async def analyze_competitor(competitor_id: str, db: Session = Depends(get_db)):
    """Analyze a competitor's market position"""
    service = CompetitorService(db)
    analyzer = MarketAnalyzer()
    
    competitor = service.get_competitor(competitor_id)
    if not competitor:
        raise HTTPException(status_code=404, detail="Competitor not found")
    
    analysis = analyzer.analyze_competitor(competitor)
    return {
        "competitor": competitor,
        "analysis": analysis,
        "recommendations": analyzer.get_recommendations(analysis)
    }


@router.get("/market/overview")
async def get_market_overview(db: Session = Depends(get_db)):
    """Get market overview with all competitors"""
    service = CompetitorService(db)
    analyzer = MarketAnalyzer()
    
    competitors = service.get_competitors()
    overview = analyzer.get_market_overview(competitors)
    
    return overview
