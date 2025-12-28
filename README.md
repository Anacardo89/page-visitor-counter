# Page Visitor Counter

A simple GoLang web service to track and count unique visitors per page URL.

---

## Features

* Ingest visitor events via HTTP POST requests (one event per request).
* Serve the number of distinct visitors for a given page URL via HTTP GET.
* In-memory storage, fully concurrent-safe.
* Graceful shutdown on SIGINT/SIGTERM.
* Middleware for logging and request timeouts.
* Health check endpoint.

---

## Project Layout

```bash
cmd/                # Main entrypoint for the server
config/             # Config files and models
internal/           # Core application logic
  api/                # API handlers, interfaces, and tests
  repo/               # Repository and mocks
  middleware/         # Logging & timeout middleware
  server/             # HTTP server and router
pkg/                # Reusable packages
  validator/          # JSON validator
                       # - Repurposed from your own codebase
                       # - Used instead of the standard Go playground validator
                       # - Checks presence of required fields only, ignoring their actual values
```

---

## Requirements

- Go 1.25+

## External Dependencies

- [testify](https://github.com/stretchr/testify) - for unit tests
- [gomock](https://github.com/golang/mock) – for testing with mocks
- [gorilla/mux](https://github.com/gorilla/mux) – for routing HTTP requests
- [google/uuid](https://github.com/google/uuid) - for generating request IDs
- [godotenv](https://github.com/joho/godotenv) - for loading env vars
- [yaml](https://gopkg.in/yaml.v3) - for parsing Yaml (to load config)

---

## Installation

Clone the repository:

```bash
git clone git@github.com:Anacardo89/page-visitor-counter.git
cd page-visitor-counter
```

Install dependencies:

```bash
go mod tidy
```

---

## Running Tests

```bash
go test ./...
```

Core modules, repository and API, are fully covered.

---

## Configuration

A `config.yaml` example is provided in `/config`.  
For the purposes of this demo you can use the default values.  

Set environment variables:

```bash
export APP_HOME=$(pwd)
export CFG_PATH=config/config.yaml
```

---

## Running the Server

```bash
go run ./cmd/main/
```

The server will listen on the port specified in `config.yaml`. (Default port: `8080`)

---

## API Endpoints

### POST /visitors

Add a visitor event.

**Request Body**:

```json
{
  "url": "https://example.com",
  "visitor_id": "user1"
}
```

**Responses**:

* `204 No Content` - Event recorded successfully.
* `400 Bad Request` - Missing or invalid fields.

---

### GET /visitors?url=<page_url>

Retrieve the count of distinct visitors for a URL.

**Response Body**:

```json
{
  "visitor_count": 5
}
```

**Responses**:

* `200 OK` - Request successful.
* `400 Bad Request` - Missing `url` query parameter.

---

### GET / (Health Check)

**Response Body**:

```json
{
  "status": "OK"
}
```

### Catch-all

Any other path will return `404 Not Found` with JSON:

```json
{
  "error": "endpoint not found"
}
```

---

## Example cURL Requests

```bash
# Add visitor
curl -X POST -H "Content-Type: application/json" \
  -d '{"url":"https://example.com", "visitor_id":"user1"}' \
  http://localhost:8080/visitors

# Count visitors
curl http://localhost:8080/visitors?url=https://example.com
```
