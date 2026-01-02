from ephemeris import ts, sun, moon
from astro_utils import angular_separation, get_moon_equatorial, get_sun_equatorial
from skyfield.api import wgs84

def get_moon_altitude(dt_utc, lat, lon):
    t = ts.utc(dt_utc.year, dt_utc.month, dt_utc.day, dt_utc.hour, dt_utc.minute, dt_utc.second)
    observer = wgs84.latlon(latitude_degrees=lat, longitude_degrees=lon)
    moon_apparent = observer.at(t).observe(moon).apparent()
    alt, _, _ = moon_apparent.altaz()
    return alt.degrees

def get_sun_altitude(dt_utc, lat, lon):
    t = ts.utc(dt_utc.year, dt_utc.month, dt_utc.day, dt_utc.hour, dt_utc.minute, dt_utc.second)
    observer = wgs84.latlon(latitude_degrees=lat, longitude_degrees=lon)
    sun_apparent = observer.at(t).observe(sun).apparent()
    alt, _, _ = sun_apparent.altaz()
    return alt.degrees

def evaluate_visibility(sunset_dt_utc, lat, lon, conjunction_jd):
    jd_sunset = sunset_dt_utc.timestamp() / 86400.0 + 2440587.5

    moon_pos = get_moon_equatorial(jd_sunset)
    sun_pos = get_sun_equatorial(jd_sunset)

    elongation = angular_separation(
        moon_pos.ra.degrees, moon_pos.dec.degrees,
        sun_pos.ra.degrees, sun_pos.dec.degrees
    )

    moon_alt = get_moon_altitude(sunset_dt_utc, lat, lon)
    moon_age_hours = (jd_sunset - conjunction_jd) * 24

    visibility = {
        "moon_altitude": moon_alt,
        "elongation": elongation,
        "moon_age": moon_age_hours,
        "is_visible": moon_alt >= 3 and elongation >= 6.4 and moon_age_hours >= 8,
    }

    return visibility
