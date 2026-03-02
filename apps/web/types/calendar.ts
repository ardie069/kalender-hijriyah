import { HijriAstronomicalData, Method } from "./hijri";

export type DateSystem = "hijri" | "gregorian";
export type ViewMode = "monthly" | "yearly";

/** Data satu bulan Hijriyah dari backend /calendar/year */
export interface MonthlyCalendarData {
  month_id: number;
  month_name: string;
  total_days: 29 | 30;
  start_gregorian: string;
  day_1_weekday: number;
  visibility_at_29?: HijriAstronomicalData;
}

export interface YearlyCalendarResponse {
  year: number;
  method: Method;
  months: MonthlyCalendarData[];
}

/** Data satu bulan Masehi dari backend /calendar/gregorian/year */
export interface GregorianMonthData {
  month_id: number;
  month_name: string;
  total_days: number;
  day_1_weekday: number;
}

export interface GregorianYearResponse {
  year: number;
  months: GregorianMonthData[];
}

/**
 * Unified month data — dipakai di MonthGrid supaya satu interface
 * bisa render Hijriyah maupun Masehi.
 */
export interface UnifiedMonthData {
  month_id: number;
  month_name: string;
  total_days: number;
  day_1_weekday: number;
  start_gregorian?: string; // hanya ada di mode Hijriyah
}