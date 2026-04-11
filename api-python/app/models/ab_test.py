"""A/B Test model for experimentation."""
from sqlalchemy import Column, String, Integer, Float, Numeric, DateTime, ForeignKey, CheckConstraint
from sqlalchemy.dialects.postgresql import UUID
from datetime import datetime
import uuid
from app.core.database import Base
import enum


class ABTestVariant(str, enum.Enum):
    A = "A"
    B = "B"


class ABTestStatus(str, enum.Enum):
    RUNNING = "running"
    COMPLETED = "completed"
    DRAFT = "draft"


class ABTest(Base):
    """SQLAlchemy model for A/B test."""
    __tablename__ = "ab_tests"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    user_id = Column(UUID(as_uuid=True), ForeignKey("users.id", ondelete="CASCADE"))
    campaign_id = Column(UUID(as_uuid=True), ForeignKey("campaigns.id", ondelete="SET NULL"))
    name = Column(String(255), nullable=False)
    variant = Column(String(10), nullable=False)
    status = Column(String(50), default="draft")
    visitors = Column(Integer, default=0)
    conversions = Column(Integer, default=0)
    conversion_rate = Column(Numeric(5, 2), default=0.0)
    confidence = Column(Numeric(5, 2), default=0.0)
    start_date = Column(DateTime(timezone=True), nullable=False)
    end_date = Column(DateTime(timezone=True))
    created_at = Column(DateTime(timezone=True), default=datetime.utcnow, nullable=False)
    updated_at = Column(DateTime(timezone=True), default=datetime.utcnow, onupdate=datetime.utcnow, nullable=False)

    __table_args__ = (
        CheckConstraint("variant IN ('A', 'B')", name="ab_test_variant_check"),
        CheckConstraint("status IN ('running', 'completed', 'draft')", name="ab_test_status_check"),
    )

    def __repr__(self):
        return f"<ABTest {self.name} ({self.variant})>"

    def calculate_conversion_rate(self) -> float:
        """Calculate and update conversion rate."""
        if self.visitors == 0:
            return 0.0
        rate = (self.conversions / self.visitors) * 100
        self.conversion_rate = rate
        return rate