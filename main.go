package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	env  string
	port int
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production)")
	flag.IntVar(&cfg.port, "port", 8000, "server port")
	flag.Parse()

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 2 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sig := <-quit
	app.logger.Info(fmt.Sprintf("Received signal: %v. Shutting down server...\n", sig))
	srv.Shutdown(ctx)
}
