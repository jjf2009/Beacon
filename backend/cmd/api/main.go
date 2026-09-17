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

	"github.com/jjf2009/beacon/backend/internal/config"
	"github.com/jjf2009/beacon/backend/internal/database"
	"github.com/jjf2009/beacon/backend/internal/server"
)

func main() {
	cfg := config.MustLoad()
	db, err := database.New(cfg)
    if err != nil {
    slog.Error("failed to connect to database", "error", err)
    os.Exit(1)
   }
   defer db.Close()

	srv := server.New(cfg.HTTPServer.Addr,db)

	slog.Info("server started", "addr", cfg.HTTPServer.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server: ", err)
		}
	}()

	<-done
	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", "error", err)
	}

	slog.Info("server shutdown successfully")
}
