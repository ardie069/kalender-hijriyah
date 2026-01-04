from pydantic import BaseModel
from datetime import datetime
from typing import Literal, Optional

HijriMethod = Literal["global", "hisab", "rukyat"]


class HijriDateSchema(BaseModel):
    year: int
    month: int
    day: int


class LocationSchema(BaseModel):
    lat: float
    lon: float
    timezone: str


class HijriDateResponse(BaseModel):
    method: HijriMethod
    location: LocationSchema
    hijri_date: HijriDateSchema
    generated_at: datetime


class HijriEndMonthEstimate(BaseModel):
    hijri: Optional[HijriDateSchema] = None
    moon_age: Optional[float] = None
    moon_altitude: Optional[float] = None
    elongation: Optional[float] = None
    is_visible: Optional[bool] = None


class HijriEndMonthResponse(BaseModel):
    method: str
    location: LocationSchema
    generated_at: datetime

    today: HijriDateSchema
    estimated_end_of_month: Optional[HijriEndMonthEstimate] = None
    message: Optional[str] = None
