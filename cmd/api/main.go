package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/ardie069/kalender-hijriyah/docs"
	"github.com/ardie069/kalender-hijriyah/pkg/app"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
	"github.com/gin-contrib/gzip"
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
	// 1. Inisialisasi Engine NASA
	manager, err := cspice.NewEphemerisManager(
		"pkg/cspice/kernels/de440s.bsp",
		"pkg/cspice/kernels/naif0012.tls",
		"pkg/cspice/kernels/pck00011.tpc",
	)
	if err != nil {
		log.Fatalf("❌ NASA Engine Failure: %v", err)
	}

	// 2. Initialize AppConfig with factory function
	appConfig, err := app.NewAppConfig(manager)
	if err != nil {
		log.Fatalf("❌ Application Initialization Failure: %v", err)
	}

	// 3. Add compression middleware to the engine
	appConfig.GetEngine().Use(gzip.Gzip(gzip.DefaultCompression))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Kalender Hijriyah Engine: STANDBY ON PORT %s\n", port)
	fmt.Printf("🌍 Coordinate System: NASA SPICE Topocentric\n")
	appConfig.GetEngine().Run(":" + port)
}
