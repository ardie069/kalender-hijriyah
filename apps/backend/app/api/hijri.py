from fastapi import APIRouter, Query
from datetime import datetime
import pytz # type: ignore

from app.core.hijri_calculator import get_hijri_date
from app.core.month_predictor import predict_end_of_month
from app.core.hijri_explain import explain_hijri_decision
from app.deps.astronomy import ts, eph, sun, moon, earth

router = APIRouter()


@router.get("/hijri-date")
def hijri_date(
    lat: float = Query(...),
    lon: float = Query(...),
    method: str = Query("global"),
    timezone: str = Query("UTC"),
):
    return get_hijri_date(
        lat, lon, method, timezone, now_local=now_local, ts=ts, eph=eph, sun=sun, moon=moon, earth=earth # type: ignore
    )


@router.get("/hijri-explain")
def hijri_explain(
    lat: float,
    lon: float,
    method: str,
    timezone: str,
    now: str,
):
    now_local = datetime.fromisoformat(now).replace(tzinfo=pytz.utc)

    return explain_hijri_decision(
        lat=lat,
        lon=lon,
        method=method,
        timezone=timezone,
        now_local=now_local,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )
