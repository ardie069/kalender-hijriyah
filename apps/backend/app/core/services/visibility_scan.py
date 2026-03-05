"""
Visibility scan global — UGHC/KHGT global visibility scanning.

Dioptimasi untuk GEOCENTRIC:
- Observer geocentric (earth) → posisi bulan/matahari sama untuk semua grid points
  pada waktu sunset yang sama → pre-compute sekali, reuse altitude refraksi
- Cached grid, pre-computed conjunction, early-exit
"""

import logging
import time
from functools import lru_cache
import pytz
from datetime import datetime, time as dt_time

from ..astronomy.global_grid import generate_global_grid
from ..astronomy.criteria_registry import CRITERIA_REGISTRY
from ..astronomy.visibility import evaluate_visibility_geocentric_batch
from ..calendar.julian import jd_from_datetime
from ..cache.disk import _make_key, get_cache, set_cache
from .engine import (
    calculate_sunset,
    calculate_conjunction,
)

logger = logging.getLogger(__name__)

# Module-level cached grids
_GRID_CACHE: dict[tuple, list] = {}


def _get_grid(lat_step: int, lon_step: int) -> list:
    """Cached grid generation — hanya di-build sekali per konfigurasi step."""
    key = (lat_step, lon_step)
    if key not in _GRID_CACHE:
        _GRID_CACHE[key] = generate_global_grid(lat_step=lat_step, lon_step=lon_step)
        logger.info("Grid generated: %d sites (lat=%d, lon=%d)", len(_GRID_CACHE[key]), lat_step, lon_step)
    return _GRID_CACHE[key]


def _get_astro():
    from app.deps.astronomy import get_provider
    return get_provider()


class GlobalVisibilityRegistry:

    @staticmethod
    @lru_cache(maxsize=64)
    def _cached_scan(date_key, criteria, lat_step, lon_step):
        p = _get_astro()
        sites = _get_grid(lat_step, lon_step)

        criteria_def = CRITERIA_REGISTRY.get(criteria)
        is_geocentric = criteria_def and criteria_def.get("observer") == "geocentric"

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

        # Pre-compute conjunction SEKALI untuk seluruh grid
        noon_dt = datetime.combine(date_key, dt_time(12, 0)).replace(tzinfo=pytz.utc)
        noon_jd = jd_from_datetime(noon_dt)
        shared_conj_jd = calculate_conjunction(noon_jd)

        found_before_24 = False

        # Cache geocentric positions per-sunset-time
        # Untuk geocentric: observer = earth → posisi bulan/matahari SAMA untuk semua titik
        # yang memiliki waktu sunset UTC yang persis sama (jarang terjadi).
        # Tapi kita bisa cache per-sunset-JD (dibulatkan ke menit) untuk reuse.
        _geo_cache: dict[int, dict] = {}  # key: sunset timestamp menit

        for site in sites:
            offset = int(site["lon"] / 15)
            tz_name = f"Etc/GMT{-offset:+d}" if offset != 0 else "UTC"

            sunset_local = calculate_sunset(date_key, site["lat"], site["lon"], tz_name)
            if not sunset_local:
                continue

            sunset_utc = sunset_local.astimezone(pytz.utc)

            # Geocentric optimization: cache positions per menit
            precomputed = None
            if is_geocentric:
                # Bucket per menit (60 detik presisi) — sunset di bujur yang sama ≈ waktu UTC sama
                bucket_key = int(sunset_utc.timestamp() // 60)

                if bucket_key not in _geo_cache:
                    # Hitung geocentric positions — observer = earth, SAMA untuk semua titik
                    t = p.ts.from_datetime(sunset_utc)
                    observer = p.earth
                    moon_app = observer.at(t).observe(p.moon).apparent()
                    sun_app = observer.at(t).observe(p.sun).apparent()
                    elongation_deg = float(moon_app.separation_from(sun_app).degrees)

                    _geo_cache[bucket_key] = {
                        "moon_app": moon_app,
                        "sun_app": sun_app,
                        "elongation_deg": elongation_deg,
                    }

                precomputed = _geo_cache[bucket_key]

            # Evaluasi visibilitas
            if is_geocentric:
                vis = evaluate_visibility_geocentric_batch(
                    sunset_utc, site["lat"], site["lon"], shared_conj_jd,
                    p.ts, p.sun, p.moon, p.earth,
                    criteria_name=criteria,
                    precomputed=precomputed,
                )
            else:
                # Fallback untuk topocentric (non-UGHC criteria)
                from ..astronomy.visibility import evaluate_visibility
                vis = evaluate_visibility(
                    sunset_utc, site["lat"], site["lon"], shared_conj_jd,
                    p.ts, p.sun, p.moon, p.earth,
                    criteria=criteria,
                )

            if vis.get("is_visible"):
                if sunset_utc <= threshold_24utc:
                    result["anywhere_before_24utc"] = True
                    found_before_24 = True
                else:
                    result["anywhere_after_24utc"] = True

                if site.get("is_america", False) or (-170 <= site["lon"] <= -35):
                    result["america_visible"] = True

                score = vis["moon_altitude"] + vis["elongation"]
                if score > best_score:
                    best_score = score
                    result["best_visibility"] = {**vis, "lat": site["lat"], "lon": site["lon"]}

                # Early-exit: Butir 1 tercapai + sudah punya best → cukup
                if found_before_24 and result["best_visibility"] is not None:
                    break

        return result

    @classmethod
    def scan_global(
        cls,
        date,
        criteria,
        lat_step=5,
        lon_step=10,
    ):
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
