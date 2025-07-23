package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	type healthcheck struct {
		Environment string `json:"env"`
		Status      string `json:"status"`
	}

	health := &healthcheck{
		Environment: app.config.env,
		Status:      "available",
	}

	if err := app.toJsonResponse(w, http.StatusOK, health, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
