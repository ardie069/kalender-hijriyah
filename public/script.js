document.addEventListener("DOMContentLoaded", () => {
    const loadingText = document.getElementById("loading");
    const hijriDateText = document.getElementById("hijri-date");
    const currentTimeText = document.getElementById("current-time");
    const methodSelect = document.getElementById("method");
    const methodLabel = document.getElementById("method-label");
    const body = document.body;
    const toggleThemeButton = document.getElementById("toggle-theme");
    const box = document.getElementById("box");
    const dateBox = document.getElementById("date-box");
    const timezoneText = document.getElementById("timezone");

    function updateRealTime() {
        setInterval(() => {
            const now = new Date();
            const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
            const offset = now.getTimezoneOffset() / -60;

            currentTimeText.textContent = `🕒 ${now.toLocaleString("id-ID", {
                weekday: "long", year: "numeric", month: "long", day: "numeric",
                hour: "2-digit", minute: "2-digit", second: "2-digit"
            })}`;

            timezoneText.textContent = `🌍 Zona Waktu: ${timezone} (UTC${offset >= 0 ? "+" : ""}${offset})`;
        }, 1000);
    }

    async function fetchLocationAndHijriDate() {
        loadingText.style.display = "block";
        hijriDateText.textContent = "";

        if ("geolocation" in navigator) {
            navigator.geolocation.getCurrentPosition(
                async (position) => {
                    await fetchHijriDate(position.coords.latitude, position.coords.longitude);
                },
                async () => {
                    console.warn("⚠️ Geolocation ditolak. Mencoba lokasi berdasarkan IP...");
                    await fetchLocationByIP();
                }
            );
        } else {
            console.warn("⚠️ Geolocation tidak didukung. Menggunakan lokasi berdasarkan IP...");
            await fetchLocationByIP();
        }
    }

    async function fetchLocationByIP() {
        try {
            const response = await fetch("https://ip-api.com/json");
            if (!response.ok) throw new Error("Gagal mengambil lokasi IP.");

            const data = await response.json();
            if (data.lat && data.lon) {
                console.log(`📍 Lokasi berdasarkan IP: ${data.city}, ${data.country}`);
                await fetchHijriDate(data.lat, data.lon);
            } else {
                hijriDateText.textContent = "⚠️ Gagal mendapatkan lokasi.";
                loadingText.style.display = "none";
            }
        } catch (error) {
            console.error("❌ Error Fetching Location by IP:", error);
            hijriDateText.textContent = "❌ Gagal mendapatkan lokasi.";
            loadingText.style.display = "none";
        }
    }

    async function fetchHijriDate(lat, lon) {
        const method = methodSelect.value;
        const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;

        const API_URL = window.location.hostname.includes("localhost")
            ? "/hijri-date"
            : "/api/hijri-date";

        try {
            console.log(`📡 Fetching from ${API_URL} with location ${lat}, ${lon}, TZ: ${timezone}`);
            const response = await fetch(`${API_URL}?lat=${lat}&lon=${lon}&method=${method}&timezone=${timezone}`);

            if (!response.ok) throw new Error(`HTTP Error: ${response.status}`);

            const data = await response.json();
            console.log("✅ Data diterima:", data);

            loadingText.style.display = "none";

            if (data.hijriDate) {
                hijriDateText.textContent = `${data.hijriDate.day} ${getHijriMonthName(data.hijriDate.month)} ${data.hijriDate.year} H`;
            } else {
                hijriDateText.textContent = "⚠️ Gagal mendapatkan data.";
            }
        } catch (error) {
            console.error("❌ Error Fetching Data:", error);
            hijriDateText.textContent = "❌ Terjadi kesalahan saat mengambil data.";
        }
    }

    function getHijriMonthName(month) {
        const hijriMonths = [
            "Muharram", "Safar", "Rabi'ul Awal", "Rabi'ul Akhir",
            "Jumadil Awal", "Jumadil Akhir", "Rajab", "Sya'ban",
            "Ramadan", "Syawal", "Dzulqa'dah", "Dzulhijjah"
        ];
        return hijriMonths[month - 1] || "Tidak diketahui";
    }

    function toggleTheme() {
        body.classList.toggle("bg-gray-900");
        body.classList.toggle("bg-white");
        body.classList.toggle("text-white");
        body.classList.toggle("text-black");

        box.classList.toggle("bg-gray-800");
        box.classList.toggle("bg-gray-200");

        dateBox.classList.toggle("bg-gray-700");
        dateBox.classList.toggle("bg-gray-300");

        currentTimeText.classList.toggle("text-white");
        currentTimeText.classList.toggle("text-black");

        timezoneText.classList.toggle("text-gray-400");
        timezoneText.classList.toggle("text-black");

        methodSelect.classList.toggle("bg-gray-300");
        methodSelect.classList.toggle("bg-gray-700");

        methodSelect.classList.toggle("text-black");
        methodSelect.classList.toggle("text-white");

        methodLabel.classList.toggle("text-gray-300");
        methodLabel.classList.toggle("text-black");

        const isDarkMode = body.classList.contains("bg-gray-900");
        toggleThemeButton.textContent = isDarkMode ? "🌞 Light Mode" : "🌙 Dark Mode";
        localStorage.setItem("theme", isDarkMode ? "dark" : "light");
    }

    function loadTheme() {
        const savedTheme = localStorage.getItem("theme");
        if (savedTheme === "dark") {
            toggleTheme();
        }
    }

    updateRealTime();
    fetchLocationAndHijriDate();
    methodSelect.addEventListener("change", fetchLocationAndHijriDate);
    toggleThemeButton.addEventListener("click", toggleTheme);
    loadTheme();
});
