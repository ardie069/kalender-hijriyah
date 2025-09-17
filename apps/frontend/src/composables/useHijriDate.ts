import { ref, watch, Ref } from "vue";
import { isLocationInJava, getWeton } from "../utils/weton";

// --- Type Definitions --- //
export interface HijriDate {
  day: number;
  month: number;
  year: number;
}

export interface GregorianDate {
  day: number;
  month: number;
  year: number;
}

export interface HijriEndMonth {
  hijri?: HijriDate;
  gregorian?: GregorianDate;
  moonAge?: number;
  moonAltitude?: number;
  elongation?: number;
}

export function useHijriDate(
  selectedMethod: Ref<string>,
  userTimezone: Ref<string>,
  API_BASE_URL: string
) {
  const lat = ref<number | null>(null);
  const lon = ref<number | null>(null);
  const hijriDateText = ref<string>("");
  const hijriEndPrediction = ref<string>("");
  const showWeton = ref<boolean>(false);
  const wetonText = ref<string>("");
  const loading = ref<boolean>(false);

  // --- Core Functions --- //

  const fetchLocationAndHijriDate = async () => {
    loading.value = true;
    hijriDateText.value =
      `<p class="text-sm text-gray-500">📅 Menghitung hasil tanggal Hijriyah... 🔍</p>`;
    hijriEndPrediction.value =
      `<p class="text-sm text-gray-500">📅 Menunggu prediksi akhir bulan... ⏳</p>`;

    if ("geolocation" in navigator) {
      navigator.geolocation.getCurrentPosition(
        async ({ coords }) => {
          lat.value = coords.latitude;
          lon.value = coords.longitude;
          await fetchHijriData();
          checkAndSetWeton();
        },
        async () => {
          await fetchLocationByIP();
        }
      );
    } else {
      await fetchLocationByIP();
    }
  };

  const fetchLocationByIP = async () => {
    try {
      const res = await fetch("https://ip-api.com/json");
      const { lat: resLat, lon: resLon } = await res.json();
      if (resLat && resLon) {
        lat.value = resLat;
        lon.value = resLon;
        await fetchHijriData();
        checkAndSetWeton();
      } else {
        showError("⚠️ Gagal mendapatkan lokasi.");
      }
    } catch {
      showError("❌ Gagal mendapatkan lokasi.");
    }
  };

  const checkAndSetWeton = () => {
    if (lat.value && lon.value && isLocationInJava(lat.value, lon.value)) {
      showWeton.value = true;
      let day = new Date().toLocaleDateString("id-ID", { weekday: "long" });
      if (day === "Minggu") day = "Ahad";
      wetonText.value = `${day} ${getWeton(new Date())}`;
    } else {
      showWeton.value = false;
      wetonText.value = "";
    }
  };

  const getHijriMonthName = (month: number): string => {
    const months = [
      "Muharam",
      "Safar",
      "Rabiulawal",
      "Rabiulakhir",
      "Jumadilawal",
      "Jumadilakhir",
      "Rajab",
      "Syakban",
      "Ramadan",
      "Syawal",
      "Zulkaidah",
      "Zulhijah",
    ];
    return months[month - 1] || "Tidak diketahui";
  };

  const formatGregorianDate = ({ day, month, year }: GregorianDate): string => {
    const months = [
      "Januari",
      "Februari",
      "Maret",
      "April",
      "Mei",
      "Juni",
      "Juli",
      "Agustus",
      "September",
      "Oktober",
      "November",
      "Desember",
    ];
    return `${day} ${months[month - 1]} ${year}`;
  };

  const checkImkanurRukyat = ({
    moonAge,
    moonAltitude,
    elongation,
  }: Partial<HijriEndMonth>) => {
    const criteria = [
      { name: "Usia Bulan ≥ 8 jam", value: moonAge ?? NaN, required: 8 },
      { name: "Ketinggian Bulan ≥ 3°", value: moonAltitude ?? NaN, required: 3 },
      { name: "Elongasi ≥ 6,4°", value: elongation ?? NaN, required: 6.4 },
    ];

    const isValid = criteria.every((c) => !isNaN(c.value) && c.value >= c.required);
    const summary = isValid
      ? "✅ Memenuhi syarat Imkanur Rukyat"
      : "❌ Tidak memenuhi syarat → Bulan ini 30 hari";

    const details = criteria
      .map((c) => {
        const satuan = c.name.includes("jam") ? " jam" : "°";
        const status = !isNaN(c.value) && c.value >= c.required ? "✅" : "❌";
        return `${status} ${c.name}: ${
          isNaN(c.value) ? "Tidak tersedia" : c.value.toFixed(2) + satuan
        }`;
      })
      .join("<li>");

    return { summary, details };
  };

  const showError = (message: string) => {
    hijriDateText.value = message;
    hijriEndPrediction.value = "<p class='text-red-500'>Data tidak tersedia.</p>";
    loading.value = false;
  };

  const fetchHijriData = async () => {
    if (!lat.value || !lon.value || !userTimezone.value) {
      showError("❌ Lokasi atau zona waktu belum diatur.");
      return;
    }

    try {
      const params = `lat=${lat.value}&lon=${lon.value}&method=${selectedMethod.value}&timezone=${userTimezone.value}`;
      const [dateRes, endRes] = await Promise.all([
        // fetch(`${API_BASE_URL}/hijri-date?${params}`), // production
        // fetch(`${API_BASE_URL}/hijri-end-month?${params}`), // production
        fetch(`http://127.0.0.1:5000/api/hijri-date?${params}`), // local
        fetch(`http://127.0.0.1:5000/api/hijri-end-month?${params}`), // local
      ]);

      const dateData: { hijriDate?: HijriDate } = await dateRes.json();
      const endMonthData: {
        message?: string;
        estimatedEndOfMonth?: HijriEndMonth;
        estimatedStartOfMonth?: { gregorian: GregorianDate };
      } = await endRes.json();

      loading.value = false;

      if (!dateData?.hijriDate) {
        showError("⚠️ Gagal mendapatkan data.");
        return;
      }

      const hijriToday = dateData.hijriDate;
      const hijriText = `${hijriToday.day} ${getHijriMonthName(
        hijriToday.month
      )} ${hijriToday.year} H`;
      hijriDateText.value = hijriText;

      if (!endMonthData) return;

      const { message, estimatedEndOfMonth, estimatedStartOfMonth } =
        endMonthData;

      const showPrediction = hijriToday.day >= 29;
      if (message) {
        hijriEndPrediction.value = `<p class="text-gray-500">🔍 ${message}</p>`;
        return;
      }
      if (!showPrediction) {
        hijriEndPrediction.value =
          `<p class="text-gray-500">⚠️ Prediksi hilal hanya tersedia saat tanggal 29 Hijriyah.</p>`;
        return;
      }

      const hijri29 = estimatedEndOfMonth?.hijri
        ? `${estimatedEndOfMonth.hijri.day} ${getHijriMonthName(
            estimatedEndOfMonth.hijri.month
          )} ${estimatedEndOfMonth.hijri.year} H`
        : "Tidak tersedia";

      const formattedStartGregorian = estimatedStartOfMonth?.gregorian
        ? formatGregorianDate(estimatedStartOfMonth.gregorian)
        : "Tidak tersedia";

      const { summary, details } = checkImkanurRukyat(estimatedEndOfMonth || {});

      hijriEndPrediction.value = `
        <div>
          <h3 class="text-lg font-semibold mb-2">📅 Informasi Tanggal Hijriyah</h3>
          <p><strong>🗓️ Hari Ini:</strong> ${hijriText}</p>
          <p><strong>🔮 Perkiraan Akhir Bulan:</strong> ${hijri29}</p>
          <p><strong>🌙 Awal Bulan Baru:</strong> ${formattedStartGregorian}</p>
          <br>
          <p><strong>✅ Validasi Imkanur Rukyat:</strong> ${summary}</p>
          <ul><li>${details}</li></ul>
        </div>
      `;
    } catch {
      showError("❌ Terjadi kesalahan saat mengambil data.");
    }
  };

  // --- Watcher --- //
  watch([selectedMethod, userTimezone], () => {
    checkAndSetWeton();
    fetchHijriData();
  });

  return {
    hijriDateText,
    hijriEndPrediction,
    showWeton,
    wetonText,
    loading,
    fetchLocationAndHijriDate,
  };
}