package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func TestFizzbuzzHandler_Success(t *testing.T) {
	app := &application{}

	params := &fizzbuzzParams{
		Int1:  3,
		Int2:  5,
		Limit: 15,
		Str1:  "fizz",
		Str2:  "buzz",
	}

	ctx := context.WithValue(context.Background(), paramsCtxKey, params)
	req := httptest.NewRequest(http.MethodGet, "/fizzbuzz", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	app.fizzbuzzHandler(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	expected := `{"result":["1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz"]}` + "\n"
	if string(body) != expected {
		t.Errorf("expected response body:\n%s\ngot:\n%s", expected, string(body))
	}
}

func TestFizzbuzz(t *testing.T) {
	tests := []struct {
		name string
		fizzbuzzParams

		expRes []string
	}{
		{
			name: "standard input",
			fizzbuzzParams: fizzbuzzParams{
				Int1:  2,
				Int2:  3,
				Limit: 10,
				Str1:  "two",
				Str2:  "three",
			},
			expRes: []string{"1", "two", "three", "two", "5", "twothree", "7", "two", "three", "two"},
		},
		{
			name: "empty str1",
			fizzbuzzParams: fizzbuzzParams{
				Int1:  2,
				Int2:  3,
				Limit: 10,
				Str1:  "",
				Str2:  "three",
			},
			expRes: []string{"1", "", "three", "", "5", "three", "7", "", "three", ""},
		},
		{
			name: "limit = 1",
			fizzbuzzParams: fizzbuzzParams{
				Int1:  2,
				Int2:  3,
				Limit: 1,
				Str1:  "one",
				Str2:  "three",
			},
			expRes: []string{"1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := fizzbuzz(&tt.fizzbuzzParams)
			if !reflect.DeepEqual(tt.expRes, res) {
				t.Errorf("expected %s, but got, %s", tt.expRes, res)
			}
		})
	}
}

func TestParseQueryInt(t *testing.T) {
	tests := []struct {
		name string
		key  string
		qp   *url.Values

		expRes   int
		expError error
		expPanic bool
		isNumErr bool
	}{
		{
			name: "happy path",
			key:  "int1",
			qp: &url.Values{
				"int1": {"3"},
			},
			expRes:   3,
			expError: nil,
			expPanic: false,
		},
		{
			name:     "key is missing from query parameters",
			key:      "int1",
			qp:       &url.Values{},
			expRes:   3,
			expError: queryParamNotFoundError("int1"),
			expPanic: false,
		},
		{
			name: "value is not an integer",
			key:  "int1",
			qp: &url.Values{
				"int1": {"a"},
			},
			expRes: 0,
			expError: &strconv.NumError{
				Func: "Atoi",
				Num:  "a",
				Err:  strconv.ErrSyntax,
			},
			isNumErr: true,
			expPanic: false,
		},
		{
			name: "int1 is negative",
			key:  "int1",
			qp: &url.Values{
				"int1": {"-3"},
			},
			expRes:   -3,
			expError: nil,
			expPanic: false,
		},
		{
			name: "int1 is zero",
			key:  "int1",
			qp: &url.Values{
				"int1": {"0"},
			},
			expRes:   0,
			expError: invalidParamError("int1"),
			expPanic: false,
		},
		{
			name: "limit is not positive",
			key:  "limit",
			qp: &url.Values{
				"limit": {"-10"},
			},
			expRes:   0,
			expError: invalidParamError("limit"),
			expPanic: false,
		},
		{
			name:     "key is empty",
			key:      "",
			qp:       &url.Values{},
			expRes:   0,
			expError: nil,
			expPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expPanic {
				defer func() {
					if pv := recover(); pv == nil {
						t.Fatalf("expected test to panic, but got no panic")
					}
				}()
			}
			res, err := parseQueryInt(tt.qp, tt.key)

			if tt.expError != nil {
				if err == nil {
					t.Errorf("expected error %v, but got nil", err)
				}

				if tt.isNumErr {
					var numErr *strconv.NumError
					if !errors.As(err, &numErr) {
						t.Fatalf("expected *strconv.NumError, but got %T (%v)", err, err)
					}
				} else if !errors.Is(err, tt.expError) {
					t.Errorf("expected error:%v,\n but got %v\n", tt.expError, err)
				}
			}

			if err == nil {
				if res != tt.expRes {
					t.Errorf("expected result=%d, but got %d", tt.expRes, res)
				}
			}
		})
	}
}

func TestParseFizzbuzzReq(t *testing.T) {
	tests := []struct {
		name     string
		queryStr string

		expRes   *fizzbuzzParams
		expErr   error
		isNumErr bool
	}{
		{
			name:     "happy path",
			queryStr: "int1=2&int2=3&limit=10&str1=two&str2=three",

			expRes: &fizzbuzzParams{
				Int1:  2,
				Int2:  3,
				Limit: 10,
				Str1:  "two",
				Str2:  "three",
			},
			expErr: nil,
		},
		{
			name:     "value is not an integer",
			queryStr: "int1=a&int2=3&limit=10&str1=fizz&str2=buzz",

			expRes: nil,
			expErr: &strconv.NumError{
				Func: "Atoi",
				Num:  "a",
				Err:  strconv.ErrSyntax,
			},
			isNumErr: true,
		},
		{
			name:     "int1 is negative",
			queryStr: "int1=-2&int2=3&limit=10&str1=fizz&str2=buzz",

			expRes: &fizzbuzzParams{
				Int1:  -2,
				Int2:  3,
				Limit: 10,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expErr: nil,
		},
		{
			name:     "int1 is zero",
			queryStr: "int1=0&int2=3&limit=10&str1=fizz&str2=buzz",

			expRes: nil,
			expErr: invalidParamError("int1"),
		},
		{
			name:     "limit is not positive",
			queryStr: "int1=2&int2=3&limit=-10&str1=fizz&str2=buzz",

			expRes: &fizzbuzzParams{
				Int1:  0,
				Int2:  1,
				Limit: 1,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			expErr: invalidParamError("limit"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/fizzbuzz?"+tt.queryStr, nil)
			res, err := parseFizzbuzzReq(req)

			if tt.expErr == nil {
				if err != nil {
					t.Errorf("expected no error, but got %v", err)
				}
				if !reflect.DeepEqual(*tt.expRes, *res) {
					t.Errorf("expected %v\n, but got %v\n", *tt.expRes, *res)
				}
			}

			if tt.expErr != nil {
				if err == nil {
					t.Errorf("expected error, but got no error")
				}
				if tt.isNumErr {
					var numErr *strconv.NumError
					if !errors.As(err, &numErr) {
						t.Fatalf("expected *strconv.NumError, but got %T (%v)", err, err)
					}
				} else if !errors.Is(err, tt.expErr) {
					t.Errorf("expected error:%v,\n but got %v\n", tt.expErr, err)
				}
			}
		})
	}
}
