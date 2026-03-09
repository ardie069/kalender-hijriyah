package routes

import (
	"github.com/ardie069/kalender-hijriyah/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, hijriHandler *handlers.HijriHandler) {
	// Root ping
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
	}

	// Serve on both /api/v4 and /v4
	setupV4(router.Group("/api/v4"))
	setupV4(router.Group("/v4"))
}
