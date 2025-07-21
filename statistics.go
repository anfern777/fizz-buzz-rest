package main

import (
	"fmt"
	"net/http"
)

func (app *application) statisticsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "statistics endpoint")
}
