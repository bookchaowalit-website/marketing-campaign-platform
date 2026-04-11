"""
SEO Data Sources - Real integrations for Google Search Console, Ahrefs, SEMrush
"""
import os
import logging
from typing import Dict, Any, List, Optional
from datetime import datetime, timedelta
from abc import ABC, abstractmethod
import httpx

logger = logging.getLogger(__name__)


class SEODatasourceError(Exception):
    """Custom exception for SEO data source errors"""
    pass


class SEODataSource(ABC):
    """Abstract base class for SEO data sources"""
    
    @abstractmethod
    async def get_domain_metrics(self, domain: str) -> Dict[str, Any]:
        pass
    
    @abstractmethod
    async def get_keywords(self, domain: str, limit: int = 50) -> List[Dict[str, Any]]:
        pass
    
    @abstractmethod
    async def get_backlinks(self, domain: str, limit: int = 50) -> List[Dict[str, Any]]:
        pass
    
    @abstractmethod
    async def get_ranking_positions(self, domain: str) -> List[Dict[str, Any]]:
        pass


class GoogleSearchConsoleDataSource(SEODataSource):
    """
    Google Search Console API integration
    Provides: organic traffic, impressions, clicks, positions, queries, pages
    """
    
    def __init__(self, credentials_path: Optional[str] = None, site_url: Optional[str] = None):
        self.credentials_path = credentials_path or os.getenv("GSC_CREDENTIALS_PATH")
        self.site_url = site_url or os.getenv("GSC_SITE_URL")
        self._service = None
    
    def _get_service(self):
        """Initialize Google Search Console service"""
        if self._service is None:
            try:
                from google.oauth2 import service_account
                from googleapiclient.discovery import build
                
                credentials = service_account.Credentials.from_service_account_file(
                    self.credentials_path,
                    scopes=["https://www.googleapis.com/auth/webmasters.readonly"]
                )
                self._service = build("searchconsole", "v1", credentials=credentials)
            except ImportError:
                raise SEODatasourceError("Install google-api-python-client: pip install google-api-python-client")
            except Exception as e:
                raise SEODatasourceError(f"Failed to initialize GSC: {e}")
        return self._service
    
    async def get_domain_metrics(self, domain: str) -> Dict[str, Any]:
        """Get domain metrics from GSC"""
        try:
            service = self._get_service()
            end_date = datetime.now()
            start_date = end_date - timedelta(days=30)
            
            request = {
                "startDate": start_date.strftime("%Y-%m-%d"),
                "endDate": end_date.strftime("%Y-%m-%d"),
                "dimensions": ["date"],
                "rowLimit": 30,
            }
            
            response = service.searchanalytics().query(
                siteUrl=f"sc-domain:{domain}",
                body=request
            ).execute()
            
            rows = response.get("rows", [])
            total_clicks = sum(r.get("clicks", 0) for r in rows)
            total_impressions = sum(r.get("impressions", 0) for r in rows)
            avg_position = sum(r.get("position", 0) for r in rows) / len(rows) if rows else 0
            avg_ctr = (total_clicks / total_impressions * 100) if total_impressions > 0 else 0
            
            return {
                "domain": domain,
                "source": "google_search_console",
                "organic_traffic": total_clicks,
                "impressions": total_impressions,
                "avg_position": round(avg_position, 2),
                "ctr": round(avg_ctr, 2),
                "period_days": 30,
            }
        except Exception as e:
            logger.error(f"GSC get_domain_metrics failed: {e}")
            return {"domain": domain, "source": "google_search_console", "error": str(e)}
    
    async def get_keywords(self, domain: str, limit: int = 50) -> List[Dict[str, Any]]:
        """Get top keywords from GSC"""
        try:
            service = self._get_service()
            end_date = datetime.now()
            start_date = end_date - timedelta(days=30)
            
            request = {
                "startDate": start_date.strftime("%Y-%m-%d"),
                "endDate": end_date.strftime("%Y-%m-%d"),
                "dimensions": ["query"],
                "rowLimit": limit,
                "sortTypes": ["clicks"],
            }
            
            response = service.searchanalytics().query(
                siteUrl=f"sc-domain:{domain}",
                body=request
            ).execute()
            
            return [
                {
                    "keyword": row["keys"][0],
                    "clicks": row.get("clicks", 0),
                    "impressions": row.get("impressions", 0),
                    "ctr": round(row.get("ctr", 0) * 100, 2),
                    "position": round(row.get("position", 0), 1),
                }
                for row in response.get("rows", [])
            ]
        except Exception as e:
            logger.error(f"GSC get_keywords failed: {e}")
            return []
    
    async def get_backlinks(self, domain: str, limit: int = 50) -> List[Dict[str, Any]]:
        """GSC doesn't provide backlink data - return empty"""
        logger.info("GSC does not provide backlink data")
        return []
    
    async def get_ranking_positions(self, domain: str) -> List[Dict[str, Any]]:
        """Get average positions for pages"""
        try:
            service = self._get_service()
            end_date = datetime.now()
            start_date = end_date - timedelta(days=30)
            
            request = {
                "startDate": start_date.strftime("%Y-%m-%d"),
                "endDate": end_date.strftime("%Y-%m-%d"),
                "dimensions": ["page"],
                "rowLimit": 50,
            }
            
            response = service.searchanalytics().query(
                siteUrl=f"sc-domain:{domain}",
                body=request
            ).execute()
            
            return [
                {
                    "url": row["keys"][0],
                    "clicks": row.get("clicks", 0),
                    "impressions": row.get("impressions", 0),
                    "avg_position": round(row.get("position", 0), 1),
                }
                for row in response.get("rows", [])
            ]
        except Exception as e:
            logger.error(f"GSC get_ranking_positions failed: {e}")
            return []


class AhrefsDataSource(SEODataSource):
    """
    Ahrefs API integration (v2)
    Provides: backlinks, domain rating, organic keywords, traffic, referring domains
    """
    
    def __init__(self, api_token: Optional[str] = None):
        self.api_token = api_token or os.getenv("AHREFS_API_TOKEN")
        self.base_url = "https://apiv2.ahrefs.com"
        if not self.api_token:
            logger.warning("Ahrefs API token not configured")
    
    async def _request(self, endpoint: str, params: Dict[str, Any]) -> Dict[str, Any]:
        """Make API request to Ahrefs"""
        async with httpx.AsyncClient(timeout=30.0) as client:
            params["token"] = self.api_token
            params["output"] = "json"
            response = await client.get(f"{self.base_url}{endpoint}", params=params)
            response.raise_for_status()
            return response.json()
    
    async def get_domain_metrics(self, domain: str) -> Dict[str, Any]:
        """Get domain metrics from Ahrefs"""
        if not self.api_token:
            return {"domain": domain, "source": "ahrefs", "error": "API token not configured"}
        
        try:
            data = await self._request("", {
                "from": "domain_rank",
                "target": domain,
                "limit": 1,
            })
            
            if "results" not in data or not data["results"]:
                return {"domain": domain, "source": "ahrefs", "error": "Domain not found"}
            
            result = data["results"][0]
            return {
                "domain": domain,
                "source": "ahrefs",
                "domain_rating": result.get("domain_rating", 0),
                "url_rating": result.get("url_rating", 0),
                "backlinks": result.get("backlinks", 0),
                "referring_domains": result.get("referring_domains", 0),
                "organic_keywords": result.get("organic_keywords", 0),
                "organic_traffic": result.get("organic_traffic", 0),
                "traffic_value": result.get("traffic_value", 0),
            }
        except Exception as e:
            logger.error(f"Ahrefs get_domain_metrics failed: {e}")
            return {"domain": domain, "source": "ahrefs", "error": str(e)}
    
    async def get_keywords(self, domain: str, limit: int = 50) -> List[Dict[str, Any]]:
        """Get top organic keywords from Ahrefs"""
        if not self.api_token:
            return []
        
        try:
            data = await self._request("", {
                "from": "organic2",
                "target": domain,
                "limit": limit,
                "order_by": "traffic:desc",
            })
            
            return [
                {
                    "keyword": row.get("keyword", ""),
                    "position": row.get("position", 0),
                    "volume": row.get("volume", 0),
                    "traffic": row.get("traffic", 0),
                    "difficulty": row.get("keyword_difficulty", 0),
                    "cpc": row.get("cpc", 0),
                    "url": row.get("url", ""),
                }
                for row in data.get("results", [])
            ]
        except Exception as e:
            logger.error(f"Ahrefs get_keywords failed: {e}")
            return []
    
    async def get_backlinks(self, domain: str, limit: int = 50) -> List[Dict[str, Any]]:
        """Get backlinks from Ahrefs"""
        if not self.api_token:
            return []
        
        try:
            data = await self._request("", {
                "from": "backlinks",
                "target": domain,
                "limit": limit,
                "order_by": "domain_rating:desc",
            })
            
            return [
                {
                    "source_url": row.get("url_from", ""),
                    "target_url": row.get("url_to", ""),
                    "anchor_text": row.get("anchor_text", ""),
                    "domain_rating": row.get("domain_rating", 0),
                    "follow": row.get("link_type") == "follow",
                    "first_seen": row.get("first_seen", ""),
                }
                for row in data.get("results", [])
            ]
        except Exception as e:
            logger.error(f"Ahrefs get_backlinks failed: {e}")
            return []
    
    async def get_ranking_positions(self, domain: str) -> List[Dict[str, Any]]:
        """Get top ranking URLs"""
        if not self.api_token:
            return []
        
        try:
            data = await self._request("", {
                "from": "pages",
                "target": domain,
                "limit": 50,
                "order_by": "traffic:desc",
            })
            
            return [
                {
                    "url": row.get("url", ""),
                    "traffic": row.get("traffic", 0),
                    "keywords": row.get("keywords", 0),
                    "domain_rating": row.get("domain_rating", 0),
                }
                for row in data.get("results", [])
            ]
        except Exception as e:
            logger.error(f"Ahrefs get_ranking_positions failed: {e}")
            return []


class SEMrushDataSource(SEODataSource):
    """
    SEMrush API integration
    Provides: organic research, positions, domain overview, backlinks
    """
    
    def __init__(self, api_key: Optional[str] = None):
        self.api_key = api_key or os.getenv("SEMRUSH_API_KEY")
        self.base_url = "https://api.semrush.com"
        if not self.api_key:
            logger.warning("SEMrush API key not configured")
    
    async def _request(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Make API request to SEMrush"""
        params["key"] = self.api_key
        async with httpx.AsyncClient(timeout=30.0) as client:
            response = await client.get(self.base_url, params=params)
            response.raise_for_status()
            # SEMrush returns CSV by default, parse it
            content = response.text
            lines = content.strip().split("\n")
            if len(lines) < 2:
                return {"headers": [], "data": []}
            headers = lines[0].split(";")
            data = []
            for line in lines[1:]:
                values = line.split(";")
                data.append(dict(zip(headers, values)))
            return {"headers": headers, "data": data}
    
    async def get_domain_metrics(self, domain: str) -> Dict[str, Any]:
        """Get domain overview from SEMrush"""
        if not self.api_key:
            return {"domain": domain, "source": "semrush", "error": "API key not configured"}
        
        try:
            data = await self._request({
                "type": "domain_ranks",
                "target": domain,
                "export_columns": "OrganicKeywords,OrganicTraffic,OrganicTrafficPrice,DomainAuthority,Backlinks",
                "database": "us",
            })
            
            if not data.get("data"):
                return {"domain": domain, "source": "semrush", "error": "Domain not found"}
            
            row = data["data"][0]
            return {
                "domain": domain,
                "source": "semrush",
                "organic_keywords": int(row.get("OrganicKeywords", 0)),
                "organic_traffic": int(row.get("OrganicTraffic", 0)),
                "traffic_value": float(row.get("OrganicTrafficPrice", 0)),
                "domain_authority": int(row.get("DomainAuthority", 0)),
                "backlinks": int(row.get("Backlinks", 0)),
            }
        except Exception as e:
            logger.error(f"SEMrush get_domain_metrics failed: {e}")
            return {"domain": domain, "source": "semrush", "error": str(e)}
    
    async def get_keywords(self, domain: str, limit: int = 50) -> List[Dict[str, Any]]:
        """Get organic keywords from SEMrush"""
        if not self.api_key:
            return []
        
        try:
            data = await self._request({
                "type": "phrase_organic",
                "target": domain,
                "display_limit": limit,
                "export_columns": "Phrase,Position,SearchVolume,Cpc,Competition,Traffic,TrafficCost,Difficulty",
                "database": "us",
            })
            
            return [
                {
                    "keyword": row.get("Phrase", ""),
                    "position": int(row.get("Position", 0)),
                    "volume": int(row.get("SearchVolume", 0)),
                    "cpc": float(row.get("Cpc", 0)),
                    "competition": row.get("Competition", ""),
                    "traffic": int(row.get("Traffic", 0)),
                    "difficulty": int(row.get("Difficulty", 0)),
                }
                for row in data.get("data", [])
            ]
        except Exception as e:
            logger.error(f"SEMrush get_keywords failed: {e}")
            return []
    
    async def get_backlinks(self, domain: str, limit: int = 50) -> List[Dict[str, Any]]:
        """Get backlinks from SEMrush"""
        if not self.api_key:
            return []
        
        try:
            data = await self._request({
                "type": "backlinks",
                "target": domain,
                "display_limit": limit,
                "export_columns": "BacklinkUrl,BacklinkDomain,Anchor,LinkType,FirstSeen,DomainAuthority",
                "database": "us",
            })
            
            return [
                {
                    "source_url": row.get("BacklinkUrl", ""),
                    "source_domain": row.get("BacklinkDomain", ""),
                    "anchor_text": row.get("Anchor", ""),
                    "follow": "follow" in row.get("LinkType", "").lower(),
                    "first_seen": row.get("FirstSeen", ""),
                    "domain_authority": int(row.get("DomainAuthority", 0)),
                }
                for row in data.get("data", [])
            ]
        except Exception as e:
            logger.error(f"SEMrush get_backlinks failed: {e}")
            return []
    
    async def get_ranking_positions(self, domain: str) -> List[Dict[str, Any]]:
        """Get top pages from SEMrush"""
        if not self.api_key:
            return []
        
        try:
            data = await self._request({
                "type": "pages",
                "target": domain,
                "display_limit": 50,
                "export_columns": "Url,OrganicKeywords,OrganicTraffic,TrafficCost",
                "database": "us",
            })
            
            return [
                {
                    "url": row.get("Url", ""),
                    "keywords": int(row.get("OrganicKeywords", 0)),
                    "traffic": int(row.get("OrganicTraffic", 0)),
                    "traffic_value": float(row.get("TrafficCost", 0)),
                }
                for row in data.get("data", [])
            ]
        except Exception as e:
            logger.error(f"SEMrush get_ranking_positions failed: {e}")
            return []


class SEODataService:
    """
    Unified SEO data service that aggregates multiple sources
    """
    
    def __init__(self):
        self.sources: Dict[str, SEODataSource] = {}
        self._initialize_sources()
    
    def _initialize_sources(self):
        """Initialize configured SEO data sources"""
        # Google Search Console
        if os.getenv("GSC_CREDENTIALS_PATH"):
            try:
                self.sources["google_search_console"] = GoogleSearchConsoleDataSource()
                logger.info("Google Search Console source initialized")
            except Exception as e:
                logger.warning(f"Failed to initialize GSC: {e}")
        
        # Ahrefs
        if os.getenv("AHREFS_API_TOKEN"):
            try:
                self.sources["ahrefs"] = AhrefsDataSource()
                logger.info("Ahrefs source initialized")
            except Exception as e:
                logger.warning(f"Failed to initialize Ahrefs: {e}")
        
        # SEMrush
        if os.getenv("SEMRUSH_API_KEY"):
            try:
                self.sources["semrush"] = SEMrushDataSource()
                logger.info("SEMrush source initialized")
            except Exception as e:
                logger.warning(f"Failed to initialize SEMrush: {e}")
        
        if not self.sources:
            logger.warning("No SEO data sources configured")
    
    async def get_domain_metrics(self, domain: str, source: Optional[str] = None) -> Dict[str, Any]:
        """Get domain metrics from specified source or all sources"""
        if source:
            if source not in self.sources:
                return {"error": f"Source '{source}' not available"}
            return await self.sources[source].get_domain_metrics(domain)
        
        # Aggregate from all sources
        results = {}
        for name, src in self.sources.items():
            results[name] = await src.get_domain_metrics(domain)
        return results
    
    async def get_keywords(self, domain: str, source: Optional[str] = None, limit: int = 50) -> Dict[str, List[Dict[str, Any]]]:
        """Get keywords from specified source or all sources"""
        if source:
            if source not in self.sources:
                return {"error": f"Source '{source}' not available"}
            return {source: await self.sources[source].get_keywords(domain, limit)}
        
        results = {}
        for name, src in self.sources.items():
            results[name] = await src.get_keywords(domain, limit)
        return results
    
    async def get_backlinks(self, domain: str, source: Optional[str] = None, limit: int = 50) -> Dict[str, List[Dict[str, Any]]]:
        """Get backlinks from specified source or all sources"""
        if source:
            if source not in self.sources:
                return {"error": f"Source '{source}' not available"}
            return {source: await self.sources[source].get_backlinks(domain, limit)}
        
        results = {}
        for name, src in self.sources.items():
            results[name] = await src.get_backlinks(domain, limit)
        return results
    
    async def get_ranking_positions(self, domain: str, source: Optional[str] = None) -> Dict[str, List[Dict[str, Any]]]:
        """Get ranking positions from specified source or all sources"""
        if source:
            if source not in self.sources:
                return {"error": f"Source '{source}' not available"}
            return {source: await self.sources[source].get_ranking_positions(domain)}
        
        results = {}
        for name, src in self.sources.items():
            results[name] = await src.get_ranking_positions(domain)
        return results
    
    def get_available_sources(self) -> List[str]:
        """List configured sources"""
        return list(self.sources.keys())


# Singleton instance
_seo_data_service: Optional[SEODataService] = None


def get_seo_service() -> SEODataService:
    """Get or create SEO data service singleton"""
    global _seo_data_service
    if _seo_data_service is None:
        _seo_data_service = SEODataService()
    return _seo_data_service
