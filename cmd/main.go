package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/anfern777/fizz-buzz-rest/internal/vcs"
)

type config struct {
	env     string
	port    int
	limiter struct {
		rate    float64
		burst   int
		enabled bool
	}
	trustedOrigins []string
}

type statsStore struct {
	requests map[fizzbuzzParams]int
	mu       sync.RWMutex
}

type application struct {
	config config
	logger *slog.Logger
	stats  *statsStore
}

var (
	version = vcs.Version()
)

func main() {
	var cfg config
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	flag.StringVar(&cfg.env, "env", "development", "Environment (development | production)")
	flag.IntVar(&cfg.port, "port", 8080, "server port")

	flag.Float64Var(&cfg.limiter.rate, "limiter-rate", 2, "Rate limiter requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter burst limit")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Func("trusted-origins", "Trusted CORS origins", func(val string) error {
		cfg.trustedOrigins = strings.Fields(val)
		return nil
	})
	showVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	app := &application{
		config: cfg,
		logger: logger,
		stats: &statsStore{
			requests: make(map[fizzbuzzParams]int),
		},
	}

	if err := app.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
