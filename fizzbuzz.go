package main

import (
	"fmt"
	"net/http"
)

func (app *application) fizzbuzzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "fizzbuzz endpoint")
}
