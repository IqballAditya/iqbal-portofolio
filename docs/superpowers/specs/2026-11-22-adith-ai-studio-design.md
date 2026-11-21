# Design Document: ADITH AI STUDIO
**Date:** 2026-11-22  
**Author:** Iqbal Aditya & Antigravity  
**Status:** Ready for Approval (Post-Reviewer Check)  
**Location Target:** `C:\Users\iqbal\adith-ai-studio`

---

## 1. Intent and The Why
Modern software engineering is transitioning from manual coding to agentic AI collaboration. Tools like Cursor, GitHub Copilot, and Claude Code have demonstrated the power of AI-assisted development. However, many existing tools either lock users into heavy desktop IDE forks (Electron) or lack deep contextual understanding of whole-codebase abstract syntax trees (AST).

**ADITH AI STUDIO** is designed as a next-generation, high-performance Web-based AI Coding Studio. It combines:
1. The conversational flexibility and reasoning transparency of frontier chat interfaces (ChatGPT, Claude, DeepSeek-R1).
2. The precision of a full-fledged web code editor (Monaco Editor).
3. Codebase intelligence powered by AST parsing (Tree-sitter) and semantic vector search (Qdrant).
4. An interactive terminal execution environment for automated testing and diff validation.

---

## 2. Goals, Non-Goals, and Anti-Goals

### Goals
- **G1: End-to-End Agentic Loop**: A user issues natural language instructions in the AI chat; the agent inspects the project AST, retrieves relevant symbol chunks via Qdrant, drafts a code patch (diff), and presents it for one-click Accept/Reject.
- **G2: Real-time Streaming & Reasoning**: Display model reasoning traces (e.g. DeepSeek-R1 / Claude thinking blocks) and streaming token responses via Server-Sent Events (SSE).
- **G3: Polyglot Clean Architecture**:
  - Frontend: TypeScript + Modern Web IDE (Tailwind CSS, Monaco Editor).
  - API Gateway: Golang (HTTP/SSE router, session handling, WebSocket/SSE dispatcher).
  - Agent Engine: Python (Agent loop, LLM adapters, MCP tool server).
  - Memory Engine: Tree-sitter (AST parser) + Qdrant (Vector indexing).
- **G4: Zero-Friction Local Development (Option A - Selected)**:
  - Local embedded databases (SQLite in WAL mode + Qdrant embedded storage) requiring zero Docker overhead during early development.
  - Abstraction layers enabling instant upgrade to PostgreSQL + Dockerized Qdrant + Redis when scaling.

### Non-Goals
- **NG1: Pre-training Foundation Models**: We will not train base LLMs from scratch. We leverage frontier LLM APIs (Anthropic Claude 3.5 Sonnet, DeepSeek R1, OpenAI) or local open-weights models via Ollama.
- **NG2: Immediate Multi-Tenant SaaS Billing**: Billing, stripe webhooks, and team organizations are deferred until single-user and local workspace workflows are robust.

### Anti-Goals
- **AG1: Monolithic Bloat**: Avoid dumping the entire project into the portfolio repository (`netflix_iqbal`). The studio lives in its own dedicated repository `adith-ai-studio`.
- **AG2: Desktop Electron Lock-in**: Do not build a heavy, non-portable Electron app first. Build a responsive, fast Web Cloud IDE first; packaging with Tauri or Electron can come later if desired.

---

## 3. Approaches Considered & Trade-Offs

### Approach 1: Heavy Dockerized Enterprise Monolith (Option B)
- **Architecture**: Docker Compose orchestrating PostgreSQL, Qdrant Server, Redis, Go Gateway, Python Agent, and Next.js frontend.
- **Pros**: Matches production cloud topology from day 1.
- **Cons**: High memory consumption (2-4GB background RAM), requires Windows Subsystem for Linux (WSL2) + Docker Desktop installation, friction for rapid prototyping.

### Approach 2: Lightweight Embedded Polyglot Core (Option A) — *SELECTED*
- **Architecture**:
  - Relational: SQLite in WAL mode (designed with PostgreSQL-compatible schema).
  - Vector: Qdrant Embedded file-based client (`QdrantClient(path="./qdrant_storage")`).
  - Caching & Events: Go channels + in-memory pub/sub (clean interface `EventBroker` allowing hot-swap with Redis).
- **Pros**: Instant startup, zero RAM bloat, zero dependencies to install (no Docker required), clean modular codebase.
- **Cons**: Concurrency limited to single-node filesystem (sufficient for local/developer studio).

---

## 4. Architectural Decisions & Builder Clarifications (Post-Review)

### Clarification 1: Inter-Service Communication Protocol & Responsibilities (Builder Call)
- **Protocol**: Standard **HTTP REST + SSE (JSON)** over `http://localhost:5050`.
- **Responsibilities**:
  - **Golang Gateway (Port 8090)**: Owns client sessions, workspace file tree CRUD, SQLite persistence, and public SSE streaming to the browser.
  - **Python Agent (Port 5050)**: Owns LLM interactions, Tree-sitter parsing, Qdrant vector memory, and diff generation. Exposes `/v1/chat/stream` and `/v1/index/codebase` to Go.
- **Why**: Zero code-generation dependencies (avoids gRPC/protoc dual-compiler overhead on Windows), highly inspectable, easy to debug with curl.

### Clarification 2: Terminal Execution Model in Option A (Builder Call)
- In **Option A (No Docker)**, the terminal runs as a directory-scoped local shell session (PowerShell on Windows / Bash on Linux) executed via WebSocket PTY bridge. Process execution has configurable timeouts and path restriction to the current workspace root.
- True hardware microVM sandboxing (Firecracker / Docker container) is designated for **Option B**, unlocked when Docker is introduced later.

### Clarification 3: Diff Validation & Workspace Mutation (Confirmed Design)
- **User-in-the-Loop Safe Flow**:
  1. Agent inspects AST and generates a structured patch (`unified diff`).
  2. The UI renders the diff preview card in the AI chat with syntax highlights.
  3. The user reviews the diff and clicks **"Accept Changes"** (or presses `Tab`).
  4. Only upon acceptance, Go applies the diff in-place to the file on disk and updates the Monaco editor model.
  5. User can then run test commands in the terminal with zero risk of dirtying speculative code.

### Clarification 4: Process Lifecycle & Developer Workflow (Builder Call)
- The services operate as independent decoupled daemons.
- A root PowerShell script `run-dev.ps1` boots both Go Gateway and Python Agent in tandem, ensuring zero coupling in the Go binary itself.

---

## 5. Data Models & Schemas

### Relational Schema (SQLite / PostgreSQL Compatible)
1. `workspaces`: `id`, `name`, `root_path`, `default_branch`, `created_at`.
2. `chat_threads`: `id`, `workspace_id`, `title`, `model_name`, `created_at`.
3. `messages`: `id`, `thread_id`, `role`, `content`, `reasoning_content`, `prompt_tokens`, `completion_tokens`, `created_at`.
4. `tool_executions`: `id`, `message_id`, `tool_name`, `tool_arguments` (JSON), `tool_output`, `status`, `duration_ms`.
5. `codebase_files`: `id`, `workspace_id`, `relative_path`, `file_hash`, `ast_parsed`, `indexed_at`.

### Qdrant Vector Payload Structure (`codebase_chunks`)
- `id`: UUID
- `vector`: Dense embeddings (e.g. `fastembed` BGE-small-en or OpenAI embeddings)
- `payload`:
  - `workspace_id`: string
  - `file_path`: string
  - `language`: string
  - `symbol_name`: string (e.g., `HandleCheckout`)
  - `symbol_type`: string (`function`, `struct`, `class`, `interface`)
  - `start_line`: int
  - `end_line`: int
  - `raw_code`: string

---

## 6. Directory Structure (`C:\Users\iqbal\adith-ai-studio`)

```
adith-ai-studio/
├── cmd/
│   └── gateway/
│       └── main.go                 # Golang Gateway entrypoint
├── internal/
│   ├── api/                        # HTTP & SSE handlers
│   ├── config/                     # Environment configuration
│   ├── database/                   # SQLite database & migrations
│   └── workspace/                  # File system & project workspace ops
├── agent/                          # Python Agent Engine
│   ├── orchestrator.py             # Agent reasoning loop & LLM router
│   ├── parser/                     # Tree-sitter AST extraction
│   ├── memory/                     # Qdrant client & vector indexer
│   ├── tools/                      # File reader, diff patcher
│   └── requirements.txt            # Python dependencies (fastapi, qdrant-client, tree-sitter)
├── web/                            # TypeScript Frontend (Web IDE)
│   ├── index.html                  # Main SPA entrypoint
│   ├── src/
│   │   ├── components/             # ChatPanel, MonacoEditor, TerminalDrawer
│   │   └── api/                    # SSE client & Gateway bindings
│   └── package.json
├── run-dev.ps1                     # One-click unified dev bootstrapper
├── docs/                           # Architectural guides & specs
└── README.md
```

---

## 7. What's Left to The Builder
- Tree-sitter grammar bindings for target languages (Go, Python, TypeScript).
- FastEmbed local embedding model setup for 100% free vector indexing.
- Monaco Editor inline diff decoration rendering.
