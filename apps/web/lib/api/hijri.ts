import type {
  HijriDateResponse,
  HijriEndMonthResponse,
  Method,
} from "@/types/hijri";
import { request } from "./client";

/**
 * /hijri-date
 */
export function fetchHijriDate(
  lat: number,
  lon: number,
  method: Method,
  timezone: string,
) {
  const params = new URLSearchParams({
    lat: String(lat),
    lon: String(lon),
    method,
    timezone,
  });

  return request<HijriDateResponse>(`/hijri-date?${params.toString()}`);
}

/**
 * /hijri-end-month
 */
export function fetchHijriEndMonth(
  lat: number,
  lon: number,
  method: Method,
  timezone: string,
) {
  const params = new URLSearchParams({
    lat: String(lat),
    lon: String(lon),
    method,
    timezone,
  });

  return request<HijriEndMonthResponse>(
    `/hijri-end-month?${params.toString()}`,
  );
}
