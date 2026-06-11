package notification

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/calendar"
	"github.com/ardie069/kalender-hijriyah/internal/hijri"
	"github.com/ardie069/kalender-hijriyah/internal/notification/telegram"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
)

// Service handles notification logic
type Service struct {
	notifier *telegram.Notifier
	dateSvc  *hijri.DateService
	astro    *cspice.Adapter
}

func NewService(notifier *telegram.Notifier, dateSvc *hijri.DateService, astro *cspice.Adapter) *Service {
	return &Service{
		notifier: notifier,
		dateSvc:  dateSvc,
		astro:    astro,
	}
}

func (s *Service) GenerateNotifications(ctx context.Context) {
	log.Println("Cron job triggered: GenerateNotifications")

	// Get current UTC date (since sunset checking is typically evening)
	now := time.Now().UTC()
	hDate := calendar.GetTabularHijri(now)

	// Check if today is the 29th of the Hijri month
	if hDate.Day == 29 {
		log.Printf("Today is %d %s %d. Generating predictions for the end of the month...", hDate.Day, hDate.MonthName, hDate.Year)

		// Parse CHAT_ID from string to int64
		var chatID int64
		if chatIDStr := os.Getenv("CHAT_ID"); chatIDStr != "" {
			parsedID, err := strconv.ParseInt(chatIDStr, 10, 64)
			if err == nil {
				chatID = parsedID
			} else {
				log.Printf("Failed to parse CHAT_ID: %v", err)
			}
		}

		// TODO: Fetch users from repository. Here we mock it:
		users := []User{
			{Username: "Pengguna (Jakarta)", Latitude: -6.2, Longitude: 106.8, Timezone: "WIB", TelegramID: chatID},
		}

		for _, user := range users {
			sunset, err := s.astro.GetSunsetFast(now, user.Latitude, user.Longitude)
			if err != nil {
				log.Printf("Failed to get sunset for user %d: %v", user.TelegramID, err)
				continue
			}

			moonset, _ := s.astro.GetMoonset(sunset, user.Latitude, user.Longitude)
			alt, elong, _, width := s.astro.CalculateTopocentricParamsGlobal(sunset, user.Latitude, user.Longitude)

			// Simple visibility rule based on MABIMS (Altitude > 3, Elongation > 6.4)
			isVisible := alt > 3.0 && elong > 6.4

			pred := Prediction{
				Date:         now,
				Sunset:       sunset,
				Moonset:      moonset,
				Altitude:     alt,
				Azimuth:      0, // Optional or calculated separately
				Elongation:   elong,
				AgeHours:     0,                      // Would need ijtima time to calculate accurately
				Illumination: (width * 60.0) / 100.0, // rough proxy for demonstration
				Visible:      isVisible,
			}

			msg := FormatPredictionMessage(user, pred)
			err = s.notifier.Send(ctx, user.TelegramID, msg)
			if err != nil {
				log.Printf("Failed to send notification to user %d: %v\n", user.TelegramID, err)
			} else {
				log.Printf("Successfully sent Hilal prediction to user %d", user.TelegramID)
			}
		}
	} else {
		log.Printf("Today is %d %s %d (Not the 29th). Skipping notification.", hDate.Day, hDate.MonthName, hDate.Year)
	}
}

func (s *Service) StartListening(ctx context.Context) {
	log.Println("Bot is now listening for incoming messages...")
	s.notifier.StartListener(ctx, func(chatID int64, text string) {
		log.Printf("Received message from %d: %s\n", chatID, text)

		if strings.HasPrefix(text, "/start") || strings.HasPrefix(text, "/hilal") {
			// Mock data
			user := User{Username: "Pengguna", Latitude: -6.2, Longitude: 106.8, Timezone: "WIB", TelegramID: chatID}
			pred := Prediction{
				Date:         time.Now(),
				Sunset:       time.Now().Add(time.Hour * 6),
				Moonset:      time.Now().Add(time.Hour * 7),
				Visible:      true,
				Altitude:     5.5,
				Azimuth:      270.5,
				Elongation:   7.2,
				AgeHours:     15.3,
				Illumination: 0.02,
			}

			msg := FormatPredictionMessage(user, pred)
			err := s.notifier.Send(ctx, chatID, msg)
			if err != nil {
				log.Printf("Failed to reply to user %d: %v\n", chatID, err)
			}
		} else {
			// Default response
			msg := "Silakan ketik /hilal untuk mendapatkan prediksi visibilitas hilal hari ini."
			_ = s.notifier.Send(ctx, chatID, msg)
		}
	})
}
