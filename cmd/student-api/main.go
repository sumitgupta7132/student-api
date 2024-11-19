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

	"github.com/sumitgupta7132/student-api/internal/config"
	"github.com/sumitgupta7132/student-api/internal/http/handlers/student"
	"github.com/sumitgupta7132/student-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("Failed to connect to DB ", err)
	}
	slog.Info("storage intialised", slog.String("env", cfg.Env), slog.String("version", "1.0"))
	// setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.Create(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetStudentById(storage))
	router.HandleFunc("GET /api/students", student.GetStudents(storage))
	// setup server

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	slog.Info("Server started at", slog.String("address", cfg.Address))
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	slog.Info("Shutting the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to Shutdowm server", slog.String("error", err.Error()))
	}
	slog.Info("Server Shutdown successfully")

}
