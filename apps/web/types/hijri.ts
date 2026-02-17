export type Method = "global" | "hisab" | "rukyat";

/**
 * Representasi tanggal Hijriyah
 */
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

/**
 * Data astronomis hilal untuk penjelasan logika
 */
export interface HijriAstronomicalData {
  moon_altitude: number;
  elongation: number;
  conjunction_before_sunset: boolean;
  is_visible: boolean;
}

/**
 * Penjelasan mendalam mengenai penentuan tanggal (Scientific Soul)
 */
export interface HijriExplanation {
  method: Method;
  after_sunset: boolean;
  criteria_used: string;
  reasoning: string[];
  astronomical_data?: HijriAstronomicalData;
}

/**
 * Response utama dari endpoint /hijri-date
 */
export interface HijriDateResponse {
  method: Method;
  location: Location;
  hijri_date: HijriDate;
  explanation?: HijriExplanation;
  generated_at: string;
}

/**
 * Data visibilitas untuk prediksi akhir bulan
 */
export interface HijriVisibility {
  moon_altitude: number;
  elongation: number;
  moon_age: number;
  is_visible: boolean;
}

/**
 * Response dari endpoint /hijri-end-month
 */
export interface HijriEndMonthResponse {
  method: Method;
  today: HijriDate;
  estimated_end_of_month: HijriDate | null;
  visibility?: HijriVisibility;
  message?: string | null;
}
