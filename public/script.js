document.addEventListener("DOMContentLoaded", () => {
    const elements = {
        loadingText: document.getElementById("loading"),
        hijriDateText: document.getElementById("hijri-date"),
        currentTimeText: document.getElementById("current-time"),
        methodSelect: document.getElementById("method"),
        methodLabel: document.getElementById("method-label"),
        body: document.body,
        toggleThemeButton: document.getElementById("toggle-theme"),
        box: document.getElementById("box"),
        dateBox: document.getElementById("date-box"),
        timezoneText: document.getElementById("timezone"),
        hijriEndPrediction: document.getElementById("hijri-end-prediction")
    };

    let lat = 0, lon = 0;
    let userTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone;

    const API_BASE_URL = window.location.hostname.includes("localhost")
        ? "http://localhost:3000"
        : "/api";

    function updateRealTime() {
        setInterval(() => {
            const now = new Date();
            const offset = now.getTimezoneOffset() / -60;

            elements.currentTimeText.textContent = `🕒 ${now.toLocaleString("id-ID", {
                weekday: "long", year: "numeric", month: "long", day: "numeric",
                hour: "2-digit", minute: "2-digit", second: "2-digit"
            })}`;

            elements.timezoneText.textContent = `🌍 Zona Waktu: ${userTimezone} (UTC${offset >= 0 ? "+" : ""}${offset})`;
        }, 1000);
    }

    async function fetchLocationAndHijriDate() {
        elements.loadingText.style.display = "block";
        elements.hijriDateText.textContent = "📅 Menghitung hasil tanggal Hijriyah... 🔍";
        elements.hijriEndPrediction.innerHTML = "📅 Menunggu prediksi akhir bulan... ⏳";

        if ("geolocation" in navigator) {
            navigator.geolocation.getCurrentPosition(
                async (position) => {
                    lat = position.coords.latitude;
                    lon = position.coords.longitude;
                    await fetchHijriData();
                },
                async () => {
                    await fetchLocationByIP();
                }
            );
        } else {
            await fetchLocationByIP();
        }
    }

    async function fetchLocationByIP() {
        try {
            const response = await fetch("https://ip-api.com/json");
            if (!response.ok) throw new Error("Gagal mengambil lokasi IP.");
            const data = await response.json();
            if (data.lat && data.lon) {
                lat = data.lat;
                lon = data.lon;
                await fetchHijriData();
            } else {
                elements.hijriDateText.textContent = "⚠️ Gagal mendapatkan lokasi.";
                elements.loadingText.style.display = "none";
            }
        } catch (error) {
            elements.hijriDateText.textContent = "❌ Gagal mendapatkan lokasi.";
            elements.loadingText.style.display = "none";
        }
    }

    // Fungsi untuk mengubah format tanggal Gregorian menjadi "dd mm (nama bulan) yyyy"
    function formatGregorianDate(dateString) {
        const dateObj = new Date(dateString);
        const day = dateObj.getDate().toString().padStart(2, '0'); // Pastikan dua digit
        const monthNames = [
            "Januari", "Februari", "Maret", "April", "Mei", "Juni",
            "Juli", "Agustus", "September", "Oktober", "November", "Desember"
        ];
        const month = monthNames[dateObj.getMonth()];
        const year = dateObj.getFullYear();
        return `${day} ${month} ${year}`;
    }

    async function fetchHijriData() {
        const selectedMethod = elements.methodSelect.value;
        const API_URL = `${API_BASE_URL}/hijri-date`;
        const END_MONTH_URL = `${API_BASE_URL}/hijri-end-month`;
    
        if (!lat || !lon || !userTimezone) {
            elements.hijriDateText.textContent = "❌ Lokasi atau zona waktu belum diatur.";
            elements.hijriEndPrediction.innerHTML = "<p class='text-red-500'>Data tidak tersedia.</p>";
            return;
        }
    
        try {
            const [dateResponse, endMonthResponse] = await Promise.all([
                fetch(`${API_URL}?lat=${lat}&lon=${lon}&method=${selectedMethod}&timezone=${userTimezone}`),
                fetch(`${END_MONTH_URL}?lat=${lat}&lon=${lon}&method=${selectedMethod}&timezone=${userTimezone}`)
            ]);
    
            if (!dateResponse.ok || !endMonthResponse.ok) throw new Error("HTTP Error saat mengambil data.");
    
            const dateData = await dateResponse.json();
            const endMonthData = await endMonthResponse.json();
    
            elements.loadingText.style.display = "none";
    
            // Ambil tanggal Hijriyah hari ini
            let hijriToday;
            if (dateData?.hijriDate) {
                hijriToday = dateData.hijriDate;
                elements.hijriDateText.textContent =
                    `${hijriToday.day} ${getHijriMonthName(hijriToday.month)} ${hijriToday.year} H`;
            } else {
                elements.hijriDateText.textContent = "⚠️ Gagal mendapatkan data.";
                return;
            }
    
            // Jika data akhir bulan tersedia
            if (endMonthData?.estimatedEndOfMonth) {
                const {
                    hijri: estimatedHijri,
                    moonAltitude,
                    elongation,
                    sunAltitude,
                    moonAge,
                    conjunctionBeforeSunset,
                    explanation,
                } = endMonthData.estimatedEndOfMonth;
    
                // Pastikan tanggal 29 diambil dari perkiraan akhir bulan
                const hijri29 = estimatedHijri
                    ? `${estimatedHijri.day} ${getHijriMonthName(estimatedHijri.month)} ${estimatedHijri.year} H`
                    : "Tidak tersedia";
    
                // Format angka dengan 2 desimal
                const formattedMoonAltitude = moonAltitude !== "Tidak tersedia" ? parseFloat(moonAltitude).toFixed(2) : "Tidak tersedia";
                const formattedElongation = elongation !== "Tidak tersedia" ? parseFloat(elongation).toFixed(2) : "Tidak tersedia";
                const formattedSunAltitude = sunAltitude !== "Tidak tersedia" ? parseFloat(sunAltitude).toFixed(2) : "Tidak tersedia";
    
                const parsedMoonAge = parseFloat(moonAge);
                const formattedMoonAge = !isNaN(parsedMoonAge) && moonAge !== null ? `${parsedMoonAge.toFixed(2)} jam` : "Tidak diketahui";
    
                // Validasi Imkanur Rukyat
                let imkanurRukyat = "❌ Tidak memenuhi syarat → Bulan ini 30 hari";
                if (
                    moonAltitude !== "Tidak tersedia" &&
                    elongation !== "Tidak tersedia" &&
                    !isNaN(parsedMoonAge) && parsedMoonAge >= 8 &&
                    parseFloat(moonAltitude) >= 3 && parseFloat(elongation) >= 6.4
                ) {
                    imkanurRukyat = `✅ Memenuhi syarat Imkanur Rukyat → ${explanation}`;
                }
    
                // Format tanggal Gregorian awal bulan Hijriyah jika tersedia
                const formattedStartGregorian = endMonthData.estimatedStartOfMonth?.gregorian
                    ? formatGregorianDate(endMonthData.estimatedStartOfMonth.gregorian)
                    : "Tidak tersedia";
    
                // Hanya tampilkan prediksi jika hari ini **tanggal 29 Hijriyah atau setelahnya**
                const now = new Date();
                const userHour = now.getHours();

                if (hijriToday.day === 29 && userHour >= 6 || hijriToday.day === 30 && userHour < 18) {
                    elements.hijriEndPrediction.innerHTML = `
                        <div class="bg-gray-100 p-2 rounded-lg shadow-md">
                            <h3 class="text-lg font-semibold">📅 Informasi Tanggal Hijriyah</h3>
                            <div class="mt-4">
                                <p class="text-gray-700"><strong>🗓️ Hari Ini:</strong> <span class="text-blue-600">${hijriToday.day} ${getHijriMonthName(hijriToday.month)} ${hijriToday.year} H</span></p>
                                <p class="text-gray-700"><strong>🔮 Perkiraan Akhir Bulan:</strong> <span class="text-green-600">${hijri29}</span></p>
                                <p class="text-gray-700"><strong>🌙 Prediksi Awal Bulan Hijriyah:</strong> <span class="text-purple-600">${formattedStartGregorian}</span></p>
                            </div>
                        </div>
                        <br>
                        <p><strong>📆 Perkiraan Durasi Bulan:</strong> ${explanation}</p>
                        <br>
                        <p><strong>🌙 Posisi Bulan:</strong></p>
                        <ul>
                            <li>- Ketinggian: ${formattedMoonAltitude}°</li>
                            <li>- Elongasi: ${formattedElongation}°</li>
                        </ul>
                        <br>
                        <p><strong>☀️ Posisi Matahari:</strong></p>
                        <ul>
                            <li>- Ketinggian: ${formattedSunAltitude}°</li>
                        </ul>
                        <br>
                        <p><strong>🔭 Konjungsi Terjadi Sebelum Matahari Terbenam:</strong> ${conjunctionBeforeSunset ? "Ya" : "Tidak"}</p>
                        <p><strong>⏳ Usia Bulan:</strong> ${formattedMoonAge}</p>
                        <br>
                        <p class="mb-4"><strong>✅ Validasi Imkanur Rukyat:</strong> ${imkanurRukyat}</p>
                    `;
                } else {
                    // Jika belum tanggal 29, jangan tampilkan prediksi
                    elements.hijriEndPrediction.innerHTML = `<p class='text-gray-500'>🔍 Prediksi hilal hanya tersedia saat tanggal 29 Hijriyah.</p>`;
                }
            } else {
                elements.hijriEndPrediction.innerHTML = "<p class='text-red-500'>❌ Gagal mendapatkan data.</p>";
            }
        } catch (error) {
            elements.hijriDateText.textContent = "❌ Terjadi kesalahan saat mengambil data.";
            elements.hijriEndPrediction.innerHTML = "<p class='text-red-500'>❌ Gagal mengambil data.</p>";
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
        elements.body.classList.toggle("bg-gray-900");
        elements.body.classList.toggle("bg-white");
        elements.body.classList.toggle("text-white");
        elements.body.classList.toggle("text-black");

        elements.loadingText.classList.toggle("text-yellow-400");
        elements.loadingText.classList.toggle("text-yellow-800");

        elements.box.classList.toggle("bg-gray-800");
        elements.box.classList.toggle("bg-gray-200");

        elements.dateBox.classList.toggle("bg-gray-700");
        elements.dateBox.classList.toggle("bg-gray-300");

        elements.hijriEndPrediction.classList.toggle("bg-gray-700");
        elements.hijriEndPrediction.classList.toggle("bg-gray-300");

        elements.currentTimeText.classList.toggle("text-white");
        elements.currentTimeText.classList.toggle("text-black");

        elements.timezoneText.classList.toggle("text-gray-400");
        elements.timezoneText.classList.toggle("text-black");

        elements.methodSelect.classList.toggle("bg-gray-300");
        elements.methodSelect.classList.toggle("bg-gray-700");

        elements.methodSelect.classList.toggle("text-black");
        elements.methodSelect.classList.toggle("text-white");

        elements.methodLabel.classList.toggle("text-gray-300");
        elements.methodLabel.classList.toggle("text-black");

        const isDarkMode = elements.body.classList.contains("bg-gray-900");
        elements.toggleThemeButton.textContent = isDarkMode ? "🌞 Light Mode" : "🌙 Dark Mode";
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
    elements.methodSelect.addEventListener("change", fetchLocationAndHijriDate);
    elements.toggleThemeButton.addEventListener("click", toggleTheme);
    loadTheme();
});
