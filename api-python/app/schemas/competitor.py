from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime
from enum import Enum


class CompetitorType(str, Enum):
    DIRECT = "direct"
    INDIRECT = "indirect"
    EMERGING = "emerging"


class CompetitorBase(BaseModel):
    competitor_name: str = Field(..., min_length=1, max_length=200)
    domain: str = Field(..., pattern=r'^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]?\.[a-zA-Z]{2,}$')
    type: CompetitorType = CompetitorType.DIRECT


class CompetitorCreate(CompetitorBase):
    """Schema for creating a competitor"""
    pass


class CompetitorUpdate(BaseModel):
    """Schema for updating a competitor"""
    competitor_name: Optional[str] = None
    domain: Optional[str] = None
    type: Optional[CompetitorType] = None
    market_share: Optional[float] = None
    seo_score: Optional[int] = None
    sem_spend: Optional[float] = None
    keywords: Optional[int] = None
    backlinks: Optional[int] = None
    social_followers: Optional[int] = None


class CompetitorResponse(CompetitorBase):
    """Schema for competitor response"""
    id: str
    market_share: float
    seo_score: int
    sem_spend: float
    keywords: int
    backlinks: int
    social_followers: int
    last_analyzed: datetime
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True
