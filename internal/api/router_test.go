package api

import (
	"bytes"
	"embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/khw315/calendarr/internal/config"
	"github.com/khw315/calendarr/internal/models"
	"github.com/khw315/calendarr/internal/services/calendar"
	"github.com/khw315/calendarr/internal/services/formatter"
	"github.com/khw315/calendarr/internal/services/platform"
	"github.com/khw315/calendarr/internal/services/scheduler"
)

func TestAPIRoutes(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "calendarr.json")

	cfgMgr := config.NewManager(configPath)
	calSvc := calendar.NewService()
	fmtSvc := formatter.NewService()
	platSvc := platform.NewService()
	schedSvc := scheduler.NewService(cfgMgr, calSvc, fmtSvc, platSvc)

	var emptyFS embed.FS
	router := NewRouter(cfgMgr, schedSvc, calSvc, emptyFS)
	handler := router.Setup()

	// 1. GET /api/status
	reqStatus := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	recStatus := httptest.NewRecorder()
	handler.ServeHTTP(recStatus, reqStatus)
	if recStatus.Code != http.StatusOK {
		t.Errorf("Expected status 200 on /api/status, got %d", recStatus.Code)
	}

	// 1b. GET /api/version
	reqVer := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	recVer := httptest.NewRecorder()
	handler.ServeHTTP(recVer, reqVer)
	if recVer.Code != http.StatusOK {
		t.Errorf("Expected status 200 on /api/version, got %d", recVer.Code)
	}

	// 2. GET /api/events
	reqEvents := httptest.NewRequest(http.MethodGet, "/api/events", nil)
	recEvents := httptest.NewRecorder()
	handler.ServeHTTP(recEvents, reqEvents)
	if recEvents.Code != http.StatusOK {
		t.Errorf("Expected status 200 on /api/events, got %d", recEvents.Code)
	}

	// 3. GET /api/config
	reqCfg := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	recCfg := httptest.NewRecorder()
	handler.ServeHTTP(recCfg, reqCfg)
	if recCfg.Code != http.StatusOK {
		t.Errorf("Expected status 200 on /api/config, got %d", recCfg.Code)
	}

	// 4. POST /api/config
	newCfg := models.DefaultConfig()
	newCfg.Timezone = "Asia/Jakarta"
	bodyBytes, _ := json.Marshal(newCfg)
	reqPostCfg := httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewBuffer(bodyBytes))
	recPostCfg := httptest.NewRecorder()
	handler.ServeHTTP(recPostCfg, reqPostCfg)
	if recPostCfg.Code != http.StatusOK {
		t.Errorf("Expected status 200 on POST /api/config, got %d", recPostCfg.Code)
	}

	// 4b. POST /api/config (invalid JSON)
	reqBadPost := httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewBufferString("{invalid"))
	recBadPost := httptest.NewRecorder()
	handler.ServeHTTP(recBadPost, reqBadPost)
	if recBadPost.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 on bad JSON POST /api/config, got %d", recBadPost.Code)
	}

	// 5. GET /api/languages
	reqLang := httptest.NewRequest(http.MethodGet, "/api/languages", nil)
	recLang := httptest.NewRecorder()
	handler.ServeHTTP(recLang, reqLang)
	if recLang.Code != http.StatusOK {
		t.Errorf("Expected status 200 on /api/languages, got %d", recLang.Code)
	}

	// 6. GET /api/schedule
	reqSched := httptest.NewRequest(http.MethodGet, "/api/schedule", nil)
	recSched := httptest.NewRecorder()
	handler.ServeHTTP(recSched, reqSched)
	if recSched.Code != http.StatusOK {
		t.Errorf("Expected status 200 on /api/schedule, got %d", recSched.Code)
	}

	// 7. GET /api/releases
	reqRel := httptest.NewRequest(http.MethodGet, "/api/releases?days=3", nil)
	recRel := httptest.NewRecorder()
	handler.ServeHTTP(recRel, reqRel)
	if recRel.Code != http.StatusOK {
		t.Errorf("Expected status 200 on /api/releases, got %d", recRel.Code)
	}

	// 7b. GET /api/past-releases
	reqPastRel := httptest.NewRequest(http.MethodGet, "/api/past-releases?days=3", nil)
	recPastRel := httptest.NewRecorder()
	handler.ServeHTTP(recPastRel, reqPastRel)
	if recPastRel.Code != http.StatusOK {
		t.Errorf("Expected status 200 on /api/past-releases, got %d", recPastRel.Code)
	}

	// 8. GET /api/timezones
	reqTZ := httptest.NewRequest(http.MethodGet, "/api/timezones", nil)
	recTZ := httptest.NewRecorder()
	handler.ServeHTTP(recTZ, reqTZ)
	if recTZ.Code != http.StatusOK {
		t.Errorf("Expected status 200 on /api/timezones, got %d", recTZ.Code)
	}

	// 9. OPTIONS /api/config (CORS preflight)
	reqOpt := httptest.NewRequest(http.MethodOptions, "/api/config", nil)
	recOpt := httptest.NewRecorder()
	handler.ServeHTTP(recOpt, reqOpt)
	if recOpt.Code != http.StatusOK {
		t.Errorf("Expected status 200 on OPTIONS /api/config, got %d", recOpt.Code)
	}

	// 10. POST /api/trigger
	reqTrigger := httptest.NewRequest(http.MethodPost, "/api/trigger", nil)
	recTrigger := httptest.NewRecorder()
	handler.ServeHTTP(recTrigger, reqTrigger)
	if recTrigger.Code != http.StatusOK {
		t.Errorf("Expected status 200 on POST /api/trigger, got %d", recTrigger.Code)
	}
}
