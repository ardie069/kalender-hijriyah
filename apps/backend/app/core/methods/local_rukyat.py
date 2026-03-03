from ..calendar.hijri_date import start_new_month
from ..config import REGIONAL_RUKYAT_CONFIG
from ._shared import prepare_hijri_baseline, handle_before_sunset, handle_normal_day
from .base import BaseHijriMethod, HijriResult
from ..services.engine import (
    calculate_conjunction,
    calculate_visibility,
    calculate_sunset,
)
import pytz
from ..calendar.julian import jd_from_datetime


class LocalRukyatMethod(BaseHijriMethod):
    def __init__(self, mode="individual", region="indonesia"):
        self.mode = mode
        self.region = region
    
    def calculate(self, context):
        # 1. Siapkan baseline tanggal Hijriyah dan cek apakah sudah melewati Maghrib
        baseline, after_sunset, sunset_local = prepare_hijri_baseline(
            context,
            criteria="MABIMS", # Default kriteria rukyat lokal menggunakan MABIMS
        )

        # 2. Jika belum Maghrib, gunakan tanggal hari ini (sebelum ganti hari)
        if not after_sunset:
            return handle_before_sunset(baseline, "local_rukyat")

        # 3. Jika bukan tanggal 29, lakukan increment hari secara normal (bukan fase penentuan bulan baru)
        if baseline["day"] != 29:
            return handle_normal_day(baseline, "local_rukyat")

        # 4. Tentukan daftar titik observasi (sites) berdasarkan mode
        if self.mode == "individual":
            # Mode individual: hanya cek di lokasi user
            sites = [
                {
                    "name": "User Location",
                    "lat": context.lat,
                    "lon": context.lon,
                    "timezone": context.timezone,
                }
            ]
            criteria = "MABIMS"
        else:
            # Mode regional: cek di beberapa titik sesuai konfigurasi wilayah (misal: Indonesia)
            region_config = REGIONAL_RUKYAT_CONFIG.get(
                self.region, REGIONAL_RUKYAT_CONFIG["indonesia"]
            )
            sites = region_config["sites"]
            criteria = region_config["criteria"]

        best_visibility = None
        best_score = -999
        any_visible = False

        # 5. Iterasi setiap titik observasi untuk mengecek visibilitas hilal
        for i, site in enumerate(sites):
            if self.mode == "individual" and i == 0:
                sunset_site = sunset_local
            else:
                # Hitung waktu sunset untuk koordinat spesifik site
                sunset_site = calculate_sunset(
                    context.now_local.date(), site["lat"], site["lon"], site["timezone"]
                )

            if not sunset_site:
                continue

            # Konversi ke UTC dan Julian Date untuk perhitungan astronomi
            sunset_utc = sunset_site.astimezone(pytz.utc)
            sunset_jd = jd_from_datetime(sunset_utc)
            conj_jd = calculate_conjunction(sunset_jd)

            # Hitung parameter visibilitas (ketinggian, elongasi, dll)
            vis = calculate_visibility(
                sunset_utc, site["lat"], site["lon"], conj_jd, criteria=criteria
            )

            # Simpan data visibilitas terbaik untuk keperluan metadata/debug
            score = vis["moon_altitude"] + vis["elongation"]
            if score > best_score:
                best_score = score
                best_visibility = {
                    **vis,
                    "lat": site["lat"],
                    "lon": site["lon"],
                    "site_name": site["name"],
                    "moon_age_hours": vis.get("moon_age_hours", vis.get("moon_age", 0)),
                }

            # Jika salah satu titik berhasil melihat hilal (sesuai kriteria), tandai True
            if vis["is_visible"]:
                any_visible = True

        # 6. Pengambilan keputusan akhir bulan
        if any_visible:
            final_date = start_new_month(baseline) # Hilal terlihat: Besok tanggal 1
            decision = "new_month"
        else:
            final_date = {**baseline, "day": 30} # Hilal tidak terlihat: Istikmal (genapkan 30 hari)
            decision = "istikmal_30"

        # 7. Kembalikan hasil beserta metadata untuk transparansi perhitungan
        return HijriResult(
            hijri_date=final_date,
            metadata={
                "type": "local_rukyat",
                "decision": decision,
                "visibility": best_visibility,
                "region": self.region,
                "mode": self.mode,
            },
        )
