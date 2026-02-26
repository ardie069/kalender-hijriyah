from fastapi import APIRouter, Depends, Query
from datetime import datetime, timezone
from app.schemas.moon import MoonInfoResponse
from app.core.services.lunar_telemetry import LunarTelemetryService
from app.core.astronomy.skyfield_adapter import SkyfieldAdapter
from app.core.config import REGIONAL_RUKYAT_CONFIG

router = APIRouter()


def get_lunar_service():
    adapter = SkyfieldAdapter(ephemeris_path="data/de421.bsp")
    return LunarTelemetryService(adapter)


@router.get("/moon/info", response_model=MoonInfoResponse)
async def get_moon_info(
    lat: float = Query(..., description="Latitude observer"),
    lon: float = Query(..., description="Longitude observer"),
    method: str = Query("local_rukyat"),
    region: str | None = Query(None),
    service: LunarTelemetryService = Depends(get_lunar_service),
):
    if method == "local_rukyat":
        if region and region in REGIONAL_RUKYAT_CONFIG:
            criteria = REGIONAL_RUKYAT_CONFIG[region]["criteria"]
        else:
            criteria = "MABIMS"
    else:
        criteria = "MABIMS"

    now = datetime.now(timezone.utc)

    analysis = service.get_full_analysis(lat, lon, now, criteria=criteria)

    return analysis
