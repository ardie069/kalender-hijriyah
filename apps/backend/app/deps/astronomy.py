from pathlib import Path
from skyfield.api import load  # type: ignore

BASE_DIR = Path(__file__).resolve().parents[2]
DATA_DIR = BASE_DIR / "data"
DATA_DIR.mkdir(exist_ok=True)

EPHEMERIS_FILE = DATA_DIR / "de421.bsp"

ts = load.timescale()

if EPHEMERIS_FILE.exists():
    eph = load(str(EPHEMERIS_FILE))
else:
    eph = load("data/de421.bsp")
    eph.save(str(EPHEMERIS_FILE))  # type: ignore

earth = eph["earth"]
sun = eph["sun"]
moon = eph["moon"]
