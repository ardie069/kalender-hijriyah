import { ref, watch } from "vue";
import { isLocationInJava, getWeton } from "../utils/weton";

export function useHijriDate(selectedMethod, userTimezone, API_BASE_URL) {
  const lat = ref(null);
  const lon = ref(null);
  const hijriDateText = ref("");
  const hijriEndPrediction = ref("");
  const showWeton = ref(false);
  const wetonText = ref("");
  const loading = ref(false);

  const fetchLocationAndHijriDate = async () => {
    loading.value = true;
    hijriDateText.value = `<p class="text-sm text-gray-500">📅 Menghitung hasil tanggal Hijriyah... 🔍</p>`;
    hijriEndPrediction.value = `<p class="text-sm text-gray-500">📅 Menunggu prediksi akhir bulan... ⏳</p>`;

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
    if (isLocationInJava(lat.value, lon.value)) {
      showWeton.value = true;
      // Ambil nama hari dan ganti "Minggu" menjadi "Ahad"
      let day = new Date().toLocaleDateString("id-ID", { weekday: "long" });
      if (day === "Minggu") {
        day = "Ahad";
      }

      // Gabungkan nama hari dengan weton
      wetonText.value = `${day} ${getWeton(new Date())}`;
    } else {
      showWeton.value = false;
      wetonText.value = "";
    }
  };

  const getHijriMonthName = (month) => {
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

  const formatGregorianDate = ({ day, month, year }) => {
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

  const checkImkanurRukyat = ({ moonAge, moonAltitude, elongation }) => {
    const criteria = [
      { name: "Usia Bulan ≥ 8 jam", value: parseFloat(moonAge), required: 8 },
      {
        name: "Ketinggian Bulan ≥ 3°",
        value: parseFloat(moonAltitude),
        required: 3,
      },
      { name: "Elongasi ≥ 6,4°", value: parseFloat(elongation), required: 6.4 },
    ];

    const isValid = criteria.every(
      (c) => !isNaN(c.value) && c.value >= c.required
    );
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

  const showError = (message) => {
    hijriDateText.value = message;
    hijriEndPrediction.value =
      "<p class='text-red-500'>Data tidak tersedia.</p>";
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
        fetch(`http://localhost:5000/api/hijri-date?${params}`), // development
        fetch(`http://localhost:5000/api/hijri-end-month?${params}`), // development
      ]);

      const dateData = await dateRes.json();
      const endMonthData = await endRes.json();

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
        hijriEndPrediction.value = `<p class="text-gray-500">⚠️ Prediksi hilal hanya tersedia saat tanggal 29 Hijriyah.</p>`;
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

      const { summary, details } = checkImkanurRukyat(estimatedEndOfMonth);

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

  // Watch changes
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
