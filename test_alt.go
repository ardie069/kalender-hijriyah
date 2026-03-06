package main

import (
	"fmt"
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
)

func main() {
	// Let's use Jakarta coordinates as a default if not specified elsewhere (-6.2, 106.8)
	// Or maybe the user meant generally Indonesia. Let's use Malang (-7.98, 112.63) as an example.
	lat := -7.98
	lon := 112.63

	kernels := []string{"data/de440.bsp", "data/naif0012.tls", "data/pck00011.tpc"}
	em, err := astronomy.NewEphemerisManager(kernels...)
	if err != nil {
		fmt.Println("Error loading kernel:", err)
		return
	}
	adapter := astronomy.GetAdapter(em)

	fmt.Println("Testing GetSunset on March 19...")
	searchDate := time.Date(2026, 3, 19, 10, 0, 0, 0, time.UTC)
	sunset, err := adapter.GetSunset(searchDate, lat, lon)
	fmt.Println("Sunset at:", sunset, err)

	sunsetMoon, err := adapter.GetMoonTelemetry(sunset, lat, lon)
	fmt.Printf("Moon at Sunset -> Altitude: %.4f, Azimuth: %.4f, Illumination: %.2f%%\n",
		sunsetMoon.Altitude, sunsetMoon.Azimuth, sunsetMoon.Illumination)

	fmt.Println("Testing Sun and Moon altitude on March 19...")
	for hour := 4; hour <= 14; hour++ { // 4 UTC (11 WIB) to 14 UTC (21 WIB)
		dt := time.Date(2026, 3, 19, hour, 0, 0, 0, time.UTC)
		et, _ := astronomy.Str2et(dt.Format("2006-01-02 15:04:05 UTC"))

		sunPos, _ := em.GetTopocentricPosition(astronomy.Sun, et, lat, lon)
		sunAlt, _ := em.GetLocalAltAz(sunPos, lat, lon)

		moonPos, _ := em.GetTopocentricPosition(astronomy.Moon, et, lat, lon)
		moonAlt, _ := em.GetLocalAltAz(moonPos, lat, lon)

		fmt.Printf("Time: %s, Sun Alt: %8.4f, Moon Alt: %8.4f\n", dt.Format("15:04 UTC"), sunAlt, moonAlt)
	}
}
