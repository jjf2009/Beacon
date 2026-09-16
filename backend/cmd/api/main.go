package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jjf2009/beacon/backend/internal/config"
)

func main(){
	// load config 
	cfg := config.MustLoad()

	router :=http.NewServeMux()
	router.HandleFunc("GET /",func (w http.ResponseWriter, r *http.Request)  {
		fmt.Fprintln(w,"Beacon server is running")
	})
	router.HandleFunc("GET /api/projects",GetProjects)
	router.HandleFunc("POST /api/projects",CreateProjects)

	server := http.Server{
		Addr:cfg.HTTPServer.Addr,
		Handler:router,
	}
   	slog.Info("server started", "addr", cfg.HTTPServer.Addr)
	fmt.Printf("Server started:%s", cfg.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

go func() {
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatal("fail to start server", err)
    }
}()


	<-done
	slog.Info("shutting done the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown sucessffully")

}