# Calendarr Agent Guidelines & Repository Architecture

This document provides guidelines, technical standards, code conventions, and workflows for AI agents and human developers contributing to the **Calendarr** repository. This application is a high-performance Go service that parses Sonarr & Radarr calendar feeds (iCal) and dispatches upcoming TV show and movie release notifications to Discord and Slack.

---

## 1. Repository Identity & Technology Stack

- **Primary Language**: Go 1.27.2 (Modern standard library, concurrency patterns, embedded filesystem).
- **HTTP Routing & API**: `github.com/go-chi/chi/v5` with logging, recovery, CORS, and Single-Page Application (SPA) fallback middleware.
- **iCal Parser**: `github.com/arran4/golang-ical`.
- **Background Scheduler**: `github.com/robfig/cron/v3` supporting dynamic cron expressions.
- **Frontend SPA**: Next.js 15 (App Router) + React 19 + Tailwind CSS + Lucide React, statically built into `public/` and embedded into the Go binary via `//go:embed public/*`.
- **Webhook Integrations**:
  - Discord Webhooks (Embeds, dynamic markdown timestamps `<t:TIMESTAMP:STYLE>`, bulk episode grouping, premiere celebration).
  - Slack Webhooks (Block Kit sections, context, divider layouts).
- **Internationalization (i18n)**: Embedded JSON locales in `internal/localization/locales/` (`en`, `fr`, `id`, `ja`, `ko`).

---

## 2. Git Workflow & Contribution Discipline

### 2.1 Branching Model
- **Base Branch**: `main` is the primary production branch.
- **Feature/Fix Branches**: Create isolated branches branching off `main`:
  - `feat/<feature-name>` for new features.
  - `fix/<bug-name>` for bug fixes.
  - `refactor/<change-name>` for non-functional code restructuring.
  - `chore/<task-name>` for dependency maintenance or build configurations.
- All pull requests (PRs) must target `main`.
- Working branches must be deleted promptly once merged into `main`.

### 2.2 Conventional Commits
All commits must strictly follow the **Conventional Commits** standard:
```text
<type>(<scope>): <short description in lowercase and imperative mood>
```

- **Allowed Types**:
  - `feat`: New features or enhancements.
  - `fix`: Bug fixes.
  - `refactor`: Code refactoring that neither fixes a bug nor adds a feature.
  - `style`: Formatting, whitespace, semi-colons (no logical changes).
  - `test`: Adding or correcting unit tests.
  - `chore`: Maintenance tasks, dependencies, build configurations.
  - `docs`: Documentation updates (README, AGENTS.md, etc.).
  - `build`: Build system changes, Go module updates, Dockerfile.

- **Standard Scopes**:
  - `calendar`: iCal parsing, event deduplication, HTTP fetching.
  - `formatter`: Discord embeds, Slack blocks, bulk release grouping, localization headers.
  - `platform`: Webhook HTTP dispatcher client, retry mechanism.
  - `scheduler`: Background cron scheduler, lifecycle triggers, startup runs.
  - `api`: Chi router, HTTP handlers, REST endpoints.
  - `config`: Configuration manager, file persistence, validation.
  - `localization`: i18n locales, date formatters.
  - `models`: Domain structs, helpers, validation methods.
  - `frontend`: Next.js UI, React components, Tailwind styling.
  - `deps`: Dependency updates in `go.mod` or `package.json`.

### 2.3 Atomic Commits
- Commit changes in small, logical, self-contained units.
- Do not mix mass formatting (`gofmt`), library upgrades, and bug fixes into a single commit.

### 2.4 CHANGELOG & Versioning
- Every functional change or bug fix must be documented in `CHANGELOG.md` under the `## [Unreleased]` section.
- Adhere strictly to **Keep a Changelog** standards:
  - `### Added` for new features.
  - `### Changed` for changes in existing functionality.
  - `### Fixed` for bug fixes.
  - `### Removed` for removed features.
  - `### Security` for security vulnerability fixes.
- The repository follows **Semantic Versioning** (`vMAJOR.MINOR.PATCH`).

---

## 3. Architecture & Directory Layout

```text
calendarr/
├── main.go                     # Application bootstrap, embedded static FS, graceful shutdown
├── go.mod / go.sum             # Go module dependencies
├── calendarr.json              # Runtime configuration file (auto-generated if missing)
├── CHANGELOG.md                # Release history (Keep a Changelog)
├── AGENTS.md                   # Agent technical guidelines & operational standards
├── frontend/                   # Next.js 15 frontend source code
│   ├── app/                    # Next.js App Router (pages, layout, API proxies)
│   ├── components/             # Reusable UI components
│   └── package.json            # Node.js dependencies & build scripts
├── public/                     # Static frontend build artifacts (embedded into Go binary)
└── internal/                   # Private application packages
    ├── api/                    # Chi HTTP REST router & SPA handler
    ├── config/                 # Thread-safe config manager & JSON persistence
    ├── constants/              # Centralized domain constants (colors, thresholds, limits)
    ├── localization/           # Multi-language translation & date formatters
    │   └── locales/            # Translation JSON files (en.json, id.json, etc.)
    ├── models/                 # Domain models (Event, Config, DTOs, Webhook payloads)
    ├── services/               # Core business logic services
    │   ├── calendar/           # Fetching & parsing iCal feeds
    │   ├── formatter/          # Generating Discord & Slack payloads
    │   ├── platform/           # Dispatching HTTP webhooks to Discord & Slack
    │   └── scheduler/          # Background cron runner & dynamic rescheduling
    └── tzdata/                 # IANA timezone dataset provider
```

---

## 4. Go Coding Standards & Idiomatic Practices

### 4.1 Cognitive Complexity & Readability
- **Complexity Threshold**: Maintain **Cognitive Complexity < 15** for every function.
- If a function contains deep iterations, nested conditionals, or multi-step formatting (such as Discord payload builders), break it down into single-purpose helper functions (e.g., `groupEventsByDay`, `groupTVShows`, `formatDayDiscordLines`).
- Favor linear code flow; avoid deep nesting and utilize early returns / guard clauses.

### 4.2 Avoid Shadowing Predeclared Built-ins
- Never shadow Go's predeclared built-in identifiers:
  - ❌ `min`, `max`, `len`, `cap`, `close`, `copy`, `error`, `new`, `make`.
  - ✅ Use descriptive variable names: `minute`, `minVal`, `maxCount`, `length`.

### 4.3 Explicit Error Handling & Logging
- **Explicit Handling**: Never ignore or swallow errors (`_ = err`).
- **Error Wrapping**: Always wrap errors with contextual messages using `%w`:
  ```go
  if err != nil {
      return fmt.Errorf("failed to parse iCal feed %s: %w", url, err)
  }
  ```
- **Logging Conventions**: Use standard Calendarr emoji prefixes for log readability:
  - 🚀 Process startup / initialization
  - ⚡ Calendar job execution
  - 🗓️ Fetching and parsing calendar events
  - 📅 Setting cron schedules
  - ✅ Operation success (config saved, webhook dispatched)
  - ⚠️ Non-fatal warning (fallback config, invalid timezone, partial feed error)
  - ❌ Fatal error or webhook dispatch failure

### 4.4 Concurrency & Thread Safety
- **Shared Memory**: Protect all access to shared in-memory state (`ConfigManager`) using `sync.RWMutex` (`RLock` for reads, `Lock` for writes).
- **Bounded Concurrency**: Fetch calendar feeds concurrently using goroutines and `sync.WaitGroup`, protected by mutexes when deduplicating results.
- **Context & Timeouts**: All outbound HTTP requests (iCal fetch, webhook POST) must accept a `context.Context` and enforce deadlines (`http.Client.Timeout` or `context.WithTimeout`).
- **Memory Allocation**: Pre-allocate slice capacities whenever the length is predictable:
  ```go
  events := make([]models.Event, 0, len(rawEvents))
  ```

### 4.5 Testing Standards
- **Table-Driven Tests**: Use table-driven test structures to cover edge cases and input variations:
  ```go
  tests := []struct {
      name     string
      input    InputType
      expected ExpectedType
  }{ ... }
  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) { ... })
  }
  ```
- **Race Detector**: All tests must pass race detection:
  ```bash
  go test -v ./...
  # When CGO is available (e.g., in CI or environments with GCC/Clang):
  go test -v -race ./...
  ```
- **HTTP Mocking**: Use `httptest.NewServer` for testing network calls without hitting live endpoints.
- **Filesystem Isolation**: Use `t.TempDir()` for tests that write config files or disk artifacts.

---

## 5. Webhook Domain Constraints & Rules

### 5.1 Discord Webhooks
- **Payload Limits**: Maximum 10 embeds per message; total characters across all fields must remain well below 6,000 characters (safety margin: 5,800 bytes).
- **Dynamic Timestamps**: Use standard Discord timestamp formatting:
  - `<t:UNIX:F>` for full date and time.
  - `<t:UNIX:R>` for relative time (e.g., *in 2 hours*).
  - `<t:UNIX:t>` for short time (e.g., *09:30*).
- **Multi-Episode Bulk Grouping**:
  - When a TV series releases more than the bulk threshold (`BulkThresholdDiscord = 2` episodes) on the same day, collapse them into a single concise line:
    `**Show Name** — <t:TIMESTAMP:STYLE> 🎉` (append the celebration emoji if it is a premiere).
  - For 1 or 2 episodes, render each episode individually (`Season X Episode Y - Episode Title`).
- **Premiere Detection**:
  - Use `event.IsPremiere()` to identify premiere releases (Season 1 Episode 1 or series premiere) and append `🎉`.

### 5.2 Slack Webhooks
- Use Slack Block Kit:
  - Header block for the schedule title.
  - Section blocks with markdown for daily release listings.
  - Divider blocks between dates.
  - Context blocks for metadata summaries.

---

## 6. Frontend Architecture & Embedded Static Assets

- The UI source code resides in `frontend/`.
- During production builds, assets are compiled into `public/` (`npm run build` in `frontend/`).
- `main.go` embeds the `public/` directory:
  ```go
  //go:embed public/*
  var embeddedPublicFS embed.FS
  ```
- The HTTP handler in `internal/api/router.go` serves static files from `embeddedPublicFS` and routes non-API paths back to `index.html` (Single-Page Application fallback).
- Backend APIs are hosted under `/api/*` (`/api/config`, `/api/run`, `/api/timezones`, `/api/languages`).

---

## 7. Pre-Commit Quality Gates

Before creating a commit or opening a pull request, AI agents and developers MUST execute and verify all the following steps:

1. **Formatting**:
   ```powershell
   gofmt -w .
   ```
2. **Static Analysis**:
   ```powershell
   go vet ./...
   ```
3. **Unit Tests**:
   ```powershell
   go test -v ./...
   # If CGO is enabled:
   go test -v -race ./...
   ```
4. **Git Workspace Cleanliness**:
   Ensure no untracked scratch files, temp files, or unintended changes exist in `git status`.

---

## 8. AI Agent Operational Principles

1. **Inspect Before Execution**: Always inspect existing code and understand dependencies before refactoring or implementing features.
2. **Preserve Code Integrity**: Retain existing docstrings, naming structures, and translation dictionaries in `internal/localization/locales/`.
3. **Credential Security**: Never commit or log real webhook URLs, API keys, or sensitive credentials.
4. **Structured Communication**: Always provide clear change summaries, clickable markdown file links (`[filename](file:///path/to/file)`), and easily verifiable instructions.
