# API

Provides api service to other downstream services, i.e sports-service, racing-service.

## Protobuf

To regenerate the protobuf bindings:

```bash
go mod tidy
go generate ./...
```


## Manual Testing

### v1/list-races

All races:

```
curl -X POST http://localhost:8000/v1/list-races \
  -H "Content-Type: application/json" \
  -d '{}' | jq
```

### v1/list-events

All sports events:

```
curl -X POST http://localhost:8000/v1/list-events \
  -H "Content-Type: application/json" \
  -d '{}' | jq
```






