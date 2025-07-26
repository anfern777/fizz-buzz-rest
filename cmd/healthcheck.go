package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	type healthcheck struct {
		Status      string `json:"status"`
		Environment string `json:"env"`
		Version     string `json:"version"`
	}

	health := &healthcheck{
		Status:      "available",
		Environment: app.config.env,
		Version:     version,
	}

	err := app.toJSONResponse(w, http.StatusOK, &wrapper{"health": health}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
