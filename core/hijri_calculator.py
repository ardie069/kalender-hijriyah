from datetime import datetime
import pytz
from julian import jd_from_datetime, julian_to_hijri
from conjunction import get_conjunction_time
from visibility import evaluate_visibility
from sun_times import get_sunset_time
from config import DEFAULT_LOCATION


def get_hijri_date(lat, lon, method, timezone, *, ts, eph, sun, moon, earth):
    tz = pytz.timezone(timezone)
    now_local = datetime.now(tz)

    # Lokasi referensi berdasarkan metode
    if method == "global":
        ref_lat, ref_lon, ref_zone = DEFAULT_LOCATION["global"]
    else:
        ref_lat, ref_lon, ref_zone = lat, lon, timezone

    # Waktu matahari terbenam di lokasi referensi
    sunset_local = get_sunset_time(
        now_local.date(), ref_lat, ref_lon, ref_zone, ts, eph
    )
    if sunset_local is None:
        return {"error": "Gagal menghitung waktu matahari terbenam."}

    # Apakah sudah lewat maghrib?
    now_ref = now_local.astimezone(pytz.timezone(ref_zone))
    after_sunset = now_ref >= sunset_local

    # Julian Date saat maghrib (UTC)
    sunset_utc = sunset_local.astimezone(pytz.utc)
    sunset_jd = jd_from_datetime(sunset_utc, ts)
    effective_jd = sunset_jd + 1 if after_sunset else sunset_jd

    # Konversi ke tanggal Hijriyah
    hijri_date = julian_to_hijri(effective_jd)

    # Evaluasi hanya saat tanggal 29 atau 30
    if hijri_date["day"] == 29:
        conj_jd = get_conjunction_time(effective_jd - 1, ts, earth, sun, moon)
        is_new_month = evaluate_new_month(
            method, lat, lon, conj_jd, sunset_utc, now_local, ts=ts, eph=eph, sun=sun, moon=moon, earth=earth
        )

        if not is_new_month and after_sunset:
            # Genapkan jadi 30 hari
            hijri_date["day"] = 30

    elif hijri_date["day"] == 30 and after_sunset:
        # Otomatis masuk bulan baru setelah maghrib hari ke-30
        effective_jd += 1
        hijri_date = julian_to_hijri(effective_jd)

    return hijri_date


def evaluate_new_month(
    method, lat, lon, conj_jd, sunset_utc, now_local, *, ts, eph, sun, moon, earth
):
    """
    Evaluasi apakah bulan baru dimulai menurut metode tertentu.
    """
    if method == "global":
        mekkah_lat, mekkah_lon, mekkah_zone = DEFAULT_LOCATION["global"]
        sunset_mekkah = get_sunset_time(
            now_local.date(), mekkah_lat, mekkah_lon, mekkah_zone, ts, eph
        )
        if sunset_mekkah is None:
            return False
        jd_mekkah_sunset = jd_from_datetime(sunset_mekkah.astimezone(pytz.utc), ts)
        return conj_jd < jd_mekkah_sunset

    elif method == "hisab":
        jd_sunset = jd_from_datetime(sunset_utc, ts)
        return conj_jd < jd_sunset

    elif method == "rukyat":
        visibility = evaluate_visibility(sunset_utc, lat, lon, conj_jd, ts, sun, moon)
        return visibility["is_visible"]

    return False
