import pytz
from datetime import datetime, time, timedelta
from functools import lru_cache

from ..calendar.julian import jd_from_datetime, jd_to_datetime, julian_to_hijri
from ..astronomy.conjunction import get_conjunction_time
from ..astronomy.visibility import evaluate_visibility
from ..astronomy.sun_times import get_sunset_time

_TS_REGISTRY = {}
_EPH_REGISTRY = {}
_EARTH_REGISTRY = {}
_SUN_REGISTRY = {}
_MOON_REGISTRY = {}


def calculate_baseline_hijri(now_local, timezone, ts, lat=None, lon=None, eph=None):
    """
    Baseline Hijriyah berbasis noon,
    tapi mengikuti pergantian hari saat Maghrib.
    """

    tz = pytz.timezone(timezone)
    now_local = now_local.astimezone(tz)

    # Hitung sunset lokal dulu
    if lat is not None and lon is not None and eph is not None:
        sunset = calculate_sunset(now_local.date(), lat, lon, timezone, ts, eph)
    else:
        sunset = None

    after_sunset = sunset and now_local >= sunset

    # Jika sudah lewat Maghrib, geser ke hari Gregorian berikutnya
    target_date = now_local.date()
    if after_sunset:
        target_date = target_date + timedelta(days=1)

    noon_local = tz.localize(datetime.combine(target_date, time(12, 0)))
    noon_jd = jd_from_datetime(noon_local.astimezone(pytz.utc), ts)

    return julian_to_hijri(noon_jd), noon_jd


@lru_cache(maxsize=512)
def _calculate_sunset_cached(date, lat, lon, timezone, ts_id, eph_id):
    return get_sunset_time(
        date, lat, lon, timezone, _TS_REGISTRY[ts_id], _EPH_REGISTRY[eph_id]
    )


def calculate_sunset(date, lat, lon, timezone, ts, eph):
    _TS_REGISTRY[id(ts)] = ts
    _EPH_REGISTRY[id(eph)] = eph

    return _calculate_sunset_cached(date, lat, lon, timezone, id(ts), id(eph))


@lru_cache(maxsize=256)
def _calculate_conjunction_cached(sunset_jd, ts_id, earth_id, sun_id, moon_id):
    return get_conjunction_time(
        sunset_jd,
        _TS_REGISTRY[ts_id],
        _EARTH_REGISTRY[earth_id],
        _SUN_REGISTRY[sun_id],
        _MOON_REGISTRY[moon_id],
    )


def calculate_conjunction(sunset_jd, ts, earth, sun, moon):
    _TS_REGISTRY[id(ts)] = ts
    _EARTH_REGISTRY[id(earth)] = earth
    _SUN_REGISTRY[id(sun)] = sun
    _MOON_REGISTRY[id(moon)] = moon

    return _calculate_conjunction_cached(
        sunset_jd, id(ts), id(earth), id(sun), id(moon)
    )


@lru_cache(maxsize=128)
def _calculate_month_conjunction_cached(jd_key, ts_id, earth_id, sun_id, moon_id):
    ts = _TS_REGISTRY[ts_id]
    earth = _EARTH_REGISTRY[earth_id]
    sun = _SUN_REGISTRY[sun_id]
    moon = _MOON_REGISTRY[moon_id]

    return get_conjunction_time(jd_key, ts, earth, sun, moon)


def calculate_month_conjunction(now_local, ts, earth, sun, moon):
    now_utc = now_local.astimezone(pytz.utc)

    _TS_REGISTRY[id(ts)] = ts
    _EARTH_REGISTRY[id(earth)] = earth
    _SUN_REGISTRY[id(sun)] = sun
    _MOON_REGISTRY[id(moon)] = moon

    jd_key = jd_from_datetime(now_utc, ts)

    return _calculate_month_conjunction_cached(
        jd_key, id(ts), id(earth), id(sun), id(moon)
    )


def calculate_visibility(sunset_utc, lat, lon, conj_jd, ts, sun, moon, earth, criteria):
    """
    Visibilitas Hilal
    """
    return evaluate_visibility(
        sunset_utc, lat, lon, conj_jd, ts, sun, moon, earth, criteria=criteria
    )


def check_historical_lag(
    current_baseline, noon_jd, lat, lon, timezone, ts, eph, sun, moon, earth, criteria
):
    days_since_start = current_baseline["day"] - 1
    month_start_noon_jd = noon_jd - days_since_start

    day_before_first_jd = month_start_noon_jd - 1.0
    day_before_first_utc = jd_to_datetime(day_before_first_jd, ts)

    sunset_local = calculate_sunset(
        day_before_first_utc.date(), lat, lon, timezone, ts, eph
    )

    if not sunset_local:
        return False

    sunset_utc = sunset_local.astimezone(pytz.utc)
    sunset_jd = jd_from_datetime(sunset_utc, ts)

    conj_jd = calculate_conjunction(sunset_jd, ts, earth, sun, moon)

    vis = evaluate_visibility(
        sunset_utc, lat, lon, conj_jd, ts, sun, moon, earth, criteria=criteria
    )

    return not vis["is_visible"]
