from datetime import datetime
import pytz
from .config import ACEH_LOCATION as ACEH_REF
from .hijri_calculator import get_hijri_date
from .sun_times import get_sunset_time
from .conjunction import get_conjunction_time
from .visibility import evaluate_visibility
from .julian import jd_from_datetime


def predict_end_of_month(lat, lon, method, timezone, *, ts, eph, sun, moon, earth):
    tz = pytz.timezone(timezone)
    now_local = datetime.now(tz)

    # 1. Ambil Data Hari Ini
    today = get_hijri_date(
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

    if today["day"] != 29:
        return {
            "method": method,
            "today": today,
            "estimated_end_of_month": None,
            "visibility": None,
            "message": "Prediksi akhir bulan hanya tersedia pada tanggal 29 Hijriyah.",
        }

    # 2. Persiapan Data Astronomis
    visibility_data = None
    is_new_month_nationally = False

    # Logic cuma jalan buat Hisab & Rukyat
    if method != "global":
        criteria = "MABIMS" if method == "rukyat" else "Wujudul Hilal"

        # --- HITUNG LOKAL (User) ---
        sunset_local = get_sunset_time(now_local.date(), lat, lon, timezone, ts, eph)

        if sunset_local:
            sunset_utc = sunset_local.astimezone(pytz.utc)
            sunset_jd = jd_from_datetime(sunset_utc, ts)
            conj_jd = get_conjunction_time(sunset_jd, ts, earth, sun, moon)

            # Evaluasi di lokasi user
            vis_user = evaluate_visibility(
                sunset_utc, lat, lon, conj_jd, ts, sun, moon, earth, criteria=criteria
            )

            # --- HITUNG REFERENSI (Aceh) ---
            sunset_aceh = get_sunset_time(
                now_local.date(), ACEH_REF["lat"], ACEH_REF["lon"], timezone, ts, eph
            )
            sunset_aceh_utc = sunset_aceh.astimezone(pytz.utc)

            vis_aceh = evaluate_visibility(
                sunset_aceh_utc,
                ACEH_REF["lat"],
                ACEH_REF["lon"],
                conj_jd,
                ts,
                sun,
                moon,
                earth,
                criteria=criteria,
            )

            # Keputusan Nasional: Jika salah satu titik (User ATAU Aceh) "Lolos" kriteria
            is_new_month_nationally = vis_user["is_visible"] or vis_aceh["is_visible"]

            # Kita tampilkan data yang "terbaik" (biasanya Aceh lebih tinggi posisinya)
            # Biar user dapet visualisasi kenapa keputusannya 'Lolos'
            final_vis = (
                vis_aceh
                if (vis_aceh["moon_altitude"] > vis_user["moon_altitude"])
                else vis_user
            )

            visibility_data = {
                "moon_age": float(
                    (
                        sunset_utc - ts.utc_to_datetime(ts.from_jd(conj_jd))
                    ).total_seconds()
                    / 3600
                ),
                "moon_altitude": float(final_vis["moon_altitude"]),
                "elongation": float(final_vis["elongation"]),
                "is_visible": bool(is_new_month_nationally),  # Status Nasional
            }

    # 3. Penentuan Estimasi Tanggal
    if is_new_month_nationally:
        estimated_date = today
        msg = f"Hilal memenuhi kriteria {criteria} di wilayah Indonesia. Besok diperkirakan ganti bulan."
    else:
        estimated_date = {**today, "day": 30}
        msg = "Hilal belum memenuhi kriteria di seluruh wilayah Indonesia. Bulan digenapkan (Istikmal)."

    return {
        "method": method,
        "today": today,
        "estimated_end_of_month": estimated_date,
        "visibility": visibility_data,
        "message": msg,
    }


def get_next_month_name(current_month_index):
    months = [
        "Safar",
        "Rabiulawal",
        "Rabiulakhir",
        "Jumadilawal",
        "Jumadilakhir",
        "Rajab",
        "Syakban",
        "Ramadan",
        "Syawal",
        "Zulkaidah",
        "Zulhijah",
        "Muharam",
    ]
    return months[current_month_index - 1]
