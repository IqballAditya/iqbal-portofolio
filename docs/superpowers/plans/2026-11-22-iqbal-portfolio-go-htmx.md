# Iqbal Portfolio Platform (Go + HTMX + Tailwind) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a high-performance, single-binary portfolio web platform in Go + HTMX + Tailwind CSS for Iqbal Aditya, faithfully translating all 6 Netflix UI features designed in Figma into interactive web views.

**Architecture:** Pure Go standard library HTTP server (`net/http`) using Go 1.22+ routing, pure Go SQLite database (`modernc.org/sqlite`) with WAL mode, Server-Driven UI with `html/template` and `embed.FS`, asynchronous micro-interactions with HTMX, and a client-side simulated architecture walkthrough player in vanilla JS.

**Tech Stack:** Go 1.22+, `modernc.org/sqlite`, `html/template`, HTMX 1.9+, Tailwind CSS.

**Spec:** `docs/superpowers/specs/2026-11-22-iqbal-portfolio-go-htmx-design.md`

## Global Constraints

- Zero CGO dependency: must compile cleanly without gcc/cgo using `modernc.org/sqlite`.
- Single binary deployment: all templates (`web/templates/**`) and static assets (`web/static/**`) must be embedded via `embed.FS`.
- Netflix design tokens: Background `#141414`, surface `#1F1F1F`, brand red `#E50914`, match score `#46D369`, card radius `4px`.
- Code hygiene: `antislop-code` compliance (no decorative banner comments, no restating the obvious, concise technical documentation).
- Ponytail: Minimum code that solves the entire problem cleanly; no unnecessary frameworks or wrappers.

## Review Focus

1. Visitor with missing or corrupted `iqbal_persona` cookie -> Defaults safely to `guest` persona without panic or broken page render.
2. HTMX modal request for non-existent or malformed project ID (`GET /projects/{id}/modal`) -> Returns clean 404 partial card instead of crashing server or rendering an empty overlay.
3. Live search with empty query or special SQL characters (`GET /api/search?q=' OR 1=1--`) -> Safely parameterized query returning appropriate matches without injection or template errors.
4. Shortlist toggle by anonymous visitor -> Persists against visitor session cookie without losing or duplicating existing bookmarks.
5. Hire Me form submission with invalid email or empty message -> Returns validation error messages directly within the HTMX target without page reload.

---

### Task 1: Project Scaffolding & Pure Go SQLite Database Layer

**Files:**
- Create: `go.mod`
- Create: `internal/database/schema.sql`
- Create: `internal/database/db.go`
- Test: `internal/database/db_test.go`

**Interfaces:**
- Consumes: Standard library `database/sql`, `modernc.org/sqlite`
- Produces:
  ```go
  func InitDB(dataSourceName string) (*sql.DB, error)
  func Migrate(db *sql.DB) error
  ```

- [ ] **Step 1: Initialize Go module & dependencies**

Run:
```bash
go mod init iqbal-portfolio
go get modernc.org/sqlite@v1.36.0
```

- [ ] **Step 2: Write failing test in `internal/database/db_test.go`**

```go
package database_test

import (
	"os"
	"testing"
	"iqbal-portfolio/internal/database"
)

func TestInitDBAndMigrate(t *testing.T) {
	tmpFile := "test_portfolio.db"
	defer os.Remove(tmpFile)

	db, err := database.InitDB(tmpFile)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='projects'").Scan(&tableName)
	if err != nil || tableName != "projects" {
		t.Fatalf("expected projects table to exist, got %v (table: %s)", err, tableName)
	}
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./internal/database/... -v`
Expected: FAIL (compilation error, package `database` undefined)

- [ ] **Step 4: Implement `schema.sql` and `db.go`**

Create `internal/database/schema.sql` with table definitions for `projects`, `milestones`, `skills`, `shortlists`, and `inquiries`.
Create `internal/database/db.go` embedding `schema.sql` using `//go:embed schema.sql`, enabling WAL mode `PRAGMA journal_mode=WAL;`, and executing migrations.

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./internal/database/... -v`
Expected: PASS

- [ ] **Step 6: Initialize git repository and commit**

```bash
git init
git add go.mod go.sum internal/database/
git commit -m "feat(database): setup pure Go SQLite connection and schema migrations"
```

---

### Task 2: Data Models, Repositories & Comprehensive Seed Data

**Files:**
- Create: `internal/models/project.go`
- Create: `internal/models/skill.go`
- Create: `internal/models/inquiry.go`
- Create: `internal/repository/project_repo.go`
- Create: `internal/repository/skill_repo.go`
- Create: `internal/repository/shortlist_repo.go`
- Create: `internal/repository/inquiry_repo.go`
- Create: `internal/repository/seed.go`
- Test: `internal/repository/repo_test.go`

**Interfaces:**
- Consumes: `internal/database`, `database/sql`
- Produces:
  ```go
  type ProjectRepository interface {
      ListProjects(persona string) ([]models.Project, error)
      GetProjectByID(id string) (*models.Project, error)
      SearchProjects(query string, category string) ([]models.Project, error)
      ListMilestones(projectID string) ([]models.Milestone, error)
  }
  type SkillRepository interface {
      ListSkills() ([]models.Skill, error)
  }
  type ShortlistRepository interface {
      ToggleShortlist(sessionID, projectID string) (bool, error)
      ListShortlist(sessionID string) ([]models.Project, error)
      IsShortlisted(sessionID, projectID string) (bool, error)
  }
  type InquiryRepository interface {
      CreateInquiry(inq *models.Inquiry) error
  }
  func SeedData(db *sql.DB) error
  ```

- [ ] **Step 1: Write failing test in `internal/repository/repo_test.go`**

Test covers:
- Seeding data into an in-memory SQLite DB (`:memory:`).
- Querying projects with persona ranking (`recruiter`, `tech_lead`, `guest`).
- Searching projects with query strings and categories.
- Toggling project shortlist status for a session.
- Submitting an inquiry.

```go
package repository_test

import (
	"testing"
	"iqbal-portfolio/internal/database"
	"iqbal-portfolio/internal/repository"
)

func TestRepositoriesAndSeed(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}

	if err := repository.SeedData(db); err != nil {
		t.Fatalf("SeedData failed: %v", err)
	}

	pRepo := repository.NewProjectRepo(db)
	projects, err := pRepo.ListProjects("tech_lead")
	if err != nil || len(projects) == 0 {
		t.Fatalf("expected projects, got %v (len: %d)", err, len(projects))
	}

	hero, err := pRepo.GetProjectByID("nexus-ai")
	if err != nil || hero.ID != "nexus-ai" {
		t.Fatalf("expected nexus-ai project, got %v", err)
	}

	sRepo := repository.NewSkillRepo(db)
	skills, err := sRepo.ListSkills()
	if err != nil || len(skills) < 10 {
		t.Fatalf("expected at least 10 skills, got %v (len: %d)", err, len(skills))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/repository/... -v`
Expected: FAIL (packages and constructors undefined)

- [ ] **Step 3: Implement models, repositories and seed data**

- Define `models.Project`, `models.Milestone`, `models.Skill`, `models.Inquiry`.
- Implement `project_repo.go`, `skill_repo.go`, `shortlist_repo.go`, `inquiry_repo.go`.
- Implement `seed.go` containing:
  - *Nexus AI Agent* (ID: `nexus-ai`, category: `ai`, match: 99%, role: "Lead Systems Architect", metrics: "50k+ req/sec, <80ms", 3 milestones with simulated terminal snippets).
  - 5+ supporting engineering projects: Distributed Key-Value Store, Cloud Cost Optimizer, Real-time Stream Engine, High-Throughput gRPC Gateway, Event Broker.
  - Top 10 skills: Go, Rust, PostgreSQL, Kubernetes, Redis, Docker, gRPC, TypeScript, Kafka, AWS.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/repository/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/models/ internal/repository/
git commit -m "feat(repo): add data models, sqlite repositories, and portfolio seed data"
```

---

### Task 3: Server Core, Static Assets & Base Netflix Layout

**Files:**
- Create: `web/static/css/style.css`
- Create: `web/static/js/htmx.min.js`
- Create: `web/templates/base.html`
- Create: `internal/server/server.go`
- Test: `internal/server/server_test.go`

**Interfaces:**
- Consumes: `html/template`, `net/http`, `embed`
- Produces:
  ```go
  type Server struct {
      router *http.ServeMux
      db     *sql.DB
      tmpl   *template.Template
  }
  func NewServer(db *sql.DB) (*Server, error)
  func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request)
  ```

- [ ] **Step 1: Write failing test in `internal/server/server_test.go`**

```go
package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"iqbal-portfolio/internal/database"
	"iqbal-portfolio/internal/server"
)

func TestServerStaticAndBaseRoute(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	srv, err := server.NewServer(db)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	req := httptest.NewRequest("GET", "/static/css/style.css", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected /static/css/style.css to return 200, got %d", rec.Code)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/server/... -v`
Expected: FAIL (compilation error, `server` package undefined)

- [ ] **Step 3: Implement style.css, htmx.min.js, base.html, and server.go**

- `web/static/css/style.css`: Netflix theme CSS variables (`--bg: #141414`, `--card: #1F1F1F`, `--red: #E50914`, `--green: #46D369`), utility classes for dark Netflix layout, card transitions, custom scrollbars, and modal backdrops.
- `web/static/js/htmx.min.js`: Vendored HTMX 1.9+ library for zero-dependency execution.
- `web/templates/base.html`: Top navbar with IQBAL logo in Netflix bold serif, links (Browse, Projects, Tech Stack, Shortlist), search bar trigger, persona switcher dropdown/link to `/gateway`, and responsive container.
- `internal/server/server.go`: Embedded template loading (`//go:embed web/templates/*`), static file server via `http.FS`, template helper functions (e.g. JSON array parser, score formatter).

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/server/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add web/ internal/server/
git commit -m "feat(server): setup embedded assets, netflix theme styles, and server foundation"
```

---

### Task 4: Persona Gateway & Home / Browse Showcase Feed

**Files:**
- Create: `internal/handlers/gateway.go`
- Create: `internal/handlers/home.go`
- Create: `web/templates/pages/gateway.html`
- Create: `web/templates/pages/index.html`
- Create: `web/templates/partials/project_card.html`
- Test: `internal/handlers/home_test.go`

**Interfaces:**
- Consumes: `internal/server`, `internal/repository`
- Produces:
  ```go
  func (s *Server) HandleGateway(w http.ResponseWriter, r *http.Request)
  func (s *Server) HandleGatewaySelect(w http.ResponseWriter, r *http.Request)
  func (s *Server) HandleHome(w http.ResponseWriter, r *http.Request)
  ```

- [ ] **Step 1: Write failing test in `internal/handlers/home_test.go`**

```go
package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"iqbal-portfolio/internal/database"
	"iqbal-portfolio/internal/repository"
	"iqbal-portfolio/internal/server"
)

func TestGatewayAndHomeRoutes(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()
	database.Migrate(db)
	repository.SeedData(db)

	srv, err := server.NewServer(db)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	// 1. Gateway view
	reqGateway := httptest.NewRequest("GET", "/gateway", nil)
	recGateway := httptest.NewRecorder()
	srv.ServeHTTP(recGateway, reqGateway)
	if recGateway.Code != http.StatusOK || !strings.Contains(recGateway.Body.String(), "Who's Watching?") {
		t.Fatalf("expected gateway view with Who's Watching, got status %d", recGateway.Code)
	}

	// 2. Persona selection
	form := strings.NewReader("persona=tech_lead")
	reqSelect := httptest.NewRequest("POST", "/gateway/select", form)
	reqSelect.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recSelect := httptest.NewRecorder()
	srv.ServeHTTP(recSelect, reqSelect)
	if recSelect.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect 303, got %d", recSelect.Code)
	}
	cookies := recSelect.Result().Cookies()
	var found bool
	for _, c := range cookies {
		if c.Name == "iqbal_persona" && c.Value == "tech_lead" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected iqbal_persona cookie to be set to tech_lead")
	}

	// 3. Home view
	reqHome := httptest.NewRequest("GET", "/", nil)
	recHome := httptest.NewRecorder()
	srv.ServeHTTP(recHome, reqHome)
	if recHome.Code != http.StatusOK || !strings.Contains(recHome.Body.String(), "Nexus AI Agent") {
		t.Fatalf("expected home view with Nexus AI Agent, got status %d", recHome.Code)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/handlers/... -v`
Expected: FAIL (handlers not implemented)

- [ ] **Step 3: Implement Gateway and Home handlers & templates**

- `web/templates/pages/gateway.html`: Full-screen Netflix profile selection ("Who's Watching?") displaying 4 avatar cards: Recruiter, Tech Lead, Founder, Guest.
- `web/templates/partials/project_card.html`: Reusable Netflix card with hover scale, match badge (`99% Match`), tags, and HTMX trigger `hx-get="/projects/{id}/modal"` targeting `#modal-container`.
- `web/templates/pages/index.html`:
  - Hero banner with *Nexus AI Agent* spotlight: title, metadata (99% Match, 2026, 18+ YRS EXP, Ultra HD, System Architect), synopsis, "Play Walkthrough" and "More Info" buttons.
  - Carousels: "Continue Reviewing Iqbal's Projects", "Top 10 Tech Stacks by Iqbal Today", "Proven Impact & Case Studies".
- `internal/handlers/gateway.go` & `home.go`: Read persona cookie, fetch projects and skills, render templates.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/handlers/... -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/handlers/ web/templates/
git commit -m "feat(handlers): implement persona gateway and netflix home showcase feed"
```

---

### Task 5: HTMX Case Study Modal & Live Search Filter

**Files:**
- Create: `internal/handlers/modal.go`
- Create: `internal/handlers/search.go`
- Create: `web/templates/partials/modal_detail.html`
- Create: `web/templates/pages/search.html`
- Create: `web/templates/partials/search_results.html`
- Test: `internal/handlers/modal_search_test.go`

**Interfaces:**
- Consumes: `internal/server`, `internal/repository`
- Produces:
  ```go
  func (s *Server) HandleProjectModal(w http.ResponseWriter, r *http.Request)
  func (s *Server) HandleSearchPage(w http.ResponseWriter, r *http.Request)
  func (s *Server) HandleSearchAPI(w http.ResponseWriter, r *http.Request)
  ```

- [ ] **Step 1: Write failing test in `internal/handlers/modal_search_test.go`**

```go
package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"iqbal-portfolio/internal/database"
	"iqbal-portfolio/internal/repository"
	"iqbal-portfolio/internal/server"
)

func TestModalAndSearchHandlers(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()
	database.Migrate(db)
	repository.SeedData(db)

	srv, err := server.NewServer(db)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	// 1. Modal partial request
	reqModal := httptest.NewRequest("GET", "/projects/nexus-ai/modal", nil)
	recModal := httptest.NewRecorder()
	srv.ServeHTTP(recModal, reqModal)
	if recModal.Code != http.StatusOK || !strings.Contains(recModal.Body.String(), "Episodes & Milestones") {
		t.Fatalf("expected modal with Episodes & Milestones, got status %d", recModal.Code)
	}

	// 2. Missing project modal (Review Focus 2)
	reqMissing := httptest.NewRequest("GET", "/projects/non-existent/modal", nil)
	recMissing := httptest.NewRecorder()
	srv.ServeHTTP(recMissing, reqMissing)
	if recMissing.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing project modal, got %d", recMissing.Code)
	}

	// 3. Live search HTMX query (Review Focus 3)
	reqSearch := httptest.NewRequest("GET", "/api/search?q=Golang", nil)
	recSearch := httptest.NewRecorder()
	srv.ServeHTTP(recSearch, reqSearch)
	if recSearch.Code != http.StatusOK {
		t.Fatalf("expected search query to succeed, got %d", recSearch.Code)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/handlers/... -run TestModalAndSearchHandlers -v`
Expected: FAIL

- [ ] **Step 3: Implement Modal and Search handlers + partial templates**

- `web/templates/partials/modal_detail.html`: Netflix "More Info" dialog with backdrop, hero preview, play CTA, episode milestone accordion/list with duration and summary, tech stack tags, and close button (`hx-on:click="document.getElementById('modal-container').innerHTML=''"`).
- `web/templates/pages/search.html`: Netflix Search Matrix page with category pills (`All`, `AI`, `Backend`, `Infra`) and search input with HTMX debounced trigger: `hx-get="/api/search" hx-trigger="keyup changed delay:200ms" hx-target="#search-results"`.
- `web/templates/partials/search_results.html`: Grid of project cards matching filter criteria.
- `internal/handlers/modal.go` & `search.go`: Handlers returning the partials with proper 404 handling.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/handlers/... -run TestModalAndSearchHandlers -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/handlers/modal.go internal/handlers/search.go web/templates/
git commit -m "feat(handlers): add htmx case study modal and instant search matrix"
```

---

### Task 6: Interactive Demo Player & Shortlist / Hire Me Engagement

**Files:**
- Create: `internal/handlers/player.go`
- Create: `internal/handlers/hire.go`
- Create: `web/templates/pages/player.html`
- Create: `web/templates/pages/hire.html`
- Create: `web/templates/partials/shortlist_item.html`
- Create: `web/static/js/player.js`
- Test: `internal/handlers/player_hire_test.go`

**Interfaces:**
- Consumes: `internal/server`, `internal/repository`
- Produces:
  ```go
  func (s *Server) HandlePlayer(w http.ResponseWriter, r *http.Request)
  func (s *Server) HandleShortlistPage(w http.ResponseWriter, r *http.Request)
  func (s *Server) HandleShortlistToggle(w http.ResponseWriter, r *http.Request)
  func (s *Server) HandleInquirySubmit(w http.ResponseWriter, r *http.Request)
  ```

- [ ] **Step 1: Write failing test in `internal/handlers/player_hire_test.go`**

```go
package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"iqbal-portfolio/internal/database"
	"iqbal-portfolio/internal/repository"
	"iqbal-portfolio/internal/server"
)

func TestPlayerAndHireHandlers(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()
	database.Migrate(db)
	repository.SeedData(db)

	srv, err := server.NewServer(db)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	// 1. Player route
	reqPlayer := httptest.NewRequest("GET", "/projects/nexus-ai/player", nil)
	recPlayer := httptest.NewRecorder()
	srv.ServeHTTP(recPlayer, reqPlayer)
	if recPlayer.Code != http.StatusOK || !strings.Contains(recPlayer.Body.String(), "TERMINAL OUTPUT") {
		t.Fatalf("expected player route with terminal output, got %d", recPlayer.Code)
	}

	// 2. Shortlist toggle (Review Focus 4)
	reqToggle := httptest.NewRequest("POST", "/shortlist/nexus-ai", nil)
	recToggle := httptest.NewRecorder()
	srv.ServeHTTP(recToggle, reqToggle)
	if recToggle.Code != http.StatusOK {
		t.Fatalf("expected shortlist toggle 200, got %d", recToggle.Code)
	}

	// 3. Inquiry submission validation error on empty email (Review Focus 5)
	inqForm := strings.NewReader("name=Alex&email=&message=")
	reqInq := httptest.NewRequest("POST", "/inquiries", inqForm)
	reqInq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recInq := httptest.NewRecorder()
	srv.ServeHTTP(recInq, reqInq)
	if recInq.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty inquiry email, got %d", recInq.Code)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/handlers/... -run TestPlayerAndHireHandlers -v`
Expected: FAIL

- [ ] **Step 3: Implement Player, Shortlist, and Hire Me handlers + templates + JS**

- `web/static/js/player.js`: Scrubber controller, play/pause state, chapter markers, simulated typewriter terminal logs showing microsecond-level execution of Go services.
- `web/templates/pages/player.html`: Netflix video player interface with back-to-browse arrow, scrub bar, chapter selector, and interactive terminal display.
- `web/templates/pages/hire.html`: Netflix "Choose Your Plan" pricing/engagement tiers:
  - Tier 1: Full-Time Senior / Lead Engineer (ID: `staff`)
  - Tier 2: Advisory / Architecture Sprint (ID: `sprint`)
  - Tier 3: Fractional / Contract Consulting (ID: `fractional`)
  - With shortlisted projects list and contact form submitting to `/inquiries` via HTMX.
- `internal/handlers/player.go` & `hire.go`: Implementing session cookie handling for shortlists and inquiry form validation.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/handlers/... -run TestPlayerAndHireHandlers -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/handlers/player.go internal/handlers/hire.go web/
git commit -m "feat(handlers): add interactive demo player and shortlist hire me system"
```

---

### Task 7: Application Entry Point, End-to-End Tests & Single Binary Build Verification

**Files:**
- Create: `cmd/server/main.go`
- Create: `test/e2e_test.go`
- Create: `README.md`
- Delete: `test.jpg` (cleanup)
- Delete: `upload_job.py` (cleanup)

**Interfaces:**
- Consumes: All internal packages (`internal/database`, `internal/repository`, `internal/server`, `internal/handlers`)
- Produces: Compiled executable binary (`bin/iqbal-portfolio.exe` / `bin/iqbal-portfolio`)

- [ ] **Step 1: Write E2E integration test in `test/e2e_test.go`**

```go
package test_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"iqbal-portfolio/internal/database"
	"iqbal-portfolio/internal/repository"
	"iqbal-portfolio/internal/server"
)

func TestEndToEndVisitorFlow(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()
	database.Migrate(db)
	repository.SeedData(db)

	srv, err := server.NewServer(db)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	ts := httptest.NewServer(srv)
	defer ts.Close()

	client := ts.Client()

	// 1. Visit Home
	res, err := client.Get(ts.URL + "/")
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("home request failed: %v", err)
	}

	// 2. Select Persona
	form := url.Values{"persona": {"founder"}}
	res, err = client.PostForm(ts.URL+"/gateway/select", form)
	if err != nil || res.StatusCode != http.StatusOK { // Follows redirect to /
		t.Fatalf("gateway select failed: %v", err)
	}

	// 3. Open Modal
	res, err = client.Get(ts.URL + "/projects/nexus-ai/modal")
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("modal request failed: %v", err)
	}

	// 4. Submit Valid Inquiry
	inqForm := url.Values{
		"name":            {"Recruiter Alex"},
		"email":           {"alex@acme.corp"},
		"engagement_type": {"staff"},
		"message":         {"Let's build high-throughput systems together."},
	}
	res, err = client.PostForm(ts.URL+"/inquiries", inqForm)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("inquiry submission failed: %v", err)
	}
}
```

- [ ] **Step 2: Implement `cmd/server/main.go` and `README.md`**

- `cmd/server/main.go`:
  - Reads `PORT` env var (default: `8080`).
  - Reads `DB_PATH` env var (default: `portfolio.db`).
  - Initializes DB with `modernc.org/sqlite`, runs migrations, runs seeds.
  - Starts HTTP server with graceful shutdown listening to `os.Interrupt` and `syscall.SIGTERM`.
- `README.md`: Document build instructions (`go build -o bin/iqbal-portfolio.exe ./cmd/server`), environment variables, features, and Figma design alignment.
- Clean up obsolete temporary files: `test.jpg` and `upload_job.py`.

- [ ] **Step 3: Run all unit and integration tests**

Run: `go test ./... -v`
Expected: PASS across all packages.

- [ ] **Step 4: Build single standalone executable binary**

Run: `go build -o bin/iqbal-portfolio.exe ./cmd/server`
Expected: Zero compilation errors, single executable generated in `bin/`.

- [ ] **Step 5: Commit final deliverable**

```bash
git add cmd/ test/ README.md
git rm -f test.jpg upload_job.py 2>$null
git commit -m "feat: complete iqbal portfolio platform implementation with tests and build artifact"
```
