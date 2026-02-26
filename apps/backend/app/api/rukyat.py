from fastapi import APIRouter, Query
from datetime import datetime
import pytz

from app.core.methods.base import HijriContext
from app.core.astronomy.skyfield_adapter import SkyfieldAdapter
from app.core.services.rukyat_service import RukyatService
from app.schemas.rukyat import RukyatEvaluateResponse

router = APIRouter()

adapter = SkyfieldAdapter(ephemeris_path="data/de421.bsp")
service = RukyatService()


def _build_context(lat: float, lon: float, timezone: str) -> HijriContext:
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


@router.get("/rukyat/evaluate", response_model=RukyatEvaluateResponse)
async def evaluate_rukyat(
    lat: float = Query(...),
    lon: float = Query(...),
    timezone: str = Query("Asia/Jakarta"),
    region: str = Query("indonesia"),
    mode: str = Query("local"),
):
    context = _build_context(lat, lon, timezone)
    return service.evaluate(context, region, mode)
