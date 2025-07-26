package main

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type queryParamNotFoundError string

func (e queryParamNotFoundError) Error() string {
	return strconv.Quote(string(e)) + " query parameter not found in request"
}

type invalidParamError string

func (e invalidParamError) Error() string {
	val := strconv.Quote(string(e))
	var errStr string
	switch e {
	case "limit":
		errStr = " parameter has to be a positive integer"
	case "int1", "int2":
		errStr = val + " parameter cannot be zero"
	}
	return errStr
}

type fizzbuzzParams struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

func (app *application) fizzbuzzHandler(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(paramsCtxKey).(*fizzbuzzParams)

	res := fizzbuzz(params)

	err := app.toJSONResponse(w, http.StatusOK, &wrapper{"result": res}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func fizzbuzz(params *fizzbuzzParams) []string {
	var res = make([]string, 0, params.Limit)
	for i := 1; i <= params.Limit; i++ {
		var added bool
		var strBuilder strings.Builder
		if i%params.Int1 == 0 {
			strBuilder.WriteString(params.Str1)
			added = true
		}
		if i%params.Int2 == 0 {
			strBuilder.WriteString(params.Str2)
			added = true
		}
		if !added {
			res = append(res, strconv.Itoa(i))
		} else {
			res = append(res, strBuilder.String())
		}
	}
	return res
}

// parses and validates client request input parameters for
// fizzbuzz handler. It assumes that empty str1 and str2
// are valid use cases and, if omitted, are taken as empty.
func parseFizzbuzzReq(r *http.Request) (*fizzbuzzParams, error) {
	if r.URL.RawQuery == "" {
		return nil, errors.New("query string is not present")
	}

	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	int1, err := parseQueryInt(&queryParams, "int1")
	if err != nil {
		return nil, err
	}
	int2, err := parseQueryInt(&queryParams, "int2")
	if err != nil {
		return nil, err
	}
	limit, err := parseQueryInt(&queryParams, "limit")
	if err != nil {
		return nil, err
	}

	str1 := queryParams.Get("str1")
	str2 := queryParams.Get("str2")

	return &fizzbuzzParams{
		Int1:  int1,
		Int2:  int2,
		Limit: limit,
		Str1:  str1,
		Str2:  str2,
	}, nil
}

func parseQueryInt(qp *url.Values, key string) (int, error) {
	if key == "" {
		panic("key must be defined")
	}

	val := qp.Get(key)
	if val == "" {
		return 0, queryParamNotFoundError(key)
	}

	res, err := strconv.Atoi(val)
	if err != nil {
		return res, err
	}

	if key == "limit" && res < 1 || res == 0 {
		return 0, invalidParamError(key)
	}

	return res, nil
}
