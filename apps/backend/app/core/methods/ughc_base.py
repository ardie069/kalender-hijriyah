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
        # 1. Baseline Aritmatika (Jangkar)
        baseline, noon_jd = calculate_baseline_hijri(
            context.now_local, context.timezone, context.ts
        )

        # 2. Maghrib Lokal User (Ganti Hari)
        sunset_local = calculate_sunset(
            context.now_local.date(),
            context.lat,
            context.lon,
            context.timezone,
            context.ts,
            context.eph,
        )
        after_sunset = context.now_local >= sunset_local if sunset_local else False

        if not after_sunset:
            return HijriResult(
                hijri_date=baseline, metadata={"decision": "before_maghrib"}
            )

        # 3. Hari Biasa (1-28) -> Increment
        if baseline["day"] != 29:
            return HijriResult(
                hijri_date=increment_hijri_day(baseline),
                metadata={"decision": "normal_increment"},
            )

        # 4. Evaluasi Akhir Bulan (Tanggal 29)
        # Scan fakta material dunia
        scan = GlobalVisibilityRegistry.scan_global(
            context.now_local.date(),
            self.CRITERIA,
            context.ts,
            context.eph,
            context.sun,
            context.moon,
            context.earth,
        )

        # --- LOGIKA KEPUTUSAN KHGT MUHAMMADIYAH ---
        is_new_month = False

        # A. Butir 1: Terpenuhi sebelum 24:00 UTC di mana pun
        if scan["anywhere_before_24utc"]:
            is_new_month = True

        # B. Butir 2: Terpenuhi setelah 24:00 UTC
        elif scan["anywhere_after_24utc"] or scan["america_visible"]:
            # 2b. Di Benua Amerika -> Sah
            if scan["america_visible"]:
                is_new_month = True
            # 2a. Di tempat lain + Konjungsi sebelum Fajr NZ
            else:
                conj_jd = calculate_conjunction(
                    noon_jd, context.ts, context.earth, context.sun, context.moon
                )
                conj_dt_utc = jd_to_datetime(conj_jd, context.ts)
                nz_fajr = get_fajr_time(
                    conj_dt_utc.date(),
                    NZ_LOCATION["lat"],
                    NZ_LOCATION["lon"],
                    "UTC",
                    context.ts,
                    context.eph,
                )

                if nz_fajr and conj_dt_utc < nz_fajr.astimezone(pytz.utc):
                    is_new_month = True

        # 5. Final Decision
        if is_new_month:
            final_date = start_new_month(baseline)
            decision = "new_month"
        else:
            final_date = {**baseline, "day": 30}
            decision = "istikmal_30"

        return HijriResult(
            hijri_date=final_date,
            metadata={
                "type": self.TYPE,
                "decision": decision,
                "global_visible": is_new_month,
                "visibility_data": scan["best_visibility"],
            },
        )
