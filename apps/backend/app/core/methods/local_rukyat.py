# app/core/methods/local_rukyat.py

from ..calendar.hijri_date import start_new_month
from ..config import REGIONAL_RUKYAT_CONFIG
from ._shared import prepare_hijri_baseline, handle_before_sunset, handle_normal_day
from .base import BaseHijriMethod, HijriResult
from ..services.engine import (
    calculate_sunset,
    calculate_conjunction,
    calculate_visibility,
)
import pytz
from ..calendar.julian import jd_from_datetime


class LocalRukyatMethod(BaseHijriMethod):
    def __init__(self, mode="individual", region="indonesia"):
        self.mode = mode
        self.region = region

    def calculate(self, context):
        baseline, after_sunset, sunset_local = prepare_hijri_baseline(
            context,
            criteria="MABIMS",  # Baseline awal
        )

        if not after_sunset:
            return handle_before_sunset(baseline, "local_rukyat")

        if baseline["day"] != 29:
            return handle_normal_day(baseline, "local_rukyat")

        # ─── TENTUKAN TARGET SCANNING ───
        if self.mode == "individual":
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
            region_config = REGIONAL_RUKYAT_CONFIG.get(
                self.region, REGIONAL_RUKYAT_CONFIG["indonesia"]
            )
            sites = region_config["sites"]
            criteria = region_config["criteria"]

        # ─── PROSES SCANNING (Hanya di Lokasi Terdaftar) ───
        best_visibility = None
        best_score = -999
        any_visible = False

        for site in sites:
            # Hitung sunset di lokasi spesifik site
            sunset_site = calculate_sunset(
                context.now_local.date(), site["lat"], site["lon"], site["timezone"]
            )
            if not sunset_site:
                continue

            sunset_utc = sunset_site.astimezone(pytz.utc)
            sunset_jd = jd_from_datetime(sunset_utc)
            conj_jd = calculate_conjunction(sunset_jd)

            vis = calculate_visibility(
                sunset_utc, site["lat"], site["lon"], conj_jd, criteria=criteria
            )

            # Cari yang paling tinggi hilalnya buat metadata
            score = vis["moon_altitude"] + vis["elongation"]
            if score > best_score:
                best_score = score
                best_visibility = {
                    **vis,
                    "lat": site["lat"],
                    "lon": site["lon"],
                    "site_name": site["name"],
                }

            if vis["is_visible"]:
                any_visible = True

        # ─── KEPUTUSAN FINAL ───
        if any_visible:
            final_date = start_new_month(baseline)
            decision = "new_month"
        else:
            final_date = {**baseline, "day": 30}
            decision = "istikmal_30"

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
