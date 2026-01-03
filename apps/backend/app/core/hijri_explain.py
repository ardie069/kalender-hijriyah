import pytz  # type: ignore
from datetime import datetime
from functools import lru_cache
from .julian import jd_from_datetime, julian_to_hijri
from .conjunction import get_conjunction_time
from .visibility import evaluate_visibility
from .sun_times import get_sunset_time
from .config import DEFAULT_LOCATION

_TS = None
_EPH = None
_SUN = None
_MOON = None
_EARTH = None


def setup_explain_dependencies(ts, eph, sun, moon, earth):
    global _TS, _EPH, _SUN, _MOON, _EARTH
    _TS = ts
    _EPH = eph
    _SUN = sun
    _MOON = moon
    _EARTH = earth


@lru_cache(maxsize=512)
def _cached_explain(
    year: int,
    month: int,
    day: int,
    lat: float,
    lon: float,
    method: str,
    timezone: str,
    after_sunset: bool,
):
    tz = pytz.timezone(timezone)
    hour = 20 if after_sunset else 10
    now_local = tz.localize(datetime(year, month, day, hour, 0))  # type: ignore

    return _explain_hijri_decision_internal(
        lat,
        lon,
        method,
        timezone,
        now_local=now_local,
        ts=_TS,
        eph=_EPH,
        sun=_SUN,
        moon=_MOON,
        earth=_EARTH,
    )


def _explain_hijri_decision_internal(
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

    sunset_jd = None

    if sunset_local is not None:
        sunset_utc = sunset_local.astimezone(pytz.utc)
        sunset_jd = jd_from_datetime(sunset_utc, ts)
        after_sunset = now_local >= sunset_local
        effective_jd = sunset_jd + 1 if after_sunset else sunset_jd
    else:
        effective_jd = jd_from_datetime(now_local.astimezone(pytz.utc), ts)
        after_sunset = False

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
        explanation["method_relation"] = _method_relation_label(method, baseline["day"])
        return explanation

    if baseline["day"] != 29:
        explanation["decision"] = "no_evaluation_needed"
        explanation["final_hijri_date"] = baseline
        explanation["method_relation"] = _method_relation_label(method, baseline["day"])
        return explanation

    # Konjungsi
    conj_jd = get_conjunction_time(effective_jd - 1, ts, earth, sun, moon)
    explanation["conjunction"] = {
        "jd": conj_jd,
        "before_sunset": conj_jd < sunset_jd,  # type: ignore
    }

    # Hisab
    if method == "hisab":
        if conj_jd < sunset_jd:  # type: ignore
            explanation["decision"] = "new_month_by_hisab"
            explanation["final_hijri_date"] = {
                "year": baseline["year"],
                "month": baseline["month"] + 1 if baseline["month"] < 12 else 1,
                "day": 1,
            }
        else:
            explanation["decision"] = "istikmal"
            explanation["final_hijri_date"] = {**baseline, "day": 30}

        explanation["method_relation"] = _method_relation_label(method, baseline["day"])
        return explanation

    if sunset_jd is None:
        explanation["decision"] = "sunset_not_available"
        explanation["final_hijri_date"] = baseline
        explanation["method_relation"] = _method_relation_label(method, baseline["day"])
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

    explanation["method_relation"] = _method_relation_label(method, baseline["day"])
    return explanation


def _method_relation_label(method, baseline_day):
    if method == "global":
        return {
            "status": "not_applicable",
            "note": "Metode global menggunakan kalender baku Mekkah.",
        }

    if baseline_day == 29:
        return {
            "status": "diverged",
            "note": (
                "Perbedaan metode dapat terjadi pada akhir bulan "
                "tergantung hasil hisab dan rukyat."
            ),
        }

    return {
        "status": "converged",
        "note": (
            "Perbedaan metode terjadi di awal bulan. "
            "Pada tanggal ini, hisab dan rukyat menghasilkan tanggal yang sama."
        ),
    }


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
    if _TS is None:
        setup_explain_dependencies(ts, eph, sun, moon, earth)

    tz = pytz.timezone(timezone)
    now_local = now_local.astimezone(tz)

    sunset_local = get_sunset_time(
        now_local.date(),
        lat if method != "global" else DEFAULT_LOCATION["global"][0],
        lon if method != "global" else DEFAULT_LOCATION["global"][1],
        timezone if method != "global" else DEFAULT_LOCATION["global"][2],
        ts,
        eph,
    )

    after_sunset = sunset_local is not None and now_local >= sunset_local

    return _cached_explain(
        now_local.year,
        now_local.month,
        now_local.day,
        round(lat, 4),
        round(lon, 4),
        method,
        timezone,
        after_sunset,
    )
