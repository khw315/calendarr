package scheduler

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/khw315/calendarr/internal/config"
	"github.com/khw315/calendarr/internal/models"
	"github.com/khw315/calendarr/internal/services/calendar"
	"github.com/khw315/calendarr/internal/services/formatter"
	"github.com/khw315/calendarr/internal/services/platform"
)

func TestSchedulerServiceLifecycle(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "calendarr.json")

	cfgMgr := config.NewManager(configPath)
	cfg := cfgMgr.Get()
	cfg.ScheduleSettings.ScheduleType = "DAILY"
	cfg.ScheduleSettings.RunTime = "08:30"
	cfg.ScheduleSettings.RunOnStartup = true
	_ = cfgMgr.Save(cfg)

	calSvc := calendar.NewService()
	fmtSvc := formatter.NewService()
	platSvc := platform.NewService()

	svc := NewService(cfgMgr, calSvc, fmtSvc, platSvc)
	svc.Start()
	defer svc.Stop()

	// Update to Custom Cron
	cfg.ScheduleSettings.CronSchedule = "0 10 * * *"
	_ = cfgMgr.Save(cfg)
	svc.UpdateSchedule()

	status := svc.GetStatus()
	if status == nil {
		t.Fatalf("Expected status map to be non-nil")
	}

	events, err := svc.TriggerRun(context.Background())
	if err != nil {
		t.Fatalf("TriggerRun failed: %v", err)
	}

	cached := svc.GetCachedEvents()
	if len(cached) != len(events) {
		t.Errorf("Mismatch in cached events count")
	}
}

func TestBuildCronExpr(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	s1 := models.ScheduleSettings{ScheduleType: "WEEKLY", RunTime: "invalid", ScheduleDay: ""}
	expr1 := svc.buildCronExpr(s1)
	if expr1 != "0 9 * * 1" {
		t.Errorf("Expected fallback weekly cron 0 9 * * 1, got %s", expr1)
	}

	s2 := models.ScheduleSettings{CronSchedule: "*/5 * * * *"}
	expr2 := svc.buildCronExpr(s2)
	if expr2 != "*/5 * * * *" {
		t.Errorf("Expected custom cron expression, got %s", expr2)
	}
}
