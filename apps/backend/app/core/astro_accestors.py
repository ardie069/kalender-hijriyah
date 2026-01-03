def get_moon_equatorial(jd, ts, moon):
    t = ts.tt(jd)
    return moon.at(t).apparent()


def get_sun_equatorial(jd, ts, sun):
    t = ts.tt(jd)
    return sun.at(t).apparent()
