package main

import (
	"fmt"
	"log"

	"github.com/ardie069/kalender-hijriyah/core/api/handlers"
	"github.com/ardie069/kalender-hijriyah/core/api/routes"
	"github.com/ardie069/kalender-hijriyah/core/astronomy"
	"github.com/ardie069/kalender-hijriyah/core/calendar"
	"github.com/ardie069/kalender-hijriyah/core/services"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	// 1. Inisialisasi Engine NASA (Pake NewEphemerisManager biar aman)
	// Fungsi ini bakal nge-loop dan mastiin semua kernel ke-load tanpa skip error
	manager, err := astronomy.NewEphemerisManager(
		"data/de440s.bsp",
		"data/naif0012.tls",
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

	r := gin.Default()

	// 3. Setup Handlers dan Routes
	prayerService := services.NewPrayerService(adapter)
	prayerHandler := handlers.NewPrayerHandler(prayerService, service)
	hijriHandler := handlers.NewHijriHandler(service, adapter)

	routes.SetupRoutes(r, hijriHandler, prayerHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Kalender Hijriyah Engine: STANDBY ON PORT %s\n", port)
	fmt.Printf("🌍 Coordinate System: NASA SPICE Topocentric\n")
	r.Run(":" + port)
}
