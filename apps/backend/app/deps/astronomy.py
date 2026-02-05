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

    earth = eph["earth"]
    sun = eph["sun"]
    moon = eph["moon"]
except Exception as e:
    print(f"Load Astronomi Gagal: {e}")
