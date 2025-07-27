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

	var res wrapper
	if maxHits == 0 {
		res = wrapper{"most_frequent_request": "no stats available yet"}
	} else {
		res = wrapper{"most_frequent_request": response{
			Parameters: topParams,
			Hits:       maxHits,
		}}
	}

	err := app.toJSONResponse(w, http.StatusOK, res, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
