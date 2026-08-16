package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/khw315/calendarr/internal/config"
	"github.com/khw315/calendarr/internal/constants"
	"github.com/khw315/calendarr/internal/models"
	"github.com/khw315/calendarr/internal/services/calendar"
	"github.com/khw315/calendarr/internal/services/formatter"
	"github.com/khw315/calendarr/internal/services/platform"
)

type Service struct {
	mu           sync.Mutex
	cfgMgr       *config.Manager
	calSvc       *calendar.Service
	fmtSvc       *formatter.Service
	platSvc      *platform.Service
	cronRunner   *cron.Cron
	entryID      cron.EntryID
	lastRun      time.Time
	lastStatus   string
	lastErr      error
	cachedEvents []*models.Event
}

func NewService(cfgMgr *config.Manager, calSvc *calendar.Service, fmtSvc *formatter.Service, platSvc *platform.Service) *Service {
	return &Service{
		cfgMgr:     cfgMgr,
		calSvc:     calSvc,
		fmtSvc:     fmtSvc,
		platSvc:    platSvc,
		cronRunner: cron.New(),
	}
}

func (s *Service) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	cfg := s.cfgMgr.Get()
	s.updateScheduleLocked(cfg)

	s.cronRunner.Start()

	if cfg.ScheduleSettings.RunOnStartup {
		log.Println("🚀 'Run On Startup' enabled. Triggering initial run...")
		go func() {
			_, _ = s.TriggerRun(context.Background())
		}()
	}
}

func (s *Service) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cronRunner.Stop()
}

func (s *Service) UpdateSchedule() {
	s.mu.Lock()
	defer s.mu.Unlock()
	cfg := s.cfgMgr.Get()
	s.updateScheduleLocked(cfg)
}

func (s *Service) updateScheduleLocked(cfg *models.Config) {
	if s.entryID != 0 {
		s.cronRunner.Remove(s.entryID)
		s.entryID = 0
	}

	cronExpr := s.buildCronExpr(cfg.ScheduleSettings)
	if cronExpr == "" {
		log.Println("⚠️ Invalid or empty schedule expression")
		return
	}

	id, err := s.cronRunner.AddFunc(cronExpr, func() {
		log.Println("⏰ Scheduled run triggered...")
		_, _ = s.TriggerRun(context.Background())
	})

	if err != nil {
		log.Printf("❌ Failed to schedule cron [%s]: %v", cronExpr, err)
		return
	}

	s.entryID = id
	log.Printf("📅 Cron schedule set: %s (Entry ID: %d)", cronExpr, id)
}

func (s *Service) buildCronExpr(sched models.ScheduleSettings) string {
	if sched.CronSchedule != "" {
		return sched.CronSchedule
	}

	runTime := sched.RunTime
	if runTime == "" {
		runTime = "09:00"
	}

	var hour, min int
	_, err := fmt.Sscanf(runTime, "%d:%d", &hour, &min)
	if err != nil {
		hour = 9
		min = 0
	}

	switch sched.ScheduleType {
	case "DAILY":
		return fmt.Sprintf("%d %d * * *", min, hour)
	case "WEEKLY":
		day := sched.ScheduleDay
		if day == "" {
			day = "1" // Monday
		}
		return fmt.Sprintf("%d %d * * %s", min, hour, day)
	default:
		return fmt.Sprintf("%d %d * * *", min, hour)
	}
}

func (s *Service) TriggerRun(ctx context.Context) ([]*models.Event, error) {
	cfg := s.cfgMgr.Get()
	loc := cfg.TimezoneLocation
	if loc == nil {
		loc = time.UTC
	}

	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	var startDate, endDate time.Time
	if cfg.ScheduleSettings.ScheduleType == "DAILY" {
		// For DAILY schedule type, set window from start to end of current day (00:00 - 23:59:59)
		startDate = startOfDay
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)
	} else {
		// 7 days window for WEEKLY / default starting from start of today
		startDate = startOfDay
		endDate = startOfDay.AddDate(0, 0, 7).Add(-1 * time.Nanosecond)
	}

	log.Printf("⚡ Executing calendar run (%s) from %s to %s", cfg.ScheduleSettings.ScheduleType, startDate.Format("2006-01-02 15:04"), endDate.Format("2006-01-02 15:04"))

	events, err := s.calSvc.FetchEvents(ctx, cfg, startDate, endDate)

	s.mu.Lock()
	s.lastRun = now
	s.cachedEvents = events
	if err != nil {
		s.lastStatus = "ERROR"
		s.lastErr = err
		s.mu.Unlock()
		return nil, err
	}
	s.lastStatus = "SUCCESS"
	s.lastErr = nil
	s.mu.Unlock()

	formatted := s.fmtSvc.Format(events, cfg, startDate, endDate)
	_ = s.platSvc.Dispatch(ctx, cfg, formatted.Discord, formatted.Slack)

	return events, nil
}

func (s *Service) GetStatus() map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	nextRun := time.Time{}
	if s.entryID != 0 {
		entry := s.cronRunner.Entry(s.entryID)
		nextRun = entry.Next
	}

	errStr := ""
	if s.lastErr != nil {
		errStr = s.lastErr.Error()
	}

	return map[string]interface{}{
		"last_run":     s.lastRun,
		"last_status":  s.lastStatus,
		"last_error":   errStr,
		"next_run":     nextRun,
		"events_count": len(s.cachedEvents),
		"version":      constants.Version,
	}
}

func (s *Service) GetCachedEvents() []*models.Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cachedEvents
}
