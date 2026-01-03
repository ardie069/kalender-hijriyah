from datetime import datetime
import pytz # type: ignore

from app.deps.astronomy import ts, eph, sun, moon, earth
from app.core.hijri_calculator import get_hijri_date

now_local = datetime.now(pytz.timezone("Asia/Jakarta"))

result = get_hijri_date(
    lat=-7.9666,
    lon=112.6326,
    method="hisab",
    timezone="Asia/Jakarta",
    now_local=now_local,
    ts=ts,
    eph=eph,
    sun=sun,
    moon=moon,
    earth=earth,
)

print(result)
