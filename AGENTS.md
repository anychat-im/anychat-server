# Repository Guidelines

## Project Structure & Module Organization
This repository is a Go microservice backend for AnyChat. Service entrypoints live in `cmd/*-service`, domain code in `internal/<domain>`, and reusable helpers in `pkg/` (config, database, JWT, Redis, logging, middleware). API contracts are in `api/proto/`, migrations in `migrations/`, local setup in `deployments/` and `configs/`, and generated binaries in `bin/`. Contributor docs live in `docs/`, while HTTP integration scripts live in `tests/api/`.

## Build, Test, and Development Commands
Use Mage for day-to-day workflows:

- `mage deps` - download and tidy Go modules.
- `mage docker:up` - start PostgreSQL, Redis, NATS, MinIO, and other local dependencies from `deployments/docker-compose.yml`.
- `mage db:up` - apply local migrations.
- `mage build:all` - compile every service into `bin/`.
- `mage dev:gateway` or `mage dev:auth` - run a single service locally with `go run`.
- `mage test:all` - run all Go tests with race detection and coverage output.
- `./tests/api/test-all.sh` - run end-to-end HTTP API scripts against a running gateway stack.
- `mage lint` / `mage fmt` - run `golangci-lint` and format Go code.

## Coding Style & Naming Conventions
Follow standard Go formatting; run `mage fmt` before opening a PR. Keep package names lowercase, keep service directories named `*-service`, and match domain folders between `cmd/` and `internal/`. Linting is enforced through `.golangci.yml` with `gofmt`, `govet`, `errcheck`, `staticcheck`, and complexity checks. Do not hand-edit generated protobuf files such as `*.pb.go`; regenerate them with `mage proto`.

## Testing Guidelines
Unit tests use Go's `testing` package with `stretchr/testify`; place them next to the code as `*_test.go`. API coverage uses shell scripts named `test-<domain>-api.sh` under `tests/api/`. For new handlers or repository logic, add at least one focused Go test and update the relevant API script when behavior changes. Use `mage test:coverage` to produce `coverage.out` and `coverage.html`.

## Commit & Pull Request Guidelines
Recent history follows Conventional Commit style, for example `feat(auth): ...`, `fix(docs): ...`, `test: ...`, and `docs: ...`. Keep commits scoped to one service or concern. PRs should summarize affected services, note config or migration changes, link related issues, and include example requests/responses when API contracts change. If docs or specs change, update them in the same PR.

## Security & Configuration Tips
Keep secrets out of Git; prefer environment overrides and local config files. Treat the hard-coded migration DSN in `magefile.go` as a local default only, and verify Docker-based dependencies are running before executing destructive migration commands.
