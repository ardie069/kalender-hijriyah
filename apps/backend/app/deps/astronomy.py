from pathlib import Path
from skyfield.api import Loader

BASE_DIR = Path(__file__).resolve().parents[2]
DATA_DIR = BASE_DIR / "data"
DATA_DIR.mkdir(exist_ok=True)

loader = Loader(str(DATA_DIR))

ts = loader.timescale()
eph = loader("de421.bsp")

earth = eph["earth"]
sun = eph["sun"]
moon = eph["moon"]
