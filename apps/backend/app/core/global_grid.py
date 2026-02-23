def generate_global_grid(
    lat_step=10,
    lon_step=15,
    min_lat=-60,
    max_lat=60,
):
    """
    Generate grid lokasi global untuk scan visibilitas hilal.
    """

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
                }
            )
            lon += lon_step
        lat += lat_step

    return sites
