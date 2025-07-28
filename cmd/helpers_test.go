package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"testing"
)

func TestToJsonResponse(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		data       any
		headers    http.Header
		expStatus  int
		expResp    string
		expHeaders http.Header
	}{
		{
			name:   "happy path",
			status: http.StatusOK,
			data: map[string]any{
				"Test": "test",
			},
			headers: http.Header{
				"Test": {"test"},
			},
			expStatus: http.StatusOK,
			expResp:   `{"Test":"test"}` + "\n",
			expHeaders: http.Header{
				"Test":         {"test"},
				"Content-Type": {"application/json"},
			},
		},
	}

	for _, tt := range tests {
		mockApp := application{}
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			err := mockApp.toJSONResponse(rr, tt.status, tt.data, tt.headers)

			if err != nil {
				t.Errorf("expected no error, but got %v", err)
			}

			if rr.Code != tt.expStatus {
				t.Errorf("expected status %d, got %d", tt.expStatus, rr.Code)
			}

			if rr.Body.String() != tt.expResp {
				t.Errorf("expected body %q, got %q", tt.expResp, rr.Body.String())
			}

			for k, v := range tt.expHeaders {
				if got := rr.Header().Values(k); !reflect.DeepEqual(got, v) {
					t.Errorf("expected header %q to be %v, got %v", k, v, got)
				}
			}

		})
	}
}

func TestComputeStats(t *testing.T) {
	params1 := fizzbuzzParams{Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz"}
	params2 := fizzbuzzParams{Int1: 2, Int2: 4, Limit: 50, Str1: "a", Str2: "b"}
	params3 := fizzbuzzParams{Int1: 7, Int2: 11, Limit: 200, Str1: "c", Str2: "d"}

	tests := []struct {
		name         string
		initialStats map[fizzbuzzParams]int
		expHits      int
		expParams    []fizzbuzzParams
	}{
		{
			name:         "empty stats",
			initialStats: map[fizzbuzzParams]int{},

			expHits:   0,
			expParams: []fizzbuzzParams{{}},
		},
		{
			name: "single entry",
			initialStats: map[fizzbuzzParams]int{
				params1: 10,
			},

			expHits:   10,
			expParams: []fizzbuzzParams{params1},
		},
		{
			name: "Multiple entries with a clear winner",
			initialStats: map[fizzbuzzParams]int{
				params1: 10,
				params2: 50,
				params3: 5,
			},

			expHits:   50,
			expParams: []fizzbuzzParams{params2},
		},
		{
			name: "Multiple entries with a tie",
			initialStats: map[fizzbuzzParams]int{
				params1: 50,
				params2: 10,
				params3: 50,
			},

			expHits:   50,
			expParams: []fizzbuzzParams{params1, params3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{
				stats: &statsStore{
					requests: tt.initialStats,
				},
			}

			topParams, maxHits := app.computeStats()

			if maxHits != tt.expHits {
				t.Errorf("expected max hits %d; got %d", tt.expHits, maxHits)
			}

			if !slices.Contains(tt.expParams, topParams) {
				t.Errorf("expected top params to be one of %+v; got %+v", tt.expParams, topParams)
			}
		})
	}
}
