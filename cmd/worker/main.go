package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/ardie069/kalender-hijriyah/internal/notification"
	"github.com/ardie069/kalender-hijriyah/internal/notification/telegram"
	"github.com/ardie069/kalender-hijriyah/pkg/app"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it")
	}

	// Initialize NASA SPICE Engine
	manager, err := cspice.NewEphemerisManager(
		"pkg/cspice/kernels/de440s.bsp",
		"pkg/cspice/kernels/naif0012.tls",
		"pkg/cspice/kernels/pck00011.tpc",
	)
	if err != nil {
		log.Fatalf("❌ NASA Engine Failure: %v", err)
	}

	// Initialize AppConfig to get DateSvc and Astro
	appConfig, err := app.NewAppConfig(manager)
	if err != nil {
		log.Fatalf("❌ Application Initialization Failure: %v", err)
	}

	c := cron.New()

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set in environment")
	}

	tgNotifier, err := telegram.New(token)
	if err != nil {
		log.Fatalf("Failed to initialize telegram notifier: %v", err)
	}

	// Initialize Notification Service with dependencies
	adapter := cspice.GetAdapter(manager)
	notifSvc := notification.NewService(tgNotifier, appConfig.DateSvc, adapter)

	c.AddFunc(
		"0 8 * * *",
		func() {
			notifSvc.GenerateNotifications(context.Background())
		},
	)

	c.Start()

	// Start bot listener
	go notifSvc.StartListening(context.Background())

	select {}
}
