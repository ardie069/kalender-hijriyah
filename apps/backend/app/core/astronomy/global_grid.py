def is_america_continent(lon: float) -> bool:
    """Check jika bujur berada di Benua Amerika (-170 s/d -30)."""
    return -170 <= lon <= -30


def generate_global_grid(lat_step=10, lon_step=15, min_lat=-60, max_lat=60):
    sites = []
    lat = min_lat
    while lat <= max_lat:
        lon = -180
        while lon < 180:
            sites.append(
                {
                    "name": f"{lat}_{lon}",
                    "lat": lat,
                    "lon": lon,
                    "timezone": "UTC",
                    "is_america": is_america_continent(lon),
                }
            )
            lon += lon_step
        lat += lat_step
    return sites
