package main

import "net/http"

func (app *application) logError(r *http.Request, err error) {
	var (
		uri    = r.URL.RequestURI()
		method = r.Method
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	wr := &wraper{"error": message}
	if err := app.toJsonResponse(w, status, wr, nil); err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process the request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}
