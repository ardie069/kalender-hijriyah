"""
Sun times — sunset & fajr calculation.

Astronomy objects diakses dari AstronomyProvider singleton.
"""

from functools import lru_cache
from skyfield import almanac
from skyfield.api import wgs84
import pytz


def _get_astro():
    from app.deps.astronomy import get_provider
    return get_provider()


@lru_cache(maxsize=512)
def _cached_sunset(year: int, month: int, day: int, lat: float, lon: float, timezone: str):
    p = _get_astro()
    observer = wgs84.latlon(lat, lon)

    t0 = p.ts.utc(year, month, day, 0, 0)
    t1 = p.ts.utc(year, month, day, 23, 59)

    f = almanac.sunrise_sunset(p.eph, observer)
    times, events = almanac.find_discrete(t0, t1, f)

    for ti, ev in zip(times, events):
        if ev == 0:
            return (
                ti.utc_datetime()
                .replace(tzinfo=pytz.utc)
                .astimezone(pytz.timezone(timezone))
            )

    return None


def get_sunset_time(year_or_date, month_or_lat=None, day_or_lon=None,
                    lat_or_tz=None, lon_or_ts=None, timezone_or_eph=None,
                    ts=None, eph=None):
    """
    Flexible API:
      - get_sunset_time(year, month, day, lat, lon, timezone)  # new style
      - get_sunset_time(date, lat, lon, timezone, ts, eph)     # old style (ts/eph ignored)
    """
    # Detect old-style call: first arg is a date object
    if hasattr(year_or_date, 'year') and hasattr(year_or_date, 'month') and hasattr(year_or_date, 'day') and not isinstance(year_or_date, int):
        date = year_or_date
        lat = round(month_or_lat, 4)
        lon = round(day_or_lon, 4)
        tz = lat_or_tz
        return _cached_sunset(date.year, date.month, date.day, lat, lon, tz)

    # New-style call
    return _cached_sunset(
        year_or_date, month_or_lat, day_or_lon,
        round(lat_or_tz, 4), round(lon_or_ts, 4),
        timezone_or_eph,
    )


@lru_cache(maxsize=512)
def _cached_fajr(year: int, month: int, day: int, lat: float, lon: float, timezone: str):
    p = _get_astro()
    observer = wgs84.latlon(lat, lon)

    t0 = p.ts.utc(year, month, day, 0, 0)
    t1 = p.ts.utc(year, month, day, 12, 0)

    sun = p.eph["sun"]

    f = almanac.risings_and_settings(
        p.eph,
        sun,
        observer,
        horizon_degrees=-18.0,
    )

    times, events = almanac.find_discrete(t0, t1, f)

    for ti, ev in zip(times, events):
        if ev == 1:
            return (
                ti.utc_datetime()
                .replace(tzinfo=pytz.utc)
                .astimezone(pytz.timezone(timezone))
            )

    return None


def get_fajr_time(date, lat, lon, timezone, ts=None, eph=None):
    """ts and eph params kept for backward compat but are no longer used."""
    return _cached_fajr(
        date.year,
        date.month,
        date.day,
        round(lat, 4),
        round(lon, 4),
        timezone,
    )
