from fastapi import APIRouter, Query, Request
from datetime import datetime
import pytz  # type: ignore
from app.main import limiter
from app.core.hijri_calculator import get_hijri_date
from app.core.month_predictor import predict_end_of_month
from app.core.hijri_explain import explain_hijri_decision
from app.schemas.hijri import (
    HijriDateResponse,
    HijriDateSchema,
    LocationSchema,
    HijriEndMonthResponse,
)
from app.deps.astronomy import ts, eph, sun, moon, earth

router = APIRouter()


@router.get(
    "/hijri-date",
    response_model=HijriDateResponse,
)
@limiter.limit("30/minute")
def hijri_date(
    lat: float,
    lon: float,
    method: str,
    timezone: str,
    request: Request,
):

    now_local = datetime.now(pytz.timezone(timezone))

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

    return HijriDateResponse(
        method=method,  # type: ignore
        location=LocationSchema(
            lat=lat,
            lon=lon,
            timezone=timezone,
        ),
        hijri_date=HijriDateSchema(**hijri),  # type: ignore
        generated_at=now_local,
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


@router.get("/hijri-end-month", response_model=HijriEndMonthResponse)
@limiter.limit("30/minute")
def hijri_predict_end(
    request: Request,
    lat: float,
    lon: float,
    method: str,
    timezone: str,
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

    return HijriEndMonthResponse(
        method=method,
        location=LocationSchema(
            lat=lat,
            lon=lon,
            timezone=timezone,
        ),
        generated_at=now_local,
        **result,
    )
