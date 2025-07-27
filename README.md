# FizzBuzz REST API

This is a REST server that implements an extended version of the classic **FizzBuzz** problem, designed for a technical assessment.

The classic FizzBuzz task is to print numbers from 1 to 100, replacing multiples of 3 with "fizz", multiples of 5 with "buzz", and multiples of both with "fizzbuzz". This version exposes a flexible and parameterized REST API that lets you customize this behavior.

This version exposes a flexible and parameterized REST API that lets you customize the numbers and the replacement strings.

## Tech Stack
- Go: 1.24+
- Docker
- GitHub Actions

## Project Structure
The project follows the standard Go project layout.
```bash
├── cmd/ # Application entry point
├── internal/ # Application core logic
├── deploy/ # Configuration files required for production deployment
├── Makefile # Helper commands for development
└── .github/workflows/ci.yml # Build/Test/Validation Workflow
```

## API Endpoints

### `GET /fizzbuzz`

Generates a custom FizzBuzz sequence based on the provided parameters.

#### Query Parameters:
- `int1` (integer): First multiple
- `int2` (integer): Second multiple
- `limit` (integer): Upper bound of the sequence
- `str1` (string): Replacement string for `int1` multiples
- `str2` (string): Replacement string for `int2` multiples

**Example**:
```shell 
GET /fizzbuzz?int1=2&int2=3&limit=5&str1=fizz&str2=buzz
```
**Response**
```json
{
  "result": [
    "1",
    "fizz",
    "buzz",
    "fizz",
    "5",
    "fizzbuzz",
    "7",
    "fizz",
    "buzz",
    "fizz"
  ]
}
```

### `GET /stats`

Returns the parameters corresponding to the most used request, as well as the number of hits for this request

```shell 
GET /stats
```
**Response**
```json
{
  "healthcheck": {
    "status": "available",
    "environment": "development",
    "version": "1.0.0"
  }
}
```

### `GET /healthcheck`

Returns the current running state of the server.

```shell 
GET /healthcheck
```
**Response**
```json
{
  "healthcheck": {
    "status": "available",
    "environment": "development",
    "version": "1.0.0"
  }
}
```

## Configuration
You can configure the application using command-line flags at runtime.

| Flag | Description | Default |
| ---  |     ---     |  ----   |
|-port	| The port the server listens on.| 8080 |
| -env	| The runtime environment.	| development |
| -limiter-enabled	| Enable or disable the rate limiter.	| true |
| -limiter-rate	| Max number of requests per second for the rate limiter.|	2 |
| -limiter-burst |	The burst allowance for the rate limiter.	| 4 |
| -trusted-origins	| A space-separated list of trusted CORS origins. |	"" (none)
| -version	| Displays the application version and exits. |	false |

## Local Setup

### Prerequisites
**Make**: A build automation tool.
**Docker**: For running the application in a containerized environment.

### Running the application
To run the server locally with live-reload enabled for development:
```shell
make dev
```
The application will be available at http://localhost:8080.

To run the application using the standard Go toolchain:
```shell
% make dev
```

## Developer Commands
Run `make help` to see a list of all available developer commands.
```shell
% fizz-buzz-rest % make                             
Usage:
  run     run the fuzzbuzz api with go cli
  dev     run containerized fizzbuzz api with live-reload
  build   build the application binaries (strips out both symbol tables and DWARF debugging information to reduce size)
  tidy    tidy module dependencies and format .go files
  audit   run tests and quality control checks
```

## Continuous Integration

This project uses GitHub Actions for CI. The workflow is triggered on every pull request to the main branch and performs the following checks:

- Builds the application.
- Runs unit tests with race detection.
- Performs static analysis and security vulnerability audits.

> Note: All checks must pass before a pull request can be merged. This is enforced by branch protection rules.

The CI configuration can be found in .github/workflows/pr-build-test.yml.

> **Tests must pass before a PR can be merged** (enforced via GitHub branch protection rules).


CI configuration lives in [`/.github/workflows/ci.yml`](.github/workflows/ci.yml).

## Key Behaviors
- If str1 or str2 parameters are submitted empty, they are treated as valid empty strings.
- int1 and int2 can be positive or negative integers, and are mandatory query parameters
- Statistics gathered by the /stats endpoint are ephemeral and will be reset if the application restarts
