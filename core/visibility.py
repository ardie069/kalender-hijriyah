from skyfield.api import wgs84
from .astro_utils import angular_separation
from .astro_accestors import get_moon_equatorial, get_sun_equatorial


def get_moon_altitude(dt_utc, lat, lon, ts, moon):
    t = ts.utc(
        dt_utc.year,
        dt_utc.month,
        dt_utc.day,
        dt_utc.hour,
        dt_utc.minute,
        dt_utc.second,
    )
    observer = wgs84.latlon(lat, lon)
    alt, _, _ = observer.at(t).observe(moon).apparent().altaz()  # type: ignore
    return alt.degrees


def evaluate_visibility(
    sunset_dt_utc,
    lat,
    lon,
    conjunction_jd,
    ts,
    sun,
    moon,
):
    jd_sunset = sunset_dt_utc.timestamp() / 86400.0 + 2440587.5

    moon_pos = get_moon_equatorial(jd_sunset, ts, moon)
    sun_pos = get_sun_equatorial(jd_sunset, ts, sun)

    elongation = angular_separation(
        moon_pos.ra.degrees,
        moon_pos.dec.degrees,
        sun_pos.ra.degrees,
        sun_pos.dec.degrees,
    )

    moon_alt = get_moon_altitude(sunset_dt_utc, lat, lon, ts, moon)
    moon_age = (jd_sunset - conjunction_jd) * 24

    return {
        "moon_altitude": moon_alt,
        "elongation": elongation,
        "moon_age": moon_age,
        "is_visible": moon_alt >= 3 and elongation >= 6.4 and moon_age >= 8, # type: ignore
    }
