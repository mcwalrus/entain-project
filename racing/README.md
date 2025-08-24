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

## Protobuf

On updating `racing.proto`, you will need to regenerate the bindings:

```bash
$ go mod tidy
$ go generate ./...
```


## Testing

Using `grpcurl` to test the endpoints manually. To install:

```bash
brew install grpcurl
```

To perform unit tests:

```bash
go test -run . git.neds.sh/matty/entain/racing/db
```

### GetRace

By race id:

```bash
$ grpcurl -plaintext -emit-defaults -proto proto/racing/racing.proto -d '{"id": 1}' localhost:9000 racing.Racing/GetRace
```

### ListRaces

By meeting ids:

```bash
$ grpcurl -plaintext -emit-defaults -proto proto/racing/racing.proto -d '{"filter": {"meeting_ids": [1, 2, 3]}}' localhost:9000 racing.Racing/ListRaces
```

By meeting ids, visible only:

```bash
$ grpcurl -plaintext -emit-defaults -proto proto/racing/racing.proto -d '{"filter": {"meeting_ids": [1, 2, 3], "visible_only": true}}' localhost:9000 racing.Racing/ListRaces
```

Default sorted order should be by ascending `advertised_start_time`:

```bash
$ grpcurl -plaintext -emit-defaults -proto proto/racing/racing.proto -d '{}' localhost:9000 racing.Racing/ListRaces
```

With sorting by advertised start time descending:

```bash
$ grpcurl -plaintext -emit-defaults -proto proto/racing/racing.proto -d '{"filter": {"sort_by": "ADVERTISED_START_TIME_DESC"}}' localhost:9000 racing.Racing/ListRaces
```

With sorting by name ascending:

```bash
$ grpcurl -plaintext -emit-defaults -proto proto/racing/racing.proto -d '{"filter": {"sort_by": "NAME_ASC"}}' localhost:9000 racing.Racing/ListRaces
```

**Manual validation**

Verify different sets of results with counting returned entries from races list:

```
 $ grpcurl ... | jq '.races | length'
```

