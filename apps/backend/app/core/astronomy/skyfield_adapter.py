from skyfield.api import load, wgs84
from skyfield.almanac import fraction_illuminated
import numpy as np
from datetime import timezone


class SkyfieldAdapter:
    """
    Astronomical Adapter Layer.
    Bertanggung jawab hanya pada perhitungan astronomi murni.
    Tidak mengandung logika kalender atau keputusan fiqh.
    """

    _instance = None

    def __new__(cls, ephemeris_path="data/de421.bsp"):
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._initialize(ephemeris_path)
        return cls._instance

    def _initialize(self, ephemeris_path):
        self.ts = load.timescale()
        self.eph = load(ephemeris_path)

        self.earth = self.eph["earth"]
        self.moon = self.eph["moon"]
        self.sun = self.eph["sun"]

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
