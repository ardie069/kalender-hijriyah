package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/ardie069/kalender-hijriyah/pkg/app"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice/kernels"
)

var (
	appConfig *app.AppConfig
	initError error
)

func init() {
	gin.SetMode(gin.ReleaseMode)

	// 1. Extract embedded NASA kernels ke temp directory
	kernelPaths, err := extractKernels()
	if err != nil {
		initError = err
		log.Printf("❌ Kernel Extraction Failure: %v", err)
		setupErrorHandler(kernelPaths)
		return
	}

	// 2. Initialize NASA Engine
	manager, err := cspice.NewEphemerisManager(kernelPaths...)
	if err != nil {
		initError = err
		log.Printf("❌ NASA Engine Failure: %v", err)
		setupErrorHandler(kernelPaths)
		return
	}

	// 3. Initialize AppConfig with all services and handlers
	appConfig, err = app.NewAppConfig(manager)
	if err != nil {
		initError = err
		log.Printf("❌ Application Initialization Failure: %v", err)
		setupErrorHandler(kernelPaths)
		return
	}

	// 4. Add debug endpoint
	appConfig.GetEngine().GET("/debug", func(c *gin.Context) {
		cwd, _ := os.Getwd()
		files, _ := os.ReadDir(cwd)
		var fileList []string
		for _, f := range files {
			fileList = append(fileList, f.Name())
		}

		c.JSON(200, gin.H{
			"status":       "Ready",
			"cwd":          cwd,
			"files":        fileList,
			"init_error":   "none",
			"kernel_paths": kernelPaths,
		})
	})
}

func setupErrorHandler(kernelPaths []string) {
	errorEngine := gin.New()

	errorEngine.GET("/debug", func(c *gin.Context) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":       "error",
			"init_error":   initError.Error(),
			"kernel_paths": kernelPaths,
		})
	})

	errorEngine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Application initialization failed",
			"detail":  initError.Error(),
		})
	})

	appConfig = &app.AppConfig{Engine: errorEngine}
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
	if appConfig == nil {
		http.Error(w, "Application not initialized", http.StatusInternalServerError)
		return
	}
	appConfig.ServeHTTP(w, r)
}
