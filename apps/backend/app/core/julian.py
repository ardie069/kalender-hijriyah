from convertdate import islamic


def jd_from_datetime(dt, ts):
    """
    Convert UTC datetime to Julian Date (UTC Scale).
    """
    t = ts.utc(dt.year, dt.month, dt.day, dt.hour, dt.minute, dt.second)
    return t.ut1


def jd_to_datetime(jd, ts):
    """
    Convert Julian Date (UT1 Scale) to a timezone-aware UTC datetime object.
    """
    t = ts.ut1_jd(jd)
    return t.utc_datetime()


def julian_to_hijri(jd_float: float):
    y, m, d = islamic.from_jd(jd_float)
    return {"year": y, "month": m, "day": d}
