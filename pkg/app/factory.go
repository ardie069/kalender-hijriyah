package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ardie069/kalender-hijriyah/internal/delivery/http/handlers"
	"github.com/ardie069/kalender-hijriyah/internal/delivery/http/routes"
	"github.com/ardie069/kalender-hijriyah/internal/hijri"
	"github.com/ardie069/kalender-hijriyah/internal/prayer"
	"github.com/ardie069/kalender-hijriyah/internal/usecase/timezone"
	"github.com/ardie069/kalender-hijriyah/internal/astronomy/ephemeris"
	"github.com/ardie069/kalender-hijriyah/internal/visibility"
	"github.com/ardie069/kalender-hijriyah/internal/visibility/scan"
	"github.com/ardie069/kalender-hijriyah/internal/calendar/ummalqura"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
)

// AppConfig holds the application configuration and dependencies
type AppConfig struct {
	Manager    *cspice.EphemerisManager
	TzService  *timezone.Service
	DateSvc    *hijri.DateService
	CalSvc     *hijri.CalendarService
	PrayerSvc  *prayer.Service
	VisSvc     *visibility.Service
	Engine     *gin.Engine
}

// NewAppConfig creates and initializes a new AppConfig with all required services
func NewAppConfig(manager *cspice.EphemerisManager) (*AppConfig, error) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Initialize Timezone Service
	tzSvc, err := timezone.NewService()
	if err != nil {
		return nil, err
	}

	// Initialize Calendar Logic
	adapter := cspice.GetAdapter(manager)
	ephemSvc := ephemeris.NewService(adapter, manager)
	scanSvc := &scan.Scanner{Astro: adapter}
	visSvc := visibility.NewService(adapter, ephemSvc)
	ummAlQuraSvc := ummalqura.NewService(adapter, ephemSvc)

	// Initialize Services
	dateSvc := hijri.NewDateService(adapter, ephemSvc, scanSvc, ummAlQuraSvc, tzSvc)
	calSvc := hijri.NewCalendarService(dateSvc)
	prayerSvc := prayer.NewService(adapter)

	// Initialize Handlers
	hHandler := handlers.NewHijriHandler(dateSvc, calSvc, visSvc, adapter)
	pHandler := handlers.NewPrayerHandler(prayerSvc, dateSvc, tzSvc)

	// Setup Routes
	routes.SetupRoutes(engine, hHandler, pHandler)

	return &AppConfig{
		Manager:    manager,
		TzService:  tzSvc,
		DateSvc:    dateSvc,
		CalSvc:     calSvc,
		PrayerSvc:  prayerSvc,
		VisSvc:     visSvc,
		Engine:     engine,
	}, nil
}

// GetEngine returns the Gin engine
func (ac *AppConfig) GetEngine() *gin.Engine {
	return ac.Engine
}

// ServeHTTP implements http.Handler interface for Vercel
func (ac *AppConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ac.Engine.ServeHTTP(w, r)
}
