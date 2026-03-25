package routes

import (
	"net/http"
	"time"

	"github.com/ardie069/kalender-hijriyah/core/api/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, hijriHandler *handlers.HijriHandler, prayerHandler *handlers.PrayerHandler) {
	// CORS configuration using gin-contrib/cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))	// Safety Check: Avoid nil pointer dereference if NASA Engine failure occurred during init
	if hijriHandler == nil || prayerHandler == nil {
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"message": "NASA Engine Gagal Start / Data Kernel Ga Ketemu. Cek log Vercel lu, Bang!",
			})
		})
		return
	}

	// Root ping
	router.GET("/ping", hijriHandler.Ping)

	// API Version 4 endpoints setup function
	setupV4 := func(g *gin.RouterGroup) {
		hijri := g.Group("/hijri")
		{
			hijri.GET("/date", hijriHandler.GetHijriDate)
			hijri.GET("/search", hijriHandler.SearchDate)
			hijri.GET("/calendar", hijriHandler.GetYearlyCalendar)
		}

		moon := g.Group("/moon")
		{
			moon.GET("/telemetry", hijriHandler.GetTelemetry)
		}

		prayer := g.Group("/prayer")
		{
			prayer.GET("/times", prayerHandler.GetPrayerTimes)
		}
	}

	// Serve on both /api/v4 and /v4
	setupV4(router.Group("/api/v4"))
	setupV4(router.Group("/v4"))
}
