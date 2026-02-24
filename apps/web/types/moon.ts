export interface MoonTelemetry {
  altitude: number;
  azimuth: number;
  illumination: number;
  elongation: number;
  distance_km: string;
  age_days: number;
}

export interface MoonStatus {
  phase_name: string;
  is_waning: boolean;
  is_mabims_met: boolean;
  is_rukyat_time: boolean;
  observation_ref: string;
}

export interface MoonDataResponse {
  telemetry: MoonTelemetry;
  status: MoonStatus;
  timestamp: string;
}
