"""
Conjunction calculation — waktu konjungsi (ijtimak) bulan-matahari.

Astronomy objects diakses dari AstronomyProvider singleton.
"""

from skyfield import almanac
from datetime import timedelta, timezone
from functools import lru_cache


def _get_astro():
    from app.deps.astronomy import get_provider
    return get_provider()


@lru_cache(maxsize=256)
def _cached_conjunction(jd_key: float) -> float:
    p = _get_astro()
    return _compute_conjunction(jd_key, p.ts, p.earth, p.sun, p.moon)


def get_conjunction_time(jd_start, ts=None, earth=None, sun=None, moon=None):
    """
    Public API — menerima parameter opsional untuk backward compatibility.
    Kunci cache menggunakan JD yang dibulatkan.
    """
    jd_key = round(jd_start, 3)
    return _cached_conjunction(jd_key)


@lru_cache(maxsize=128)
def _get_previous_conjunction_cached(jd_bucket: float):
    """Cache berdasarkan bucket JD agar tidak bergantung bulan Gregorian."""
    p = _get_astro()
    dt_utc = p.ts.ut1_jd(jd_bucket).utc_datetime()

    t1 = p.ts.from_datetime(dt_utc - timedelta(days=32))
    t2 = p.ts.from_datetime(dt_utc)

    f = almanac.moon_phases(p.eph)
    times, phases = almanac.find_discrete(t1, t2, f)

    new_moons = [t for t, ph in zip(times, phases) if ph == 0]

    if not new_moons:
        raise RuntimeError("No conjunction found")

    return new_moons[-1].utc_datetime()


def get_previous_conjunction(dt_utc, adapter=None):
    """
    Get previous conjunction before dt_utc.
    adapter parameter kept for backward compat but no longer needed.
    """
    p = _get_astro()

    if dt_utc.tzinfo is None:
        raise ValueError("dt_utc must be timezone-aware")

    dt_utc = dt_utc.astimezone(timezone.utc)

    jd_now = p.ts.from_datetime(dt_utc).ut1
    jd_bucket = round(jd_now, 0)

    return _get_previous_conjunction_cached(jd_bucket)


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
