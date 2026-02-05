from functools import lru_cache

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


def _get_ra_dec(body, jd, ts, earth):
    t = ts.tt(jd=jd)
    astrometric = earth.at(t).observe(body).apparent()
    ra, dec, _ = astrometric.radec()
    return ra.hours, dec.degrees


def _compute_conjunction(
    jd_start,
    ts,
    earth,
    sun,
    moon,
    max_days=2,
    step=0.01,
    precision=1e-5,
):
    jd = jd_start
    limit = jd_start + max_days
    min_diff = float("inf")
    best_jd = None

    while jd < limit:
        ra_moon, _ = _get_ra_dec(moon, jd, ts, earth)
        ra_sun, _ = _get_ra_dec(sun, jd, ts, earth)

        diff = abs(ra_moon - ra_sun)
        diff = min(diff, 24 - diff)

        if diff < min_diff:
            min_diff = diff
            best_jd = jd

        jd += step

    left = best_jd - step
    right = best_jd + step

    while right - left > precision:
        mid = (left + right) / 2
        ra_moon, _ = _get_ra_dec(moon, mid, ts, earth)
        ra_sun, _ = _get_ra_dec(sun, mid, ts, earth)

        diff = abs(ra_moon - ra_sun)
        diff = min(diff, 24 - diff)

        if diff < min_diff:
            min_diff = diff
            best_jd = mid

        if ra_moon > ra_sun:
            right = mid
        else:
            left = mid

    return best_jd
