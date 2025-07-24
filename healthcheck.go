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

	if err := app.toJSONResponse(w, http.StatusOK, &wrapper{"health": health}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
