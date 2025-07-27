# FizzBuzz REST API

This is a REST server that implements an extended version of the classic **FizzBuzz** problem, designed for a technical assessment.

The classic FizzBuzz task is to print numbers from 1 to 100, replacing:
- multiples of 3 with `"fizz"`,
- multiples of 5 with `"buzz"`,
- multiples of both 3 and 5 with `"fizzbuzz"`.

This version exposes a flexible and parameterized REST API that lets you customize the numbers and the replacement strings.

## Tech Stack
- Go >= 1.24
- Docker
- GitHub Actions

## Project Structure
```bash
├── cmd/ # application entry point
├── internal/ # application core logic
├── deploy/ # configuration files required for production deployment
├── Makefile 
└── .github/workflows/ci.yml # CI Pipeline 
```

## API Endpoints

### `GET /fizzbuzz`

Generates a FizzBuzz sequence using the parameters you provide.

#### Query Parameters:
- `int1` (integer): First multiple
- `int2` (integer): Second multiple
- `limit` (integer): Upper bound of the sequence (inclusive)
- `str1` (string): Replacement for `int1` multiples
- `str2` (string): Replacement for `int2` multiples

**Example**:
```shell 
GET /fizzbuzz?int1=2&int2=3&limit=5&str1=fizz&str2=buzz
```
**Response**
```json
{"fizzbuzz":["1","fizz","buzz","fizz","5","fizzbuzz","7","fizz","buzz","fizz"]}
```

### `GET /healthcheck`

Returns the current running state of the server.

```shell 
GET /healthcheck
```
**Response**
```json
{"healthcheck":{"status":"available","env":"development","version":"(devel)"}}
```

### `GET /stats`

Returns the parameters corresponding to the most used request, as well as the number of hits for this request

```shell 
GET /stats
```
**Response**
```json
{"stats":{"parameters":{"int1":2,"int2":3,"limit":5,"str1":"fizz","str2":"buzz"},"hits":4}}
```

## Local Setup

### Running locally, with make
```shell
$ make run
```

### Running locally, with docker (with live-reload)
```shell
$ make dev
```

## Developer Commands
Run `make` (or `make help`) to display available commands and their description:
```shell
$ make 
```

## Continuous Integration

This project uses **GitHub Actions** for continuous integration.

### Pull Request Checks

All code changes must go through a Pull Request into the `main` branch.  
CI will:

- Build the application
- Run unit tests with race detection
- Perform static analysis and dependency audits

> **Tests must pass before a PR can be merged** (enforced via GitHub branch protection rules).


CI configuration lives in [`/.github/workflows/ci.yml`](.github/workflows/ci.yml).

## Assumptions
- If str1 or str2 parameters are empty, the query is valid, and they are taken as empty strings `""`
- int1 and int2 can be positive or negative integers, and are mandatory query parameters
- stats endpoint statistics are ephemeral - they are not persisted in a database and are lost if the application goes down

## TODOs
- version - (devel)
- add flags to setup instructions
