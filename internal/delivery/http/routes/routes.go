package routes

import (
	"net/http"
	"time"

	hijriCtrl "github.com/ardie069/kalender-hijriyah/internal/delivery/http/controllers/hijri"
	prayerCtrl "github.com/ardie069/kalender-hijriyah/internal/delivery/http/controllers/prayer"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, hijriHandler interface{}, prayerHandler interface{}, hijriCtrl *hijriCtrl.Controller, prayerCtrl *prayerCtrl.Controller) {
	// CORS configuration using gin-contrib/cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))	// Safety Check: Avoid nil pointer dereference if NASA Engine failure occurred during init
	if hijriCtrl == nil || prayerCtrl == nil {
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "NASA Engine Gagal Start / Data Kernel Ga Ketemu. Cek log Vercel lu, Bang!",
			})
		})
		return
	}

	// Root ping
	router.GET("/ping", hijriCtrl.Ping)

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API Version 4 endpoints setup function
	setupV4 := func(g *gin.RouterGroup) {
		hijri := g.Group("/hijri")
		{
			hijri.GET("/date", hijriCtrl.GetHijriDate)
			hijri.GET("/search", hijriCtrl.SearchDate)
			hijri.GET("/calendar", hijriCtrl.GetYearlyCalendar)
			hijri.GET("/month", hijriCtrl.GetGregorianMonth)
			hijri.GET("/visibility-map", hijriCtrl.GetVisibilityMap)
		}

		moon := g.Group("/moon")
		{
			moon.GET("/telemetry", hijriCtrl.GetTelemetry)
		}

		prayer := g.Group("/prayer")
		{
			prayer.GET("/times", prayerCtrl.GetPrayerTimes)
		}
	}

	// Serve on both /api/v4 and /v4
	setupV4(router.Group("/api/v4"))
	setupV4(router.Group("/v4"))
}
