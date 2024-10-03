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

	"github.com/SandipKurmi/students-api/internal/config"
	"github.com/SandipKurmi/students-api/internal/http/handlers/student"
	"github.com/SandipKurmi/students-api/internal/storage/sqlite"
)

func main() {
	// Load config
	cfg := config.MustLoad()

	// Database setup
	storage, err := sqlite.New(cfg)
	// post gres

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Database connected", slog.String("path", cfg.StoragePath))


	
	// Create a new HTTP router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/v1/students", student.New(storage))
	// get by id
	router.HandleFunc("GET /api/v1/students/{id}", student.GetById(storage))


	// Start server
	server := http.Server{
		Addr:    cfg.HttpServer.Address,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.HttpServer.Address))

	// Signal handling
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start the server in a separate goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-done
	slog.Info("Server stopped", slog.String("address", cfg.HttpServer.Address))

	// Graceful shutdown with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the existing 'err' variable instead of redeclaring it
	if err = server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("Server Shutdown successfully")
}