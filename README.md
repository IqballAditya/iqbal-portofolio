# IQBAL - Lead Systems Architect & Full-Stack Engineer Portfolio

High-throughput systems architecture, AI orchestration, and full-stack engineering portfolio presented in a cinematic Netflix Dark Mode user interface. Built with Go 1.22+, HTMX, Tailwind-inspired design tokens, and pure-Go SQLite—compiled into a single deployable binary with zero CGO dependencies.

Designed from Figma specification `BF1A9Lq5hNSw01evMchdbC`.

---

## Architecture Highlights

- **Single Deployable Binary**: All HTML templates, stylesheets, JavaScript, and SVG/JPEG assets are embedded into the compiled binary via `embed.FS`. No runtime file dependencies.
- **Zero CGO SQLite**: Uses pure-Go SQLite (`modernc.org/sqlite`) configured in WAL mode (`PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`), enabling cross-compilation to any platform without gcc/clang.
- **HTMX Interactivity**: Delivers Single Page Application (SPA) interactions (real-time filtering, live search, dynamic modal popups, reactive shortlist bookmarking) with standard HTML over the wire and zero JavaScript framework overhead.
- **Sub-Millisecond Response Times**: Lightweight Go standard library router (`net/http.ServeMux`), in-memory SQLite querying, and pre-parsed template cloning achieve p99 response times under 1ms.
- **Minimal Footprint**: Operates with an RSS memory footprint under 25MB and zero external runtime dependencies.

---

## Features & UI Modules

### 1. Home Showcase Feed (`/`)
- **Billboard Hero Banner**: Highlighted flagship project (`Adith AI Agent`) with dynamic match scores, metadata badges, and direct action triggers.
- **Top 10 Tech Stacks Carousel**: Ranked Netflix Top 10 badges featuring key languages, databases, networking, and distributed systems tools.
- **Curated Rows**: Categorized carousels (AI & Distributed Agents, High-Throughput Infrastructure, Core Backend & Storage) with hover scale effects and quick-preview triggers.

### 2. "Who's Watching?" Persona Gateway (`/gateway`)
- Profile switcher allowing visitors to browse through tailored lenses:
  - **Recruiter**: Focuses on impact metrics, core accomplishments, and hire details.
  - **Tech Lead**: Prioritizes low-level architectures, throughput benchmarks, and concurrency models.
  - **Founder**: Emphasizes business ROI, engineering velocity, and end-to-end case studies.
  - **Guest**: Standard holistic portfolio showcase.
- Persisted via `iqbal_persona` cookie with dynamic avatar badges and tailored carousel sorting.

### 3. Netflix-Style Case Study Modal (`/projects/{id}/modal`)
- Full-screen modal overlay requested dynamically via HTMX.
- Includes architectural blueprints, system overview, role breakdown, measurable metrics, and episode-by-episode milestone guides with syntax-highlighted code snippets.

### 4. Interactive Demo Player (`/projects/{id}/player`)
- Cinematic terminal-based execution simulator mimicking a Netflix video player.
- Features simulated console boot sequences, ASCII banner art, interactive command prompt, step-by-step milestone log navigation, and real-time playback controls.

### 5. Live Search Matrix (`/search`, `/api/search`)
- Instant keystroke search filtering across projects, technical skills, tags, categories, and architectural concepts.
- Debounced HTMX queries with SQL injection-safe parameterized queries.

### 6. My Shortlist & Hire Me Tiers (`/shortlist`)
- **Session-Persisted Shortlist**: Bookmarking engine backed by cookie sessions (`iqbal_session`) allowing visitors to curate relevant projects.
- **Three Engagement Models**:
  - *Full-Time Staff / Lead*: Principal / Staff Backend Engineer & Systems Architect.
  - *Advisory / Architecture Sprint*: 2-4 week intensive bottleneck audit, concurrency review, and migration blueprints.
  - *Contract / Fractional*: Quarterly or monthly retainer for fractional CTO and senior bar evaluation.
- **Hire Me Contact Form**: Instant inquiry dispatch with input validation and inline feedback.

---

## Quickstart & Build

### Prerequisites
- Go 1.22 or newer

### Run in Development
```bash
go run ./cmd/server
```
The server will start on `http://localhost:8080` with an automatic SQLite database bootstrap and idempotent data seeding.

### Build Standalone Binary
```bash
# Windows
go build -o bin/iqbal-portfolio.exe ./cmd/server

# Linux / macOS
go build -o bin/iqbal-portfolio ./cmd/server
```

Run the built executable directly:
```bash
./bin/iqbal-portfolio.exe
```

---

## Configuration

The application is configured via environment variables:

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP port the server binds to |
| `DB_PATH` | `portfolio.db` | Path to SQLite database file |
| `WEBHOOK_URL` | *(optional)* | Webhook endpoint for forwarding contact form inquiries |

---

## Test Suite

Run all unit tests and end-to-end integration scenarios:
```bash
# Run full suite
go test ./... -v

# Run only E2E integration test
go test ./test/... -v
```

All 8 user journey stages (Home feed, Persona selection, Case study modal, Shortlist toggle, Shortlist page, Hire Me inquiry, Demo player, and Live search API) are automated and verified in `test/e2e_test.go`.
