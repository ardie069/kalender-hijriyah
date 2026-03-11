package models

import "time"

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
	Method struct {
		Name    string `json:"name"`
		Madhab  string `json:"madhab"`
		HighLat string `json:"high_lat,omitempty"`
	} `json:"method"`
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

type PrayerTimes struct {
	Fajr       time.Time `json:"fajr"`
	Sunrise    time.Time `json:"sunrise"`
	Dhuhr      time.Time `json:"dhuhr"`
	Asr        time.Time `json:"asr"`
	Maghrib    time.Time `json:"maghrib"`
	Isha       time.Time `json:"isha"`
	Midnight   time.Time `json:"midnight"`
	ThirdNight time.Time `json:"third_night"`
}
