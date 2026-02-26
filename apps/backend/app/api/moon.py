from fastapi import APIRouter, Query
from datetime import datetime, timezone
from app.schemas.moon import MoonInfoResponse
from app.core.services.lunar_telemetry import LunarTelemetryService
from app.core.astronomy.skyfield_adapter import SkyfieldAdapter
from app.core.config import REGIONAL_RUKYAT_CONFIG

router = APIRouter()

# Singleton — SkyfieldAdapter sekarang otomatis menggunakan AstronomyProvider
_service = LunarTelemetryService(SkyfieldAdapter())


@router.get("/moon/info", response_model=MoonInfoResponse)
async def get_moon_info(
    lat: float = Query(..., description="Latitude observer"),
    lon: float = Query(..., description="Longitude observer"),
    method: str = Query("local_rukyat"),
    region: str | None = Query(None),
):
    if method == "local_rukyat":
        if region and region in REGIONAL_RUKYAT_CONFIG:
            criteria = REGIONAL_RUKYAT_CONFIG[region]["criteria"]
        else:
            criteria = "MABIMS"
    else:
        criteria = "MABIMS"

    now = datetime.now(timezone.utc)

    analysis = _service.get_full_analysis(lat, lon, now, criteria=criteria)

    return analysis
