from fastapi import APIRouter, Request, Query, HTTPException
from datetime import datetime
import pytz
from app.deps.rate_limit import limiter
from app.core.services.hijri_explain import explain_hijri_decision
from app.schemas.hijri import (
    HijriDateResponse,
    HijriDateSchema,
    LocationSchema,
    HijriMethod,
    HijriEndMonthResponse,
)
from app.deps.astronomy import ts, eph, sun, moon, earth
from app.core.methods.factory import get_method_instance
from app.core.methods.base import HijriContext

router = APIRouter()


def resolve_now(timezone: str, now: str | None) -> datetime:
    try:
        tz = pytz.timezone(timezone)
    except Exception:
        raise HTTPException(status_code=400, detail="Timezone invalid")

    if now:
        dt = datetime.fromisoformat(now)

        if dt.tzinfo is None:
            dt = tz.localize(dt)

        return dt

    return datetime.now(tz)


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
    now_local = resolve_now(timezone, now)

    context = HijriContext(
        lat=lat,
        lon=lon,
        timezone=timezone,
        now_local=now_local,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    method_instance = get_method_instance(method)
    result = method_instance.calculate(context)

    explanation = result.explanation
    if not explanation:
        explanation = explain_hijri_decision(method, result)

    return HijriDateResponse(
        method=method,
        location=LocationSchema(lat=lat, lon=lon, timezone=timezone),
        hijri_date=HijriDateSchema(**result.hijri_date),
        explanation=explanation,
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
    now_local = resolve_now(timezone, None)

    context = HijriContext(
        lat=lat,
        lon=lon,
        timezone=timezone,
        now_local=now_local,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    method_instance = get_method_instance(method)
    result = method_instance.calculate(context)

    meta = result.metadata
    decision = meta.get("decision")

    if decision in ("new_month", "new_year"):
        message = "Kriteria terpenuhi. Besok diperkirakan masuk bulan baru."
    elif decision == "istikmal_30":
        message = "Hilal tidak memenuhi kriteria. Bulan digenapkan 30 hari."
    else:
        message = "Belum masuk fase evaluasi akhir bulan."

    return HijriEndMonthResponse(
        location=LocationSchema(lat=lat, lon=lon, timezone=timezone),
        generated_at=now_local,
        method=method,
        today=result.hijri_date,
        estimated_end_of_month=result.hijri_date,
        visibility=meta.get("visibility") or meta.get("visibility_data"),
        message=message,
    )
