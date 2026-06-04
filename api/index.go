package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/ardie069/kalender-hijriyah/internal/delivery/http/handlers"
	"github.com/ardie069/kalender-hijriyah/internal/delivery/http/routes"
	"github.com/ardie069/kalender-hijriyah/internal/usecase/calendar"
	"github.com/ardie069/kalender-hijriyah/internal/usecase/hijri"
	"github.com/ardie069/kalender-hijriyah/internal/usecase/prayer"
	"github.com/ardie069/kalender-hijriyah/internal/usecase/timezone"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice/kernels"
)

var (
	app       *gin.Engine
	initError error
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	app = gin.New()

	// 1. Extract embedded NASA kernels ke temp directory
	kernelPaths, err := extractKernels()
	if err != nil {
		initError = err
		log.Printf("❌ Kernel Extraction Failure: %v", err)
	}

	// 2. Initialize NASA Engine
	var manager *cspice.EphemerisManager
	if err == nil {
		manager, err = cspice.NewEphemerisManager(kernelPaths...)
		if err != nil {
			initError = err
			log.Printf("❌ NASA Engine Failure: %v", err)
		}
	}

	// 3. Debug endpoint (selalu available)
	app.GET("/debug", func(c *gin.Context) {
		cwd, _ := os.Getwd()
		files, _ := os.ReadDir(cwd)
		var fileList []string
		for _, f := range files {
			fileList = append(fileList, f.Name())
		}

		errMsg := "none"
		if initError != nil {
			errMsg = initError.Error()
		}

		c.JSON(200, gin.H{
			"status":       "Ready",
			"cwd":          cwd,
			"files":        fileList,
			"init_error":   errMsg,
			"kernel_paths": kernelPaths,
		})
	})

	// 4. Setup Handlers & Routes
	if err != nil {
		app.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "NASA Engine Offline",
				"detail":  initError.Error(),
			})
		})
		return
	}

	tzSvc, err := timezone.NewService()
	if err != nil {
		initError = err
		app.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "Timezone Service Offline",
				"detail":  initError.Error(),
			})
		})
		return
	}

	adapter := cspice.GetAdapter(manager)
	logic := calendar.NewLogic(adapter, manager)

	dateSvc := hijri.NewDateService(adapter, logic, tzSvc)
	calSvc := hijri.NewCalendarService(dateSvc)
	prayerCalc := prayer.NewCalculator(adapter)

	hHandler := handlers.NewHijriHandler(dateSvc, calSvc, adapter)
	pHandler := handlers.NewPrayerHandler(prayerCalc, dateSvc, tzSvc)

	routes.SetupRoutes(app, hHandler, pHandler)
}

// extractKernels writes embedded kernel files to a temp directory
// and returns the absolute paths.
func extractKernels() ([]string, error) {
	tmpDir, err := os.MkdirTemp("", "nasa-kernels-*")
	if err != nil {
		return nil, err
	}

	kernelNames := []string{"de440s.bsp", "naif0012.tls", "pck00011.tpc"}
	var paths []string

	for _, name := range kernelNames {
		data, err := kernels.FS.ReadFile(name)
		if err != nil {
			return nil, err
		}

		dest := filepath.Join(tmpDir, name)
		if err := os.WriteFile(dest, data, 0644); err != nil {
			return nil, err
		}
		paths = append(paths, dest)
	}

	return paths, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if app == nil {
		http.Error(w, "Engine not initialized", http.StatusInternalServerError)
		return
	}
	app.ServeHTTP(w, r)
}
