package main

import (
	"flag"
	"log/slog"
	"os"
)

type config struct {
	env     string
	port    int
	limiter struct {
		rate    float64
		burst   int
		enabled bool
	}
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

	flag.Float64Var(&cfg.limiter.rate, "limiter-rate", 2, "Rate limiter requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter burst limit")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()

	app := &application{
		config: cfg,
		logger: logger,
	}

	if err := app.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
