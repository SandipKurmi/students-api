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

	"github.com/SandipKurmi/students-api/internal/config"
	"github.com/SandipKurmi/students-api/internal/http/handlers/student"
)

func main() {
	// fmt.Println("Welcome to go project")

	cfg := config.MustLoad()

	fmt.Printf("%+v\n", cfg)

	// load config
	// data base setup
	// setup router

	router := http.NewServeMux()

	router.HandleFunc("POST /api/v1/students", student.New())

	// start server

	server := http.Server{
		Addr: cfg.HttpServer.Address,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.HttpServer.Address))


	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)


	go func() {
		if err := server.ListenAndServe(); 
		
		err != nil {
			log.Fatal(err)
		}
	}()

	<-done


	slog.Info("Server stopped", slog.String("address", cfg.HttpServer.Address))


	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Server forced to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("Server Shutdown successfully")
}