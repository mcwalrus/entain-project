# Sports Service

A gRPC microservice for managing sport events data with SQLite storage.

## Quick Start
```bash
go run .
```
Server runs on `localhost:9001` by default.

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

### ListEvents

List all:

```bash
$ grpcurl -plaintext -emit-defaults -proto proto/sports/sports.proto -d '{}' localhost:9001 sports.Sports/ListEvents
```

Filter by sport-type:

```bash
$ grpcurl -plaintext -emit-defaults -proto proto/sports/sports.proto -d '{"filter": {"sport_types": ["Hockey"]}}' localhost:9001 sports.Sports/ListEvents
```

Filter by league:

```bash
$ grpcurl -plaintext -emit-defaults -proto proto/sports/sports.proto -d '{"filter": {"sport_types": ["Hockey"], "leagues": ["Armstrong-Williamson League"]}}' localhost:9001 sports.Sports/ListEvents
```

Visible only:

```bash
$ grpcurl -plaintext -emit-defaults -proto proto/sports/sports.proto -d '{"filter": {"visible_only": true}}' localhost:9001 sports.Sports/ListEvents
```

