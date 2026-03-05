package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/astronomy"
	"github.com/gin-gonic/gin"
)

func main() {
	manager, err := astronomy.NewEphemerisManager(
		"data/naif0012.tls",
		"data/de440.bsp",
	)
	if err != nil {
		log.Fatalf("Gagal inisialisasi NASA: %v", err)
	}

	adapter := astronomy.GetAdapter(manager)

	fmt.Println("🚀 Sistem Kalender Hijriyah: ONLINE (NASA Engine)")

	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "Alive",
			"engine":  "Go 1.26.0",
			"message": "Siap porting dari Python!",
		})
	})

	r.GET("/api/v1/moon", func(ctx *gin.Context) {
		tel, _ := adapter.GetMoonTelemetry(time.Now(), -7.98, 112.62, 444)
		ctx.JSON(200, tel)
	})

	r.Run(":8080")
}
