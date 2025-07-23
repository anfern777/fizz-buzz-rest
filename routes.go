package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/fizzbuzz", app.fizzbuzzHandler)
	router.HandlerFunc(http.MethodGet, "/stats", app.statisticsHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimiter(router)))
}
