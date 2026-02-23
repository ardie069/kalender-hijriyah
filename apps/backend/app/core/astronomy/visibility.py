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

    criteria_name = criteria
    criteria_def = CRITERIA_REGISTRY.get(criteria_name)

    if not criteria_def:
        return {
            "moon_altitude": None,
            "elongation": None,
            "moon_age": None,
            "is_visible": False,
            "criteria_used": criteria_name,
        }

    t = ts.from_datetime(sunset_dt_utc)

    if criteria_def.get("observer") == "geocentric":
        observer = earth
    else:
        observer = earth + wgs84.latlon(lat, lon)

    moon_app = observer.at(t).observe(moon).apparent()
    sun_app = observer.at(t).observe(sun).apparent()

    if criteria_def.get("observer") == "geocentric":
        alt, _, _ = moon_app.altaz(wgs84.latlon(lat, lon))
    else:
        alt, _, _ = moon_app.altaz()

    alt_deg = float(alt.degrees)
    elongation_deg = float(moon_app.separation_from(sun_app).degrees)
    moon_age_hours = float((t.tt - conjunction_jd) * 24)

    if criteria_def.get("requires_conjunction") and moon_age_hours <= 0:
        return {
            "moon_altitude": alt_deg,
            "elongation": elongation_deg,
            "moon_age": moon_age_hours,
            "is_visible": False,
            "criteria_used": criteria_name,
        }

    if criteria_def["type"] == "wujud":
        is_visible = moon_age_hours > 0 and alt_deg > 0

    elif criteria_def["type"] == "numeric":
        is_visible = (
            alt_deg >= criteria_def["altitude_min"]
            and elongation_deg >= criteria_def["elongation_min"]
        )
    else:
        is_visible = False

    return {
        "moon_altitude": alt_deg,
        "elongation": elongation_deg,
        "moon_age": moon_age_hours,
        "is_visible": bool(is_visible),
        "criteria_used": criteria_name,
    }
