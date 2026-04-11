from sqlalchemy import Column, String, Integer, Float, DateTime, Enum as SQLEnum
from sqlalchemy.dialects.postgresql import UUID
from datetime import datetime
import uuid
from app.core.database import Base
import enum


class CompetitorType(str, enum.Enum):
    DIRECT = "direct"
    INDIRECT = "indirect"
    EMERGING = "emerging"


class Competitor(Base):
    """SQLAlchemy model for competitor"""
    __tablename__ = "competitors"
    
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    competitor_name = Column(String(200), nullable=False)
    domain = Column(String(200), nullable=False, unique=True)
    type = Column(SQLEnum(CompetitorType), default=CompetitorType.DIRECT)
    
    # Metrics
    market_share = Column(Float, default=0.0)
    seo_score = Column(Integer, default=0)
    sem_spend = Column(Float, default=0.0)
    keywords = Column(Integer, default=0)
    backlinks = Column(Integer, default=0)
    social_followers = Column(Integer, default=0)
    
    # Timestamps
    last_analyzed = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)
    created_at = Column(DateTime, default=datetime.utcnow, nullable=False)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow, nullable=False)
    
    def __repr__(self):
        return f"<Competitor {self.competitor_name} ({self.domain})>"
