# Design Document: Iqbal Portfolio Platform (Go + HTMX + Tailwind)

## 1. Intent & Purpose
Build a high-performance, single-binary portfolio web application for Iqbal Aditya using the Netflix UI design system ("IQBAL") created in Figma (`BF1A9Lq5hNSw01evMchdbC`). The application showcases full-stack and systems engineering craftsmanship by running on pure Go (`net/http`), Server-Driven UI with `html/template`, dynamic asynchronous interactions via HTMX, and SQLite persistence.

## 2. Goals, Non-Goals, and Anti-Goals

### Goals
- Recreate all 6 core Netflix features designed in Figma as interactive web pages:
  1. **Home / Browse Showcase Feed**: Hero spotlight (*Nexus AI Agent*), project carousels, ranked tech stacks, case studies, and footer.
  2. **Profile Gateway ("Who's Watching?")**: Persona switchers (Recruiter, Tech Lead, Founder, Guest) tailoring highlighted content.
  3. **Case Study Modal ("More Info")**: HTMX-driven modal popups loading project architectures, performance metrics, and engineering milestone episodes.
  4. **Interactive Demo Player ("Video Player")**: Interactive scrubber with chapter timelines, terminal output simulation, and playback controls.
  5. **Live Search & Filter**: Real-time project search by stack/category without full-page reloads.
  6. **Shortlist & Hire Me**: Bookmarked projects, engagement model selection, and contact inquiry submission.
- Zero CGO dependency: Uses pure Go SQLite (`modernc.org/sqlite`) for portable compilation across Windows, Linux, and macOS.
- Single deployable executable: All templates, static assets, and database migrations can be compiled directly into the binary with Go 1.16+ `embed`.
- Strict code hygiene: Conforms to `antislop-code` (no decorative banner comments, no restating the obvious, concise technical documentation).

### Non-Goals
- Complex multi-tenant authentication or JWT authorization: Portfolio is public-facing with visitor persona cookies.
- Video streaming backend / HLS transcoders: The "Video Player" simulates project architecture walkthroughs, terminal demos, and case study chapters.

### Anti-Goals
- No Node.js runtime required in production.
- No heavy JavaScript frameworks (React/Vue/Angular); HTMX and micro-scripts handle all dynamic client behavior.
- No generic AI-slop comments in Go code.

## 3. Technology Stack & Approaches

### Approaches Considered
1. **Go Backend + Next.js/React Frontend (Decoupled)**:
   - *Trade-off*: Requires maintaining two codebases, Node.js tooling, and more complex hosting. Ruled out to favor unified Go engineering.
2. **Go + html/template + HTMX + Tailwind CSS (Selected)**:
   - *Rationale*: Delivers server-driven speed (<5ms response), 1 compiled binary, native SEO, sub-20MB memory footprint, and showcases modern Go craftsmanship.

### Core Stack
- **Language**: Go 1.22+ (utilizing the enhanced `http.ServeMux` routing syntax).
- **Server**: Standard library `net/http`.
- **Database**: Embedded SQLite (`modernc.org/sqlite`) with WAL mode for fast concurrent reads.
- **Frontend Interactivity**: HTMX 1.9+ (CDN / vendored local script) for partial DOM swaps.
- **Styling**: Tailwind CSS (Netflix dark theme: `#141414` base, `#1F1F1F` cards, `#E50914` brand red, `#46D369` match green).
- **Design Source of Truth**: Figma File `BF1A9Lq5hNSw01evMchdbC` (Pages 01, 02, and 03).

## 4. Architecture & Data Model

### Builder Decisions & Clarifications
1. **Persona Content Tailoring**:
   - Stored in a `target_personas` column (CSV or JSON array) on `projects`.
   - The repository query dynamically boosts projects matching the visitor's selected persona (`recruiter`, `tech_lead`, `founder`, `guest`).
   - For example: `tech_lead` prioritizes `infra` and `ai` with architecture badges; `recruiter` highlights `fullstack` with business metrics.
2. **First-Time Visitors & Gateway Entry**:
   - `GET /` serves the portfolio showcase directly with default persona (`guest`) to ensure 100% SEO indexability and immediate recruiter access.
   - A subtle header badge and banner allows switching personas at any time, leading to `GET /gateway`.
   - `GET /gateway` is an explicit, full-screen Netflix persona selection gateway ("Who's Watching?").
3. **Demo Player Terminal Simulation**:
   - Self-contained client-side interactive player (`player.js`) driven by structured milestone chapter data passed from the Go template.
   - Simulates terminal output, interactive logs, and code tabs per chapter without needing heavy backend websocket streams.
4. **Inquiry Delivery & Notifications**:
   - Inquiries write to the SQLite `inquiries` table and log to server stdout.
   - An optional `WEBHOOK_URL` environment variable (Telegram/Discord/Slack webhook) can be set for instant notification dispatch, falling back gracefully if unconfigured.

### Database Schema (SQLite)
- `projects`:
  - `id` (TEXT PRIMARY KEY)
  - `title` (TEXT)
  - `category` (TEXT) - `ai`, `fullstack`, `infra`, `mobile`
  - `tagline` (TEXT)
  - `synopsis` (TEXT)
  - `tech_stack` (TEXT) - JSON array of strings
  - `target_personas` (TEXT) - e.g. `tech_lead,recruiter`
  - `match_score` (INTEGER) - e.g. 99
  - `role` (TEXT)
  - `metrics` (TEXT) - e.g. "50k+ req/sec, <80ms"
  - `github_url` (TEXT)
  - `demo_url` (TEXT)
  - `image_url` (TEXT)
  - `is_featured` (BOOLEAN)
  - `sort_order` (INTEGER)
- `milestones` (Episodes for Case Studies):
  - `id` (INTEGER PRIMARY KEY AUTOINCREMENT)
  - `project_id` (TEXT REFERENCES projects(id))
  - `episode_number` (INTEGER)
  - `title` (TEXT)
  - `duration` (TEXT)
  - `summary` (TEXT)
  - `code_snippet` (TEXT)
- `skills` (Top 10 Stacks):
  - `rank` (INTEGER PRIMARY KEY)
  - `name` (TEXT)
  - `category` (TEXT)
  - `description` (TEXT)
  - `image_url` (TEXT)
- `shortlists`:
  - `session_id` (TEXT)
  - `project_id` (TEXT)
  - `created_at` (DATETIME)
- `inquiries` (Hire Me submissions):
  - `id` (INTEGER PRIMARY KEY AUTOINCREMENT)
  - `name` (TEXT)
  - `email` (TEXT)
  - `engagement_type` (TEXT)
  - `message` (TEXT)
  - `created_at` (DATETIME)

### Routing Specification
- `GET /` - Home browse showcase feed
- `GET /gateway` - Who's Watching persona picker
- `POST /gateway/select` - Set persona cookie and redirect to `/`
- `GET /projects/{id}/modal` - HTMX endpoint returning the Case Study Modal partial
- `GET /projects/{id}/player` - Interactive Demo Player view
- `GET /search` - Search & filter page
- `GET /api/search` - HTMX endpoint returning live search results cards
- `GET /shortlist` - My Shortlist & Hire Me page
- `POST /shortlist/{id}` - HTMX endpoint toggling shortlist state
- `POST /inquiries` - HTMX endpoint handling hire/contact submission

## 5. Directory Structure
```
netflix_iqbal/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── database/
│   │   ├── db.go
│   │   └── schema.sql
│   ├── models/
│   │   ├── project.go
│   │   ├── skill.go
│   │   └── inquiry.go
│   ├── repository/
│   │   ├── project_repo.go
│   │   ├── skill_repo.go
│   │   ├── shortlist_repo.go
│   │   └── seed.go
│   ├── handlers/
│   │   ├── home.go
│   │   ├── gateway.go
│   │   ├── modal.go
│   │   ├── player.go
│   │   ├── search.go
│   │   └── hire.go
│   └── server/
│       └── server.go
├── web/
│   ├── templates/
│   │   ├── base.html
│   │   ├── pages/
│   │   │   ├── index.html
│   │   │   ├── gateway.html
│   │   │   ├── player.html
│   │   │   └── hire.html
│   │   └── partials/
│   │       ├── project_card.html
│   │       ├── modal_detail.html
│   │       ├── search_results.html
│   │       └── shortlist_item.html
│   └── static/
│       ├── css/
│       │   └── style.css
│       └── js/
│           ├── htmx.min.js
│           └── player.js
├── go.mod
├── go.sum
└── README.md
```

## 6. Implementation Guardrails
- **Antislop Code**:
  - No decorative separator comments (`// ======== ROUTES ========`).
  - No stating the obvious (`// User struct`).
  - No workflow narration (`// Step 1: parse form`).
  - Comments must only explain non-obvious constraints, architecture invariants, or protocol rules.
- **Reliability & Testing**:
  - Self-contained seed data initialized automatically on first startup.
  - Unit and HTTP handler tests in Go standard `testing` package verifying routing, database queries, and template rendering.
