package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestStatsHandler_Success(t *testing.T) {
	app := &application{
		stats: &statsStore{
			requests: map[fizzbuzzParams]int{
				{Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz"}: 10,
			},
			mu: sync.RWMutex{},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	rr := httptest.NewRecorder()

	app.statsHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	expected := `{"stats":{"parameters":{"int1":3,"int2":5,"limit":100,"str1":"fizz","str2":"buzz"},"hits":10}}` + "\n"
	if string(body) != expected {
		t.Errorf("expected body:\n%s\ngot:\n%s", expected, string(body))
	}
}
