import logging
import time

from fastapi import APIRouter, Request, Query, HTTPException
import pytz
import calendar as cal_mod
from datetime import datetime

logger = logging.getLogger(__name__)

from app.deps.rate_limit import limiter
from app.core.services.hijri_explain import explain_hijri_decision
from app.schemas.hijri import (
    HijriDateResponse,
    HijriDateSchema,
    LocationSchema,
    HijriMethod,
    HijriEndMonthResponse,
    GregorianMonthData,
    GregorianYearResponse,
)
from app.core.methods.factory import get_method_instance, ContextFactory
from app.core.services.month_predictor import MonthPredictor, predict_end_of_month

router = APIRouter()

GREGORIAN_MONTH_NAMES = [
    "", "Januari", "Februari", "Maret", "April",
    "Mei", "Juni", "Juli", "Agustus",
    "September", "Oktober", "November", "Desember",
]


def resolve_now(timezone: str, now: str | None) -> datetime:
    try:
        tz = pytz.timezone(timezone)
    except Exception:
        raise HTTPException(status_code=400, detail="Timezone invalid")

    if now:
        try:
            dt = datetime.fromisoformat(now)
            if dt.tzinfo is None:
                dt = tz.localize(dt)
            return dt
        except ValueError:
            raise HTTPException(status_code=400, detail="Invalid ISO date format")

    return datetime.now(tz)


def get_ready_factory():
    """
    Menghasilkan ContextFactory.
    Karena deps.astronomy pake lazy loading, kita nggak butuh get_provider().
    """
    return ContextFactory()


# ── ENDPOINTS ─────────────────────────────────────────────────────


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

    # Pake factory buat bikin konteks yang lengkap (ts, eph, dll)
    factory = get_ready_factory()
    context = factory.create_context(lat, lon, timezone, now_local)

    method_instance = get_method_instance(method)
    result = method_instance.calculate(context)

    explanation = result.explanation or explain_hijri_decision(method, result)

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
    factory = get_ready_factory()

    prediction = predict_end_of_month(
        lat=lat,
        lon=lon,
        method=method,
        timezone=timezone,
        context_factory=factory,
    )

    return HijriEndMonthResponse(
        location=LocationSchema(lat=lat, lon=lon, timezone=timezone),
        generated_at=datetime.now(pytz.timezone(timezone)),
        method=method,
        today=HijriDateSchema(**prediction["current_hijri"]),
        estimated_end_of_month=HijriDateSchema(
            **prediction["prediction"]["estimated_next_month_1"]
        ),
        visibility=prediction["visibility_at_29"],
        message=prediction.get("message", "Success"),
    )


@router.get("/calendar/year")
@limiter.limit("10/minute")
async def get_yearly_calendar(
    request: Request,
    year: int,
    lat: float,
    lon: float,
    method: HijriMethod,
    timezone: str = "Asia/Jakarta",
):
    """
    Endpoint utama buat nyuapin data 12 bulan ke Frontend Year View.
    """
    t0 = time.perf_counter()
    logger.info(
        "GET /calendar/year: year=%d method=%s lat=%.2f lon=%.2f",
        year, method, lat, lon,
    )

    factory = get_ready_factory()
    predictor = MonthPredictor(factory)

    data = predictor.predict_full_year(year, lat, lon, timezone, method)

    elapsed = time.perf_counter() - t0
    logger.info(
        "GET /calendar/year DONE: year=%d method=%s (%.3fs)",
        year, method, elapsed,
    )

    return {"year": year, "method": method, "months": data}


@router.get("/calendar/gregorian/year", response_model=GregorianYearResponse)
@limiter.limit("30/minute")
async def get_gregorian_calendar(
    request: Request,
    year: int,
):
    """
    Endpoint kalender Masehi — pure calendar, tanpa astronomi.
    Return 12 bulan dengan nama, total hari, dan weekday hari pertama.
    """
    months = []
    for m in range(1, 13):
        total_days = cal_mod.monthrange(year, m)[1]
        # weekday() → 0=Senin, 6=Ahad (Python convention)
        day_1_weekday = cal_mod.weekday(year, m, 1)
        months.append(
            GregorianMonthData(
                month_id=m,
                month_name=GREGORIAN_MONTH_NAMES[m],
                total_days=total_days,
                day_1_weekday=day_1_weekday,
            )
        )
    return GregorianYearResponse(year=year, months=months)
