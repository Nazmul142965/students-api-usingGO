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

	"github.com/nazmul14296/students-api/internal/config"
	"github.com/nazmul14296/students-api/internal/http/handlers/student"
	"github.com/nazmul14296/students-api/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup,

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initializes", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))

	//sertup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("server started ", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGALRM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to Srtart server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shhutdown succesfully")
}
