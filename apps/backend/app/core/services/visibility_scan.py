"""
Visibility scan global — UGHC/KHGT global visibility scanning.

Registry anti-pattern dihapus. Astronomy objects diakses dari provider.
"""

from functools import lru_cache
import pytz
from datetime import datetime, time

from ..astronomy.global_grid import generate_global_grid
from ..calendar.julian import jd_from_datetime
from ..cache.disk import _make_key, get_cache, set_cache
from .engine import (
    calculate_sunset,
    calculate_conjunction,
    calculate_visibility,
)


def _get_astro():
    from app.deps.astronomy import get_provider
    return get_provider()


class GlobalVisibilityRegistry:

    @staticmethod
    @lru_cache(maxsize=64)
    def _cached_scan(date_key, criteria, lat_step, lon_step):
        sites = generate_global_grid(lat_step=lat_step, lon_step=lon_step)

        result = {
            "anywhere_before_24utc": False,
            "america_visible": False,
            "anywhere_after_24utc": False,
            "best_visibility": None,
        }

        best_score = -999
        end_of_day_utc = datetime.combine(date_key, time(23, 59, 59)).replace(
            tzinfo=pytz.utc
        )

        for site in sites:
            sunset_local = calculate_sunset(
                date_key, site["lat"], site["lon"], "UTC",
            )
            if not sunset_local:
                continue

            sunset_utc = sunset_local.astimezone(pytz.utc)
            if sunset_utc.date() != date_key:
                continue

            sunset_jd = jd_from_datetime(sunset_utc)
            conj_jd = calculate_conjunction(sunset_jd)

            vis = calculate_visibility(
                sunset_utc,
                site["lat"],
                site["lon"],
                conj_jd,
                criteria=criteria,
            )

            if vis["is_visible"]:
                if sunset_utc <= end_of_day_utc:
                    result["anywhere_before_24utc"] = True
                else:
                    result["anywhere_after_24utc"] = True

                if site["is_america"]:
                    result["america_visible"] = True

                score = vis["moon_altitude"] + vis["elongation"]
                if score > best_score:
                    best_score = score
                    result["best_visibility"] = {
                        **vis,
                        "lat": site["lat"],
                        "lon": site["lon"],
                    }

        return result

    @classmethod
    def scan_global(
        cls,
        date,
        criteria,
        ts=None,
        eph=None,
        sun=None,
        moon=None,
        earth=None,
        lat_step=10,
        lon_step=15,
    ):
        """
        ts/eph/sun/moon/earth params kept for backward compat
        but are no longer used — accessed from provider.
        """

        cache_key = _make_key(date, criteria, lat_step, lon_step)

        cached = get_cache(cache_key)
        if cached:
            return cached

        result = cls._cached_scan(date, criteria, lat_step, lon_step)

        set_cache(cache_key, result)

        return result
