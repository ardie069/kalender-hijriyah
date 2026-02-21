from pathlib import Path
from skyfield.api import Loader

CURRENT_DIR = Path(__file__).resolve().parent
DATA_DIR = CURRENT_DIR.parent.parent / "data"

if not DATA_DIR.exists():
    DATA_DIR.mkdir(parents=True, exist_ok=True)

loader = Loader(str(DATA_DIR), expire=False)

try:
    ts = loader.timescale()
    eph = loader("de421.bsp")
    earth, sun, moon = eph["earth"], eph["sun"], eph["moon"]
except Exception as e:
    raise RuntimeError(f"Gagal memuat data astronomi: {e}")
