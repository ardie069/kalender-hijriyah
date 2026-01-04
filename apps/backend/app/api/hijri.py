from fastapi import APIRouter, Query, Request
from datetime import datetime
import pytz  # type: ignore
from app.main import limiter
from app.core.hijri_calculator import get_hijri_date
from app.core.month_predictor import predict_end_of_month
from app.core.hijri_explain import explain_hijri_decision
from app.deps.astronomy import ts, eph, sun, moon, earth

router = APIRouter()


@router.get("/hijri-date")
@limiter.limit("30/minute")
def hijri_date(
    request: Request,
    lat: float = Query(...),
    lon: float = Query(...),
    method: str = Query("global"),
    timezone: str = Query("UTC"),
):

    now_local = datetime.now(pytz.timezone(timezone))

    return get_hijri_date(
        lat,
        lon,
        method,
        timezone,
        now_local=now_local,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )


@router.get("/hijri-explain")
@limiter.limit("10/minute")
def hijri_explain(
    request: Request,
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


@router.get("/hijri-predict-end")
@limiter.limit("30/minute")
def hijri_predict_end(
    lat: float,
    lon: float,
    method: str = "hisab",
    timezone: str = "UTC",
):
    """
    Prediksi apakah bulan hijriyah berakhir di hari ke-29 atau ke-30.
    """
    tz = pytz.timezone(timezone)
    now_local = datetime.now(tz)

    result = predict_end_of_month(
        lat,
        lon,
        method,
        timezone,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    return {
        "method": method,
        "location": {
            "lat": lat,
            "lon": lon,
            "timezone": timezone,
        },
        "generated_at": now_local.isoformat(),
        **result,
    }
