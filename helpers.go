package main

import (
	"encoding/json"
	"maps"
	"net"
	"net/http"
)

type wraper map[string]any

// Returns the remote IP of the client. It gives precedence to
// X-Forwarded-For and X-Real-IP request headers. If these are
// not present, it gets the remote IP from RemoteAddr.
func getRemoteIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return xff
	}

	xrIP := r.Header.Get("X-Real-IP")
	if xrIP != "" {
		return xrIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ip
	}

	return ""
}

// converts data into a json string, adds the necessary
// headers and writes to the ResponseWriter
func (app *application) toJsonResponse(w http.ResponseWriter, status int, data any, headers http.Header) error {
	maps.Copy(w.Header(), headers)

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
