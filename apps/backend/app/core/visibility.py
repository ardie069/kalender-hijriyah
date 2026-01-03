from skyfield.api import wgs84 # type: ignore


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
    earth,
):
    """
    Evaluate lunar visibility using geometric criteria.
    All times are UTC calendar times.
    Astronomical calculations use Skyfield ephemeris (TT internally).
    """
    

    # Time object
    t = ts.utc(
        sunset_dt_utc.year,
        sunset_dt_utc.month,
        sunset_dt_utc.day,
        sunset_dt_utc.hour,
        sunset_dt_utc.minute,
        sunset_dt_utc.second,
    )

    # Topos (lokasi pengamat)
    topos = wgs84.latlon(
        latitude_degrees=lat,
        longitude_degrees=lon,
    )

    # Observer yang BENAR: earth + topos
    observer = earth + topos

    # Observasi benda langit (AMAN)
    moon_app = observer.at(t).observe(moon).apparent()
    sun_app = observer.at(t).observe(sun).apparent()

    # Altitude bulan
    alt_moon, _, _ = moon_app.altaz()

    # Elongasi
    elongation = moon_app.separation_from(sun_app).degrees

    # Usia bulan (jam)
    moon_age_hours = (t.tt - conjunction_jd) * 24

    return {
        "moon_altitude": alt_moon.degrees,
        "elongation": elongation,
        "moon_age": moon_age_hours,
        "is_visible": (
            alt_moon.degrees >= 3 and elongation >= 6.4 and moon_age_hours >= 8
        ),
    }
