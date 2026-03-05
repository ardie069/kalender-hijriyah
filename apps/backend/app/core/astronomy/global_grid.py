from functools import lru_cache


def is_america_continent(lon: float) -> bool:
    """Check jika bujur berada di Benua Amerika (-170 s/d -35)."""
    return -170 <= lon <= -35


@lru_cache(maxsize=8)
def generate_global_grid(lat_step=5, lon_step=10, min_lat=-60, max_lat=60):
    """
    Grid titik observasi global.
    Di-cache via lru_cache agar tidak di-regenerate setiap pemanggilan.
    Return sebagai tuple of tuples (hashable) agar bisa di-cache oleh lru_cache.
    """
    sites = []
    lat = min_lat
    while lat <= max_lat:
        lon = -180
        while lon < 180:
            offset = round(lon / 15)
            tz_name = f"Etc/GMT{-offset:+d}" if offset != 0 else "UTC"

            sites.append(
                {
                    "name": f"{lat}_{lon}",
                    "lat": lat,
                    "lon": lon,
                    "timezone": tz_name,
                    "is_america": -170 <= lon <= -35,
                }
            )
            lon += lon_step
        lat += lat_step
    return sites
