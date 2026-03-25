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
	CheckDateUTC     time.Time     `json:"check_date_utc"`
	CheckDateLocal   string        `json:"check_date_local"`
	TimezoneName     string        `json:"timezone_name"`
	IsNewMonth       bool          `json:"is_new_month"`
	IjtimaTime       time.Time     `json:"ijtima_time"`
	Altitude         float64       `json:"altitude"`
	AltitudeApparent *float64      `json:"altitude_apparent,omitempty"`
	Elongation       float64       `json:"elongation"`
	AgeHours         float64       `json:"age_hours"`
	Location         *LocationInfo `json:"location,omitempty"`
	GlobalLocation   *LocationInfo `json:"global_location,omitempty"`
}

type HijriResponse struct {
	GregorianDate time.Time               `json:"gregorian_date"`
	Location      LocationInfo            `json:"location"`
	Methods       map[string]MethodResult `json:"methods"`
}

type CalendarMonth struct {
	MonthID        int    `json:"month_id"`
	MonthName      string `json:"month_name"`
	TotalDays      int    `json:"total_days"`
	Day1Weekday    int    `json:"day_1_weekday"`
	StartGregorian string `json:"start_gregorian"`
}

type YearlyCalendarResponse struct {
	Year   int             `json:"year"`
	Method string          `json:"method"`
	Months []CalendarMonth `json:"months"`
}
