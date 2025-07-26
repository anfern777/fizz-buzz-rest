package main

import (
	"encoding/json"
	"maps"
	"net"
	"net/http"
	"strings"
)

type wrapper map[string]any

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

func (app *application) computeStats() (fizzbuzzParams, int) {
	var maxHits int
	var topParams fizzbuzzParams

	app.stats.mu.RLock()
	defer app.stats.mu.RUnlock()

	for params, hits := range app.stats.requests {
		if hits > maxHits {
			maxHits = hits
			topParams = params
		}
	}

	return topParams, maxHits
}
