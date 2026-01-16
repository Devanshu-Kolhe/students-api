package main

import (
	"context"
	// "fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Devanshu-Kolhe/students-api/internal/config"
	"github.com/Devanshu-Kolhe/students-api/internal/http/handlers/student"

	// "github.com/Devanshu-Kolhe/students-api/internal/storage"
	"github.com/Devanshu-Kolhe/students-api/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage Initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.Create(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetByID(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	router.HandleFunc("PUT /api/students/{id}", student.UpdateStudentById(storage))
	router.HandleFunc("DELETE /api/students/{id}",student.DeleteStudentById(storage))
	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Server Started at Address", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}

	}()

	<-done

	slog.Info("Shutting Down server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server ShutDown Successfully")
}
