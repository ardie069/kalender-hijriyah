"""
Visibility scan global — UGHC/KHGT global visibility scanning.

Registry anti-pattern dihapus. Astronomy objects diakses dari provider.
"""

import logging
import time
from functools import lru_cache
import pytz
from datetime import datetime, time as dt_time

from ..astronomy.global_grid import generate_global_grid
from ..calendar.julian import jd_from_datetime
from ..cache.disk import _make_key, get_cache, set_cache
from .engine import (
    calculate_sunset,
    calculate_conjunction,
    calculate_visibility,
)

logger = logging.getLogger(__name__)


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

        # 24:00 UTC pada hari H (Batas penentuan KHGT)
        threshold_24utc = datetime.combine(date_key, dt_time(23, 59, 59)).replace(
            tzinfo=pytz.utc
        )

        for site in sites:
            # 1. Ambil sunset lokal untuk tanggal tersebut
            sunset_local = calculate_sunset(
                date_key, site["lat"], site["lon"], site.get("timezone", "UTC")
            )
            if not sunset_local:
                continue

            sunset_utc = sunset_local.astimezone(pytz.utc)

            # ─── FIX: JANGAN SKIP JIKA TANGGAL BERBEDA ───
            # Kita hanya perlu pastikan sunset ini terjadi dalam rentang
            # wajar setelah ijtima' (misal dalam 24 jam dari date_key)

            sunset_jd = jd_from_datetime(sunset_utc)
            conj_jd = calculate_conjunction(sunset_jd)

            # 2. Hitung visibilitas
            vis = calculate_visibility(
                sunset_utc,
                site["lat"],
                site["lon"],
                conj_jd,
                criteria=criteria,
            )

            if vis["is_visible"]:
                # Cek Batas 24:00 UTC
                if sunset_utc <= threshold_24utc:
                    result["anywhere_before_24utc"] = True
                else:
                    result["anywhere_after_24utc"] = True

                # Cek Sinyal Amerika (Kunci UGHC)
                if site.get("is_america", False):
                    result["america_visible"] = True

                # Tracking Hilal Terbaik
                score = vis["moon_altitude"] + vis["elongation"]
                if score > best_score:
                    best_score = score
                    result["best_visibility"] = {
                        **vis,
                        "lat": site["lat"],
                        "lon": site["lon"],
                        "sunset_utc": sunset_utc.isoformat(),
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

        # TTL 7 hari — data visibility per tanggal stabil
        cached = get_cache(cache_key, ttl_seconds=604800)
        if cached:
            logger.debug("scan_global CACHE HIT: date=%s criteria=%s", date, criteria)
            return cached

        t0 = time.perf_counter()
        result = cls._cached_scan(date, criteria, lat_step, lon_step)
        elapsed = time.perf_counter() - t0

        set_cache(cache_key, result)

        logger.info(
            "scan_global COMPUTED: date=%s criteria=%s visible=%s (%.3fs)",
            date, criteria, result.get("anywhere_before_24utc"), elapsed,
        )

        return result
