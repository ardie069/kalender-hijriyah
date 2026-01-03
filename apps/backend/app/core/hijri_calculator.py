import pytz  # type: ignore
from .julian import jd_from_datetime, julian_to_hijri
from .conjunction import get_conjunction_time
from .visibility import evaluate_visibility
from .sun_times import get_sunset_time
from .config import DEFAULT_LOCATION


def get_hijri_date(
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
    if now_local is None:
        raise ValueError("now_local must be a timezone-aware datetime")

    tz = pytz.timezone(timezone)
    now_local = now_local.astimezone(tz)

    # Lokasi referensi
    if method == "global":
        ref_lat, ref_lon, ref_zone = DEFAULT_LOCATION["global"]
    else:
        ref_lat, ref_lon, ref_zone = lat, lon, timezone

    # Sunset
    sunset_local = get_sunset_time(
        now_local.date(), ref_lat, ref_lon, ref_zone, ts, eph
    )
    if sunset_local is None:
        return {"error": "Gagal menghitung waktu matahari terbenam."}

    after_sunset = now_local >= sunset_local

    sunset_utc = sunset_local.astimezone(pytz.utc)
    sunset_jd = jd_from_datetime(sunset_utc, ts)

    # Hari hijriyah dimulai setelah maghrib
    effective_jd = sunset_jd + 1 if after_sunset else sunset_jd

    # Baseline kalender
    draft = julian_to_hijri(effective_jd)

    # GLOBAL: langsung selesai
    if method == "global":
        return draft

    # HISAB & RUKYAT: keputusan di hari ke-29
    if draft["day"] == 29:
        conj_jd = get_conjunction_time(effective_jd - 1, ts, earth, sun, moon)

        is_new_month = evaluate_new_month(
            method,
            lat,
            lon,
            conj_jd,
            sunset_utc,
            ts=ts,
            sun=sun,
            moon=moon,
            earth=earth,
        )

        if is_new_month:
            return _next_hijri_day(draft)

        return {**draft, "day": 30}

    return draft


def _next_hijri_day(date):
    year, month = date["year"], date["month"]

    if month == 12:
        return {"year": year + 1, "month": 1, "day": 1}

    return {"year": year, "month": month + 1, "day": 1}


def evaluate_new_month(method, lat, lon, conj_jd, sunset_utc, *, ts, sun, moon, earth):
    if method == "hisab":
        jd_sunset = jd_from_datetime(sunset_utc, ts)
        return conj_jd < jd_sunset

    if method == "rukyat":
        visibility = evaluate_visibility(
            sunset_utc, lat, lon, conj_jd, ts, sun, moon, earth
        )
        return visibility["is_visible"]

    return False
