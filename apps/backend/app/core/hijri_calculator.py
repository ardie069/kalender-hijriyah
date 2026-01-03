import pytz  # type: ignore
from .julian import jd_from_datetime, julian_to_hijri
from .conjunction import get_conjunction_time
from .visibility import evaluate_visibility
from .sun_times import get_sunset_time
from .rukyat_epoch import (
    get_rukyat_date,
    get_rukyat_day,
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
    # Tentukan sunset & JD
    # ==========================
    if method == "global":
        ref_lat, ref_lon, ref_zone = DEFAULT_LOCATION["global"]
    else:
        ref_lat, ref_lon, ref_zone = lat, lon, timezone

    sunset_local = get_sunset_time(
        now_local.date(), ref_lat, ref_lon, ref_zone, ts, eph
    )
    if sunset_local is None:
        return {"error": "Sunset not available"}

    after_sunset = now_local >= sunset_local
    sunset_utc = sunset_local.astimezone(pytz.utc)
    sunset_jd = jd_from_datetime(sunset_utc, ts)
    effective_jd = sunset_jd + 0.5 if after_sunset else sunset_jd - 0.5

    # ==========================
    # GLOBAL
    # ==========================
    if method == "global":
        return julian_to_hijri(effective_jd)

    # ==========================
    # HISAB
    # ==========================
    if method == "hisab":
        raw = julian_to_hijri(effective_jd)

        if raw["day"] == 29:
            conj_jd = get_conjunction_time(effective_jd - 1, ts, earth, sun, moon)
            if conj_jd < sunset_jd:
                return _next_hijri_day(raw)
            return {**raw, "day": 30}

        return raw

    # ==========================
    # RUKYAT (INI YANG PENTING)
    # ==========================
    if method == "rukyat":
        # ambil tanggal rukyat dari epoch
        rukyat_date = get_rukyat_date(effective_jd)
        day = rukyat_date["day"]

        # belum hari ke-29 → selesai
        if day < 29:
            return rukyat_date

        # hari ke-29 → evaluasi hilal
        if day == 29:
            conj_jd = get_conjunction_time(effective_jd - 1, ts, earth, sun, moon)

            visibility = evaluate_visibility(
                sunset_utc, lat, lon, conj_jd, ts, sun, moon, earth
            )

            if visibility["is_visible"]:
                # buat epoch bulan baru
                epoch = get_active_rukyat_epoch(effective_jd)
                next_year, next_month = _next_month(epoch.year, epoch.month)

                RUKYAT_EPOCHS[(next_year, next_month)] = RukyatMonthEpoch(
                    year=next_year,
                    month=next_month,
                    jd_start=sunset_jd,
                )

                return {"year": next_year, "month": next_month, "day": 1}

            # hilal tidak terlihat → tetap 29
            return rukyat_date

        # hari ke-30 → istikmal otomatis
        epoch = get_active_rukyat_epoch(effective_jd)
        next_year, next_month = _next_month(epoch.year, epoch.month)

        RUKYAT_EPOCHS[(next_year, next_month)] = RukyatMonthEpoch(
            year=next_year,
            month=next_month,
            jd_start=sunset_jd,
        )

        return {"year": next_year, "month": next_month, "day": 1}

    raise ValueError(f"Unknown method: {method}")


def _next_hijri_day(date):
    year, month = date["year"], date["month"]
    return _next_month(year, month) | {"day": 1}  # type: ignore


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
