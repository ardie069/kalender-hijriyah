package hijri

import (
	"time"

	"github.com/ardie069/kalender-hijriyah/internal/hijri"
	"github.com/ardie069/kalender-hijriyah/internal/visibility"
	"github.com/ardie069/kalender-hijriyah/pkg/cspice"
)

type Controller struct {
	DateSvc   *hijri.DateService
	CalendarSvc *hijri.CalendarService
	VisSvc    *visibility.Service
	Adapter   *cspice.Adapter
	startTime time.Time
}

func NewController(dateSvc *hijri.DateService, calSvc *hijri.CalendarService, visSvc *visibility.Service, adapter *cspice.Adapter) *Controller {
	return &Controller{
		DateSvc:   dateSvc,
		CalendarSvc: calSvc,
		VisSvc:    visSvc,
		Adapter:   adapter,
		startTime: time.Now(),
	}
}
