export type Method = "global" | "hisab" | "rukyat";

/* =======================
   BASIC
======================= */
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

/* =======================
   /hijri-date
======================= */
export interface HijriDateResponse {
  method: Method;
  location: Location;
  hijri_date: HijriDate;
  generated_at: string; // ISO datetime
}

/* =======================
   /hijri-end-month
======================= */
export interface HijriVisibility {
  moon_altitude: number;
  elongation: number;
  moon_age: number;
  is_visible: boolean;
}

export interface HijriEndMonthResponse {
  method: Method;
  today: HijriDate;
  estimated_end_of_month: HijriDate | null;
  visibility?: HijriVisibility;
  message?: string | null;
}

/* =======================
   /hijri-explain
======================= */
export interface HijriExplainResponse {
  method: Method;
  decision: string;
  after_sunset: boolean;
  final_hijri_date: HijriDate;
}
