# Banka-4-Backend

# Usage

This project uses a `Makefile` for common tasks.

## Docker

```bash
make docker-up-build    # Build and start services
make docker-up          # Start services
make docker-down        # Stop services
make docker-down-rm-vol # Stop services and remove volumes
```

## Swagger

```bash
make swagger-docs
```

## Protobuf

```bash
make proto
```

## Testing

```bash
make test               # Run unit tests
make test-integration   # Run integration tests
```

## Coverage

```bash
make coverage-profile   # Generate coverage data
make coverage           # Show total coverage
make coverage-report    # Coverage by layer
make coverage-html      # HTML report
```

