from skyfield.api import wgs84
from skyfield import almanac
from datetime import timedelta
import pytz

from ..services.engine import (
    calculate_baseline_hijri,
    calculate_sunset,
    calculate_conjunction,
)
from ..calendar.hijri_date import increment_hijri_day, start_new_month
from ..calendar.julian import jd_from_datetime
from ..config import AL_HARAM_LOCATION
from app.deps.astronomy import get_provider
from .base import BaseHijriMethod, HijriResult


class UmmAlQuraMethod(BaseHijriMethod):
    def calculate(self, context):
        # 1. Hitung baseline tanggal Hijriyah berdasarkan kalender aritmatik (jangkar)
        baseline, _ = calculate_baseline_hijri(context.now_local, context.timezone)

        # 2. Referensi Umm Al-Qura selalu menggunakan koordinat Ka'bah (Makkah)
        ref_lat, ref_lon, ref_tz = (
            AL_HARAM_LOCATION["lat"], 
            AL_HARAM_LOCATION["lon"],
            AL_HARAM_LOCATION["timezone"],
        )
        makkah_tz = pytz.timezone(ref_tz)
        local_now_makkah = context.now_local.astimezone(makkah_tz)

        sunset_makkah = calculate_sunset(
            local_now_makkah.date(), ref_lat, ref_lon, ref_tz 
        )

        # 3. Cek apakah waktu sekarang sudah melewati Maghrib di Makkah
        if not sunset_makkah or local_now_makkah < sunset_makkah:
            # Jika belum Maghrib, gunakan tanggal hari ini (baseline)
            return HijriResult(
                hijri_date=baseline,
                metadata={"type": "umm_al_qura", "decision": "before_maghrib"},
            )

        if baseline["day"] == 29:
            p = get_provider()
            sunset_utc = sunset_makkah.astimezone(pytz.utc) 
            sunset_jd = jd_from_datetime(sunset_utc)

            # Kriteria 1: Konjungsi terjadi sebelum matahari terbenam di Makkah
            conj_jd = calculate_conjunction(sunset_jd)
            is_conj_before_sunset = conj_jd < sunset_jd

            # Kriteria 2: Bulan terbenam setelah matahari terbenam di Makkah
            t0 = p.ts.from_datetime(sunset_utc - timedelta(hours=1))
            t1 = p.ts.from_datetime(sunset_utc + timedelta(hours=1))

            # Cari event moonset di sekitar waktu sunset
            f = almanac.risings_and_settings(
                p.eph, p.moon, wgs84.latlon(ref_lat, ref_lon)
            )
            times, events = almanac.find_discrete(t0, t1, f)

            moonset_jd = None
            for t, event in zip(times, events):
                if event == 0:
                    moonset_jd = t.tt
                    break

            # Fallback jika event moonset tidak ditemukan dalam window 2 jam
            if moonset_jd is None:
                t_sunset = p.ts.from_datetime(sunset_utc)
                observer = p.earth + wgs84.latlon(ref_lat, ref_lon)
                alt, _, _ = observer.at(t_sunset).observe(p.moon).apparent().altaz()
                # Jika altitude > 0 (dengan koreksi refraksi sederhana), berarti belum terbenam
                is_moonset_after_sunset = (
                    alt.degrees + 0.25
                ) > 0
            else:
                is_moonset_after_sunset = moonset_jd > sunset_jd

            # Keputusan: Jika kedua kriteria terpenuhi, besok bulan baru
            if is_conj_before_sunset and is_moonset_after_sunset:
                final_date = start_new_month(baseline)
                decision = "new_month"
            else:
                final_date = {**baseline, "day": 30}
                decision = "istikmal_30"

            metadata = {
                "type": "umm_al_qura",
                "decision": decision,
                "is_conjunction_before_sunset": is_conj_before_sunset,
                "is_moonset_after_sunset": is_moonset_after_sunset,
            }
        else:
            # Jika bukan tanggal 29, cukup increment hari secara normal
            final_date = increment_hijri_day(baseline)
            metadata = {"type": "umm_al_qura", "decision": "normal_increment"}

        return HijriResult(hijri_date=final_date, metadata=metadata)
