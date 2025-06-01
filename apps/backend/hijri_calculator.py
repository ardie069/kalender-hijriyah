from datetime import datetime, timedelta, timezone
from skyfield.api import load, N, W
from math import acos, cos, sin, pi
import pytz

# Load ephemeris untuk posisi planet/bulan/matahari
eph = load('data/de421.bsp')
ts = load.timescale()

def angular_separation(ra1, dec1, ra2, dec2):
    rad = pi / 180
    cos_theta = (
        sin(dec1 * rad) * sin(dec2 * rad) +
        cos(dec1 * rad) * cos(dec2 * rad) * cos((ra1 - ra2) * rad)
    )
    return acos(cos_theta) / rad

def get_julian_date(dt):
    return ts.utc(dt).utc_julian_date

def get_moon_position(jd):
    t = ts.tt_jd(jd)
    moon = eph['moon']
    earth = eph['earth']
    astrometric = earth.at(t).observe(moon).apparent()
    ra, dec, _ = astrometric.radec()
    return {'ra': ra.hours, 'dec': dec.degrees}

def get_sun_position(jd):
    t = ts.tt_jd(jd)
    sun = eph['sun']
    earth = eph['earth']
    astrometric = earth.at(t).observe(sun).apparent()
    ra, dec, _ = astrometric.radec()
    return {'ra': ra.hours, 'dec': dec.degrees}

def get_elongation(moon, sun):
    return angular_separation(moon['ra'], moon['dec'], sun['ra'], sun['dec'])

def julian_to_hijri(jd):
    jd_epoch = 1948439.5
    days_since_epoch = int(jd - jd_epoch)
    hijri_year = (30 * days_since_epoch + 10646) // 10631
    start_year_jd = jd_epoch + (10631 * hijri_year) // 30
    remaining_days = int(jd - start_year_jd)
    if remaining_days < 0:
        remaining_days += 354
    hijri_month = int(remaining_days / 29.5) + 1
    hijri_day = int(remaining_days - 29.5 * (hijri_month - 1)) + 1
    if hijri_month > 12:
        hijri_month = 1
        hijri_year += 1
    return {'year': hijri_year, 'month': hijri_month, 'day': hijri_day}
