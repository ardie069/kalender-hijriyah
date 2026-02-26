import pytz
from typing import Dict, Any
from .engine import (
    calculate_baseline_hijri,
    calculate_sunset,
    calculate_conjunction,
    calculate_visibility,
)
from ..calendar.julian import jd_from_datetime
from ..config import REGIONAL_RUKYAT_CONFIG, HIJRI_MONTHS


class RukyatService:
    """
    Unified Rukyat Service — Menangani evaluasi lokal & nasional.
    Sudah terstandarisasi untuk Frontend & Pydantic Schema.
    """

    def evaluate(
        self, context, region: str = "indonesia", mode: str = "local"
    ) -> Dict[str, Any]:
        baseline, _ = calculate_baseline_hijri(
            context.now_local, context.timezone,
        )

        baseline["month_name"] = HIJRI_MONTHS[baseline["month"] - 1]

        if baseline["day"] != 29:
            return {
                "is_rukyat_day": False,
                "hijri_date": baseline,
                "message": "Hari ini bukan waktu pengamatan hilal (bukan 29 Hijriyah).",
            }

        region_config = REGIONAL_RUKYAT_CONFIG.get(region)
        if not region_config:
            raise ValueError(f"Invalid region: {region}")

        criteria = region_config["criteria"]

        if mode == "national":
            return self._evaluate_national(context, region_config, criteria)

        return self._evaluate_local(context, criteria)

    # ── Local (Single-Site) ───────────────────────────────────────────

    def _evaluate_local(self, context, criteria: str) -> Dict[str, Any]:
        sunset_local = calculate_sunset(
            context.now_local.date(),
            context.lat,
            context.lon,
            context.timezone,
        )

        if not sunset_local:
            return {"is_rukyat_day": True, "error": "Sunset gagal dihitung."}

        vis = self._process_visibility(
            context.lat, context.lon, sunset_local, criteria,
        )

        return {
            "is_rukyat_day": True,
            "sunset_time": sunset_local.isoformat(),
            # Data untuk Simulator V4
            "altitude_at_sunset": vis["moon_altitude"],
            "azimuth_at_sunset": vis.get("azimuth_moon", 0),
            "azimuth_diff": vis.get("azimuth_diff", 0),
            "elongation_at_sunset": vis["elongation"],
            "moon_age_hours": vis["moon_age_hours"],
            "is_visible": vis["is_visible"],
            "criteria_used": criteria,
            "site_name": "Local Positioning System",
        }

    # ── National (Multi-Site) ─────────────────────────────────────────

    def _evaluate_national(
        self, context, region_config: dict, criteria: str
    ) -> Dict[str, Any]:
        sites = region_config["sites"]
        best_visibility = None
        best_score = -999
        any_visible = False

        for site in sites:
            sunset_local = calculate_sunset(
                context.now_local.date(),
                site["lat"],
                site["lon"],
                site["timezone"],
            )
            if not sunset_local:
                continue

            vis = self._process_visibility(
                site["lat"], site["lon"], sunset_local, criteria,
            )

            score = vis["moon_altitude"] + vis["elongation"]
            if score > best_score:
                best_score = score
                best_visibility = {**vis, "site": site["name"], "sunset": sunset_local}

            if vis["is_visible"]:
                any_visible = True

        if not best_visibility:
            return {"is_rukyat_day": True, "error": "Gagal menghitung data nasional."}

        return {
            "is_rukyat_day": True,
            "is_visible": any_visible,
            "sunset_time": best_visibility["sunset"].isoformat(),
            "altitude_at_sunset": best_visibility["moon_altitude"],
            "azimuth_at_sunset": best_visibility.get("azimuth_moon", 0),
            "azimuth_diff": best_visibility.get("azimuth_diff", 0),
            "elongation_at_sunset": best_visibility["elongation"],
            "moon_age_hours": best_visibility["moon_age_hours"],
            "criteria_used": criteria,
            "site_name": best_visibility["site"],
        }

    # ── Helper: Perhitungan Visibilitas ──────────────────────────────

    def _process_visibility(self, lat, lon, sunset_local, criteria):
        sunset_utc = sunset_local.astimezone(pytz.utc)
        sunset_jd = jd_from_datetime(sunset_utc)

        conj_jd = calculate_conjunction(sunset_jd)

        return calculate_visibility(
            sunset_utc, lat, lon, conj_jd, criteria=criteria,
        )
