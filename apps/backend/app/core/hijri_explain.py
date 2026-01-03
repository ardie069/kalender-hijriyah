import pytz  # type: ignore
from datetime import datetime
from typing import Any, Dict

from .julian import jd_from_datetime, julian_to_hijri
from .conjunction import get_conjunction_time
from .visibility import evaluate_visibility
from .sun_times import get_sunset_time
from .config import DEFAULT_LOCATION
from .hijri_calculator import get_hijri_date


def explain_hijri_decision(
    lat: float,
    lon: float,
    method: str,
    timezone: str,
    *,
    now_local: datetime,
    ts: Any,
    eph: Any,
    sun: Any,
    moon: Any,
    earth: Any,
) -> Dict[str, Any]:

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

    if sunset_local is None:
        return {
            "error": "sunset_not_available",
            "method": method,
        }

    after_sunset = now_local >= sunset_local
    sunset_utc = sunset_local.astimezone(pytz.utc)
    sunset_jd = jd_from_datetime(sunset_utc, ts)

    # Hari Hijriah dimulai setelah maghrib
    effective_jd = sunset_jd + 1 if after_sunset else sunset_jd

    # Baseline (tanpa keputusan metode)
    baseline = julian_to_hijri(effective_jd)

    # Final date dari engine utama
    final_date = get_hijri_date(
        lat,
        lon,
        method,
        timezone,
        now_local=now_local,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )

    explanation: Dict[str, Any] = {
        "method": method,
        "location": {"lat": lat, "lon": lon, "timezone": timezone},
        "after_sunset": after_sunset,
        "baseline_date": baseline,
        "final_hijri_date": final_date,
        "method_relation": _get_method_note(method),
    }

    # GLOBAL
    if method == "global":
        explanation["decision"] = "global_standard"
        return explanation

    # Evaluasi hanya di hari 29
    if baseline["day"] != 29:
        explanation["decision"] = "no_evaluation_needed"
        return explanation

    # Konjungsi dicari SEBELUM hari evaluasi
    conj_jd = get_conjunction_time(effective_jd - 1, ts, earth, sun, moon)
    explanation["conjunction"] = {
        "jd": conj_jd,
        "before_sunset": conj_jd < sunset_jd,
    }

    if method == "hisab":
        explanation["decision"] = (
            "new_month_by_hisab" if conj_jd < sunset_jd else "istikmal"
        )

    elif method == "rukyat":
        visibility = evaluate_visibility(
            sunset_utc, lat, lon, conj_jd, ts, sun, moon, earth
        )
        explanation["visibility"] = {
            "moon_altitude_deg": visibility["moon_altitude"],
            "elongation_deg": visibility["elongation"],
            "is_visible": visibility["is_visible"],
        }
        explanation["decision"] = (
            "new_month_by_rukyat" if visibility["is_visible"] else "istikmal"
        )

    return explanation


def _get_method_note(method: str) -> Dict[str, str]:
    return {
        "global": {
            "status": "not_applicable",
            "note": "Kalender baku berbasis standar Mekkah (Umm al-Qura-like).",
        },
        "hisab": {
            "status": "deterministic",
            "note": "Penentuan bulan berbasis konjungsi astronomis.",
        },
        "rukyat": {
            "status": "observational",
            "note": "Penentuan bulan berbasis visibilitas hilal.",
        },
    }.get(method, {"status": "unknown", "note": ""})
