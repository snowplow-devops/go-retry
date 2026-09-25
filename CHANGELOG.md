# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-09-25

### Added

- `Lint` CI job running `golangci-lint` via `golangci/golangci-lint-action` v9.3.0.
- Package doc comment for `retry`.
- Unit tests for `Exponential` and `ExponentialWithInterface` (100% statement coverage).
- `make test` target running the tests with the race detector and a coverage report.
- `Run unit tests` step in the CI `Compile & Test` job across the Go version matrix.

### Changed

- Convert `CHANGELOG` to `CHANGELOG.md` following the Keep a Changelog 1.1.0 format.
- Raise the minimum Go version from 1.13 to 1.26 (required by `golang.org/x/sys` v0.48.0).
- CI now tests against Go 1.26.x and 1.27.x instead of Go 1.13.15.
- Upgrade GitHub Actions to Node 24 releases, ahead of the Node 20 runtime deprecation:
  `actions/checkout` v2 → v7.0.1 and `actions/setup-go` v2 → v7.0.0.
- Pin all GitHub Actions to full commit SHAs instead of mutable tags.
- Replace the separate `actions/cache` step with the built-in module caching in `actions/setup-go`.
- Replace the deprecated `golint` with `golangci-lint` v2.14.0 (configured in `.golangci.yml`),
  used by both `make lint` and `make format`.

### Removed

- Snyk workflow (`.github/workflows/snyk.yml`), which is no longer used.

### Fixed

- README usage example now compiles (missing returns, unused variable, import grouping).
- README release badge now uses HTTPS and shows the latest tag automatically.
- `Exponential` and `ExponentialWithInterface` no longer panic when the sleep duration is
  zero or negative; they now retry immediately without sleeping.

### Security

- Upgrade `github.com/sirupsen/logrus` from v1.7.0 to v1.10.2, fixing
  [GO-2025-4188](https://pkg.go.dev/vuln/GO-2025-4188).
- Upgrade `golang.org/x/sys` from v0.0.0-20200930185726-fdedc70b468f to v0.48.0, fixing
  [GO-2022-0493](https://pkg.go.dev/vuln/GO-2022-0493) and
  [GO-2026-5024](https://pkg.go.dev/vuln/GO-2026-5024).

## [0.1.0] - 2021-01-06

### Added

- Initial release.

[unreleased]: https://github.com/snowplow-devops/go-retry/compare/0.2.0...HEAD
[0.2.0]: https://github.com/snowplow-devops/go-retry/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/snowplow-devops/go-retry/releases/tag/0.1.0
