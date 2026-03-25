package routes

import (
	"net/http"

	"github.com/ardie069/kalender-hijriyah/core/api/handlers"
	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func SetupRoutes(router *gin.Engine, hijriHandler *handlers.HijriHandler, prayerHandler *handlers.PrayerHandler) {
	// Safety Check: Avoid nil pointer dereference if NASA Engine failure occurred during init
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
	router.Use(corsMiddleware())
	router.GET("/ping", hijriHandler.Ping)

	// API Version 4 endpoints setup function
	setupV4 := func(g *gin.RouterGroup) {
		hijri := g.Group("/hijri")
		{
			hijri.GET("/date", hijriHandler.GetHijriDate)
			hijri.GET("/search", hijriHandler.SearchDate)
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
