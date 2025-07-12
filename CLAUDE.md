# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**github.com/bborbe/kv** is a Go library providing unified interfaces for key-value stores. It abstracts different KV database implementations (BadgerDB, BoltDB, in-memory) behind common interfaces, enabling easy switching between backends without code changes.

## Development Commands

### Primary Development Workflow
```bash
make precommit  # Full pre-commit pipeline: ensure format generate test check addlicense
```

### Individual Commands
```bash
# Dependencies
make ensure     # go mod tidy, verify, remove vendor

# Code formatting
make format     # gofmt + goimports-reviser with project name

# Code generation (mocks)
make generate   # go generate (creates counterfeiter mocks)

# Testing
make test       # Tests with race detection and coverage
go test -run TestSpecificTest  # Run specific test

# Code quality
make check      # vet + errcheck + vulncheck
make vet        # go vet
make errcheck   # Check unchecked errors (ignores Close/Write/Fprint)
make vulncheck  # Security vulnerability check

# Licensing
make addlicense # Add BSD license headers to all Go files
```

## Core Architecture

### Interface Hierarchy
The library uses a layered interface approach:

1. **DB Interface** (`kv_db.go`) - Database lifecycle management
   - `Update()` - Write transactions
   - `View()` - Read-only transactions
   - `Sync()`, `Close()`, `Remove()` - Database management

2. **Transaction Interface** (`kv_tx.go`) - Bucket operations within transactions
   - `Bucket()`, `CreateBucket()`, `DeleteBucket()` - Bucket management
   - `ListBucketNames()` - Bucket enumeration

3. **Bucket Interface** (`kv_bucket.go`) - Key-value operations
   - `Put()`, `Get()`, `Delete()` - Basic operations
   - `Iterator()`, `IteratorReverse()` - Iteration support

4. **Generic Store Layer** (`kv_store.go`) - Type-safe operations
   - Uses Go generics: `Store[KEY ~[]byte | ~string, OBJECT any]`
   - Automatic JSON marshaling/unmarshaling
   - Composed interfaces: `StoreAdder`, `StoreGetter`, `StoreRemover`, etc.

### Advanced Components

- **Relation Store** (`kv_relation-store.go`) - Bidirectional 1:N relationships
- **Metrics Wrapper** (`kv_db-metrics.go`) - Prometheus monitoring
- **Provider Pattern** (`kv_provider.go`) - Database factory interface
- **Benchmarking** (`benchmark/`) - Performance testing framework

## Testing Framework

### Test Structure
- **Ginkgo/Gomega BDD** - Behavior-driven testing framework
- **Reusable Test Suites** - `BasicTestSuite()`, `BucketTestSuite()`, etc.
- **Provider Pattern** - Tests accept `Provider` interface for different implementations

### Mock Generation
All interfaces use counterfeiter for mock generation:
```go
//counterfeiter:generate -o mocks/db.go --fake-name DB . DB
```
Run `make generate` to regenerate mocks.

### Test Execution
- Tests run with race detection (`-race`) and coverage reporting
- Parallel execution controlled by `GO_TEST_PARALLEL` environment variable
- Use `go test -run TestName` for specific tests

## Key Dependencies

### Benjamin's Ecosystem
- `github.com/bborbe/errors` - Error wrapping and handling
- `github.com/bborbe/collection` - Use `collection.Ptr()` for pointer utilities
- `github.com/bborbe/time` - Time handling with dependency injection patterns

### External Libraries
- `github.com/onsi/ginkgo/v2` & `github.com/onsi/gomega` - BDD testing
- `github.com/prometheus/client_golang` - Metrics collection
- `github.com/gorilla/mux` - HTTP routing (for benchmarking server)

## Code Patterns

### Generic Type Safety
The library extensively uses Go 1.18+ generics for type-safe operations:
```go
Store[KEY ~[]byte | ~string, OBJECT any]
```

### Error Handling
All errors are wrapped using `github.com/bborbe/errors` for consistent error handling and stack traces.

### Context Usage
All operations accept and properly propagate `context.Context` for cancellation and deadlines.

### Interface Segregation
Small, focused interfaces that can be composed together rather than large monolithic interfaces.

## Implementation Variants

This library defines interfaces implemented by three concrete packages:
- `github.com/bborbe/badgerkv` - BadgerDB implementation
- `github.com/bborbe/boltkv` - BoltDB implementation  
- `github.com/bborbe/memorykv` - In-memory implementation