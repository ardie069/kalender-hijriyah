from functools import lru_cache
from skyfield import almanac
from datetime import timedelta, timezone

_TS = None
_EARTH = None
_SUN = None
_MOON = None


def setup_conjunction_dependencies(ts, earth, sun, moon):
    global _TS, _EARTH, _SUN, _MOON
    _TS = ts
    _EARTH = earth
    _SUN = sun
    _MOON = moon


@lru_cache(maxsize=256)
def _cached_conjunction(jd_key: float) -> float:
    return _compute_conjunction(
        jd_key,
        _TS,
        _EARTH,
        _SUN,
        _MOON,
    )


def get_conjunction_time(
    jd_start,
    ts,
    earth,
    sun,
    moon,
):
    if _TS is None:
        setup_conjunction_dependencies(ts, earth, sun, moon)

    jd_key = round(jd_start, 3)
    return _cached_conjunction(jd_key)


def get_previous_conjunction(dt_utc, adapter):

    if dt_utc.tzinfo is None:
        raise ValueError("dt_utc must be timezone-aware")

    dt_utc = dt_utc.astimezone(timezone.utc)

    ts = adapter.ts
    eph = adapter.eph

    t1 = ts.from_datetime(dt_utc - timedelta(days=32))
    t2 = ts.from_datetime(dt_utc)

    f = almanac.moon_phases(eph)
    times, phases = almanac.find_discrete(t1, t2, f)

    new_moons = [t for t, p in zip(times, phases) if p == 0]

    if not new_moons:
        raise RuntimeError("No conjunction found in 32-day window")

    return new_moons[-1].utc_datetime()


def _get_diff_lon(jd, ts, earth, sun, moon):
    t = ts.tt(jd=jd)
    e = earth.at(t)
    _, slon, _ = e.observe(sun).apparent().ecliptic_latlon()
    _, mlon, _ = e.observe(moon).apparent().ecliptic_latlon()

    return (mlon.degrees - slon.degrees + 180) % 360 - 180


def _compute_conjunction(
    jd_start, ts, earth, sun, moon, max_days=2, step=0.1, precision=1e-6
):
    """
    Mencari waktu konjungsi (saat selisih bujur ekliptika Matahari dan Bulan nol).
    """
    jd = jd_start - max_days
    limit = jd_start + max_days

    last_diff = _get_diff_lon(jd, ts, earth, sun, moon)
    found_jd = jd

    while jd < limit:
        jd += step
        current_diff = _get_diff_lon(jd, ts, earth, sun, moon)

        if last_diff < 0 and current_diff >= 0:
            found_jd = jd
            break
        last_diff = current_diff

    left = found_jd - step
    right = found_jd

    while right - left > precision:
        mid = (left + right) / 2
        mid_diff = _get_diff_lon(mid, ts, earth, sun, moon)

        if mid_diff > 0:
            right = mid
        else:
            left = mid

    return (left + right) / 2
