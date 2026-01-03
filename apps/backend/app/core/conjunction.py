def get_ra_dec(body, jd, ts, earth):
    t = ts.tt(jd=jd)
    astrometric = earth.at(t).observe(body).apparent()
    ra, dec, _ = astrometric.radec()
    return ra.hours, dec.degrees


def get_conjunction_time(
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

    # linear scan
    while jd < limit:
        ra_moon, _ = get_ra_dec(moon, jd, ts, earth)
        ra_sun, _ = get_ra_dec(sun, jd, ts, earth)

        # handle RA wrap
        diff = abs(ra_moon - ra_sun)
        diff = min(diff, 24 - diff)

        if diff < min_diff:
            min_diff = diff
            best_jd = jd

        jd += step

    # refinement
    left = best_jd - step  # type: ignore
    right = best_jd + step  # type: ignore

    while right - left > precision:
        mid = (left + right) / 2
        ra_moon, _ = get_ra_dec(moon, mid, ts, earth)
        ra_sun, _ = get_ra_dec(sun, mid, ts, earth)

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
