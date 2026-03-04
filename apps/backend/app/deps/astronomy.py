"""
Unified Astronomy Provider — Single source of truth untuk semua object Skyfield.

Menggabungkan bekas `deps/astronomy.py` (module-level globals) dan
`SkyfieldAdapter` (singleton class) menjadi satu provider.
"""

from __future__ import annotations
from dataclasses import dataclass, field
from pathlib import Path
from skyfield.api import Loader


_PROVIDER: AstronomyProvider | None = None

CURRENT_DIR = Path(__file__).resolve().parent
DATA_DIR = CURRENT_DIR.parent.parent / "data"

if not DATA_DIR.exists():
    DATA_DIR.mkdir(parents=True, exist_ok=True)


@dataclass(frozen=True)
class AstronomyProvider:
    """Immutable container for all Skyfield objects."""

    ts: object
    eph: object
    earth: object
    sun: object
    moon: object


def _init_provider() -> AstronomyProvider:
    loader = Loader(str(DATA_DIR), expire=False)
    try:
        ts = loader.timescale()
        eph = loader("de440.bsp")
        earth, sun, moon = eph["earth"], eph["sun"], eph["moon"]
    except Exception as e:
        raise RuntimeError(f"Gagal memuat data astronomi: {e}")

    return AstronomyProvider(ts=ts, eph=eph, earth=earth, sun=sun, moon=moon)


def get_provider() -> AstronomyProvider:
    """Get or create the singleton AstronomyProvider."""
    global _PROVIDER
    if _PROVIDER is None:
        _PROVIDER = _init_provider()
    return _PROVIDER


# ── Backward-compatible module-level exports ─────────────────────
# Some modules still import `from app.deps.astronomy import ts, eph, ...`
# These lazy properties keep them working during the migration.

def __getattr__(name: str):
    """Lazy module-level access for backward compat."""
    if name in ("ts", "eph", "earth", "sun", "moon"):
        p = get_provider()
        return getattr(p, name)
    raise AttributeError(f"module {__name__!r} has no attribute {name!r}")
