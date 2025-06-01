from datetime import datetime, timedelta
import pytz
from julian import jd_from_datetime, julian_to_hijri
from conjunction import get_conjunction_time
from visibility import evaluate_visibility
from sun_times import get_sunset_time
from astro_utils import DEFAULT_LOCATION


def get_hijri_date(lat, lon, method="global", timezone="UTC", jd=None):
    tz = pytz.timezone(timezone)
    now_local = datetime.now(tz)

    # Lokasi referensi
    if method == "global":
        ref_lat, ref_lon, ref_zone = DEFAULT_LOCATION["global"]
    else:
        ref_lat, ref_lon, ref_zone = lat, lon, timezone

    # Hitung waktu matahari terbenam
    sunset_local = get_sunset_time(now_local.date(), ref_lat, ref_lon, ref_zone)

    if sunset_local is None:
        # Jika waktu matahari terbenam tidak dapat dihitung, kembalikan error atau gunakan waktu default
        return {"error": "Gagal menghitung waktu matahari terbenam."}

    sunset_utc = sunset_local.astimezone(pytz.utc)
    sunset_jd = jd_from_datetime(sunset_utc)

    # Apakah sekarang sudah lewat maghrib?
    now_ref = now_local.astimezone(pytz.timezone(ref_zone))
    after_sunset = now_ref >= sunset_local

    # Gunakan JD saat maghrib, dan tambah 1 jika sudah malam
    effective_jd = sunset_jd + 1 if after_sunset else sunset_jd

    # Tentukan tanggal Hijriyah awal
    hijri_date = julian_to_hijri(effective_jd)

    # Jika tanggal 29 → Evaluasi kemungkinan bulan berganti
    if hijri_date["day"] == 29:
        conj_jd = find_conjunction(effective_jd - 1)
        is_conj_valid = False

        if method == "global":
            mekkah_sunset = get_sunset_time(
                now_local.date(), *DEFAULT_LOCATION["global"], DEFAULT_LOCATION["zone"]
            )
            if mekkah_sunset is None:
                return {"error": "Gagal menghitung waktu matahari terbenam di Mekkah."}
            jd_mekkah_sunset = jd_from_datetime(mekkah_sunset.astimezone(pytz.utc))
            is_conj_valid = conj_jd < jd_mekkah_sunset
        elif method == "hisab":
            is_conj_valid = conj_jd < sunset_jd
        elif method == "rukyat":
            visibility = evaluate_visibility(sunset_utc, lat, lon, conj_jd)
            is_conj_valid = visibility["is_visible"]

        if not is_conj_valid and after_sunset:
            effective_jd += 1
            hijri_date = julian_to_hijri(effective_jd)

    return hijri_date
