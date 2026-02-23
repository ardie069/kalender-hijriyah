import pytz
from datetime import datetime

from ..services.engine import (
    calculate_baseline_hijri,
    calculate_sunset,
    calculate_conjunction,
)
from ..calendar.hijri_date import (
    increment_hijri_day,
    start_new_month,
)

from ..calendar.julian import jd_from_datetime, jd_to_datetime
from ..astronomy.sun_times import get_fajr_time
from ..astronomy.global_grid import generate_global_grid
from ..services.visibility_scan import GlobalVisibilityRegistry
from ..config import NZ_LOCATION
from .base import BaseHijriMethod, HijriResult


class BaseUGHCMethod(BaseHijriMethod):
    CRITERIA = None
    TYPE = None

    @staticmethod
    def is_america(lon: float) -> bool:
        return -170 <= lon <= -30

    def calculate(self, context):

        # ==================================================
        # 1️⃣ BASELINE (Arithmetic Sync)
        # ==================================================
        baseline, noon_jd = calculate_baseline_hijri(
            context.now_local,
            context.timezone,
            context.ts,
        )

        # ==================================================
        # 2️⃣ MAGHRIB LOKAL USER
        # ==================================================
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
                hijri_date=baseline,
                metadata={
                    "type": self.TYPE,
                    "decision": "before_maghrib",
                    "baseline": baseline,
                },
            )

        # ==================================================
        # 3️⃣ HARI BUKAN 29 → INCREMENT NORMAL
        # ==================================================
        if baseline["day"] != 29:
            return HijriResult(
                hijri_date=increment_hijri_day(baseline),
                metadata={
                    "type": self.TYPE,
                    "decision": "normal_increment",
                    "baseline": baseline,
                },
            )

        # ==================================================
        # 4️⃣ KONJUNGSI GLOBAL (STABIL)
        # ==================================================
        conj_jd = calculate_conjunction(
            noon_jd,
            context.ts,
            context.earth,
            context.sun,
            context.moon,
        )

        conj_dt_utc = jd_to_datetime(conj_jd, context.ts)

        # ==================================================
        # 5️⃣ GLOBAL SCAN
        # ==================================================
        now_utc = context.now_local.astimezone(pytz.utc)

        end_of_day_utc = datetime(
            now_utc.year,
            now_utc.month,
            now_utc.day,
            23,
            59,
            59,
            tzinfo=pytz.utc,
        )

        sites = generate_global_grid(lat_step=10, lon_step=15)

        visibility_before_midnight = False
        visibility_after_midnight = False
        visibility_after_midnight_america = False

        result = GlobalVisibilityRegistry.scan_global(
            context.now_local.date(),
            sites,
            self.CRITERIA,
            context.ts,
            context.eph,
            context.sun,
            context.moon,
            context.earth,
        )

        # ==================================================
        # 6️⃣ KHGT GLOBAL DECISION
        # ==================================================
        global_visible = False

        # A. Sebelum 24:00 UTC
        if visibility_before_midnight:
            global_visible = True

        # B. Setelah 24:00 UTC
        elif visibility_after_midnight:

            # 1. Jika di Amerika → sah
            if visibility_after_midnight_america:
                global_visible = True

            # 2. Jika bukan Amerika → cek konjungsi sebelum Fajr NZ
            else:
                nz_fajr = get_fajr_time(
                    conj_dt_utc.date(),
                    NZ_LOCATION["lat"],
                    NZ_LOCATION["lon"],
                    "UTC",
                    context.ts,
                    context.eph,
                )

                if nz_fajr:
                    nz_fajr_utc = nz_fajr.astimezone(pytz.utc)
                    if conj_dt_utc < nz_fajr_utc:
                        global_visible = True

        # ==================================================
        # 7️⃣ FINAL DECISION
        # ==================================================
        if global_visible:
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
                "baseline": baseline,
                "global_visible": global_visible,
                "visibility_data": result["best_visibility"],
                "conjunction_before_maghrib": (
                    conj_jd
                    < jd_from_datetime(
                        sunset_local.astimezone(pytz.utc),
                        context.ts,
                    )
                    if sunset_local
                    else None
                ),
            },
        )
