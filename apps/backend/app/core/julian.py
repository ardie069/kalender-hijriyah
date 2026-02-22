from convertdate import islamic
from functools import lru_cache

_TS_REGISTRY = {}


@lru_cache(maxsize=1024)
def _jd_from_datetime_cached(year, month, day, hour, minute, second, ts_id):
    ts = _TS_REGISTRY[ts_id]
    t = ts.utc(year, month, day, hour, minute, second)
    return float(t.ut1)


def jd_from_datetime(dt, ts):
    _TS_REGISTRY[id(ts)] = ts

    return _jd_from_datetime_cached(
        dt.year,
        dt.month,
        dt.day,
        dt.hour,
        dt.minute,
        dt.second,
        id(ts),
    )


@lru_cache(maxsize=1024)
def _jd_to_datetime_cached(jd, ts_id):
    ts = _TS_REGISTRY[ts_id]
    t = ts.ut1_jd(jd)
    return t.utc_datetime()


def jd_to_datetime(jd, ts):
    _TS_REGISTRY[id(ts)] = ts
    return _jd_to_datetime_cached(jd, id(ts))


@lru_cache(maxsize=2048)
def julian_to_hijri(jd_float: float):
    y, m, d = islamic.from_jd(jd_float)
    return {"year": y, "month": m, "day": d}
