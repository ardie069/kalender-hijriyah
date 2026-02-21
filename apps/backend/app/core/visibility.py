from skyfield.api import wgs84


def get_moon_altitude(dt_utc, lat, lon, ts, moon):
    """Mengambil tinggi hilal dalam derajat (standard float)."""
    t = ts.from_datetime(dt_utc)
    observer = wgs84.latlon(lat, lon)
    alt, _, _ = observer.at(t).observe(moon).apparent().altaz()
    return float(alt.degrees)


def evaluate_visibility(
    sunset_dt_utc,
    lat,
    lon,
    conjunction_jd,
    ts,
    sun,
    moon,
    earth,
    criteria="MABIMS",
):
    """
    Evaluasi visibilitas hilal berdasarkan kriteria tertentu.
    Mengonversi semua output NumPy ke tipe data Python standar.
    """
    t = ts.from_datetime(sunset_dt_utc)
    topos = wgs84.latlon(latitude_degrees=lat, longitude_degrees=lon)
    observer = earth + topos

    moon_app = observer.at(t).observe(moon).apparent()
    sun_app = observer.at(t).observe(sun).apparent()

    alt, _, _ = moon_app.altaz()
    alt_deg = float(alt.degrees)

    elongation_deg = float(moon_app.separation_from(sun_app).degrees)

    moon_age_hours = float((t.tt - conjunction_jd) * 24)

    if criteria == "Wujudul Hilal":
        is_visible = moon_age_hours > 0 and alt_deg > 0
    else:
        meets_geometry = alt_deg >= 3.0 and elongation_deg >= 6.4
        meets_age = moon_age_hours >= 8.0

        is_visible = meets_geometry and meets_age

    return {
        "moon_altitude": alt_deg,
        "elongation": elongation_deg,
        "moon_age": moon_age_hours,
        "is_visible": bool(is_visible),
        "criteria_used": str(criteria),
    }
