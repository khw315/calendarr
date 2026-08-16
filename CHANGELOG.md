# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.0.1] - 2026-08-16

### Fixed
- Fixed container `permission denied` runtime error on `/app/config/calendarr.json` when mounting host volumes.
- Quoted variable expansions in `Dockerfile` to satisfy SonarQube `docker:S6570` compliance rules.

### Performance
- Optimized Docker `ARM64` multi-arch build speeds up to 10x using native Go cross-compilation (`--platform=$BUILDPLATFORM` and `TARGETARCH`).

### Added
- Added standard Open Container Initiative (OCI) image metadata labels (`title`, `source`, `documentation`, `licenses`) for GitHub Container Registry (GHCR).

## [2.0.0] - 2026-08-16

### Added
- Complete application rewrite from Python (Flask) to Golang (Go 1.24) using `go-chi/chi/v5` REST router and `robfig/cron/v3` scheduler.
- Embedded React Web UI static assets directly into a single static Go binary (`//go:embed`).
- Multi-architecture Docker container with non-root user security compliance (`alpine:3.21`, ~25MB image size).
- SonarQube quality gate scanning (`sonar-project.properties`) and automated GitHub Actions CI/CD workflows (`build.yml`, `golint.yml`, `docker_push.yml`).
- Comprehensive unit test suite across all `internal/` packages (>82% statement coverage).
- Visual GFM Mermaid architecture and data flow diagram in `README.md`.

### Changed
- Standardized DockerHub release secret names to `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN`.
- Updated Dependabot configuration to cover `gomod`, `npm`, `docker`, and `github-actions` ecosystems.

### Security
- Configured container runtime to execute as non-root user (`calendarr`).
- Added `--ignore-scripts` to `npm ci` step in Dockerfile frontend builder.

### Removed
- Removed legacy Python codebase (`src/`, `tests/`, `pyproject.toml`, `requirements.txt`, `entrypoint.sh`).
- Removed obsolete workflow files (`codeql.yml`, `add-to-project.yml`).

## [1.6.0] - 2026-03-20

### Added
- Web UI, API endpoints, and past releases view.
- Settings UI and persistent configuration system.
- Dark mode (toggle).
- Localization support with per-language files (EN, ID, KO, JA).
- Timezone API and dynamic timezone selection.
- Grouping of TV episodes into bulk events.
- Display of end time and current airing status.
- Localized empty state for no releases, including random selection messages.
- Configurable past releases lookup.

### Changed
- Refactored configuration logic, centralized bulk thresholds, and standardized timezone keys.
- Refactored shared Header and Footer; reorganized settings structure.
- Changed license from MIT to GPL v3.0.
- Updated Web UI design to Neobrutalist aesthetic.
- Improved event countdowns, airing states, and overall UI polish.
- Optimized Docker image and updated GHCR login configurations.

[Unreleased]: https://github.com/khw315/calendarr/compare/v2.0.1...HEAD
[2.0.1]: https://github.com/khw315/calendarr/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/khw315/calendarr/compare/v1.6.1...v2.0.0
[1.6.0]: https://github.com/khw315/calendarr/compare/v1.5.0...v1.6.0