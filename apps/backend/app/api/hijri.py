from fastapi import APIRouter, Request, Query
from datetime import datetime
import pytz
from app.main import limiter
from app.core.hijri_calculator import get_hijri_date
from app.core.month_predictor import predict_end_of_month
from app.core.hijri_explain import explain_hijri_decision
from app.schemas.hijri import (
    HijriDateResponse,
    HijriDateSchema,
    LocationSchema,
    HijriMethod,
    HijriEndMonthResponse,
    HijriExplanationSchema,
)
from app.deps.astronomy import ts, eph, sun, moon, earth

router = APIRouter()


@router.get("/hijri-date", response_model=HijriDateResponse)
@limiter.limit("30/minute")
def hijri_date(
    lat: float,
    lon: float,
    method: HijriMethod,
    timezone: str,
    request: Request,
    now: str = Query(None),
):

    tz = pytz.timezone(timezone)

    if now:
        now_local = datetime.fromisoformat(now)
        if now_local.tzinfo is None:
            now_local = tz.localize(now_local)
    else:
        now_local = datetime.now(tz)

    hijri = get_hijri_date(
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

    explanation_dict = explain_hijri_decision(
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

    return HijriDateResponse(
        method=method,
        location=LocationSchema(lat=lat, lon=lon, timezone=timezone),
        hijri_date=HijriDateSchema(**hijri),
        explanation=explanation_dict,
        generated_at=now_local,
    )


@router.get("/hijri-end-month", response_model=HijriEndMonthResponse)
@limiter.limit("30/minute")
def hijri_predict_end(
    request: Request,
    lat: float,
    lon: float,
    method: HijriMethod,
    timezone: str,
):
    tz = pytz.timezone(timezone)
    now_local = datetime.now(tz)

    result = predict_end_of_month(
        lat, lon, method, timezone, ts=ts, eph=eph, sun=sun, moon=moon, earth=earth
    )

    return HijriEndMonthResponse(
        location=LocationSchema(lat=lat, lon=lon, timezone=timezone),
        generated_at=now_local,
        **result
    )
