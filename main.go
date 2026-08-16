package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/khw315/calendarr/internal/api"
	"github.com/khw315/calendarr/internal/config"
	"github.com/khw315/calendarr/internal/services/calendar"
	"github.com/khw315/calendarr/internal/services/formatter"
	"github.com/khw315/calendarr/internal/services/platform"
	"github.com/khw315/calendarr/internal/services/scheduler"
)

//go:embed public/*
var embeddedPublicFS embed.FS

func main() {
	log.Println("🚀 Starting Calendarr (Golang)...")

	// 1. Config Manager
	cfgMgr := config.NewManager("")

	// 2. Services
	calSvc := calendar.NewService()
	fmtSvc := formatter.NewService()
	platSvc := platform.NewService()
	schedSvc := scheduler.NewService(cfgMgr, calSvc, fmtSvc, platSvc)

	// 3. Start Background Scheduler
	schedSvc.Start()

	// 4. HTTP Router
	router := api.NewRouter(cfgMgr, schedSvc, embeddedPublicFS)
	handler := router.Setup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🌐 Calendarr Web UI & API server running on http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ HTTP server error: %v", err)
		}
	}()

	<-stop
	log.Println("🛑 Shutting down Calendarr gracefully...")

	schedSvc.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("⚠️ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Calendarr process exited safely.")
}
