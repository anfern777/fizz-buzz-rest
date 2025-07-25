package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthcheck", app.healthcheckHandler)
	mux.HandleFunc("GET /stats", app.statsHandler)
	mux.Handle("GET /fizzbuzz", app.collectStats(http.HandlerFunc(app.fizzbuzzHandler)))

	return app.recoverPanic(app.enableCORS(app.rateLimiter(mux)))
}
