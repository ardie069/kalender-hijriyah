package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/ardie069/kalender-hijriyah/docs"
	"github.com/ardie069/kalender-hijriyah/internal/delivery/http/handlers"
	"github.com/ardie069/kalender-hijriyah/internal/delivery/http/routes"
	"github.com/ardie069/kalender-hijriyah/internal/usecase/calendar"
	"github.com/ardie069/kalender-hijriyah/internal/usecase/hijri"
	"github.com/ardie069/kalender-hijriyah/internal/usecase/prayer"
	"github.com/ardie069/kalender-hijriyah/internal/usecase/timezone"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// @title Kalender Hijriyah API
// @version 4.0
// @description API untuk konversi tanggal Masehi ke Hijriyah, informasi fase bulan, dan jadwal sholat menggunakan perhitungan NASA SPICE.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email ardie069@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v4
func main() {
	// 1. Inisialisasi Engine NASA (Pake NewEphemerisManager biar aman)
	// Fungsi ini bakal nge-loop dan mastiin semua kernel ke-load tanpa skip error
	manager, err := cspice.NewEphemerisManager(
		"pkg/cspice/kernels/de440s.bsp",
		"pkg/cspice/kernels/naif0012.tls",
		"pkg/cspice/kernels/pck00011.tpc",
	)
	if err != nil {
		log.Fatalf("❌ NASA Engine Failure: %v", err)
	}

	// 2. Setup Layers (Gunakan Factory Function biar konsisten)
	tzSvc, err := timezone.NewService()
	if err != nil {
		log.Fatalf("❌ Timezone Service Failure: %v", err)
	}

	adapter := cspice.GetAdapter(manager)
	logic := calendar.NewLogic(adapter, manager)

	dateSvc := hijri.NewDateService(adapter, logic, tzSvc)
	calSvc := hijri.NewCalendarService(dateSvc)
	prayerCalc := prayer.NewCalculator(adapter)

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// 3. Setup Handlers dan Routes
	prayerHandler := handlers.NewPrayerHandler(prayerCalc, dateSvc, tzSvc)
	hijriHandler := handlers.NewHijriHandler(dateSvc, calSvc, adapter)

	routes.SetupRoutes(r, hijriHandler, prayerHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Kalender Hijriyah Engine: STANDBY ON PORT %s\n", port)
	fmt.Printf("🌍 Coordinate System: NASA SPICE Topocentric\n")
	r.Run(":" + port)
}
