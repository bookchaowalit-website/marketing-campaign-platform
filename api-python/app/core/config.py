from pydantic_settings import BaseSettings
from typing import Optional


class Settings(BaseSettings):
    """Application settings"""

    # Application
    APP_NAME: str = "Marketing SaaS API (Python)"
    APP_VERSION: str = "1.0.0"
    DEBUG: bool = True
    LOG_LEVEL: str = "debug"

    # Database
    DATABASE_URL: str = "postgresql://postgres:postgres@localhost:5432/marketing_saas?sslmode=disable"

    # Server
    PORT: int = 8081
    HOST: str = "0.0.0.0"

    # Security
    SECRET_KEY: str = "your-secret-key-change-in-production"
    JWT_SECRET: str = "your-secret-key-change-in-production"
    ALGORITHM: str = "HS256"
    ACCESS_TOKEN_EXPIRE_MINUTES: int = 30

    # External Services
    GO_API_URL: str = "http://localhost:8080"
    GRPC_SERVER_URL: str = "localhost:9091"

    class Config:
        env_file = ".env"
        case_sensitive = True


settings = Settings()
