# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based AWS Lambda application for the SGA (Sistema de Gestión de Almacenes) generator. The application provides API endpoints to retrieve and manage product/element tags and descriptions.

## Architecture

The codebase follows a clean architecture pattern with dependency injection:

```
Handler → Service → Repository
```

- **Handler Layer** (`internal/tags/handler/`): AWS Lambda handlers that process events and return responses. Depends on service interfaces.
- **Service Layer** (`internal/tags/service/`): Business logic layer. Depends on repository interfaces.
- **Repository Layer** (`internal/tags/repository/`): Data access layer. Contains hardcoded element mappings in `fillElementsDic()`.
- **Models** (`internal/models/`): Shared data structures used across layers.

### Dependency Injection Pattern

Each layer defines interfaces for its dependencies and receives them via constructors:
- `handler.New(service TagsService)` - Handler depends on TagsService interface
- `service.New(repository TagRepository)` - Service depends on TagRepository interface
- Wire-up happens in `cmd/tags/getTags/main.go`

## Development Commands

### Build
```bash
# Build getTags Lambda
go build -o bin/getTags cmd/tags/getTags/main.go

# Build getTagById Lambda
go build -o bin/getTagById cmd/tags/getTagById/main.go
```

### Run locally
```bash
# Run getTags Lambda
go run cmd/tags/getTags/main.go

# Run getTagById Lambda
go run cmd/tags/getTagById/main.go
```

### Test
```bash
go test ./...
```

### Test specific package
```bash
go test ./internal/tags/handler
```

### Install dependencies
```bash
go mod download
```

### Update dependencies
```bash
go mod tidy
```

## Lambda Functions

Each Lambda function is in its own subdirectory under `cmd/tags/`:
- `getTags/main.go` - Retrieves all element tags from the repository
- `getTagById/main.go` - Retrieves a specific tag by its ID (code)

## Adding New Features

### Adding a new Lambda function
1. Create new directory under `cmd/tags/[functionName]/`
2. Create `main.go` with handler initialization following the DI pattern
3. Wire dependencies: Repository → Service → Handler
4. Register with `lambda.Start(h.HandlerMethod)`

### Adding new elements to repository
Modify `fillElementsDic()` in `internal/tags/repository/repository.go`

### Adding new handler methods
1. Add method to `Handler` struct in `internal/tags/handler/tags.go`
2. Add corresponding interface method to `TagsService` if needed
3. Implement service method in `internal/tags/service/service.go`
4. Add corresponding repository interface method if needed
5. Implement repository method in `internal/tags/repository/repository.go`

## Code Conventions

- JSON struct tags use Spanish field names (`codigo`, `descripcion`)
- Error handling: Return errors up the stack; handler layer returns Lambda-compatible errors
- Pointer returns: Functions return pointers to slices (`*[]models.Element`)
