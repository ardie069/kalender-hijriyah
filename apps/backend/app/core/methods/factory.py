# app/core/methods/factory.py
from app.deps import astronomy  # Langit-langit (ts, eph, sun, moon, earth)
from .base import HijriContext
from .umm_al_qura import UmmAlQuraMethod
from .local_hisab import LocalHisabMethod
from .local_rukyat import LocalRukyatMethod
from .ughc_topocentric import UGHCMethod


class ContextFactory:
    @staticmethod
    def create_context(lat: float, lon: float, timezone: str, dt) -> HijriContext:
        return HijriContext(
            lat=lat,
            lon=lon,
            timezone=timezone,
            now_local=dt,
            ts=astronomy.ts,
            eph=astronomy.eph,
            sun=astronomy.sun,
            moon=astronomy.moon,
            earth=astronomy.earth,
        )


_instances = {}


def get_method_instance(method: str):
    classes = {
        "umm_al_qura": UmmAlQuraMethod,
        "local_hisab": LocalHisabMethod,
        "local_rukyat": LocalRukyatMethod,
        "ughc": UGHCMethod,
    }
    if method not in classes:
        raise ValueError(f"Metode '{method}' tidak ditemukan.")
    if method not in _instances:
        _instances[method] = classes[method]()
    return _instances[method]
