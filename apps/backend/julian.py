# julian.py
from convertdate import islamic
from datetime import datetime
from skyfield.api import load
from ephemeris import ts

def jd_from_datetime(dt: datetime):
    """
    Menghitung Julian Date dari objek datetime.
    """
    t = ts.utc(dt.year, dt.month, dt.day, dt.hour, dt.minute, dt.second)
    return t.tt  # Julian Date (termasuk waktu universal)


def julian_to_hijri(jd_float: float):
    """
    Konversi Julian Date ke tanggal Hijriyah menggunakan convertdate.
    """
    return {
        "year": islamic.from_jd(jd_float)[0],
        "month": islamic.from_jd(jd_float)[1],
        "day": islamic.from_jd(jd_float)[2],
    }
