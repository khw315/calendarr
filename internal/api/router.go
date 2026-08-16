package api

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/khw315/calendarr/internal/config"
	"github.com/khw315/calendarr/internal/models"
	"github.com/khw315/calendarr/internal/services/scheduler"
)

type Router struct {
	cfgMgr     *config.Manager
	schedSvc   *scheduler.Service
	embeddedFS embed.FS
}

func NewRouter(cfgMgr *config.Manager, schedSvc *scheduler.Service, embeddedFS embed.FS) *Router {
	return &Router{
		cfgMgr:     cfgMgr,
		schedSvc:   schedSvc,
		embeddedFS: embeddedFS,
	}
}

func (r *Router) Setup() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// API Routes
	router.Route("/api", func(api chi.Router) {
		api.Get("/status", r.handleGetStatus)
		api.Get("/events", r.handleGetEvents)
		api.Get("/config", r.handleGetConfig)
		api.Post("/config", r.handlePostConfig)
		api.Post("/trigger", r.handlePostTrigger)
	})

	// Static Web UI serving
	publicSub, err := fs.Sub(r.embeddedFS, "public")
	if err != nil {
		log.Printf("⚠️ Could not sub public embed directory: %v", err)
	} else {
		fileServer := http.FileServer(http.FS(publicSub))
		router.Get("/*", func(w http.ResponseWriter, req *http.Request) {
			path := strings.TrimPrefix(req.URL.Path, "/")
			if path != "" {
				if _, err := publicSub.Open(path); err == nil {
					fileServer.ServeHTTP(w, req)
					return
				}
			}
			// Fallback to SPA index.html
			req.URL.Path = "/"
			fileServer.ServeHTTP(w, req)
		})
	}

	return router
}

func (r *Router) handleGetStatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := r.schedSvc.GetStatus()
	json.NewEncoder(w).Encode(status)
}

func (r *Router) handleGetEvents(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	events := r.schedSvc.GetCachedEvents()
	if events == nil {
		events = []*models.Event{}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"events": events,
		"count":  len(events),
	})
}

func (r *Router) handleGetConfig(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := r.cfgMgr.Get()
	json.NewEncoder(w).Encode(cfg)
}

func (r *Router) handlePostConfig(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newCfg models.Config
	if err := json.NewDecoder(req.Body).Decode(&newCfg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request JSON payload"})
		return
	}

	if err := r.cfgMgr.Save(&newCfg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	r.schedSvc.UpdateSchedule()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Configuration saved and schedule updated",
		"config":  newCfg,
	})
}

func (r *Router) handlePostTrigger(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	go func() {
		_, _ = r.schedSvc.TriggerRun(req.Context())
	}()

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "triggered",
		"message": "Manual run triggered successfully in background",
	})
}
