# Racing Service

A gRPC microservice for managing racing data with SQLite storage.

## Features
- **ListRaces API**: Retrieve races with optional filtering by meeting IDs
- **SQLite Database**: Local storage with seeded dummy data
- **Protocol Buffers**: Type-safe API definitions

## Quick Start
```bash
go run .
```
Server runs on `localhost:9000` by default.

## Architecture
- `service/`: Business logic and gRPC handlers
- `db/`: Database repository layer with SQLite
- `proto/`: Protocol buffer definitions and generated code

## Testing

Using `grpcurl` to test the endpoints manually. To install:

```bash
brew install grpcurl
```


**ListRaces:**

By meeting ids:

```bash
$ grpcurl -plaintext -proto racing/proto/racing/racing.proto -d '{"filter": {"meeting_ids": [1, 2, 3]}}' localhost:9000 racing.Racing/ListRaces
```

By meeting ids, visible only:

```bash
$ grpcurl -plaintext -proto racing/proto/racing/racing.proto -d '{"filter": {"meeting_ids": [1, 2, 3], "visible_only": true}}' localhost:9000 racing.Racing/ListRaces
```