from __future__ import annotations


class HijriContext:
    """
    Context object yang membawa semua parameter yang dibutuhkan
    oleh method kalkulasi Hijriyah.
    """

    def __init__(self, lat, lon, timezone, now_local, ts, eph, sun, moon, earth):
        self.lat = lat
        self.lon = lon
        self.timezone = timezone
        self.now_local = now_local
        self.ts = ts
        self.eph = eph
        self.sun = sun
        self.moon = moon
        self.earth = earth

    @classmethod
    def from_request(cls, lat: float, lon: float, timezone: str, now_local) -> HijriContext:
        """
        Factory method: buat HijriContext dari parameter request.
        Astronomy objects diambil otomatis dari AstronomyProvider singleton.
        """
        from app.deps.astronomy import get_provider

        p = get_provider()
        return cls(
            lat=lat,
            lon=lon,
            timezone=timezone,
            now_local=now_local,
            ts=p.ts,
            eph=p.eph,
            sun=p.sun,
            moon=p.moon,
            earth=p.earth,
        )


class HijriResult:
    def __init__(self, hijri_date, explanation=None, metadata=None):
        self.hijri_date = hijri_date
        self.explanation = explanation or {}
        self.metadata = metadata or {}


class BaseHijriMethod:
    def calculate(self, context: HijriContext) -> HijriResult:
        raise NotImplementedError
