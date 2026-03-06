package routes

import (
	"github.com/ardie069/kalender-hijriyah/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, hijriHandler *handlers.HijriHandler) {
	// Root ping
	router.GET("/ping", hijriHandler.Ping)

	// API Version 4 Group
	v4 := router.Group("/api/v4")
	{
		// Hijri Calendar endpoints
		hijri := v4.Group("/hijri")
		{
			hijri.GET("/date", hijriHandler.GetHijriDate)
		}

		// Moon Telemetry endpoints
		moon := v4.Group("/moon")
		{
			moon.GET("/telemetry", hijriHandler.GetTelemetry)
		}
	}
}
