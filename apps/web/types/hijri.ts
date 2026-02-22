export type Method = "umm_al_qura" | "local_hisab" | "local_rukyat" | "ughc";

export interface HijriDate {
  year: number;
  month: number;
  day: number;
}

export interface Location {
  lat: number;
  lon: number;
  timezone: string;
}

export interface HijriAstronomicalData {
  moon_altitude?: number;
  elongation?: number;
  moon_age?: number;
  is_visible?: boolean;
  conjunction_before_sunset?: boolean;
  global_visible?: boolean;
}

export interface HijriExplanation {
  method: Method;
  criteria_used?: string;
  reasoning: string[];
  decision?: string;
  astronomical_data?: HijriAstronomicalData;
}

export interface HijriDateResponse {
  method: Method;
  location: Location;
  hijri_date: HijriDate;
  explanation?: HijriExplanation;
  generated_at: string;
}

export interface HijriEndMonthResponse {
  method: Method;
  location: Location;
  generated_at: string;
  today: HijriDate;
  estimated_end_of_month: HijriDate | null;
  estimated_end_of_month_gregorian?: string | null;
  visibility?: HijriAstronomicalData;
  message?: string;
}
