from fastapi import APIRouter, Query
from datetime import datetime
import pytz

from app.core.methods.base import HijriContext
from app.core.services.rukyat_service import RukyatService
from app.schemas.rukyat import RukyatEvaluateResponse

router = APIRouter()

service = RukyatService()


@router.get("/rukyat/evaluate", response_model=RukyatEvaluateResponse)
async def evaluate_rukyat(
    lat: float = Query(...),
    lon: float = Query(...),
    timezone: str = Query("Asia/Jakarta"),
    region: str = Query("indonesia"),
    mode: str = Query("local"),
):
    tz = pytz.timezone(timezone)
    now_local = datetime.now(tz)
    context = HijriContext.from_request(lat, lon, timezone, now_local)
    return service.evaluate(context, region, mode)
