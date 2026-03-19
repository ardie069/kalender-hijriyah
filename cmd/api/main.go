package main

import (
	"fmt"
	"log"

	"github.com/ardie069/kalender-hijriyah/internal/api/handlers"
	"github.com/ardie069/kalender-hijriyah/internal/api/routes"
	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
	"github.com/ardie069/kalender-hijriyah/internal/calendar"
	"github.com/ardie069/kalender-hijriyah/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Engine NASA (Pake NewEphemerisManager biar aman)
	// Fungsi ini bakal nge-loop dan mastiin semua kernel ke-load tanpa skip error
	manager, err := astronomy.NewEphemerisManager(
		"data/naif0012.tls",
		"data/de440.bsp",
		"data/pck00011.tpc",
	)
	if err != nil {
		log.Fatalf("❌ NASA Engine Failure: %v", err)
	}

	// 2. Setup Layers (Gunakan Factory Function biar konsisten)
	adapter := astronomy.GetAdapter(manager)
	logic := calendar.NewLogic(adapter, manager)
	service, err := services.NewHijriService(adapter, logic)
	if err != nil {
		log.Fatalf("❌ Hijri Service Failure: %v", err)
	}

	fmt.Println("🚀 Kalender Hijriyah Engine: STANDBY ON PORT 8080")
	fmt.Println("🌍 Coordinate System: NASA SPICE Topocentric")

	r := gin.Default()

	// 3. Setup Handlers dan Routes
	prayerService := services.NewPrayerService(adapter)
	prayerHandler := handlers.NewPrayerHandler(prayerService, service)
	hijriHandler := handlers.NewHijriHandler(service, adapter)

	routes.SetupRoutes(r, hijriHandler, prayerHandler)

	r.Run(":8080")
}
