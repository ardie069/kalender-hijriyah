from functools import lru_cache
from skyfield import almanac
from skyfield.api import wgs84
import pytz

_TS = None
_EPH = None


def setup_sun_times_dependencies(ts, eph):
    global _TS, _EPH
    _TS = ts
    _EPH = eph


@lru_cache(maxsize=512)
def _cached_sunset(
    year: int,
    month: int,
    day: int,
    lat: float,
    lon: float,
    timezone: str,
):
    observer = wgs84.latlon(lat, lon)

    t0 = _TS.utc(year, month, day, 0, 0)
    t1 = _TS.utc(year, month, day, 23, 59)

    f = almanac.sunrise_sunset(_EPH, observer)
    times, events = almanac.find_discrete(t0, t1, f)

    for ti, ev in zip(times, events):
        if ev == 0:
            return (
                ti.utc_datetime()
                .replace(tzinfo=pytz.utc)
                .astimezone(pytz.timezone(timezone))
            )

    return None


def get_sunset_time(date, lat, lon, timezone, ts, eph):
    if _TS is None:
        setup_sun_times_dependencies(ts, eph)

    return _cached_sunset(
        date.year,
        date.month,
        date.day,
        round(lat, 4),
        round(lon, 4),
        timezone,
    )
