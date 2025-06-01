from skyfield.api import load

# Load efemeris dan timescale sekali saja
eph = load('data/de421.bsp')
ts = load.timescale()

# Akses benda langit umum
earth = eph['earth']
sun = eph['sun']
moon = eph['moon']
