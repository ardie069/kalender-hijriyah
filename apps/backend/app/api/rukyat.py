from fastapi import APIRouter, Query
from datetime import datetime
import pytz

from app.core.methods.base import HijriContext
from app.core.astronomy.skyfield_adapter import SkyfieldAdapter
from app.core.services.rukyat_service import RukyatService
from app.core.services.rukyat_national_service import RukyatNationalService

router = APIRouter()

adapter = SkyfieldAdapter(ephemeris_path="data/de421.bsp")


def create_context(lat: float, lon: float, timezone: str):
    tz = pytz.timezone(timezone)
    now_local = datetime.now(tz)

    return HijriContext(
        lat=lat,
        lon=lon,
        timezone=timezone,
        now_local=now_local,
        ts=adapter.ts,
        eph=adapter.eph,
        sun=adapter.sun,
        moon=adapter.moon,
        earth=adapter.earth,
    )


@router.get("/rukyat/evaluate")
async def evaluate_rukyat(
    lat: float = Query(...),
    lon: float = Query(...),
    timezone: str = Query("Asia/Jakarta"),
    region: str = Query("indonesia"),
):
    context = create_context(lat, lon, timezone)

    service = RukyatService()
    return service.evaluate_local_rukyat(context, region)


@router.get("/rukyat/national")
async def evaluate_rukyat_national(
    lat: float = Query(...),
    lon: float = Query(...),
    timezone: str = Query("Asia/Jakarta"),
    region: str = Query("indonesia"),
):
    context = create_context(lat, lon, timezone)

    service = RukyatNationalService()
    return service.evaluate(context, region)
