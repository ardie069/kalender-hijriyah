import pytz
from datetime import datetime, timedelta

from app.core.methods.factory import get_method_instance
from app.core.methods.base import HijriContext


def predict_end_of_month(lat, lon, method, timezone, *, ts, eph, sun, moon, earth):
    tz = pytz.timezone(timezone)
    now_local = datetime.now(tz)

    # 1. Hitung kondisi hari ini dulu
    method_instance = get_method_instance(method)
    context_today = HijriContext(
        lat=lat,
        lon=lon,
        timezone=timezone,
        now_local=now_local,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )
    result_today = method_instance.calculate(context_today)
    hijri_now = result_today.hijri_date

    # 2. Cari kapan tanggal 29 di bulan berjalan
    # Rumus Madilog: (29 - hari_sekarang) adalah sisa hari menuju hari evaluasi
    days_to_29 = 29 - hijri_now["day"]

    # Target simulasi: Maghrib di hari ke-29 Hijriyah
    simulation_local = now_local + timedelta(days=days_to_29)

    # 3. Jalankan Simulasi di Masa Depan
    context_sim = HijriContext(
        lat=lat,
        lon=lon,
        timezone=timezone,
        now_local=simulation_local,
        ts=ts,
        eph=eph,
        sun=sun,
        moon=moon,
        earth=earth,
    )
    result_sim = method_instance.calculate(context_sim)

    # 4. Ambil Keputusan dari Hasil Simulasi
    meta_sim = result_sim.metadata
    decision = meta_sim.get("decision")

    # Tentukan 1 Syawal (atau bulan berikutnya)
    year = hijri_now["year"]
    month = hijri_now["month"]
    next_month = 1 if month == 12 else month + 1
    next_year = year + 1 if month == 12 else year

    # Jika decision adalah new_month, berarti bulan ini cuma 29 hari
    # Jika decision adalah istikmal, berarti bulan ini genap 30 hari
    is_istikmal = decision == "istikmal_30"
    offset_from_29 = 2 if is_istikmal else 1

    estimated_gregorian = simulation_local + timedelta(days=offset_from_29)

    return {
        "method": method,
        "current_hijri": hijri_now,
        "target_simulation_date": simulation_local.date().isoformat(),
        "prediction": {
            "is_istikmal": is_istikmal,
            "total_days_in_month": 30 if is_istikmal else 29,
            "estimated_next_month_1": {
                "year": next_year,
                "month": next_month,
                "day": 1,
                "gregorian": estimated_gregorian.date().isoformat(),
            },
        },
        "visibility_at_29": meta_sim.get("visibility_data")
        or meta_sim.get("visibility"),
        "message": f"Berdasarkan metode {method}, bulan ini diperkirakan berjumlah {30 if is_istikmal else 29} hari. "
        f"1 {next_month} {next_year} H jatuh pada {estimated_gregorian.strftime('%d %B %Y')}.",
    }
