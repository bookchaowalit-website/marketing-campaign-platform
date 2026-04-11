from fastapi import FastAPI, Depends
from fastapi.middleware.cors import CORSMiddleware
from app.core.config import settings
from app.core.auth import get_current_user_id
from app.api import competitors, seo, analytics, grpc_client


def create_app() -> FastAPI:
    """Create and configure FastAPI application"""

    app = FastAPI(
        title=settings.APP_NAME,
        version=settings.APP_VERSION,
        description="Python FastAPI backend for Marketing SaaS - ML/AI powered analytics and research",
    )

    # CORS middleware
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    # Include routers with auth protection
    app.include_router(
        competitors.router, 
        prefix="/api/v1/competitors", 
        tags=["Competitors"],
        dependencies=[Depends(get_current_user_id)]
    )
    app.include_router(
        seo.router, 
        prefix="/api/v1/seo", 
        tags=["SEO"],
        dependencies=[Depends(get_current_user_id)]
    )
    app.include_router(
        analytics.router, 
        prefix="/api/v1/analytics", 
        tags=["Analytics"],
        dependencies=[Depends(get_current_user_id)]
    )
    app.include_router(
        grpc_client.router, 
        prefix="/api/v1/grpc", 
        tags=["gRPC Client"],
        dependencies=[Depends(get_current_user_id)]
    )

    # Health check (no auth required)
    @app.get("/health")
    async def health_check():
        return {"status": "ok", "service": "python-api"}

    @app.get("/")
    async def root():
        return {
            "name": settings.APP_NAME,
            "version": settings.APP_VERSION,
            "docs": "/docs",
        }

    return app


app = create_app()
