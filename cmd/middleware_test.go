package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

func TestCollectStats(t *testing.T) {
	app := &application{
		stats: &statsStore{
			requests: make(map[fizzbuzzParams]int),
			mu:       sync.RWMutex{},
		},
	}

	calledNext := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calledNext = true
		params, ok := r.Context().Value(paramsCtxKey).(*fizzbuzzParams)
		if !ok {
			t.Fatal("paramsCtxKey not found in context")
		}
		if params.Int1 != 3 || params.Int2 != 5 || params.Limit != 15 || params.Str1 != "fizz" || params.Str2 != "buzz" {
			t.Errorf("unexpected fizzbuzzParams: %+v", params)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz", nil)
	rr := httptest.NewRecorder()

	wrappedHandler := app.collectStats(nextHandler)
	wrappedHandler.ServeHTTP(rr, req)

	if !calledNext {
		t.Error("expected next handler to be called")
	}

	expectedParams := fizzbuzzParams{Int1: 3, Int2: 5, Limit: 15, Str1: "fizz", Str2: "buzz"}

	app.stats.mu.Lock()
	defer app.stats.mu.Unlock()
	if hits := app.stats.requests[expectedParams]; hits != 1 {
		t.Errorf("expected 1 hit for params %+v, got %d", expectedParams, hits)
	}
}

func TestEnableCORS(t *testing.T) {
	tests := []struct {
		name                 string
		originHeader         string
		trustedOrigins       []string
		expAllowOriginHeader string
	}{
		{
			name:                 "Trusted origin",
			originHeader:         "http://trusted.com",
			trustedOrigins:       []string{"http://trusted.com", "http://another.com"},
			expAllowOriginHeader: "http://trusted.com",
		},
		{
			name:                 "Untrusted origin",
			originHeader:         "http://untrusted.com",
			trustedOrigins:       []string{"http://trusted.com"},
			expAllowOriginHeader: "",
		},
		{
			name:                 "No origin header",
			originHeader:         "",
			trustedOrigins:       []string{"http://trusted.com"},
			expAllowOriginHeader: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{
				config: config{
					trustedOrigins: tt.trustedOrigins,
				},
			}

			calledNext := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calledNext = true
				w.Write([]byte("OK"))
			})

			wrappedHandler := app.enableCORS(nextHandler)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.originHeader != "" {
				req.Header.Set("Origin", tt.originHeader)
			}

			rr := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rr, req)

			if !calledNext {
				t.Error("expected next handler to be called")
			}

			if vary := rr.Header().Get("Vary"); vary != "Origin" {
				t.Errorf("expected Vary header to be 'Origin', but got %q", vary)
			}

			if allowOrigin := rr.Header().Get("Access-Control-Allow-Origin"); allowOrigin != tt.expAllowOriginHeader {
				t.Errorf("expected Access-Control-Allow-Origin header to be %q, but got %q",
					tt.expAllowOriginHeader, allowOrigin)
			}
		})
	}
}

func TestRecoverPanic(t *testing.T) {
	app := application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	expHeaders := http.Header{
		"Connection":   {"close"},
		"Content-Type": {"application/json"},
	}

	calledNext := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calledNext = true
		panic("test")
	})

	wrappedHandler := app.recoverPanic(nextHandler)

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, r)

	headers := w.Header()
	code := w.Code

	if !calledNext {
		t.Error("expected next handler to be called")
	}

	if !reflect.DeepEqual(headers, expHeaders) {
		t.Errorf("expected headers:\n%v\n but got:\n%v\n", expHeaders, headers)
	}

	if code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, but got %d", http.StatusInternalServerError, code)
	}
}

func TestRateLimiter(t *testing.T) {
	app := &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		config: config{
			limiter: struct {
				rate    float64
				burst   int
				enabled bool
			}{
				enabled: true,
				rate:    1,
				burst:   2,
			},
		},
	}
	calledNext := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calledNext = true
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := app.rateLimiter(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	for i := 1; i <= 2; i++ {
		rr := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("request %d: expected status OK; got %d", i, rr.Code)
		}
	}

	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)

	if !calledNext {
		t.Error("expected next handler to be called")
	}

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("request 3: expected status TooManyRequests; got %d", rr.Code)
	}
}
