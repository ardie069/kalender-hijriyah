from pydantic import BaseModel


class MoonTelemetry(BaseModel):
    altitude: float
    azimuth: float
    illumination: float
    elongation: float
    distance_km: str
    age_days: float


class MoonStatus(BaseModel):
    phase_name: str
    is_waning: bool
    criteria_used: str
    is_visible: bool
    is_rukyat_time: bool
    observation_ref: str


class MoonInfoResponse(BaseModel):
    telemetry: MoonTelemetry
    status: MoonStatus
    timestamp: str
