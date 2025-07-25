package main

import (
	"net/http"
)

func (app *application) statsHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Parameters fizzbuzzParams `json:"parameters"`
		Hits       int            `json:"hits"`
	}

	topParams, maxHits := app.computeStats()

	res := wrapper{"stats": response{
		Parameters: topParams,
		Hits:       maxHits,
	}}

	err := app.toJSONResponse(w, http.StatusOK, res, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
