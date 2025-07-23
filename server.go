package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 2 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	shutdownErr := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		sig := <-quit
		app.logger.Info("graceful shutdown - handling remaining connections", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownErr <- srv.Shutdown(ctx)
	}()

	app.logger.Info("server started", "addr", srv.Addr, "env", app.config.env)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErr
	if err != nil {
		return err
	}

	app.logger.Info("server stopped", "addr", srv.Addr)

	return nil
}
