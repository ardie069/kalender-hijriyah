"""
Visibility evaluation — Evaluasi visibilitas Hilal.

Dua path:
1. evaluate_visibility(): Full evaluation per-titik (untuk local_rukyat, local_hisab)
2. evaluate_visibility_geocentric(): Optimized untuk UGHC geocentric scan
   → Pre-computed geocentric positions dipakai ulang di seluruh grid
"""
from .criteria_registry import CRITERIA_REGISTRY
from skyfield.api import wgs84
import math

def get_geocentric_altaz(app_pos, t, lat, lon):
    """
    Hitung Geocentric Altitude dan Azimuth secara matematis.
    Skyfield altaz() menolak observer geocentric karena horizon adalah konsep topocentric.
    Namun KHGT geosentrik mendefinisikan altitude berbasis sudut dari geocentric zenith.
    """
    ra, dec, _ = app_pos.radec()
    
    gast = t.gast
    lst = gast + (lon / 15.0)
    ha_hours = lst - ra.hours
    ha_rad = math.radians(ha_hours * 15.0)
    
    lat_rad = math.radians(lat)
    dec_rad = dec.radians
    
    # Altitude
    sin_h = (math.sin(lat_rad) * math.sin(dec_rad)) + \
            (math.cos(lat_rad) * math.cos(dec_rad) * math.cos(ha_rad))
    h_rad = math.asin(max(-1.0, min(1.0, sin_h)))
    alt_deg = math.degrees(h_rad)

    return alt_deg

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
    """Full visibility evaluation — topocentric atau geocentric tergantung criteria."""

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
        alt_deg = get_geocentric_altaz(moon_app, t, lat, lon)
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


def evaluate_visibility_geocentric_batch(
    sunset_dt_utc,
    lat,
    lon,
    conjunction_jd,
    ts,
    sun,
    moon,
    earth,
    criteria_name,
    precomputed=None,
):
    """
    Optimized geocentric evaluation untuk UGHC global scan.

    Untuk observer geocentric (earth), posisi bulan/matahari SAMA di semua titik grid
    pada waktu yang sama. Yang berbeda hanya:
    1. Waktu sunset (berbeda per lokasi)
    2. Altitude refraksi (perlu lat/lon untuk altaz)

    precomputed: dict berisi {moon_app, sun_app, elongation_deg, moon_age_hours}
                 yang sudah dihitung di caller per-sunset-time.
                 Jika None, hitung dari awal (fallback ke path normal).
    """
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
    moon_age_hours = float((t.tt - conjunction_jd) * 24)

    if precomputed:
        moon_app = precomputed["moon_app"]
        sun_app = precomputed["sun_app"]
        elongation_deg = precomputed["elongation_deg"]
    else:
        # Fallback: hitung dari Geocentric observer
        observer = earth
        moon_app = observer.at(t).observe(moon).apparent()
        sun_app = observer.at(t).observe(sun).apparent()
        elongation_deg = float(moon_app.separation_from(sun_app).degrees)

    # Altitude matematis untuk geocentric observer
    alt_deg = get_geocentric_altaz(moon_app, t, lat, lon)

    if criteria_def.get("requires_conjunction") and moon_age_hours <= 0:
        return {
            "moon_altitude": alt_deg,
            "elongation": elongation_deg,
            "moon_age": moon_age_hours,
            "is_visible": False,
            "criteria_used": criteria_name,
        }

    if criteria_def["type"] == "numeric":
        is_visible = (
            alt_deg >= criteria_def["altitude_min"]
            and elongation_deg >= criteria_def["elongation_min"]
        )
    elif criteria_def["type"] == "wujud":
        is_visible = moon_age_hours > 0 and alt_deg > 0
    else:
        is_visible = False

    return {
        "moon_altitude": alt_deg,
        "elongation": elongation_deg,
        "moon_age": moon_age_hours,
        "is_visible": bool(is_visible),
        "criteria_used": criteria_name,
    }
