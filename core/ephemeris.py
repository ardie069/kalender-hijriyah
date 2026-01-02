from pathlib import Path
from skyfield.api import load

BASE_DIR = Path(__file__).resolve().parent.parent
DATA_DIR = BASE_DIR / "data"

eph = load(str(DATA_DIR / "de421.bsp"))
ts = load.timescale()

earth = eph["earth"]
sun = eph["sun"]
moon = eph["moon"]
