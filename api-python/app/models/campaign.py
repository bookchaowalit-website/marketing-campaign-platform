"""Campaign model for marketing campaigns."""
from sqlalchemy import Column, String, Integer, Float, Numeric, DateTime, ForeignKey, CheckConstraint
from sqlalchemy.dialects.postgresql import UUID
from datetime import datetime
import uuid
from app.core.database import Base
import enum


class CampaignType(str, enum.Enum):
    SEO = "SEO"
    SEM = "SEM"
    SOCIAL = "Social"
    EMAIL = "Email"
    DISPLAY = "Display"


class CampaignStatus(str, enum.Enum):
    ACTIVE = "active"
    PAUSED = "paused"
    COMPLETED = "completed"
    DRAFT = "draft"


class Campaign(Base):
    """SQLAlchemy model for marketing campaign."""
    __tablename__ = "campaigns"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    user_id = Column(UUID(as_uuid=True), ForeignKey("users.id", ondelete="CASCADE"))
    name = Column(String(255), nullable=False)
    type = Column(String(50), nullable=False)
    status = Column(String(50), default="draft")
    budget = Column(Numeric(15, 2), default=0.0)
    spent = Column(Numeric(15, 2), default=0.0)
    impressions = Column(Integer, default=0)
    clicks = Column(Integer, default=0)
    conversions = Column(Integer, default=0)
    start_date = Column(DateTime(timezone=True), nullable=False)
    end_date = Column(DateTime(timezone=True))
    created_at = Column(DateTime(timezone=True), default=datetime.utcnow, nullable=False)
    updated_at = Column(DateTime(timezone=True), default=datetime.utcnow, onupdate=datetime.utcnow, nullable=False)

    __table_args__ = (
        CheckConstraint("type IN ('SEO', 'SEM', 'Social', 'Email', 'Display')", name="campaign_type_check"),
        CheckConstraint("status IN ('active', 'paused', 'completed', 'draft')", name="campaign_status_check"),
    )

    def __repr__(self):
        return f"<Campaign {self.name} ({self.type})>"

    @property
    def ctr(self) -> float:
        """Calculate click-through rate."""
        if self.impressions == 0:
            return 0.0
        return (self.clicks / self.impressions) * 100

    @property
    def conversion_rate(self) -> float:
        """Calculate conversion rate."""
        if self.clicks == 0:
            return 0.0
        return (self.conversions / self.clicks) * 100

    @property
    def budget_utilization(self) -> float:
        """Calculate budget utilization percentage."""
        if self.budget == 0:
            return 0.0
        return float(self.spent / self.budget) * 100