from fastapi import APIRouter, Depends, Query
from datetime import datetime, timezone
from app.schemas.moon import MoonInfoResponse
from app.core.services.lunar_telemetry import LunarTelemetryService
from app.core.astronomy.skyfield_adapter import SkyfieldAdapter

router = APIRouter()


def get_lunar_service():
    adapter = SkyfieldAdapter(ephemeris_path="data/de421.bsp")
    return LunarTelemetryService(adapter)


@router.get("/moon/info", response_model=MoonInfoResponse)
async def get_moon_info(
    lat: float = Query(..., description="Latitude observer"),
    lon: float = Query(..., description="Longitude observer"),
    service: LunarTelemetryService = Depends(get_lunar_service),
):
    """
    Endpoint untuk mengambil data astronomi bulan lengkap untuk Dashboard.
    """
    # Secara Dialektika: Ambil waktu sekarang dalam UTC
    now = datetime.now(timezone.utc)

    # Eksekusi Logika Falak
    analysis = service.get_full_analysis(lat, lon, now)

    return analysis
