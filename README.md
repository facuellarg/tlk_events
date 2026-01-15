# Events API

A RESTful API for managing events built with Go, Echo framework, and SQLite database.

## Features

- Create, read, update, and delete events
- SQLite database for persistent storage
- UUID-based event identification
- Input validation
- RESTful API design
- JSON request/response format

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Echo v4
- **Database**: SQLite3
- **Libraries**:
  - `github.com/labstack/echo/v4` - HTTP framework
  - `github.com/mattn/go-sqlite3` - SQLite driver
  - `github.com/google/uuid` - UUID generation

## Project Structure

```
.
├── repository/
│   └── repository.go       # Database operations and models
├── models/
│   └── dto.go             # Dto definition for request
│   └── event.go           # Event model definition
├── utils/
│   └── utils.go           # Utility functions
├── service/
│   └── events.go          # Server setup and routing        
└── main.go                # Application entry point
```

## Installation

### Prerequisites

- Go 1.21 or higher
- GCC compiler (required for SQLite CGO)

### Setup

1. Clone the repository:
```bash
git clone https://github.com/facuellarg/tlk_events
cd tlk_events
```

2. Install dependencies:
```bash
go mod init challenge
go mod tidy
go get github.com/labstack/echo/v4
go get github.com/mattn/go-sqlite3
go get github.com/google/uuid
```

3. Build the application:
```bash
go build -o events-api
```

## Configuration

The application uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_PATH` | Path to SQLite database file | `./events.db` |
| `PORT` | Server port | `8080` |

## Running the Application

### Development

```bash
# Using default settings
go run .

# With custom configuration
export DB_PATH="./data/events.db"
export PORT="3000"
go run .
```

### Production

```bash
# Build binary
go build -o events-api

# Run
./events-api
```

The server will start on `http://localhost:8080` (or your configured port).

## API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Event Model

```json
{
  "id": "uuid",
  "title": "string (max 100 characters)",
  "description": "string (optional)",
  "start_time": "ISO 8601 timestamp",
  "end_time": "ISO 8601 timestamp",
  "created_at": "ISO 8601 timestamp"
}
```

## API Documentation

### 1. Create Event

Create a new event.

**Endpoint**: `POST /api/v1/events`

**Request Body**:
```json
{
  "title": "Team Meeting",
  "description": "Quarterly planning session",
  "start_time": "2026-01-20T10:00:00Z",
  "end_time": "2026-01-20T11:00:00Z"
}
```

**Response**: `201 Created`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Team Meeting",
  "description": "Quarterly planning session",
  "start_time": "2026-01-20T10:00:00Z",
  "end_time": "2026-01-20T11:00:00Z",
  "created_at": "2026-01-15T14:30:00Z"
}
```

**Validation Rules**:
- `title`: Required, non-empty, max 100 characters
- `start_time`: Required, must be before `end_time`
- `end_time`: Required
- `description`: Optional

**Error Responses**:
- `400 Bad Request`: Invalid input or validation error
- `500 Internal Server Error`: Database error

---

### 2. Get All Events

Retrieve all events ordered by start time (ascending).

**Endpoint**: `GET /api/v1/events`

**Response**: `200 OK`
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "Team Meeting",
    "description": "Quarterly planning session",
    "start_time": "2026-01-20T10:00:00Z",
    "end_time": "2026-01-20T11:00:00Z",
    "created_at": "2026-01-15T14:30:00Z"
  },
  {
    "id": "987fcdeb-51a2-43f7-b123-456789abcdef",
    "title": "Project Review",
    "description": null,
    "start_time": "2026-01-21T14:00:00Z",
    "end_time": "2026-01-21T15:30:00Z",
    "created_at": "2026-01-15T15:00:00Z"
  }
]
```

**Error Responses**:
- `500 Internal Server Error`: Database error

---

### 3. Get Event by ID

Retrieve a specific event by its UUID.

**Endpoint**: `GET /api/v1/events/:id`

**Path Parameters**:
- `id`: UUID of the event

**Response**: `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "Team Meeting",
  "description": "Quarterly planning session",
  "start_time": "2026-01-20T10:00:00Z",
  "end_time": "2026-01-20T11:00:00Z",
  "created_at": "2026-01-15T14:30:00Z"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid UUID format
- `404 Not Found`: Event not found
- `500 Internal Server Error`: Database error

---

## cURL Examples

### Create a new event

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Team Meeting",
    "description": "Quarterly planning session",
    "start_time": "2026-01-20T10:00:00Z",
    "end_time": "2026-01-20T11:00:00Z"
  }'
```

### Create an event without description

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Quick Standup",
    "start_time": "2026-01-16T09:00:00Z",
    "end_time": "2026-01-16T09:15:00Z"
  }'
```

### Get all events

```bash
curl http://localhost:8080/api/v1/events
```

### Get all events (formatted output)

```bash
curl http://localhost:8080/api/v1/events | jq
```

### Get a specific event by ID

```bash
curl http://localhost:8080/api/v1/events/123e4567-e89b-12d3-a456-426614174000
```

### Get event with formatted output

```bash
curl http://localhost:8080/api/v1/events/123e4567-e89b-12d3-a456-426614174000 | jq
```

### Create event and save response to file

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Client Presentation",
    "description": "Q4 results presentation",
    "start_time": "2026-01-25T14:00:00Z",
    "end_time": "2026-01-25T15:00:00Z"
  }' \
  -o response.json
```

### Show full HTTP response (including headers)

```bash
curl -i http://localhost:8080/api/v1/events
```

### Verbose output for debugging

```bash
curl -v -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Debug Event",
    "start_time": "2026-01-16T10:00:00Z",
    "end_time": "2026-01-16T11:00:00Z"
  }'
```

## Testing Validation

### Test: Empty title (should fail)

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "",
    "start_time": "2026-01-20T10:00:00Z",
    "end_time": "2026-01-20T11:00:00Z"
  }'
```

**Expected Response**: `400 Bad Request`
```json
{
  "error": "title cannot be empty"
}
```

### Test: Title too long (should fail)

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "This is a very long title that exceeds the maximum allowed length of one hundred characters for event titles",
    "start_time": "2026-01-20T10:00:00Z",
    "end_time": "2026-01-20T11:00:00Z"
  }'
```

**Expected Response**: `400 Bad Request`
```json
{
  "error": "title must be 100 characters or less"
}
```

### Test: End time before start time (should fail)

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Invalid Event",
    "start_time": "2026-01-20T11:00:00Z",
    "end_time": "2026-01-20T10:00:00Z"
  }'
```

**Expected Response**: `400 Bad Request`
```json
{
  "error": "start_time must be before end_time"
}
```

### Test: Invalid UUID (should fail)

```bash
curl http://localhost:8080/api/v1/events/invalid-uuid
```

**Expected Response**: `400 Bad Request`
```json
{
  "error": "Invalid UUID format"
}
```

### Test: Event not found (should fail)

```bash
curl http://localhost:8080/api/v1/events/00000000-0000-0000-0000-000000000000
```

**Expected Response**: `404 Not Found`
```json
{
  "error": "Event not found"
}
```

## Database

The application uses SQLite with the following schema:

```sql
CREATE TABLE events (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL CHECK(length(title) <= 100),
    description TEXT,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_events_start_time ON events(start_time);
CREATE INDEX idx_events_end_time ON events(end_time);
```

### Database Management

View the database:
```bash
sqlite3 events.db
```

Query events:
```sql
SELECT * FROM events ORDER BY start_time;
```

Delete all events:
```sql
DELETE FROM events;
```

## Development

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
golangci-lint run
```

## Troubleshooting

### SQLite CGO Error

If you get a CGO-related error when building, ensure you have GCC installed:

**macOS**:
```bash
xcode-select --install
```

**Linux (Ubuntu/Debian)**:
```bash
sudo apt-get install build-essential
```

**Windows**:
Install MinGW-w64 or TDM-GCC

### Port Already in Use

If port 8080 is already in use:
```bash
export PORT=3000
go run .
```

### Database Locked Error

SQLite can have locking issues with concurrent writes. The application is configured with WAL mode to minimize this, but if you encounter issues:
1. Ensure only one instance is running
2. Check file permissions on the database file
3. Close any open SQLite connections

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions, please open an issue on the GitHub repository.
