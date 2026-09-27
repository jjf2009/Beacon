package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jjf2009/beacon/backend/internal/check"
	"github.com/jjf2009/beacon/backend/internal/config"
	"github.com/jjf2009/beacon/backend/internal/database"
	"github.com/jjf2009/beacon/backend/internal/endpoint"
	"github.com/jjf2009/beacon/backend/internal/server"
	"github.com/jjf2009/beacon/backend/internal/worker"
)

func main() {
	cfg := config.MustLoad()

	db, err := database.New(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// build repos here so they can be shared between server and worker
	endpointRepo := endpoint.NewRepository(db)
	checkRepo := check.NewRepository(db)

	// context that gets cancelled on Ctrl+C — signals both server and worker to stop
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// start background worker — checks endpoints on a ticker
	w := worker.New(endpointRepo, checkRepo)
	go w.Start(ctx)

	srv := server.New(cfg.HTTPServer.Addr, db)
	slog.Info("server started", "addr", cfg.HTTPServer.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server: ", err)
		}
	}()

	// block until context is cancelled (Ctrl+C / SIGTERM)
	<-ctx.Done()
	slog.Info("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown server", "error", err)
	}

	slog.Info("server shutdown successfully")
}
