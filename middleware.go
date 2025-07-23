package main

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"golang.org/x/time/rate"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Add("Vary", "Origin")

		if slices.Contains(app.config.trustedOrigins, origin) {
			w.Header().Add("Access-Control-Allow-Origin", origin)
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if pv := recover(); pv != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%v", pv))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimiter(next http.Handler) http.Handler {
	if !app.config.limiter.enabled {
		return next
	}

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	clients := make(map[string]*client)

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 10*time.Minute {
					delete(clients, ip)
				}
			}
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getRemoteIP(r)
		_, ok := clients[ip]
		if !ok {
			clients[ip] = &client{
				limiter:  rate.NewLimiter(rate.Limit(app.config.limiter.rate), app.config.limiter.burst),
				lastSeen: time.Now(),
			}
		}
		if !clients[ip].limiter.Allow() {
			app.errorResponse(w, r, http.StatusTooManyRequests, "rate limit exceeded")
		}

		next.ServeHTTP(w, r)
	})
}
