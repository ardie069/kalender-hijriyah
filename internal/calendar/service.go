package calendar

import "time"

// CoreService for general calendar calculations
type CoreService struct{}

func NewCoreService() *CoreService {
	return &CoreService{}
}

func (s *CoreService) DateDiff(t1, t2 time.Time) int {
	return int(t2.Sub(t1).Hours() / 24)
}
