import logging
import time
import pytz
import concurrent.futures
import os
from datetime import datetime, timedelta, date
from functools import lru_cache

from app.core.methods.factory import get_method_instance
from ..config import AL_HARAM_LOCATION, HIJRI_MONTHS
from ..services.visibility_scan import GlobalVisibilityRegistry
from ..services.ughc_decision import evaluate_ughc_decision
from ..calendar.julian import jd_from_datetime

# ─── KUNCI: SMART CACHE SWITCHER ───
if os.getenv("VERCEL") == "1":
    from ..cache.redis import _make_key, get_cache, set_cache

    cache_type = "REDIS"
else:
    from ..cache.disk import _make_key, get_cache, set_cache

    cache_type = "DISK"

from app.core.services.engine import (
    calculate_sunset,
    calculate_baseline_hijri,
)

logger = logging.getLogger(__name__)

_YEAR_CACHE_TTL = 86400
_UGHC_CRITERIA = {"ughc": "TURKEY_2016_GEOCENTRIC"}
_PREDICTOR_LAT_STEP = 5
_PREDICTOR_LON_STEP = 10


class MonthPredictor:
    def __init__(self, factory):
        self.factory = factory

    def _is_ughc_method(self, method: str) -> bool:
        return method in _UGHC_CRITERIA

    def _predict_ughc_month(self, m_idx, day_29_date, method):
        criteria = _UGHC_CRITERIA[method]
        # Noon JD untuk NZ fajr check
        noon_dt = datetime.combine(day_29_date, datetime.min.time()).replace(
            tzinfo=pytz.utc
        ) + timedelta(hours=12)
        noon_jd = jd_from_datetime(noon_dt)

        is_new_month, scan = evaluate_ughc_decision(
            day_29_date,
            noon_jd,
            criteria,
            lat_step=_PREDICTOR_LAT_STEP,
            lon_step=_PREDICTOR_LON_STEP,
        )

        return (
            29 if is_new_month else 30,
            "new_month" if is_new_month else "istikmal_30",
            scan.get("best_visibility"),
        )

    def predict_full_year(self, hijri_year: int, lat, lon, timezone, method):
        t0 = time.perf_counter()

        if method == "umm_al_qura":
            lat, lon, timezone = (
                AL_HARAM_LOCATION["lat"],
                AL_HARAM_LOCATION["lon"],
                AL_HARAM_LOCATION["timezone"],
            )

        # 1. CEK CACHE (Redis/Disk otomatis terpilih)
        cache_key = _make_key(
            "year", hijri_year, method, round(lat, 2), round(lon, 2), timezone
        )

        # NOTE: Kita hapus ttl_seconds= di sini karena redis.get_cache pake **kwargs
        cached = get_cache(cache_key)
        if cached:
            logger.info(f"predict_full_year {cache_type} HIT: {hijri_year} {method}")
            return cached

        # 2. ESTIMASI AWAL TAHUN
        current_maghrib = self._estimate_start_of_hijri_year(
            hijri_year, lat, lon, timezone, method
        )

        # 3. PRE-FETCH PARALEL (Hanya untuk UGHC)
        if self._is_ughc_method(method):
            self._pre_fetch_ughc_scans(current_maghrib.date(), _UGHC_CRITERIA[method])

        method_instance = (
            None if self._is_ughc_method(method) else get_method_instance(method)
        )
        calendar_data = []

        # 4. LOOP SEQUENTIAL
        for m_idx in range(1, 13):
            day_29_date = current_maghrib.date() + timedelta(days=28)
            maghrib_29 = calculate_sunset(day_29_date, lat, lon, timezone)

            if not maghrib_29:
                total_days, decision, visibility_data = 30, "fallback", None
            elif self._is_ughc_method(method):
                total_days, decision, visibility_data = self._predict_ughc_month(
                    m_idx, day_29_date, method
                )
            else:
                sim_time = maghrib_29
                context = self.factory.create_context(lat, lon, timezone, sim_time)
                result = method_instance.calculate(context)
                result_month = result.hijri_date.get("month", m_idx)
                total_days = 29 if result_month != m_idx else 30
                decision = result.metadata.get("decision", "unknown")
                visibility_data = result.metadata.get(
                    "visibility"
                ) or result.metadata.get("visibility_data")

            calendar_data.append(
                {
                    "month_id": m_idx,
                    "month_name": HIJRI_MONTHS[m_idx],
                    "total_days": total_days,
                    "start_gregorian": current_maghrib.date().isoformat(),
                    "day_1_weekday": current_maghrib.weekday(),
                    "visibility": visibility_data,
                    "decision": decision,
                }
            )

            current_maghrib = calculate_sunset(
                current_maghrib.date() + timedelta(days=total_days), lat, lon, timezone
            )

        # SIMPAN KE CACHE
        set_cache(cache_key, calendar_data, ttl=_YEAR_CACHE_TTL)
        logger.info(
            f"predict_full_year DONE: {time.perf_counter() - t0:.3f}s ({cache_type} saved)"
        )
        return calendar_data

    def _pre_fetch_ughc_scans(self, start_date, criteria):
        potential_dates = []
        for i in range(12):
            base = start_date + timedelta(days=int(i * 29.5))
            potential_dates.extend(
                [
                    base + timedelta(days=28),
                    base + timedelta(days=29),
                    base + timedelta(days=30),
                ]
            )

        unique_dates = sorted(list(set(potential_dates)))
        with concurrent.futures.ThreadPoolExecutor(max_workers=8) as executor:
            executor.map(
                lambda d: GlobalVisibilityRegistry.scan_global(d, criteria),
                unique_dates,
            )

    def _estimate_start_of_hijri_year(self, year, lat, lon, tz, method):
        return _cached_estimate_start(year, lat, lon, tz)


@lru_cache(maxsize=64)
def _cached_estimate_start(year, lat, lon, tz):
    days_approx = int((year - 1) * 354.36707)
    base_date = date(622, 7, 19) + timedelta(days=days_approx)
    scan_start = base_date - timedelta(days=10)
    for i in range(20):
        test_date = scan_start + timedelta(days=i)
        sunset = calculate_sunset(test_date, lat, lon, tz)
        if not sunset:
            continue
        h_date, _ = calculate_baseline_hijri(sunset, tz, lat, lon)
        if h_date["month"] == 1 and h_date["day"] == 1:
            return sunset
    return calculate_sunset(base_date, lat, lon, tz)


# --- STANDALONE FUNCTION BUAT API ---
def predict_end_of_month(lat, lon, method, timezone, context_factory):
    """
    Prediksi akhir bulan Hijriyah (istikmal atau bulan baru).

    Sudah benar: menggunakan method_instance.calculate() untuk delegasi.
    """
    tz = pytz.timezone(timezone)
    now_local = datetime.now(tz)

    method_instance = get_method_instance(method)
    context_today = context_factory.create_context(lat, lon, timezone, now_local)
    result_today = method_instance.calculate(context_today)
    hijri_now = result_today.hijri_date

    days_to_29 = 29 - hijri_now["day"]
    simulation_local = now_local + timedelta(days=days_to_29)

    context_sim = context_factory.create_context(lat, lon, timezone, simulation_local)
    result_sim = method_instance.calculate(context_sim)

    is_istikmal = result_sim.metadata.get("decision") == "istikmal_30"

    return {
        "current_hijri": hijri_now,
        "prediction": {
            "is_istikmal": is_istikmal,
            "estimated_next_month_1": {
                "year": (
                    hijri_now["year"]
                    if hijri_now["month"] < 12
                    else hijri_now["year"] + 1
                ),
                "month": hijri_now["month"] + 1 if hijri_now["month"] < 12 else 1,
                "day": 1,
            },
        },
        "visibility_at_29": result_sim.metadata.get("visibility"),
        "message": "Prediksi berhasil",
    }
