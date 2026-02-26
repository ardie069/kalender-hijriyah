from skyfield.api import wgs84
from skyfield.almanac import fraction_illuminated
import numpy as np
from datetime import timezone

from ...deps.astronomy import get_provider


class SkyfieldAdapter:
    """
    Astronomical Adapter Layer.
    Bertanggung jawab hanya pada perhitungan astronomi murni.
    Tidak mengandung logika kalender atau keputusan fiqh.

    Sekarang menggunakan AstronomyProvider singleton, sehingga tidak
    ada duplikasi loading ephemeris.
    """

    _instance = None

    def __new__(cls, ephemeris_path=None):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._initialize()
        return cls._instance

    def _initialize(self):
        provider = get_provider()
        self.ts = provider.ts
        self.eph = provider.eph
        self.earth = provider.earth
        self.moon = provider.moon
        self.sun = provider.sun

    # ==========================================================
    # INTERNAL UTILITIES
    # ==========================================================

    def _ensure_utc(self, dt):
        if dt.tzinfo is None:
            raise ValueError("Datetime must be timezone-aware UTC")
        return dt.astimezone(timezone.utc)

    def get_observer(self, lat: float, lon: float, elevation: float = 0):
        return self.earth + wgs84.latlon(lat, lon, elevation_m=elevation)

    # ==========================================================
    # MOON TELEMETRY
    # ==========================================================

    def get_moon_telemetry(self, dt_utc, lat: float, lon: float):
        """
        Menghasilkan telemetry bulan toposentris:
        - altitude
        - azimuth
        - elongation
        - illumination
        - distance
        """

        dt_utc = self._ensure_utc(dt_utc)
        t = self.ts.from_datetime(dt_utc)

        observer = self.get_observer(lat, lon)

        moon_app = observer.at(t).observe(self.moon).apparent()
        sun_app = observer.at(t).observe(self.sun).apparent()

        alt, az, distance = moon_app.altaz()

        elongation = moon_app.separation_from(sun_app).degrees

        illumination = fraction_illuminated(self.eph, "moon", t) * 100

        return {
            "altitude": float(alt.degrees),
            "azimuth": float(az.degrees),
            "elongation": float(elongation),
            "illumination": float(illumination),
            "distance_km": float(distance.km),
            "timestamp": dt_utc,
        }

    # ==========================================================
    # PHASE ANGLE (CONSISTENT APPARENT FRAME)
    # ==========================================================

    def get_moon_phase_angle(self, dt_utc):
        """
        Menghitung phase angle menggunakan apparent positions.
        Digunakan untuk menentukan waxing / waning.
        """

        dt_utc = self._ensure_utc(dt_utc)
        t = self.ts.from_datetime(dt_utc)

        sun_app = self.sun.at(t).apparent()
        moon_app = self.moon.at(t).apparent()
        earth_app = self.earth.at(t).apparent()

        sm = moon_app.position.au - sun_app.position.au
        em = moon_app.position.au - earth_app.position.au

        cos_angle = np.dot(sm, em) / (np.linalg.norm(sm) * np.linalg.norm(em))

        angle = np.degrees(np.arccos(np.clip(cos_angle, -1.0, 1.0)))

        return float(angle)

    def get_altaz_data(self, dt, lat, lon):
        """Ambil data Altitude dan Azimuth untuk Matahari dan Bulan"""
        observer = self.get_observer(lat, lon)
        t = self.ts.from_datetime(dt)

        # Hitung posisi Matahari
        sun_astrometric = observer.at(t).observe(self.sun)
        sun_alt, sun_az, _ = sun_astrometric.apparent().altaz()

        # Hitung posisi Bulan
        moon_astrometric = observer.at(t).observe(self.moon)
        moon_alt, moon_az, _ = moon_astrometric.apparent().altaz()

        return {
            "sun_az": sun_az.degrees,
            "sun_alt": sun_alt.degrees,
            "moon_az": moon_az.degrees,
            "moon_alt": moon_alt.degrees,
        }
