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
)

func main() {
	//load config
	cfg:= config.MustLoad()

	//database setup
	//setup router
	router:= http .NewServeMux()

	router.HandleFunc("POST /api/students",student.Create())


	//setup server
	server:= http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}
	slog.Info("Server Started at Address",slog.String("address",cfg.Addr))
	
	
	done := make(chan os.Signal,1)

	signal.Notify(done, os.Interrupt,syscall.SIGINT,syscall.SIGTERM)

	go func(){
		err:= server.ListenAndServe()
		if err!=nil{
			log.Fatal("Failed to start server")
		}
		
	}()
	
	<-done

	slog.Info("Shutting Down server")
	ctx,cancel := context.WithTimeout(context.Background(),time.Second*5)
	defer cancel()

	err := server.Shutdown(ctx)

	if err!=nil{
		slog.Error("Failed to shutdown server", slog.String("error",err.Error()))
	}
	slog.Info("Server ShutDown Successfully")
}