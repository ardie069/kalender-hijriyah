import pytz # type: ignore
from .julian import jd_from_datetime, julian_to_hijri
from .conjunction import get_conjunction_time
from .visibility import evaluate_visibility
from .sun_times import get_sunset_time
from .config import DEFAULT_LOCATION



def explain_hijri_decision(
    lat,
    lon,
    method,
    timezone,
    *,
    now_local,
    ts,
    eph,
    sun,
    moon,
    earth,
):
    tz = pytz.timezone(timezone)
    now_local = now_local.astimezone(tz)

    # Lokasi referensi
    if method == "global":
        ref_lat, ref_lon, ref_zone = DEFAULT_LOCATION["global"]
    else:
        ref_lat, ref_lon, ref_zone = lat, lon, timezone

    sunset_local = get_sunset_time(
        now_local.date(), ref_lat, ref_lon, ref_zone, ts, eph
    )

    sunset_utc = sunset_local.astimezone(pytz.utc) # type: ignore
    after_sunset = now_local >= sunset_local

    sunset_jd = jd_from_datetime(sunset_utc, ts)
    effective_jd = sunset_jd + 1 if after_sunset else sunset_jd

    baseline = julian_to_hijri(effective_jd)

    explanation = {
        "method": method,
        "location": {
            "lat": lat,
            "lon": lon,
            "timezone": timezone,
        },
        "baseline_date": baseline,
        "after_sunset": after_sunset,
    }

    # Global tidak dievaluasi lokal
    if method == "global":
        explanation["decision"] = "global_standard"
        explanation["final_hijri_date"] = baseline
        return explanation

    if baseline["day"] != 29:
        explanation["decision"] = "no_evaluation_needed"
        explanation["final_hijri_date"] = baseline
        return explanation

    # Konjungsi
    conj_jd = get_conjunction_time(effective_jd - 1, ts, earth, sun, moon)
    explanation["conjunction"] = {
        "jd": conj_jd,
        "before_sunset": conj_jd < sunset_jd,
    }

    # Hisab
    if method == "hisab":
        if conj_jd < sunset_jd:
            explanation["decision"] = "new_month_by_hisab"
            explanation["final_hijri_date"] = {
                "year": baseline["year"],
                "month": baseline["month"] + 1 if baseline["month"] < 12 else 1,
                "day": 1,
            }
        else:
            explanation["decision"] = "istikmal"
            explanation["final_hijri_date"] = {**baseline, "day": 30}

        return explanation

    # Rukyat
    visibility = evaluate_visibility(
        sunset_utc, lat, lon, conj_jd, ts, sun, moon, earth
    )

    explanation["visibility"] = {
        "moon_altitude_deg": visibility["moon_altitude"],
        "elongation_deg": visibility["elongation"],
        "moon_age_hours": visibility["moon_age"],
        "criteria": {
            "min_altitude": 3,
            "min_elongation": 6.4,
            "min_age_hours": 8,
        },
        "is_visible": visibility["is_visible"],
    }

    if visibility["is_visible"]:
        explanation["decision"] = "new_month_by_rukyat"
        explanation["final_hijri_date"] = {
            "year": baseline["year"],
            "month": baseline["month"] + 1 if baseline["month"] < 12 else 1,
            "day": 1,
        }
    else:
        explanation["decision"] = "istikmal"
        explanation["final_hijri_date"] = {**baseline, "day": 30}

    return explanation
