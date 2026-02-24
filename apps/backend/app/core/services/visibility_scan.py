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

_TS_REGISTRY = {}
_EPH_REGISTRY = {}
_EARTH_REGISTRY = {}
_SUN_REGISTRY = {}
_MOON_REGISTRY = {}


class GlobalVisibilityRegistry:

    @staticmethod
    def _store_dependencies(ts, eph, sun, moon, earth):
        _TS_REGISTRY[id(ts)] = ts
        _EPH_REGISTRY[id(eph)] = eph
        _SUN_REGISTRY[id(sun)] = sun
        _MOON_REGISTRY[id(moon)] = moon
        _EARTH_REGISTRY[id(earth)] = earth

    @staticmethod
    @lru_cache(maxsize=64)
    def _cached_scan(
        date_key,
        criteria,
        lat_step,
        lon_step,
        ts_id,
        eph_id,
        sun_id,
        moon_id,
        earth_id,
    ):
        ts = _TS_REGISTRY[ts_id]
        eph = _EPH_REGISTRY[eph_id]
        sun = _SUN_REGISTRY[sun_id]
        moon = _MOON_REGISTRY[moon_id]
        earth = _EARTH_REGISTRY[earth_id]

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
                date_key, site["lat"], site["lon"], "UTC", ts, eph
            )
            if not sunset_local:
                continue

            sunset_utc = sunset_local.astimezone(pytz.utc)
            # Hilal hanya valid jika sunset di titik tersebut terjadi pada hari yang sama secara UTC
            # (Penting untuk membatasi scan agar tidak meluber ke hari berikutnya)
            if sunset_utc.date() != date_key:
                continue

            vis = calculate_visibility(
                sunset_utc,
                site["lat"],
                site["lon"],
                calculate_conjunction(
                    jd_from_datetime(sunset_utc, ts), ts, earth, sun, moon
                ),
                ts,
                sun,
                moon,
                earth,
                criteria=criteria,
            )

            if vis["is_visible"]:
                # Logika KHGT Butir 1 & 2
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
        ts,
        eph,
        sun,
        moon,
        earth,
        lat_step=10,
        lon_step=15,
    ):

        cache_key = _make_key(
            date,
            criteria,
            lat_step,
            lon_step,
        )

        cached = get_cache(cache_key)
        if cached:
            return cached

        cls._store_dependencies(ts, eph, sun, moon, earth)

        result = cls._cached_scan(
            date,
            criteria,
            lat_step,
            lon_step,
            id(ts),
            id(eph),
            id(sun),
            id(moon),
            id(earth),
        )

        set_cache(cache_key, result)

        return result
