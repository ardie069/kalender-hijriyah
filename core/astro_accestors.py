def get_moon_equatorial(jd_utc, ts, moon):
    t = ts.tt(jd_utc)
    return moon.at(t).apparent()


def get_sun_equatorial(jd_utc, ts, sun):
    t = ts.tt(jd_utc)
    return sun.at(t).apparent()
