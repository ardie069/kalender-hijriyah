from .criteria_registry import CRITERIA_REGISTRY
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

    criteria = CRITERIA_REGISTRY.get(criteria)

    if not criteria:
        return {
            "moon_altitude": None,
            "elongation": None,
            "moon_age": None,
            "is_visible": False,
            "criteria_used": criteria,
        }

    t = ts.from_datetime(sunset_dt_utc)

    # Observer selection
    if criteria.get("observer") == "geocentric":
        observer = earth
    else:
        observer = earth + wgs84.latlon(lat, lon)

    moon_app = observer.at(t).observe(moon).apparent()
    sun_app = observer.at(t).observe(sun).apparent()

    if criteria.get("observer") == "geocentric":
        alt, _, _ = moon_app.altaz(wgs84.latlon(lat, lon))
    else:
        alt, _, _ = moon_app.altaz()

    alt_deg = float(alt.degrees)
    elongation_deg = float(moon_app.separation_from(sun_app).degrees)
    moon_age_hours = float((t.tt - conjunction_jd) * 24)

    # ======================
    # Decision Engine
    # ======================

    if criteria["type"] == "wujud":
        is_visible = moon_age_hours > 0 and alt_deg > 0

    elif criteria["type"] == "numeric":
        is_visible = (
            alt_deg >= criteria["altitude_min"]
            and elongation_deg >= criteria["elongation_min"]
        )

    else:
        is_visible = False

    return {
        "moon_altitude": alt_deg,
        "elongation": elongation_deg,
        "moon_age": moon_age_hours,
        "is_visible": bool(is_visible),
        "criteria_used": criteria,
    }
