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
	int1, int2, limit int
	str1, str2        string
}

func (app *application) fizzbuzzHandler(w http.ResponseWriter, r *http.Request) {
	params, err := parseFizzbuzzReq(r)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	res := fizzbuzz(params)

	app.toJSONResponse(w, http.StatusOK, &wrapper{"result": res}, nil)
}

func fizzbuzz(params *fizzbuzzParams) []string {
	var res = make([]string, 0, params.limit)
	for i := 1; i <= params.limit; i++ {
		var added bool
		var strBuilder strings.Builder
		if i%params.int1 == 0 {
			strBuilder.WriteString(params.str1)
			added = true
		}
		if i%params.int2 == 0 {
			strBuilder.WriteString(params.str2)
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
		int1:  int1,
		int2:  int2,
		limit: limit,
		str1:  str1,
		str2:  str2,
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
