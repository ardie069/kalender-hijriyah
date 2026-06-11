package notification

import (
	"fmt"
)

// FormatPredictionMessage formats the hilal prediction message for a telegram bot
func FormatPredictionMessage(user User, pred Prediction) string {
	locationStr := fmt.Sprintf("%.4f, %.4f", user.Latitude, user.Longitude)
	
	visibleStr := "❌ Hilal kemungkinan TIDAK terlihat"
	if pred.Visible {
		visibleStr = "✅ Hilal kemungkinan TERLIHAT (Imkanur Rukyat)"
	}

	dateStr := pred.Date.Format("02 January 2006")
	sunsetStr := pred.Sunset.Format("15:04 MST")
	moonsetStr := pred.Moonset.Format("15:04 MST")

	msg := fmt.Sprintf(`Assalamu'alaikum %s! 👋

Berikut adalah hasil prediksi visibilitas hilal untuk lokasi Anda. 🌙✨

📍 *Informasi Lokasi:*
• Koordinat: %s
• Zona Waktu: %s

📅 *Waktu Pantau (%s):*
• 🌅 Matahari Terbenam (Sunset): %s
• 🌜 Bulan Terbenam (Moonset): %s

📊 *Data Astronomi (Saat Sunset):*
• 📐 Ketinggian (Altitude): %.2f°
• 🧭 Azimuth: %.2f°
• 📏 Elongasi: %.2f°
• ⏳ Umur Bulan: %.2f jam
• 💡 Iluminasi: %.2f%%

🔮 *Kesimpulan:*
%s

Semoga informasi ini bermanfaat! 🙏
`, 
		user.Username,
		locationStr, user.Timezone,
		dateStr, sunsetStr, moonsetStr,
		pred.Altitude, pred.Azimuth, pred.Elongation, pred.AgeHours, pred.Illumination*100,
		visibleStr,
	)

	return msg
}
