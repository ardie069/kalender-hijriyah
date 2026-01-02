from ephemeris import ts, earth, sun, moon


def get_ra_dec(body, jd):
    """
    Ambil Right Ascension (jam) dan Declination (derajat)
    untuk benda langit (sun/moon) pada Julian Date tertentu.
    """
    t = ts.tt(jd=jd)
    astrometric = earth.at(t).observe(body).apparent()
    ra, dec, _ = astrometric.radec()
    return ra.hours, dec.degrees


def get_conjunction_time(jd_start, max_days=2, step=0.01, precision=1e-5):
    """
    Cari waktu konjungsi (ijtimak/new moon) terdekat setelah jd_start.

    Args:
        jd_start (float): Julian Date awal pencarian.
        max_days (int): Rentang hari pencarian setelah jd_start.
        step (float): Resolusi awal (dalam hari).
        precision (float): Akurasi akhir (dalam hari).

    Returns:
        float: Julian Date saat konjungsi.
    """
    jd = jd_start
    limit = jd_start + max_days
    min_diff = float("inf")
    best_jd = None

    # --- Pencarian kasar (linear scan) ---
    while jd < limit:
        ra_moon, _ = get_ra_dec(moon, jd)
        ra_sun, _ = get_ra_dec(sun, jd)
        diff = abs(ra_moon - ra_sun)

        if diff < min_diff:
            min_diff = diff
            best_jd = jd
        jd += step

    # --- Refinement (binary search) ---
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
