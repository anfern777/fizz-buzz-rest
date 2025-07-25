package main

import (
	"encoding/json"
	"maps"
	"net"
	"net/http"
	"strings"
)

type wrapper map[string]any

// Returns the remote IP of the client. It gives precedence to
// X-Forwarded-For and X-Real-IP request headers. If these are
// not present, it gets the remote IP from RemoteAddr.
func getRemoteIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}

	if xrIP := r.Header.Get("X-Real-Ip"); xrIP != "" {
		return xrIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	return ip
}

// converts data into a json string, adds the necessary
// headers and writes to the ResponseWriter
func (app *application) toJSONResponse(w http.ResponseWriter, status int, data any, headers http.Header) error {
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

// computeStats returns the most frequently requested
// fizzbuzz parameters and their hit count.
func (app *application) computeStats() (fizzbuzzParams, int) {
	app.stats.mu.RLock()
	defer app.stats.mu.RUnlock()

	var maxHits int
	var topParams fizzbuzzParams
	for params, hits := range app.stats.requests {
		if hits > maxHits {
			maxHits = hits
			topParams = params
		}
	}

	return topParams, maxHits
}
