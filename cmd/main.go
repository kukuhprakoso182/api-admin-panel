package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-admin-panel/config"
	"api-admin-panel/internal/bootstrap"
	"api-admin-panel/migrations"
	applogger "api-admin-panel/pkg/logger"
	"api-admin-panel/seeder"
)

const version = "1.0.0"

func main() {
	// ── Flags ─────────────────────────────────────────────────────────────────
	withFresh := flag.Bool("fresh", false, "drop all tables before running migrations and seeders")
	withMigrate := flag.Bool("migrate", false, "run database migrations before starting")
	withSeed := flag.Bool("seed", false, "run database seeders before starting")
	flag.Parse()

	// ── Config ────────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[FATAL] config: %v\n", err)
		os.Exit(1)
	}

	// ── Logger ────────────────────────────────────────────────────────────────
	log, err := applogger.New(cfg.Log.Dir, cfg.App.Name, version, cfg.App.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[FATAL] logger: %v\n", err)
		os.Exit(1)
	}

	// ── Bootstrap (DI) ────────────────────────────────────────────────────────
	app, err := bootstrap.New(cfg, log)
	if err != nil {
		log.Fatal("failed to bootstrap application", applogger.F("error", err))
	}

	// ── Migrations ────────────────────────────────────────────────────────────
	if *withMigrate && *withFresh {
		if err := migrations.Run(app.DB, *withFresh); err != nil {
			log.Fatal("failed to run migrations", applogger.F("error", err))
		}
		log.Info("migrations completed")
	}

	// ── Seeder ────────────────────────────────────────────────────────────────
	if *withSeed {
		if err := seeder.Run(app.DB); err != nil {
			log.Fatal("failed to run seeders", applogger.F("error", err))
		}
		log.Info("seeders completed")
	}

	// ── Banner ────────────────────────────────────────────────────────────────
	log.PrintBanner("connected", cfg.App.Port)

	// ── HTTP Server ───────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:              addr,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      35 * time.Second, // harus lebih besar dari Timeout middleware
		IdleTimeout:       60 * time.Second,
	}

	// Jalankan server di goroutine agar tidak blocking signal handler
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", applogger.F("error", err))
		}
	}()

	// ── Graceful Shutdown ─────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Tunggu sinyal shutdown (Ctrl+C atau kill)

	log.Info("shutting down server...")

	// Beri waktu 10 detik untuk menyelesaikan request yang sedang berjalan
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("forced shutdown", applogger.F("error", err))
	}

	log.Info("server exited gracefully")
}
