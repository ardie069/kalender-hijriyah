from skyfield.api import load
from astro_utils import angular_separation
from ephemeris import eph, ts, earth, sun, moon
import numpy as np

def get_ra_dec(body, jd):
    t = ts.tt(jd=jd)
    astrometric = earth.at(t).observe(body).apparent()
    ra, dec, _ = astrometric.radec()
    return ra.hours, dec.degrees

def get_conjunction_time(jd_start, max_days=2, step=0.01, precision=1e-5):
    # Cari dalam rentang 2 hari dari jd_start
    jd = jd_start
    limit = jd_start + max_days
    min_diff = float('inf')
    best_jd = None

    # Pencarian kasar (linear)
    while jd < limit:
        ra_moon, dec_moon = get_ra_dec(moon, jd)
        ra_sun, dec_sun = get_ra_dec(sun, jd)
        diff = abs(ra_moon - ra_sun)

        if diff < min_diff:
            min_diff = diff
            best_jd = jd
        jd += step

    # Refinement dengan pencarian biner
    left = best_jd - step
    right = best_jd + step
    while right - left > precision:
        mid = (left + right) / 2
        ra_moon, _ = get_ra_dec(moon, mid)
        ra_sun, _ = get_ra_dec(sun, mid)
        diff = abs(ra_moon - ra_sun)

        if diff < min_diff:
            min_diff = diff
            best_jd = mid

        if ra_moon > ra_sun:
            right = mid
        else:
            left = mid

    return best_jd
