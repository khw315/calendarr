# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- 3-mode theme system (`system`, `dark`, `light`) in Web UI, supporting dynamic browser/OS preference (`prefers-color-scheme`), 3-way toggle button in header, and Theme Mode dropdown in Settings.

## [2.1.0] - 2026-08-16

### Added
- Dynamic `?days` query parameter support on `/api/past-releases` endpoint to retrieve past release events for selectable ranges.
- Automatic background reload of iCal feeds immediately upon saving settings (`POST /api/config`).
- Comprehensive TypeScript interfaces (`ConfigState`, `CalendarUrlItem`, `DayGroup`, `EventItem`) ensuring type safety across the React Web UI.

### Changed
- Web UI day headers are now formatted exclusively in English (`Sunday, Aug 16`), while multi-language translations (French, German, Indonesian) apply strictly to Discord and Slack notifications.
- Updated Past Releases Web UI layout to group events under Day Headers and display only start times on cards (`20:45`), removing raw date badges.

### Fixed
- Fixed `DAILY` schedule mode parsing to strictly return 1-day events (today only) instead of spanning into 2 days.
- Fixed Past Releases date range to query strictly up to yesterday (`startOfDay - 1ns`), filtering out future airing events.
- Fixed 20 ESLint code quality, type safety, and React hook dependency warnings in `App.tsx` and `Settings.tsx`.
- Refactored `internal/api/router.go` and `internal/services/calendar/service.go` to eliminate code duplications and reduce Cognitive Complexity for SonarQube quality gate compliance.

## [2.0.2] - 2026-08-16

### Fixed
- Resolved host volume permission error on `/app/config/calendarr.json` by adding `entrypoint.sh` runtime permission correction with `su-exec` and assigning UID `1000:1000`.
- Fixed API routing 404 errors by explicitly mounting `/api` subrouter on Chi router ahead of SPA static file fallback.
- Fixed CORS issues when accessing `/api/config` and other endpoints from external IP addresses and web frontends.
- Fixed `context canceled` errors during background calendar feed fetching over HTTP by using an independent background timeout context.
- Quoted Dockerfile variables and sorted package names alphanumerically for SonarQube `docker:S6570`, `docker:S7018`, and `shell:S2612` rule compliance.

## [2.0.1] - 2026-08-16

### Fixed
- Resolved container volume `permission denied` error on `/app/config/calendarr.json` when mounting host directories by configuring `chmod -R 775` permissions.
- Quoted variable expansions in `Dockerfile` for SonarQube `docker:S6570` compliance.

### Security
- Resolved SonarQube `docker:S2612` security rule by restricting world-write permissions on `/app/config` and `/app/logs` volume directories.

### Performance
- Accelerated ARM64 Docker image builds up to 10x using native Go cross-compilation with `--platform=$BUILDPLATFORM` and `TARGETARCH`.

### Documentation
- Added standard GHCR OCI container labels (`org.opencontainers.image.title`, `source`, `documentation`, `licenses`).

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

[Unreleased]: https://github.com/khw315/calendarr/compare/v2.1.0...HEAD
[2.1.0]: https://github.com/khw315/calendarr/compare/v2.0.2...v2.1.0
[2.0.2]: https://github.com/khw315/calendarr/compare/v2.0.1...v2.0.2
[2.0.1]: https://github.com/khw315/calendarr/compare/v2.0.0...v2.0.1
[2.0.0]: https://github.com/khw315/calendarr/compare/v1.6.1...v2.0.0
[1.6.0]: https://github.com/khw315/calendarr/compare/v1.5.0...v1.6.0