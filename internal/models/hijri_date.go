package models

import "time"

type HijriDate struct {
	Day       int    `json:"day"`
	Month     int    `json:"month"`
	MonthName string `json:"month_name"`
	Year      int    `json:"year"`
	IsTabular bool   `json:"is_tabular"`
}

type LocationInfo struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}

type MethodResult struct {
	HijriDate         HijriDate        `json:"hijri_date"`
	CurrentAltitude   *float64         `json:"current_altitude,omitempty"`
	CurrentElongation *float64         `json:"current_elongation,omitempty"`
	Prediction        *HilalPrediction `json:"prediction,omitempty"`
}

type HilalPrediction struct {
	CheckDateUTC time.Time `json:"check_date_utc"`
	IsNewMonth   bool      `json:"is_new_month"`
	IjtimaTime   time.Time `json:"ijtima_time"`
	Altitude     float64   `json:"altitude"`
	Elongation   float64   `json:"elongation"`
	AgeHours     float64   `json:"age_hours"`
}

type HijriResponse struct {
	GregorianDate time.Time               `json:"gregorian_date"`
	Location      LocationInfo            `json:"location"`
	Methods       map[string]MethodResult `json:"methods"`
}
