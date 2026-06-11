package notification

import "time"

type User struct {
	ID         int64
	TelegramID int64

	Username   string
	Latitude   float64
	Longitude  float64
	Timezone   string	

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Prediction struct {
	Date time.Time

	Sunset time.Time
	Moonset time.Time

	Altitude float64
	Azimuth float64

	Elongation float64
	AgeHours float64

	Illumination float64

	Visible bool
}
