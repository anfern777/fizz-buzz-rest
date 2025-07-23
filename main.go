package main

import (
	"flag"
	"log/slog"
	"os"
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

	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production)")
	flag.IntVar(&cfg.port, "port", 8000, "server port")
	flag.Parse()

	app := &application{
		config: cfg,
		logger: logger,
	}

	if err := app.serve(); err != nil {
		app.logger.Error("server could not listen", "reason", err.Error())
		os.Exit(1)
	}
}
