document.addEventListener("DOMContentLoaded", () => {
    const loadingText = document.getElementById("loading");
    const hijriDateText = document.getElementById("hijri-date");
    const gregorianDateText = document.getElementById("gregorian-date");
    const currentTimeText = document.getElementById("current-time");
    const methodSelect = document.getElementById("method");
    const body = document.body;
    const toggleThemeButton = document.getElementById("toggle-theme");
    const box = document.getElementById("box");
    const dateBox = document.getElementById("date-box");

    function updateRealTime() {
        setInterval(() => {
            const now = new Date();
            currentTimeText.textContent = `🕒 ${now.toLocaleString("id-ID", {
                weekday: "long", year: "numeric", month: "long", day: "numeric",
                hour: "2-digit", minute: "2-digit", second: "2-digit"
            })}`;
        }, 1000);
    }

    async function fetchHijriDate() {
        if (!("geolocation" in navigator)) {
            hijriDateText.textContent = "⚠️ Geolocation tidak didukung.";
            return;
        }

        navigator.geolocation.getCurrentPosition(async (position) => {
            const { latitude, longitude } = position.coords;
            const method = methodSelect.value;

            try {
                const response = await fetch(`/hijri-date?lat=${latitude}&lon=${longitude}&method=${method}`);
                const response2 = await fetch(`/hijri-date?lat=${latitude}&lon=${longitude}&method=${method}`);
                const data = await response.json();
                const data2 = await response2.json();

                loadingText.style.display = "none";

                if (data.hijriDate) {
                    const today = new Date();
                    gregorianDateText.textContent = today.toLocaleDateString("id-ID", {
                        weekday: "long", year: "numeric", month: "long", day: "numeric"
                    });

                    hijriDateText.textContent = `${data.hijriDate.day} ${getHijriMonthName(data.hijriDate.month)} ${data.hijriDate.year} H`;
                } else {
                    hijriDateText.textContent = "⚠️ Gagal mendapatkan data.";
                }
                if (data2.hijriDate) {
                    const today = new Date();
                    gregorianDateText.textContent = today.toLocaleDateString("id-ID", {
                        weekday: "long", year: "numeric", month: "long", day: "numeric"
                    });

                    hijriDateText.textContent = `${data.hijriDate.day} ${getHijriMonthName(data.hijriDate.month)} ${data.hijriDate.year} H`;
                } else {
                    hijriDateText.textContent = "⚠️ Gagal mendapatkan data.";
                }
            } catch (error) {
                hijriDateText.textContent = "❌ Terjadi kesalahan.";
            }
        }, () => {
            hijriDateText.textContent = "⚠️ Izin lokasi ditolak.";
        });
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
        const isDarkMode = body.classList.contains("bg-gray-900");

        if (isDarkMode) {
            // Ganti ke Light Mode
            body.classList.replace("bg-gray-900", "bg-white");
            body.classList.replace("text-white", "text-black");
            box.classList.replace("bg-gray-800", "bg-gray-200");
            dateBox.classList.replace("bg-gray-700", "bg-gray-300");
            methodSelect.classList.replace("bg-gray-700", "bg-gray-300");
            methodSelect.classList.replace("border-gray-600", "border-gray-400");
            methodSelect.classList.replace("text-white", "text-black");
            toggleThemeButton.textContent = "🌙 Dark Mode";
            toggleThemeButton.classList.replace("bg-gray-700", "bg-gray-300");
            toggleThemeButton.classList.replace("text-white", "text-black");
            currentTimeText.classList.replace("text-white", "text-black");
            localStorage.setItem("theme", "light");
        } else {
            // Ganti ke Dark Mode
            body.classList.replace("bg-white", "bg-gray-900");
            body.classList.replace("text-black", "text-white");
            box.classList.replace("bg-gray-200", "bg-gray-800");
            dateBox.classList.replace("bg-gray-300", "bg-gray-700");
            methodSelect.classList.replace("bg-gray-300", "bg-gray-700");
            methodSelect.classList.replace("border-gray-400", "border-gray-600");
            methodSelect.classList.replace("text-black", "text-white");
            toggleThemeButton.textContent = "🌞 Light Mode";
            toggleThemeButton.classList.replace("bg-gray-300", "bg-gray-700");
            toggleThemeButton.classList.replace("text-black", "text-white");
            currentTimeText.classList.replace("text-black", "text-white");
            localStorage.setItem("theme", "dark");
        }
    }

    function loadTheme() {
        const savedTheme = localStorage.getItem("theme");
        if (savedTheme === "light") {
            body.classList.replace("bg-gray-900", "bg-white");
            body.classList.replace("text-white", "text-black");
            box.classList.replace("bg-gray-800", "bg-gray-200");
            dateBox.classList.replace("bg-gray-700", "bg-gray-300");
            methodSelect.classList.replace("bg-gray-700", "bg-gray-300");
            methodSelect.classList.replace("border-gray-600", "border-gray-400");
            methodSelect.classList.replace("text-white", "text-black");
            toggleThemeButton.textContent = "🌙 Dark Mode";
            toggleThemeButton.classList.replace("bg-gray-700", "bg-gray-300");
            toggleThemeButton.classList.replace("text-white", "text-black");
            currentTimeText.classList.replace("text-white", "text-black");
        }
    }

    updateRealTime();
    fetchHijriDate();
    methodSelect.addEventListener("change", fetchHijriDate);
    toggleThemeButton.addEventListener("click", toggleTheme);
    loadTheme();
});
