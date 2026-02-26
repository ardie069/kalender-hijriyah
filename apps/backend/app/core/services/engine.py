"""
Engine — Cached computation functions for Hijri calendar.

All astronomy objects are accessed from AstronomyProvider singleton,
eliminating the _*_REGISTRY anti-pattern.
"""

import pytz
from datetime import datetime, time, timedelta
from functools import lru_cache
from skyfield.api import wgs84

from ..calendar.julian import jd_from_datetime, jd_to_datetime, julian_to_hijri
from ..astronomy.conjunction import get_conjunction_time
from ..astronomy.visibility import evaluate_visibility
from ..astronomy.sun_times import get_sunset_time


def _get_astro():
    """Get astronomy objects from provider. Lazy import to avoid circular deps."""
    from app.deps.astronomy import get_provider
    return get_provider()


def calculate_baseline_hijri(now_local, timezone, lat=None, lon=None):
    """
    Baseline Hijriyah berbasis noon,
    tapi mengikuti pergantian hari saat Maghrib.
    """
    p = _get_astro()
    tz = pytz.timezone(timezone)
    now_local = now_local.astimezone(tz)

    # Hitung sunset lokal dulu
    if lat is not None and lon is not None:
        sunset = calculate_sunset(now_local.date(), lat, lon, timezone)
    else:
        sunset = None

    after_sunset = sunset and now_local >= sunset

    # Jika sudah lewat Maghrib, geser ke hari Gregorian berikutnya
    target_date = now_local.date()
    if after_sunset:
        target_date = target_date + timedelta(days=1)

    noon_local = tz.localize(datetime.combine(target_date, time(12, 0)))
    noon_jd = jd_from_datetime(noon_local.astimezone(pytz.utc))

    return julian_to_hijri(noon_jd), noon_jd


@lru_cache(maxsize=512)
def _calculate_sunset_cached(year, month, day, lat, lon, timezone):
    return get_sunset_time(year, month, day, lat, lon, timezone)


def calculate_sunset(date, lat, lon, timezone):
    return _calculate_sunset_cached(
        date.year, date.month, date.day,
        round(lat, 4), round(lon, 4),
        timezone,
    )


@lru_cache(maxsize=256)
def _calculate_conjunction_cached(sunset_jd):
    p = _get_astro()
    return get_conjunction_time(sunset_jd, p.ts, p.earth, p.sun, p.moon)


def calculate_conjunction(sunset_jd):
    jd_key = round(sunset_jd, 3)
    return _calculate_conjunction_cached(jd_key)


def calculate_month_conjunction(now_local):
    p = _get_astro()
    now_utc = now_local.astimezone(pytz.utc)
    jd_key = round(jd_from_datetime(now_utc), 3)
    return _calculate_conjunction_cached(jd_key)


def calculate_visibility(sunset_utc, lat, lon, conj_jd, criteria):
    """Visibilitas Hilal — enriched with azimuth data."""
    p = _get_astro()

    vis = evaluate_visibility(
        sunset_utc, lat, lon, conj_jd, p.ts, p.sun, p.moon, p.earth,
        criteria=criteria,
    )

    observer = p.earth + wgs84.latlon(lat, lon)
    t = p.ts.from_datetime(sunset_utc)

    sun_pos = observer.at(t).observe(p.sun).apparent().altaz()
    moon_pos = observer.at(t).observe(p.moon).apparent().altaz()

    return {
        **vis,
        "azimuth_moon": round(moon_pos[1].degrees, 2),
        "azimuth_sun": round(sun_pos[1].degrees, 2),
        "azimuth_diff": round(moon_pos[1].degrees - sun_pos[1].degrees, 2),
    }


def check_historical_lag(
    current_baseline, noon_jd, lat, lon, timezone, criteria,
):
    p = _get_astro()
    days_since_start = current_baseline["day"] - 1
    month_start_noon_jd = noon_jd - days_since_start

    day_before_first_jd = month_start_noon_jd - 1.0
    day_before_first_utc = jd_to_datetime(day_before_first_jd)

    sunset_local = calculate_sunset(
        day_before_first_utc.date(), lat, lon, timezone,
    )

    if not sunset_local:
        return False

    sunset_utc = sunset_local.astimezone(pytz.utc)
    sunset_jd = jd_from_datetime(sunset_utc)

    conj_jd = calculate_conjunction(sunset_jd)

    vis = evaluate_visibility(
        sunset_utc, lat, lon, conj_jd, p.ts, p.sun, p.moon, p.earth,
        criteria=criteria,
    )

    return not vis["is_visible"]
