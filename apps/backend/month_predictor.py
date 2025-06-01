from datetime import datetime, timedelta
import pytz
from julian import jd_from_datetime, julian_to_hijri
from astro_utils import DEFAULT_LOCATION
from visibility import evaluate_visibility

def predict_end_of_month(lat, lon, method, timezone):
    # Tentukan referensi lokasi berdasarkan metode
    if method == 'global':
        ref_lat, ref_lon = DEFAULT_LOCATION['global']
        ref_zone = DEFAULT_LOCATION['zone']
    else:
        ref_lat, ref_lon = lat, lon
        ref_zone = timezone
    
    # Mendapatkan tanggal Hijriyah hari ini
    now_local = datetime.now(pytz.timezone(timezone))
    jd_today = jd_from_datetime(now_local)

    # Menghitung tanggal Hijriyah hari ini
    hijri_today = julian_to_hijri(jd_today)

    # Cek jika hari ini adalah tanggal 29
    if hijri_today['day'] != 29:
        return {
            'today': hijri_today,
            'estimated_end_of_month': None,
            'message': "Prediksi hilal hanya tersedia saat tanggal 29 Hijriyah."
        }

    # Perhitungan lanjut untuk tanggal 29 Hijriyah
    estimated_end = {
        'day': 29,
        'month': hijri_today['month'],
        'year': hijri_today['year']
    }

    # Hitung waktu konjungsi
    conjunction_jd = find_conjunction(jd_today)

    # Tentukan apakah bulan baru dimulai berdasarkan konjungsi
    is_new_month = False
    if method == 'global':
        sunset_global = get_sunset_time(now_local.date(), 21.422487, 39.826206, 'Asia/Riyadh')
        sunset_global_jd = jd_from_datetime(sunset_global.astimezone(pytz.utc))
        if conjunction_jd < sunset_global_jd:
            is_new_month = True
    elif method == 'hisab':
        sunset_user = get_sunset_time(now_local.date(), lat, lon, timezone)
        sunset_user_jd = jd_from_datetime(sunset_user.astimezone(pytz.utc))
        if conjunction_jd < sunset_user_jd:
            is_new_month = True
    elif method == 'rukyat':
        visibility = evaluate_visibility(sunset_global, lat, lon, conjunction_jd)
        if visibility["is_visible"]:
            is_new_month = True

    # Tentukan hasil akhir bulan Hijriyah
    if is_new_month:
        estimated_end = {
            'day': 1,
            'month': hijri_today['month'] + 1,
            'year': hijri_today['year']
        }
        if estimated_end['month'] > 12:
            estimated_end['month'] = 1
            estimated_end['year'] += 1

    return {
        'today': hijri_today,
        'estimated_end_of_month': estimated_end,
        'message': f"Prediksi bulan ini berakhir pada {estimated_end['day']}-{estimated_end['month']}-{estimated_end['year']}."
    }
