from fastapi import APIRouter, HTTPException
from typing import Dict, Any, List
import grpc
from app.core.config import settings

# Import protobuf messages (generated from .proto files)
import app.grpc_client.campaign_pb2 as campaign_pb2
import app.grpc_client.analytics_pb2 as analytics_pb2
import app.grpc_client.competitor_pb2 as competitor_pb2

router = APIRouter()

# gRPC channel singleton for connection reuse
_channel = None
_campaign_stub = None
_analytics_stub = None
_competitor_stub = None


def get_grpc_channel():
    """Get or create gRPC channel"""
    global _channel
    if _channel is None:
        _channel = grpc.insecure_channel(settings.GRPC_SERVER_URL)
    return _channel


def get_campaign_stub():
    """Get campaign service stub"""
    # Lazy import to avoid import errors if protos not generated
    from app.grpc_client.campaign_pb2_grpc import CampaignServiceStub
    return CampaignServiceStub(get_grpc_channel())


def get_analytics_stub():
    """Get analytics service stub"""
    from app.grpc_client.analytics_pb2_grpc import AnalyticsServiceStub
    return AnalyticsServiceStub(get_grpc_channel())


def get_competitor_stub():
    """Get competitor service stub"""
    from app.grpc_client.competitor_pb2_grpc import CompetitorServiceStub
    return CompetitorServiceStub(get_grpc_channel())


@router.get("/status")
async def get_grpc_status():
    """Get gRPC connection status"""
    try:
        channel = get_grpc_channel()
        grpc.channel_ready_future(channel).result(timeout=2)
        return {
            "status": "connected",
            "server": settings.GRPC_SERVER_URL,
            "services": ["CampaignService", "AnalyticsService", "CompetitorService"]
        }
    except Exception as e:
        return {
            "status": "disconnected",
            "server": settings.GRPC_SERVER_URL,
            "error": str(e)
        }


@router.post("/campaigns/sync")
async def sync_campaigns_from_grpc():
    """Sync campaign data from Go service via gRPC"""
    try:
        stub = get_campaign_stub()
        response = stub.ListCampaigns(
            campaign_pb2.ListCampaignsRequest(limit=100, offset=0)
        )
        
        campaigns = []
        for c in response.campaigns:
            campaigns.append({
                "id": c.id,
                "name": c.name,
                "type": c.type,
                "status": c.status,
                "budget": c.budget,
                "spent": c.spent,
            })
        
        return {
            "message": "Campaign sync completed",
            "status": "success",
            "synced_count": len(campaigns),
            "campaigns": campaigns
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"gRPC sync failed: {str(e)}")


@router.get("/analytics/stream")
async def stream_analytics():
    """Stream real-time analytics from Go service"""
    try:
        stub = get_analytics_stub()
        response = stub.GetDashboard(
            analytics_pb2.GetDashboardRequest(period="7d")
        )
        
        return {
            "message": "Analytics data retrieved",
            "status": "success",
            "data": {
                "total_campaigns": response.total_campaigns,
                "active_campaigns": response.active_campaigns,
                "total_spend": response.total_spend,
                "roi": response.roi,
                "conversion_rate": response.conversion_rate,
            }
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Analytics stream failed: {str(e)}")


@router.get("/campaigns/{campaign_id}")
async def get_campaign_from_grpc(campaign_id: str):
    """Get campaign data from Go service via gRPC"""
    try:
        stub = get_campaign_stub()
        response = stub.GetCampaign(
            campaign_pb2.GetCampaignRequest(id=campaign_id)
        )
        
        return {
            "id": response.id,
            "name": response.name,
            "type": response.type,
            "status": response.status,
            "budget": response.budget,
            "spent": response.spent,
            "impressions": response.impressions,
            "clicks": response.clicks,
            "conversions": response.conversions,
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"gRPC campaign fetch failed: {str(e)}")


@router.get("/competitors/{competitor_id}")
async def get_competitor_from_grpc(competitor_id: str):
    """Get competitor data from Go service via gRPC"""
    try:
        stub = get_competitor_stub()
        response = stub.GetCompetitor(
            competitor_pb2.GetCompetitorRequest(id=competitor_id)
        )
        
        return {
            "id": response.id,
            "competitor_name": response.competitor_name,
            "domain": response.domain,
            "type": response.type,
            "market_share": response.market_share,
            "seo_score": response.seo_score,
            "sem_spend": response.sem_spend,
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"gRPC competitor fetch failed: {str(e)}")
