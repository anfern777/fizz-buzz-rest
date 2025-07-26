package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestStatsHandler(t *testing.T) {
	tests := []struct {
		name  string
		stats *statsStore

		expRes string
	}{
		{
			name: "no stats available",
			stats: &statsStore{
				requests: map[fizzbuzzParams]int{
					{}: 0,
				},
			},
			expRes: `{"stats":"no stats available yet"}` + "\n",
		},
		{
			name: "happy path",
			stats: &statsStore{
				requests: map[fizzbuzzParams]int{
					{Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz"}: 10,
				},
			},
			expRes: `{"stats":{"parameters":{"int1":3,"int2":5,"limit":100,"str1":"fizz","str2":"buzz"},"hits":10}}` + "\n",
		},
	}
	app := &application{
		stats: &statsStore{
			requests: make(map[fizzbuzzParams]int),
			mu:       sync.RWMutex{},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.stats.requests = tt.stats.requests
			rr := httptest.NewRecorder()
			app.statsHandler(rr, req)
			resp := rr.Result()
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status 200 OK, got %d", resp.StatusCode)
			}

			body, _ := io.ReadAll(resp.Body)
			expected := tt.expRes
			if string(body) != expected {
				t.Errorf("expected body:\n%s\ngot:\n%s", expected, string(body))
			}
		})
	}
}
