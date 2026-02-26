"""
Julian Day & Hijri conversion utilities.

Astronomy objects diakses dari AstronomyProvider singleton.
"""

import pytz
from convertdate import islamic
from datetime import timedelta
from functools import lru_cache


def _get_ts():
    from app.deps.astronomy import get_provider
    return get_provider().ts


@lru_cache(maxsize=1024)
def _jd_from_datetime_cached(year, month, day, hour, minute, second):
    ts = _get_ts()
    t = ts.utc(year, month, day, hour, minute, second)
    return float(t.ut1)


def jd_from_datetime(dt, ts=None):
    """ts parameter kept for backward compat but is no longer used."""
    return _jd_from_datetime_cached(
        dt.year, dt.month, dt.day,
        dt.hour, dt.minute, dt.second,
    )


@lru_cache(maxsize=1024)
def _jd_to_datetime_cached(jd):
    ts = _get_ts()
    t = ts.ut1_jd(jd)
    return t.utc_datetime()


def jd_to_datetime(jd, ts=None):
    """ts parameter kept for backward compat but is no longer used."""
    return _jd_to_datetime_cached(jd)


@lru_cache(maxsize=2048)
def julian_to_hijri(jd_float: float):
    y, m, d = islamic.from_jd(jd_float)
    return {"year": y, "month": m, "day": d}


def add_days_to_datetime(dt, days):
    return dt + timedelta(days=days)


def hijri_to_jd(year: int, month: int, day: int):
    return islamic.to_jd(year, month, day)


def jd_to_local_datetime(jd, ts=None, timezone="UTC"):
    """ts parameter kept for backward compat but is no longer used."""
    utc_dt = jd_to_datetime(jd)
    tz = pytz.timezone(timezone)
    return utc_dt.astimezone(tz)
