import pytz
from datetime import datetime, time
from .julian import jd_from_datetime, jd_to_datetime, julian_to_hijri
from .conjunction import get_conjunction_time
from .visibility import evaluate_visibility
from .sun_times import get_sunset_time
from .config import DEFAULT_LOCATION


def get_hijri_date(lat, lon, method, timezone, *, now_local, ts, eph, sun, moon, earth):
    if now_local is None:
        raise ValueError("now_local must be timezone-aware")
    tz = pytz.timezone(timezone)
    now_local = now_local.astimezone(tz)

    # 1. Lokasi & Baseline
    ref_lat, ref_lon, ref_zone = (
        DEFAULT_LOCATION["global"] if method == "global" else (lat, lon, timezone)
    )

    selected_criteria = "Wujudul Hilal" if method == "rukyat" else "MABIMS"

    noon_local = tz.localize(datetime.combine(now_local.date(), time(12, 0)))
    noon_jd = jd_from_datetime(noon_local.astimezone(pytz.utc), ts)
    current_date = julian_to_hijri(noon_jd)

    # 2. Penyesuaian Rukyat
    if method == "rukyat":
        is_lagging = _check_historical_lag(
            current_date,
            noon_jd,
            lat,
            lon,
            timezone,
            ts,
            eph,
            sun,
            moon,
            earth,
            criteria=selected_criteria,
        )
        if is_lagging:
            current_date = _decrement_hijri_day(current_date)

    # 3. Cek Sunset untuk Pergantian Hari
    sunset_local = get_sunset_time(
        now_local.date(), ref_lat, ref_lon, ref_zone, ts, eph
    )
    after_sunset = now_local >= sunset_local if sunset_local else False

    # Jika Belum Maghrib, return baseline yang sudah di-adjust
    if not after_sunset:
        return current_date

    # ==========================
    # LOGIKA SETELAH MAGHRIB
    # ==========================

    # Jika bukan tgl 29, semua metode tinggal +1 hari
    if current_date["day"] != 29:
        return _increment_hijri_day(current_date)

    # Kondisi Sore Hari Tanggal 29 (Sidang Isbat)
    sunset_utc = sunset_local.astimezone(pytz.utc)
    sunset_jd = jd_from_datetime(sunset_utc, ts)
    conj_jd = get_conjunction_time(sunset_jd, ts, earth, sun, moon)

    if method == "global":
        return _increment_hijri_day(current_date)

    if method == "hisab":
        return (
            _start_new_month(current_date)
            if conj_jd < sunset_jd
            else {**current_date, "day": 30}
        )

    if method == "rukyat":
        vis = evaluate_visibility(
            sunset_utc,
            lat,
            lon,
            conj_jd,
            ts,
            sun,
            moon,
            earth,
            criteria=selected_criteria,
        )
        return (
            _start_new_month(current_date)
            if vis["is_visible"]
            else {**current_date, "day": 30}
        )

    raise ValueError(f"Unknown method: {method}")


# ==========================
# HELPERS
# ==========================


def _check_historical_lag(
    current_baseline,
    noon_jd,
    lat,
    lon,
    timezone,
    ts,
    eph,
    sun,
    moon,
    earth,
    criteria="MABIMS",
):
    days_since_start = current_baseline["day"] - 1
    month_start_noon_jd = noon_jd - days_since_start

    day_before_1st_jd = month_start_noon_jd - 1.0
    day_before_1st_utc = jd_to_datetime(day_before_1st_jd, ts)

    sunset_local = get_sunset_time(
        day_before_1st_utc.date(), lat, lon, timezone, ts, eph
    )
    if not sunset_local:
        return False

    sunset_utc = sunset_local.astimezone(pytz.utc)
    sunset_jd = jd_from_datetime(sunset_utc, ts)

    conj_jd = get_conjunction_time(sunset_jd, ts, earth, sun, moon)
    vis = evaluate_visibility(
        sunset_utc, lat, lon, conj_jd, ts, sun, moon, earth, criteria=criteria
    )

    return not vis["is_visible"]


def _increment_hijri_day(date):
    if date["day"] < 30:
        return {**date, "day": date["day"] + 1}
    return _start_new_month(date)


def _decrement_hijri_day(date):
    d, m, y = date["day"], date["month"], date["year"]
    if d > 1:
        return {**date, "day": d - 1}
    pm, py = (m - 1, y) if m > 1 else (12, y - 1)
    return {"year": py, "month": pm, "day": 30}


def _start_new_month(date):
    m, y = date["month"], date["year"]
    nm, ny = (m + 1, y) if m < 12 else (1, y + 1)
    return {"year": ny, "month": nm, "day": 1}
