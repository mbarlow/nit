# nit

A minimal Go service that provides dynamic CRUD operations for any data type without requiring predefined structs.

## Features

- **Zero configuration** - tables are created automatically
- **Pure JSON storage** - leverages SQLite's JSON1 extension
- **Single endpoint** - handles all CRUD operations
- **UUID identifiers** - uses Google UUIDs for record IDs
- **Built-in pagination** - automatic pagination support with limit/offset
- **Minimal code** - ~150 lines total

## Quick Start

```bash
make dev
```

## API Usage

The API follows a simple pattern: `/{table}` and `/{table}/{id}`

### Create Record
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "age": 30}'

# Response:
{"id": "550e8400-e29b-41d4-a716-446655440000"}
```

### Get All Records (Paginated)
```bash
# Default pagination (limit=10, offset=0)
curl http://localhost:8080/users

# With pagination parameters
curl "http://localhost:8080/users?limit=5&offset=10"

# Response:
{
  "items": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "data": {
        "name": "John Doe",
        "email": "john@example.com",
        "age": 30
      },
      "created": "2025-01-15T10:30:00Z",
      "updated": "2025-01-15T10:30:00Z"
    }
  ],
  "total_items": 42,
  "limit": 10,
  "offset": 0,
  "has_more": true
}
```

### Get Single Record
```bash
curl http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000
```

### Update Record
```bash
curl -X PUT http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe", "email": "jane@example.com", "age": 25}'
```

### Delete Record
```bash
curl -X DELETE http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000
```

## Response Format

All records include metadata:

```json
{
  "id": 1,
  "data": {
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30
  },
  "created": "2025-01-15 10:30:00",
  "updated": "2025-01-15 10:30:00"
}
```

## Architecture

- **Tables**: Auto-created with `id`, `data` (JSON), `created`, `updated`
- **Storage**: SQLite with JSON1 extension for native JSON operations
- **Validation**: None (pure schemaless approach)
- **Performance**: Indexes can be added manually to SQLite as needed

## Dependencies

- [Echo v4](https://echo.labstack.com/) - Web framework
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver

## File Structure

```
.
├── main.go         # Complete application (150 lines)
├── go.mod          # Dependencies
├── Makefile        # Development commands
└── README.md       # This file
```

## Make Commands

- `make run` - Start the server
- `make build` - Build binary
- `make test` - Run example API calls
- `make clean` - Remove binary and database
- `make dev` - Install deps and run

## Extending

To add validation, authentication, or other features:

1. Add middleware to Echo
2. Implement JSON Schema validation in the handlers
3. Add database constraints or triggers

The minimal approach makes it easy to extend exactly what you need.
