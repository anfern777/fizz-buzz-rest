# FizzBuzz REST API

This is a REST server that implements an extended version of the classic **FizzBuzz** problem, designed for a technical assessment.

## Problem Description

The classic FizzBuzz task is to print numbers from 1 to 100, replacing:
- multiples of 3 with `"fizz"`,
- multiples of 5 with `"buzz"`,
- multiples of both 3 and 5 with `"fizzbuzz"`.

This version exposes a flexible and parameterized REST API that lets you customize the numbers and the replacement strings.


## API Endpoints

### `GET /fizzbuzz`

Generates a FizzBuzz sequence using the parameters you provide.

#### Query Parameters:
- `int1` (integer): First multiple
- `int2` (integer): Second multiple
- `limit` (integer): Upper bound of the sequence (inclusive)
- `str1` (string): Replacement for `int1` multiples
- `str2` (string): Replacement for `int2` multiples

### `GET /healcheck`

Returns the current running state of the server.

### `GET /stats`

Returns the parameters corresponding to the most used request, as well as the number of hits for this request
