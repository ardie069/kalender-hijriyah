from skyfield.api import wgs84


def get_moon_altitude(dt_utc, lat, lon, ts, moon):
    """Mengambil tinggi hilal dalam derajat (standard float)."""
    t = ts.from_datetime(dt_utc)
    observer = wgs84.latlon(lat, lon)
    alt, _, _ = observer.at(t).observe(moon).apparent().altaz()
    return float(alt.degrees)


def evaluate_visibility(
    sunset_dt_utc,
    lat,
    lon,
    conjunction_jd,
    ts,
    sun,
    moon,
    earth,
    criteria="MABIMS",
):
    """
    Evaluasi visibilitas hilal berdasarkan kriteria tertentu.
    Mengonversi semua output NumPy ke tipe data Python standar.
    """
    # 1. Inisialisasi Waktu dan Observer
    t = ts.from_datetime(sunset_dt_utc)
    topos = wgs84.latlon(latitude_degrees=lat, longitude_degrees=lon)
    observer = earth + topos

    # 2. Perhitungan Astronomi (Apparent Position)
    moon_app = observer.at(t).observe(moon).apparent()
    sun_app = observer.at(t).observe(sun).apparent()

    alt_moon_deg = alt_moon, _, _ = moon_app.altaz()
    alt_deg = float(alt_moon_deg[0].degrees)
    elongation_deg = float(moon_app.separation_from(sun_app).degrees)

    # Perhitungan umur bulan dalam jam
    moon_age_hours = float((t.tt - conjunction_jd) * 24)

    # 3. Dialektika Kriteria (Logic)
    if criteria == "Wujudul Hilal":
        # Standar: Matahari terbenam dulu, lalu bulan di atas ufuk (alt > 0)
        is_visible = alt_deg > 0
    else:
        # Standar MABIMS Baru (2022): Alt >= 3 & Elongasi >= 6.4
        # Umur bulan >= 8 jam biasanya jadi pendukung
        is_visible = (
            (alt_deg >= 3.0) and (elongation_deg >= 6.4) and (moon_age_hours >= 8.0)
        )

    # 4. Return Data (Clean & Serialized)
    return {
        "moon_altitude": alt_deg,
        "elongation": elongation_deg,
        "moon_age": moon_age_hours,
        "is_visible": bool(is_visible),  # WAJIB: Casting ke standard bool
        "criteria_used": str(criteria),
    }
