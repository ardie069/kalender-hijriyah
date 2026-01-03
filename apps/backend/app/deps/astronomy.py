from pathlib import Path
from skyfield.api import load  # type: ignore

BASE_DIR = Path(__file__).resolve().parents[2]
DATA_DIR = BASE_DIR / "data"

EPHEMERIS_FILE = DATA_DIR / "de421.bsp"

if not EPHEMERIS_FILE.exists():
    raise FileNotFoundError(f"Ephemeris file not found: {EPHEMERIS_FILE}")

ts = load.timescale()
eph = load(str(EPHEMERIS_FILE))

earth = eph["earth"]
sun = eph["sun"]
moon = eph["moon"]
