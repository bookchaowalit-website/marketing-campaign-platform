from sqlalchemy.orm import Session
from typing import List, Optional
from app.models.competitor import Competitor, CompetitorType
from app.schemas.competitor import CompetitorCreate, CompetitorUpdate
import uuid


class CompetitorService:
    """Service class for competitor operations"""
    
    def __init__(self, db: Session):
        self.db = db
    
    def get_competitors(self, skip: int = 0, limit: int = 20) -> List[Competitor]:
        """Get all competitors with pagination"""
        return self.db.query(Competitor).offset(skip).limit(limit).all()
    
    def get_competitor(self, competitor_id: str) -> Optional[Competitor]:
        """Get a competitor by ID"""
        try:
            uuid_obj = uuid.UUID(competitor_id)
            return self.db.query(Competitor).filter(Competitor.id == uuid_obj).first()
        except ValueError:
            return None
    
    def get_competitor_by_domain(self, domain: str) -> Optional[Competitor]:
        """Get a competitor by domain"""
        return self.db.query(Competitor).filter(Competitor.domain == domain).first()
    
    def create_competitor(self, competitor: CompetitorCreate) -> Competitor:
        """Create a new competitor"""
        db_competitor = Competitor(
            competitor_name=competitor.competitor_name,
            domain=competitor.domain,
            type=competitor.type,
        )
        self.db.add(db_competitor)
        self.db.commit()
        self.db.refresh(db_competitor)
        return db_competitor
    
    def update_competitor(self, competitor_id: str, competitor: CompetitorUpdate) -> Optional[Competitor]:
        """Update an existing competitor"""
        db_competitor = self.get_competitor(competitor_id)
        if not db_competitor:
            return None
        
        update_data = competitor.model_dump(exclude_unset=True)
        for field, value in update_data.items():
            setattr(db_competitor, field, value)
        
        self.db.commit()
        self.db.refresh(db_competitor)
        return db_competitor
    
    def delete_competitor(self, competitor_id: str) -> bool:
        """Delete a competitor"""
        db_competitor = self.get_competitor(competitor_id)
        if not db_competitor:
            return False
        
        self.db.delete(db_competitor)
        self.db.commit()
        return True
    
    def get_total_count(self) -> int:
        """Get total number of competitors"""
        return self.db.query(Competitor).count()
