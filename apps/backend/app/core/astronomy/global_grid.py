def is_america_continent(lon: float) -> bool:
    """Check jika bujur berada di Benua Amerika (-170 s/d -35)."""
    return -170 <= lon <= -35


def generate_global_grid(lat_step=5, lon_step=10, min_lat=-65, max_lat=65):
    """
    Grid lebih rapat (5/10) biar gak ada hilal yang nyelip.
    Timezone dinamis berdasarkan bujur.
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
