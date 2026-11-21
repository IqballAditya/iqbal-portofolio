# ADITH AI STUDIO Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build ADITH AI STUDIO, a high-performance Web-based AI Coding Studio featuring an AI Chat Copilot with reasoning traces, Monaco Code Editor with interactive diff accept/reject, Tree-sitter AST parsing, Qdrant vector memory, and a Golang API gateway.

**Architecture:** Polyglot clean architecture where a Golang API Gateway manages sessions, workspace files, and SSE streaming, communicating via HTTP/JSON with a Python Agent Engine that orchestrates Tree-sitter AST parsing and Qdrant vector embeddings, all surfaced through a responsive TypeScript Web IDE.

**Tech Stack:** Go 1.22+, SQLite (WAL mode), Python 3.10+, FastAPI, Tree-sitter, Qdrant-client (Embedded), FastEmbed, TypeScript, Monaco Editor, Tailwind CSS.

**Spec:** `docs/superpowers/specs/2026-11-22-adith-ai-studio-design.md`

## Global Constraints
- Target location: Standalone directory `C:\Users\iqbal\adith-ai-studio`.
- Zero Docker requirement during local development (Option A).
- Inter-service communication via local HTTP REST and SSE (`http://localhost:8090` for Gateway, `http://localhost:5050` for Agent).
- User-in-the-loop code mutation: no files modified on disk without explicit user acceptance.

## Review Focus
- File locked on Windows during diff application: Return clean 409 Conflict with retry suggestion.
- Malformed syntax in incomplete source files: Tree-sitter gracefully yields partial AST nodes without panicking.
- Python agent engine offline or unresponsive: Go Gateway returns 503 with SSE reconnection backoff.
- Large codebase indexing memory spike: Batch indexing capped at 100 symbols per Qdrant upsert.
- UTF-8 characters and ANSI color sequences: PTY websocket correctly forwards streaming terminal output.

---

### Task 1: Project Root Scaffolding & Go Gateway Core

**Files:**
- Create: `C:\Users\iqbal\adith-ai-studio\go.mod`
- Create: `C:\Users\iqbal\adith-ai-studio\internal\database\db.go`
- Create: `C:\Users\iqbal\adith-ai-studio\internal\database\migrations.go`
- Test: `C:\Users\iqbal\adith-ai-studio\internal\database\db_test.go`

**Interfaces:**
- Consumes: None
- Produces: `database.InitDB(dbPath string) (*sql.DB, error)`, `database.Migrate(db *sql.DB) error`

- [ ] **Step 1: Initialize Go module in `adith-ai-studio`**
```bash
mkdir -p C:\Users\iqbal\adith-ai-studio
cd C:\Users\iqbal\adith-ai-studio
go mod init adith-ai-studio
go get modernc.org/sqlite
```

- [ ] **Step 2: Write the failing DB test**
```go
package database_test

import (
    "path/filepath"
    "testing"
    "adith-ai-studio/internal/database"
)

func TestInitDBAndMigrate(t *testing.T) {
    tmp := filepath.Join(t.TempDir(), "test.db")
    db, err := database.InitDB(tmp)
    if err != nil {
        t.Fatalf("InitDB failed: %v", err)
    }
    defer db.Close()

    if err := database.Migrate(db); err != nil {
        t.Fatalf("Migrate failed: %v", err)
    }
}
```

- [ ] **Step 3: Run test to verify it fails**
Run: `go test -v ./internal/database`
Expected: FAIL with undefined packages/functions.

- [ ] **Step 4: Implement `InitDB` and `Migrate` in `internal/database/`**
Implement SQLite connection with `_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)` and tables for `workspaces`, `chat_threads`, `messages`, `tool_executions`, and `codebase_files`.

- [ ] **Step 5: Run test to verify it passes**
Run: `go test -v ./internal/database`
Expected: PASS.

- [ ] **Step 6: Commit**
```bash
git -C C:\Users\iqbal\adith-ai-studio init
git -C C:\Users\iqbal\adith-ai-studio add .
git -C C:\Users\iqbal\adith-ai-studio commit -m "feat(gateway): init project and sqlite schema"
```

---

### Task 2: Workspace File Tree & Diff Patcher in Go

**Files:**
- Create: `C:\Users\iqbal\adith-ai-studio\internal\workspace\fs.go`
- Create: `C:\Users\iqbal\adith-ai-studio\internal\workspace\diff.go`
- Test: `C:\Users\iqbal\adith-ai-studio\internal\workspace\fs_test.go`
- Test: `C:\Users\iqbal\adith-ai-studio\internal\workspace\diff_test.go`

**Interfaces:**
- Consumes: None
- Produces: `workspace.ListFileTree(root string) ([]FileNode, error)`, `workspace.ApplyDiff(root string, path string, original string, modified string) error`

- [ ] **Step 1: Write failing tests for FileTree and Diff**
Test recursive directory walking, ignoring `.git` / `node_modules` / `venv`, and safe atomic file patching with backup.

- [ ] **Step 2: Run test to verify it fails**
Run: `go test -v ./internal/workspace`
Expected: FAIL.

- [ ] **Step 3: Implement `fs.go` and `diff.go`**
Walk filesystem cleanly; provide safe replacement with rollback if write fails.

- [ ] **Step 4: Run test to verify it passes**
Run: `go test -v ./internal/workspace`
Expected: PASS.

- [ ] **Step 5: Commit**
```bash
git -C C:\Users\iqbal\adith-ai-studio add internal/workspace
git -C C:\Users\iqbal\adith-ai-studio commit -m "feat(workspace): add file tree and diff application"
```

---

### Task 3: Python Agent Engine & Tree-sitter AST Parser

**Files:**
- Create: `C:\Users\iqbal\adith-ai-studio\agent\requirements.txt`
- Create: `C:\Users\iqbal\adith-ai-studio\agent\parser\ast_parser.py`
- Test: `C:\Users\iqbal\adith-ai-studio\agent\tests\test_ast_parser.py`

**Interfaces:**
- Consumes: None
- Produces: `ASTParser.parse_code(code: str, language: str) -> list[CodeSymbol]`

- [ ] **Step 1: Setup Python virtualenv and requirements**
```bash
cd C:\Users\iqbal\adith-ai-studio\agent
python -m venv venv
.\venv\Scripts\pip install fastapi uvicorn qdrant-client tree-sitter tree-sitter-languages fastembed pytest
```

- [ ] **Step 2: Write failing test `test_ast_parser.py`**
Test parsing Go and Python snippets extracting functions and structs with start/end line numbers.

- [ ] **Step 3: Run test to verify it fails**
Run: `.\venv\Scripts\pytest tests/test_ast_parser.py`
Expected: FAIL.

- [ ] **Step 4: Implement `ast_parser.py` using `tree-sitter-languages`**
Extract `function_declaration`, `method_declaration`, `type_declaration` into structured `CodeSymbol` dataclass.

- [ ] **Step 5: Run test to verify it passes**
Run: `.\venv\Scripts\pytest tests/test_ast_parser.py`
Expected: PASS.

- [ ] **Step 6: Commit**
```bash
git -C C:\Users\iqbal\adith-ai-studio add agent
git -C C:\Users\iqbal\adith-ai-studio commit -m "feat(agent): add tree-sitter ast parser"
```

---

### Task 4: Qdrant Vector Memory & Code Indexer

**Files:**
- Create: `C:\Users\iqbal\adith-ai-studio\agent\memory\vector_store.py`
- Create: `C:\Users\iqbal\adith-ai-studio\agent\indexer.py`
- Test: `C:\Users\iqbal\adith-ai-studio\agent\tests\test_vector_store.py`

**Interfaces:**
- Consumes: `ASTParser` from Task 3
- Produces: `VectorStore.index_symbols(workspace_id: str, symbols: list[CodeSymbol])`, `VectorStore.search(workspace_id: str, query: str, limit: int = 5) -> list[SearchResult]`

- [ ] **Step 1: Write failing test `test_vector_store.py`**
Verify embedding generation with FastEmbed (`bge-small-en-v1.5`) and local Qdrant collection upsert and hybrid retrieval.

- [ ] **Step 2: Run test to verify it fails**
Run: `.\venv\Scripts\pytest tests/test_vector_store.py`
Expected: FAIL.

- [ ] **Step 3: Implement `vector_store.py`**
Initialize embedded Qdrant client (`QdrantClient(path="./qdrant_storage")`), auto-create `codebase_chunks` collection with cosine distance, and perform semantic search.

- [ ] **Step 4: Run test to verify it passes**
Run: `.\venv\Scripts\pytest tests/test_vector_store.py`
Expected: PASS.

- [ ] **Step 5: Commit**
```bash
git -C C:\Users\iqbal\adith-ai-studio add agent/memory agent/indexer.py agent/tests
git -C C:\Users\iqbal\adith-ai-studio commit -m "feat(memory): add embedded qdrant vector store"
```

---

### Task 5: Agent Reasoning Loop & SSE Streaming Service

**Files:**
- Create: `C:\Users\iqbal\adith-ai-studio\agent\orchestrator.py`
- Create: `C:\Users\iqbal\adith-ai-studio\agent\server.py`
- Test: `C:\Users\iqbal\adith-ai-studio\agent\tests\test_orchestrator.py`

**Interfaces:**
- Consumes: `ASTParser`, `VectorStore`
- Produces: `POST /v1/agent/chat` (streaming SSE with `thought` chunks and `diff` proposals)

- [ ] **Step 1: Write test for streaming agent responses**
Mock LLM call to verify SSE format emitting `event: reasoning`, `event: delta`, and `event: diff`.

- [ ] **Step 2: Run test to verify it fails**
Run: `.\venv\Scripts\pytest tests/test_orchestrator.py`
Expected: FAIL.

- [ ] **Step 3: Implement FastAPI app in `agent/server.py` and `orchestrator.py`**
Port 5050; handle prompt, query Qdrant context, generate chain-of-thought, format unified diffs.

- [ ] **Step 4: Run test to verify it passes**
Run: `.\venv\Scripts\pytest tests/test_orchestrator.py`
Expected: PASS.

- [ ] **Step 5: Commit**
```bash
git -C C:\Users\iqbal\adith-ai-studio add agent/orchestrator.py agent/server.py
git -C C:\Users\iqbal\adith-ai-studio commit -m "feat(agent): add agent reasoning loop and fastapi sse stream"
```

---

### Task 6: Go Gateway API, SSE Bridge & PTY Terminal

**Files:**
- Create: `C:\Users\iqbal\adith-ai-studio\internal\api\routes.go`
- Create: `C:\Users\iqbal\adith-ai-studio\internal\api\terminal.go`
- Create: `C:\Users\iqbal\adith-ai-studio\cmd\gateway\main.go`
- Test: `C:\Users\iqbal\adith-ai-studio\internal\api\routes_test.go`

**Interfaces:**
- Consumes: `internal/database`, `internal/workspace`, Python Agent at `:5050`
- Produces: Public HTTP & SSE endpoints on `:8090` (`GET /api/tree`, `POST /api/chat`, `POST /api/diff/accept`, `WS /api/terminal`)

- [ ] **Step 1: Write failing integration test `routes_test.go`**
Test tree listing, diff acceptance, and SSE stream proxying.

- [ ] **Step 2: Run test to verify it fails**
Run: `go test -v ./internal/api`
Expected: FAIL.

- [ ] **Step 3: Implement Gateway HTTP router and terminal WebSocket bridge**
Serve static assets from `web/dist`, proxy chat to Python `:5050`, expose PTY terminal.

- [ ] **Step 4: Run test to verify it passes**
Run: `go test -v ./internal/api`
Expected: PASS.

- [ ] **Step 5: Commit**
```bash
git -C C:\Users\iqbal\adith-ai-studio add cmd/ internal/api
git -C C:\Users\iqbal\adith-ai-studio commit -m "feat(gateway): add routes, sse bridge and terminal websocket"
```

---

### Task 7: TypeScript Web IDE Frontend & Monaco Integration

**Files:**
- Create: `C:\Users\iqbal\adith-ai-studio\web\package.json`
- Create: `C:\Users\iqbal\adith-ai-studio\web\src\App.tsx`
- Create: `C:\Users\iqbal\adith-ai-studio\web\src\components\ChatPanel.tsx`
- Create: `C:\Users\iqbal\adith-ai-studio\web\src\components\CodeEditor.tsx`
- Create: `C:\Users\iqbal\adith-ai-studio\web\src\components\TerminalView.tsx`

**Interfaces:**
- Consumes: Gateway APIs on `:8090`
- Produces: Production build in `web/dist/` embedded into Go binary.

- [ ] **Step 1: Scaffold Vite + React + Tailwind + Monaco Editor**
```bash
cd C:\Users\iqbal\adith-ai-studio\web
npm create vite@latest . -- --template react-ts
npm install @monaco-editor/react lucide-react xterm xterm-addon-fit tailwindcss
```

- [ ] **Step 2: Build ChatPanel with SSE listener and Diff accept card**
Listen to `EventSource`, display reasoning process, and emit `POST /api/diff/accept` on click.

- [ ] **Step 3: Build CodeEditor with Monaco multi-tab switching**
Load files via `/api/file?path=...`, update editor buffer dynamically upon diff acceptance.

- [ ] **Step 4: Build TerminalView with xterm.js over WebSocket**
Connect to `ws://localhost:8090/api/terminal`.

- [ ] **Step 5: Verify build passes**
Run: `npm run build`
Expected: Clean build without TypeScript errors.

- [ ] **Step 6: Commit**
```bash
git -C C:\Users\iqbal\adith-ai-studio add web
git -C C:\Users\iqbal\adith-ai-studio commit -m "feat(web): add react monaco ide frontend"
```

---

### Task 8: Unified Dev Bootstrapper & End-to-End Verification

**Files:**
- Create: `C:\Users\iqbal\adith-ai-studio\run-dev.ps1`
- Create: `C:\Users\iqbal\adith-ai-studio\test\e2e_smoke_test.go`
- Create: `C:\Users\iqbal\adith-ai-studio\README.md`

- [ ] **Step 1: Create `run-dev.ps1`**
Launches Python Agent (`python agent/server.py`), builds and starts Go Gateway (`go run cmd/gateway/main.go`), and opens default browser to `http://localhost:8090`.

- [ ] **Step 2: Write `e2e_smoke_test.go`**
Simulate complete flow: send user prompt -> receive reasoning -> receive diff -> call accept diff -> verify file updated on disk.

- [ ] **Step 3: Run E2E smoke test**
Run: `go test -v ./test/...`
Expected: PASS.

- [ ] **Step 4: Commit**
```bash
git -C C:\Users\iqbal\adith-ai-studio add run-dev.ps1 test README.md
git -C C:\Users\iqbal\adith-ai-studio commit -m "chore: add dev runner and e2e smoke test"
```
