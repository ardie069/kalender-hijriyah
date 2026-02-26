import pytz

from .engine import (
    calculate_baseline_hijri,
    calculate_sunset,
    calculate_conjunction,
    calculate_visibility,
)

from ..calendar.julian import jd_from_datetime
from ..config import REGIONAL_RUKYAT_CONFIG

HIJRI_MONTHS = [
    "Muharam",
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
]


class RukyatService:
    """
    Unified Rukyat Service — menangani evaluasi rukyat lokal (single-site)
    dan nasional (multi-site) dalam satu class modular.
    """

    def evaluate(self, context, region: str = "indonesia", mode: str = "local"):
        baseline, _ = calculate_baseline_hijri(
            context.now_local, context.timezone, context.ts
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
        else:
            return self._evaluate_local(context, criteria)

    # ── Local (single-site) ───────────────────────────────────────────

    def _evaluate_local(self, context, criteria: str):
        sunset_local = calculate_sunset(
            context.now_local.date(),
            context.lat,
            context.lon,
            context.timezone,
            context.ts,
            context.eph,
        )
        if not sunset_local:
            return {
                "is_rukyat_day": True,
                "error": "Sunset could not be calculated",
            }

        sunset_utc = sunset_local.astimezone(pytz.utc)
        sunset_jd = jd_from_datetime(sunset_utc, context.ts)

        conj_jd = calculate_conjunction(
            sunset_jd, context.ts, context.earth, context.sun, context.moon
        )

        vis = calculate_visibility(
            sunset_utc,
            context.lat,
            context.lon,
            conj_jd,
            context.ts,
            context.sun,
            context.moon,
            context.earth,
            criteria=criteria,
        )

        return {
            "is_rukyat_day": True,
            "sunset_time": sunset_local.isoformat(),
            "moon_position": {
                "altitude": vis["moon_altitude"],
                "elongation": vis["elongation"],
                "moon_age_hours": vis["moon_age_hours"],
            },
            "is_visible": vis["is_visible"],
            "criteria_used": criteria,
        }

    # ── National (multi-site) ─────────────────────────────────────────

    def _evaluate_national(self, context, region_config: dict, criteria: str):
        sites = region_config["sites"]

        any_visible = False
        best_visibility = None
        best_score = -999

        for site in sites:
            sunset_local = calculate_sunset(
                context.now_local.date(),
                site["lat"],
                site["lon"],
                site["timezone"],
                context.ts,
                context.eph,
            )
            if not sunset_local:
                continue

            sunset_utc = sunset_local.astimezone(pytz.utc)
            sunset_jd = jd_from_datetime(sunset_utc, context.ts)

            conj_jd = calculate_conjunction(
                sunset_jd,
                context.ts,
                context.earth,
                context.sun,
                context.moon,
            )

            vis = calculate_visibility(
                sunset_utc,
                site["lat"],
                site["lon"],
                conj_jd,
                context.ts,
                context.sun,
                context.moon,
                context.earth,
                criteria=criteria,
            )

            score = vis["moon_altitude"] + vis["elongation"]

            if score > best_score:
                best_score = score
                best_visibility = {
                    **vis,
                    "site": site["name"],
                }

            if vis["is_visible"]:
                any_visible = True

        # Build moon_position from best site data
        moon_position = None
        if best_visibility:
            moon_position = {
                "altitude": best_visibility["moon_altitude"],
                "elongation": best_visibility["elongation"],
                "moon_age_hours": best_visibility.get("moon_age_hours", 0),
            }

        return {
            "is_rukyat_day": True,
            "is_visible": any_visible,
            "is_visible_national": any_visible,
            "criteria_used": criteria,
            "best_site": best_visibility.get("site") if best_visibility else None,
            "moon_position": moon_position,
        }
