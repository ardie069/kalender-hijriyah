from datetime import datetime
from typing import Literal, Optional

from . import NumPyBaseModel

HijriMethod = Literal["umm_al_qura", "local_hisab", "local_rukyat", "ughc"]


class HijriDateSchema(NumPyBaseModel):
    year: int
    month: int
    day: int


class LocationSchema(NumPyBaseModel):
    lat: float
    lon: float
    timezone: str


class HijriAstronomicalDataSchema(NumPyBaseModel):
    """Data astronomis lengkap, termasuk koordinat referensi (untuk UGHC)."""

    moon_altitude: float
    elongation: float
    moon_age: Optional[float] = None
    is_visible: bool
    lat: Optional[float] = None
    lon: Optional[float] = None
    location_name: Optional[str] = None


class HijriExplanationSchema(NumPyBaseModel):
    method: HijriMethod
    after_sunset: bool
    criteria_used: str
    reasoning: list[str]
    decision: str
    astronomical_data: Optional[HijriAstronomicalDataSchema] = None


class HijriDateResponse(NumPyBaseModel):
    method: HijriMethod
    location: LocationSchema
    hijri_date: HijriDateSchema
    explanation: Optional[HijriExplanationSchema] = None
    generated_at: datetime


class HijriEndMonthResponse(NumPyBaseModel):
    method: HijriMethod
    location: LocationSchema
    generated_at: datetime
    today: HijriDateSchema
    estimated_end_of_month: Optional[HijriDateSchema] = None
    visibility: Optional[HijriAstronomicalDataSchema] = None
    message: Optional[str] = None


class GregorianMonthData(NumPyBaseModel):
    month_id: int
    month_name: str
    total_days: int
    day_1_weekday: int  # 0=Senin (Python weekday)


class GregorianYearResponse(NumPyBaseModel):
    year: int
    months: list[GregorianMonthData]
