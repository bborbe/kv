# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v1.18.3

- Update Go to 1.26.0

## v1.18.2

- Update Go to 1.25.7
- Update github.com/bborbe dependencies
- Update testing dependencies (ginkgo, gomega)
- Update google/osv-scanner to v2.3.2
- Add .update-logs/ and .mcp-* to .gitignore

## v1.18.1

- Update Go to 1.25.5
- Update golang.org/x/crypto to v0.47.0
- Update dependencies

## v1.18.0

- update go and deps

## v1.17.0

- Refactor error variables to follow Go ErrFoo naming convention
- Add ErrTransactionAlreadyOpen, ErrBucketNotFound, ErrBucketAlreadyExists with backward-compatible deprecation
- Rename StoreStream/StoreList interfaces to StoreStreamer/StoreLister with backward-compatible type aliases
- Update Go version from 1.25.2 to 1.25.4
- Update dependencies (github.com/bborbe/http, github.com/bborbe/log, github.com/bborbe/run, github.com/onsi/ginkgo/v2)
- Enhance golangci-lint configuration with additional linters (errname, unparam, bodyclose, forcetypeassert, asasalint, prealloc, nestif)
- Add exclusion rules for deprecated error variables in linter configuration
- Fix deprecated Go stdlib usage (io/ioutil) in depguard rules

## v1.16.1

- Add comprehensive test suite for benchmark package (26 tests, 73% coverage)
- Add tests for RandString and ShuffleSlice utility functions
- Add tests for Benchmark core functionality with mock validation
- Add tests for HTTP handler with parameter parsing
- Update github.com/bborbe/errors from v1.3.0 to v1.3.1
- Update github.com/bborbe/run from v1.7.7 to v1.8.0 (adds FuncRunner interface and mock)
- Update .gitignore patterns for coverage output files

## v1.16.0

- Add golangci-lint configuration and integration
- Update Makefile with new lint target and improved formatting
- Integrate golines for automatic line length formatting (max 100 chars)
- Update goimports-reviser to v3
- Update github.com/bborbe/http from v1.14.2 to v1.15.2
- Update github.com/onsi/ginkgo/v2 from v2.25.3 to v2.26.0
- Apply code formatting improvements across codebase
- Improve error checking exclusions in Makefile

## v1.15.3

- Update Go version to 1.25.2
- Update CI workflow to use Go 1.25.2

## v1.15.2

- Upgrade osv-scanner from v1 to v2
- Add config file support for osv-scanner in Makefile
- go mod update

## v1.15.1

- go mod update

## v1.15.0

- Add StoreList interface with List method for retrieving all objects as slice
- Add StoreListTx interface for transaction-based list operations
- Implement List methods in Store and StoreTx concrete implementations
- Add comprehensive test coverage for new List functionality
- Fix security warning in benchmark random data generator

## v1.14.4

- Add comprehensive GoDoc documentation for all exported interfaces, types, and functions
- Improve API documentation coverage for better developer experience
- Document bucket operations, store interfaces, transaction handling, and metrics
- Add usage examples and parameter descriptions to public APIs

## v1.14.3

- improve README with usage example and installation instructions
- go mod update

## v1.14.2

- add github workflow
- go mod update

## v1.14.1

- add tests
- go mod update

## v1.14.0

- add RunnableTx and FuncTx

## v1.13.2

- add lock to reset handlers and improve logging
- go mod update

## v1.13.1

- add NewStoreFromTx
- go mod update

## v1.13.0

- add DBWithMetrics
- remove vendor files
- go mod update

## v1.12.2

- go mod update

## v1.12.1

- go mod update
- add test for relation store mocks

## v1.12.0

- add Invert for the RelationStore and RelationStoreTx

## v1.11.5

- add MapIDRelations and MapRelationIDs

## v1.11.4

- remove performance bug in relationStoreTx delete

## v1.11.3

- move JsonHandlerTx to github.com/bborbe/http
- go mod update

## v1.11.2

- add missing license file
- go mod update

## v1.11.1

- rename NewUpdateHandlerViewTx -> NewJsonHandlerUpdateTx

## v1.11.0

- add JsonHandlerTx
- go mod update

## v1.10.0

- add ListBucketNames
- go mod update

## v1.9.1

- ignore BucketNotFoundError on Map, Remove and Exists
- go mod update

## v1.9.0

- add remove to DB to delete the complete database
- add handler for reset bucket and complete database
- go mod update

## v1.8.2

- fix replace in relationStore
- go mod update

## v1.8.1

- add simple benchmark

## v1.8.0

- add relation store
- go mod update

## v1.7.0

- expect same tx returns same bucket
- go mod update

## v1.6.0

- add stream and exists to store

## v1.5.0

- add JSON store
- go mod update

## v1.4.2

- add KeyNotFoundError

## v1.4.1

- add mocks

## v1.4.0

- expect error if transaction open second transaction

## v1.3.1

- improve iterator testsuite

## v1.3.0

- add bucket testsuite

## v1.2.0

- add provider
- improve testsuite

## v1.1.1

- add test for iterator seek not found

## v1.1.0

- Add context to update and view

## v1.0.0

- Initial Version
