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
  criteria_used: string;
  is_visible: boolean;
  is_rukyat_time: boolean;
  observation_ref: string;
}

export interface MoonDataResponse {
  telemetry: MoonTelemetry;
  status: MoonStatus;
  timestamp: string;
}

export interface MoonPosition {
  altitude: number;
  elongation: number;
  moon_age_hours: number;
}

export interface RukyatResponse {
  is_rukyat_day: boolean;
  sunset_time?: string;
  moon_position?: MoonPosition;
  is_visible?: boolean;
  criteria_used?: string;
  is_visible_national?: boolean;
  best_site?: string;
  hijri_date?: {
    year: number;
    month: number;
    day: number;
    month_name: string;
  };
  error?: string;
}
