from functools import lru_cache
import pytz

from .julian import jd_from_datetime
from .disk_cache import _make_key, get_cache, set_cache
from .astronomy_engine import (
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

        from .global_grid import generate_global_grid

        sites = generate_global_grid(
            lat_step=lat_step,
            lon_step=lon_step,
        )

        visible = False
        best_visibility = None
        best_score = -999
        first_visible_site = None

        for site in sites:

            sunset_local = calculate_sunset(
                date_key,
                site["lat"],
                site["lon"],
                site["timezone"],
                ts,
                eph,
            )

            if not sunset_local:
                continue

            sunset_utc = sunset_local.astimezone(pytz.utc)
            sunset_jd = jd_from_datetime(sunset_utc, ts)

            conj_jd = calculate_conjunction(
                sunset_jd,
                ts,
                earth,
                sun,
                moon,
            )

            vis = calculate_visibility(
                sunset_utc,
                site["lat"],
                site["lon"],
                conj_jd,
                ts,
                sun,
                moon,
                earth,
                criteria=criteria,
            )

            score = vis["moon_altitude"] + vis["elongation"]

            if score > best_score:
                best_score = score
                best_visibility = {
                    **vis,
                    "location": site.get("name"),
                }

            if vis["is_visible"] and not visible:
                visible = True
                first_visible_site = site

        return {
            "visible": visible,
            "best_visibility": best_visibility,
            "first_visible_site": first_visible_site,
        }

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
