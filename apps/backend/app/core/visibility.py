from skyfield.api import wgs84


def evaluate_visibility(
    sunset_dt_utc,
    lat,
    lon,
    conjunction_jd,
    ts,
    sun,
    moon,
    earth,
    criteria,
):
    t = ts.from_datetime(sunset_dt_utc)

    # ============================
    # OBSERVER SELECTION
    # ============================
    if criteria == "TURKEY_2016_GEOCENTRIC":
        observer = earth
    else:
        topos = wgs84.latlon(lat, lon)
        observer = earth + topos

    moon_app = observer.at(t).observe(moon).apparent()
    sun_app = observer.at(t).observe(sun).apparent()

    # ============================
    # ALTITUDE CALCULATION
    # ============================
    if criteria == "TURKEY_2016_GEOCENTRIC":
        topos = wgs84.latlon(lat, lon)
        alt, _, _ = moon_app.altaz(topos)
    else:
        alt, _, _ = moon_app.altaz()

    alt_deg = float(alt.degrees)

    elongation_deg = float(moon_app.separation_from(sun_app).degrees)

    moon_age_hours = float((t.tt - conjunction_jd) * 24)

    # ============================
    # CRITERIA
    # ============================
    if criteria == "WUJUDUL_HILAL":
        is_visible = moon_age_hours > 0 and alt_deg > 0

    elif criteria == "MABIMS":
        is_visible = alt_deg >= 3.0 and elongation_deg >= 6.4

    elif criteria in (
        "TURKEY_2016_GEOCENTRIC",
        "TURKEY_2016_TOPOCENTRIC",
    ):
        is_visible = alt_deg >= 5.0 and elongation_deg >= 8.0

    else:
        is_visible = False

    return {
        "moon_altitude": alt_deg,
        "elongation": elongation_deg,
        "moon_age": moon_age_hours,
        "is_visible": bool(is_visible),
        "criteria_used": criteria,
    }
