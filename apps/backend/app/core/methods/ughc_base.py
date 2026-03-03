import pytz

from ..services.engine import (
    calculate_baseline_hijri,
    calculate_sunset,
    calculate_conjunction,
)
from ..calendar.hijri_date import (
    increment_hijri_day,
    start_new_month,
)

from ..calendar.julian import jd_to_datetime
from ..astronomy.sun_times import get_fajr_time
from ..services.visibility_scan import GlobalVisibilityRegistry
from ..config import NZ_LOCATION
from .base import BaseHijriMethod, HijriResult


class BaseUGHCMethod(BaseHijriMethod):
    CRITERIA = None
    TYPE = None

    def calculate(self, context):
        # 1. Ambil baseline hanya untuk struktur datanya (Year, Month)
        baseline, noon_jd = calculate_baseline_hijri(
            context.now_local,
            context.timezone,
        )

        # 2. Status Maghrib Lokal (Tetap perlu buat ganti hari harian)
        sunset_local = calculate_sunset(
            context.now_local.date(),
            context.lat,
            context.lon,
            context.timezone,
        )
        after_sunset = context.now_local >= sunset_local if sunset_local else False

        if not after_sunset:
            return HijriResult(
                hijri_date=baseline, metadata={"decision": "before_maghrib"}
            )

        # ─── KUNCI PERBAIKAN DI SINI ───
        # Jangan pake baseline["day"] buat ngecek hari ke-29.
        # Kita cek apakah context ini dikirim oleh Predictor untuk evaluasi akhir bulan.
        # Jika Predictor jalan, dia bakal simulasi di hari ke-29.

        # Ambil scan global
        scan = GlobalVisibilityRegistry.scan_global(
            context.now_local.date(),
            self.CRITERIA,
        )

        is_new_month = False

        # A. Butir 1: Terpenuhi sebelum 24:00 UTC
        if scan.get("anywhere_before_24utc", False):
            is_new_month = True

        # B. Butir 2: Terpenuhi setelah 24:00 UTC (Amerika atau NZ Fajr)
        elif scan.get("anywhere_after_24utc", False) or scan.get(
            "america_visible", False
        ):
            if scan.get("america_visible", False):
                is_new_month = True
            else:
                conj_jd = calculate_conjunction(noon_jd)
                conj_dt_utc = jd_to_datetime(conj_jd)
                nz_fajr = get_fajr_time(
                    conj_dt_utc.date(),
                    NZ_LOCATION["lat"],
                    NZ_LOCATION["lon"],
                    "UTC",
                )
                if nz_fajr and conj_dt_utc < nz_fajr.astimezone(pytz.utc):
                    is_new_month = True

        # 3. Penentuan Akhir
        # Jika is_new_month TRUE, kita paksa ganti bulan.
        if is_new_month:
            final_date = start_new_month(baseline)
            decision = "new_month"
        else:
            # Jika is_new_month FALSE, kita liat:
            # Kalau ini emang simulasi hari ke-29, berarti besok tanggal 30 (Istikmal)
            # Kalau bukan simulasi hari ke-29 (hari biasa), ya cuma increment biasa.
            if baseline["day"] == 29:
                final_date = {**baseline, "day": 30}
                decision = "istikmal_30"
            else:
                final_date = increment_hijri_day(baseline)
                decision = "normal_increment"

        return HijriResult(
            hijri_date=final_date,
            metadata={
                "type": self.TYPE,
                "decision": decision,
                "global_visible": is_new_month,
                "visibility_data": scan.get("best_visibility"),
            },
        )
