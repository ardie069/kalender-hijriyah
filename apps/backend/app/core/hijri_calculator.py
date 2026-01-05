import pytz  # type: ignore
from datetime import datetime, time

from .julian import jd_from_datetime, julian_to_hijri
from .conjunction import get_conjunction_time
from .visibility import evaluate_visibility
from .sun_times import get_sunset_time
from .rukyat_epoch import (
    get_rukyat_date,
    get_active_rukyat_epoch,
    RUKYAT_EPOCHS,
    RukyatMonthEpoch,
)
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
        raise ValueError("now_local must be timezone-aware")

    tz = pytz.timezone(timezone)
    now_local = now_local.astimezone(tz)

    # ==========================
    # 1. Lokasi referensi
    # ==========================
    if method == "global":
        ref_lat, ref_lon, ref_zone = DEFAULT_LOCATION["global"]
    else:
        ref_lat, ref_lon, ref_zone = lat, lon, timezone

    # ==========================
    # 2. Sunset hari ini
    # ==========================
    sunset_local = get_sunset_time(
        now_local.date(), ref_lat, ref_lon, ref_zone, ts, eph
    )
    if sunset_local is None:
        return {"error": "Sunset not available"}

    after_sunset = now_local >= sunset_local
    sunset_utc = sunset_local.astimezone(pytz.utc)
    sunset_jd = jd_from_datetime(sunset_utc, ts)

    # ==========================
    # 3. BASELINE HIJRIYAH (NOON)
    # ==========================
    noon_local = tz.localize(datetime.combine(now_local.date(), time(12, 0)))
    noon_utc = noon_local.astimezone(pytz.utc)
    noon_jd = jd_from_datetime(noon_utc, ts)

    baseline = julian_to_hijri(noon_jd)

    # Pergantian hari Hijriyah di maghrib
    if after_sunset:
        baseline = _increment_hijri_day(baseline)

    # ==========================
    # GLOBAL
    # ==========================
    if method == "global":
        return baseline

    # ==========================
    # HISAB
    # ==========================
    if method == "hisab":
        if baseline["day"] == 29:
            conj_jd = get_conjunction_time(sunset_jd, ts, earth, sun, moon)
            if conj_jd < sunset_jd and after_sunset:
                return _next_month_day(baseline)

            return {**baseline, "day": 30}

        if baseline["day"] == 30:
            if after_sunset:
                return _next_month_day(baseline)

        return baseline

    # ==========================
    # RUKYAT
    # ==========================
    if method == "rukyat":
        current_utc = now_local.astimezone(pytz.utc)
        current_jd = jd_from_datetime(current_utc, ts)

        rukyat_date = get_rukyat_date(current_jd)

        # sebelum hari 29 → aman
        if rukyat_date["day"] < 29:
            return rukyat_date

        if rukyat_date["day"] == 29 and not after_sunset:
            return rukyat_date

        # hari 29 → evaluasi
        if rukyat_date["day"] == 29 and after_sunset:
            conj_jd = get_conjunction_time(sunset_jd, ts, earth, sun, moon)
            visibility = evaluate_visibility(
                sunset_utc, lat, lon, conj_jd, ts, sun, moon, earth
            )

            if visibility["is_visible"]:
                epoch = get_active_rukyat_epoch(noon_jd)
                ny, nm = _next_month(epoch.year, epoch.month)

                RUKYAT_EPOCHS[(ny, nm)] = RukyatMonthEpoch(
                    year=ny,
                    month=nm,
                    jd_start=sunset_jd,
                )

                return {"year": ny, "month": nm, "day": 1}

            # hilal tidak terlihat → ISTIKMAL
            return {**rukyat_date, "day": 30}

        # hari 30 → otomatis bulan baru
        if rukyat_date["day"] == 30:
            epoch = get_active_rukyat_epoch(noon_jd)
            ny, nm = _next_month(epoch.year, epoch.month)

            RUKYAT_EPOCHS[(ny, nm)] = RukyatMonthEpoch(
                year=ny,
                month=nm,
                jd_start=sunset_jd,
            )

            return {"year": ny, "month": nm, "day": 1}

        return rukyat_date

    raise ValueError(f"Unknown method: {method}")


# ==========================
# HELPERS
# ==========================


def _increment_hijri_day(date):
    if date["day"] < 30:
        return {**date, "day": date["day"] + 1}
    return _next_month_day(date)


def _next_month_day(date):
    y, m = date["year"], date["month"]
    ny, nm = _next_month(y, m)
    return {"year": ny, "month": nm, "day": 1}


def _next_month(year, month):
    if month == 12:
        return year + 1, 1
    return year, month + 1


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
