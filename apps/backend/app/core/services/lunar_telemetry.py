from datetime import datetime, timedelta
from typing import Dict, Any
from ..astronomy.skyfield_adapter import SkyfieldAdapter


class LunarTelemetryService:
    def __init__(self, adapter: SkyfieldAdapter):
        self.adapter = adapter

    def get_full_analysis(
        self, lat: float, lon: float, dt: datetime = None
    ) -> Dict[str, Any]:
        """
        Menghasilkan paket data lengkap untuk Dashboard Info Bulan.
        """
        if dt is None:
            dt = datetime.now()

        # 1. Ambil Telemetry Dasar
        telemetry = self.adapter.get_moon_telemetry(dt, lat, lon)
        last_conjunction = self.adapter.get_last_conjunction(dt)
        age_delta = dt - last_conjunction
        age_days = age_delta.total_seconds() / 86400

        # 2. Dialektika Fase: Apakah membesar atau mengecil?
        # Kita bandingkan iluminasi sekarang dengan 1 jam ke depan
        dt_future = dt + timedelta(hours=1)
        tel_future = self.adapter.get_moon_telemetry(dt_future, lat, lon)

        is_waxing = tel_future["illumination"] > telemetry["illumination"]
        phase_name = self._determine_phase_name(telemetry["illumination"], is_waxing)
        
        is_relevant_time = is_waxing and age_days < 3 and telemetry["altitude"] > 0
        
        if is_relevant_time:
            is_met = bool(telemetry["altitude"] > 3 and telemetry["elongation"] > 6.4)
        else:
            is_met = False

        return {
            "telemetry": {
                "altitude": round(telemetry["altitude"], 2),
                "azimuth": round(telemetry["azimuth"], 2),
                "illumination": round(telemetry["illumination"], 2),
                "elongation": round(telemetry["elongation"], 2),
                "distance_km": f"{int(telemetry['distance_km']):,}",
                "age_days": round(age_days, 2),
            },
            "status": {
                "phase_name": phase_name,
                "is_waning": not is_waxing,
                "is_mabims_met": is_met,
                "is_rukyat_time": is_relevant_time,
                "observation_ref": f"{'North' if lat >= 0 else 'South'} Hemisphere",
            },
            "timestamp": dt.isoformat(),
        }

    def _determine_phase_name(self, illum: float, is_waxing: bool) -> str:
        """
        Logika Kategorisasi Fase Bulan (The Terminator Logic)
        """
        if illum < 1:
            return "Bulan Baru (Ijtimak)"
        if illum > 98:
            return "Bulan Purnama"

        if is_waxing:
            if illum < 45:
                return "Sabit Muda (Waxing Crescent)"
            if illum < 55:
                return "Kuartal Pertama"
            return "Cembung Awal (Waxing Gibbous)"
        else:
            if illum > 55:
                return "Cembung Akhir (Waning Gibbous)"
            if illum > 45:
                return "Kuartal Terakhir"
            return "Sabit Tua (Waning Crescent)"

    def calculate_lag_time(self, moonset: datetime, sunset: datetime) -> int:
        """Menghitung selisih waktu antara terbenam matahari dan bulan (dalam menit)"""
        if not moonset or not sunset:
            return 0
        diff = (moonset - sunset).total_seconds() / 60
        return int(max(0, diff))
