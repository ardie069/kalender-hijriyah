import { Method } from "@/types/hijri";
import {
  YearlyCalendarResponse,
  GregorianYearResponse,
} from "@/types/calendar";
import { request } from "./client";

export async function fetchYearlyCalendar(
  year: number,
  lat: number,
  lon: number,
  method: Method,
  timezone: string = "Asia/Jakarta",
): Promise<YearlyCalendarResponse> {
  const params = new URLSearchParams({
    year: String(year),
    lat: String(lat),
    lon: String(lon),
    method,
    timezone,
  });

  return request<YearlyCalendarResponse>(`/calendar/year?${params.toString()}`);
}

export async function fetchGregorianCalendar(
  year: number,
): Promise<GregorianYearResponse> {
  const params = new URLSearchParams({ year: String(year) });
  return request<GregorianYearResponse>(
    `/calendar/gregorian/year?${params.toString()}`,
  );
}
