package models

type PrayerResponse struct {
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timezone  string  `json:"timezone"`
	} `json:"location"`
	Date struct {
		Gregorian string `json:"gregorian"`
		Hijri     string `json:"hijri"`
	} `json:"date"`
	Times struct {
		Fajr       string `json:"fajr"`
		Sunrise    string `json:"sunrise"`
		Dhuhr      string `json:"dhuhr"`
		Asr        string `json:"asr"`
		Maghrib    string `json:"maghrib"`
		Isha       string `json:"isha"`
		Midnight   string `json:"midnight"`
		ThirdNight string `json:"third_night"`
	} `json:"times"`
}
