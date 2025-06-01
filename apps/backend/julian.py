# julian.py
from datetime import datetime
from skyfield.api import load
from ephemeris import ts

def jd_from_datetime(dt: datetime):
    """
    Menghitung Julian Date dari objek datetime.
    """
    t = ts.utc(dt.year, dt.month, dt.day, dt.hour, dt.minute, dt.second)
    return t.tt  # Julian Date (termasuk waktu universal)

def julian_to_hijri(jd: float):
    """
    Menghitung tanggal Hijriyah dari Julian Date.
    """
    # Konversi Julian Date ke Hijriyah
    # Rumus sederhana, bisa diganti dengan algoritma yang lebih akurat
    jd_hijri = jd - 1948440 + 10632
    year = int((jd_hijri - 1) // 354.367)
    month = int((jd_hijri - 1 - year * 354.367) // 29.5306) + 1
    day = int(jd_hijri - 1 - year * 354.367 - (month - 1) * 29.5306) + 1

    return {
        'year': year,
        'month': month,
        'day': day
    }