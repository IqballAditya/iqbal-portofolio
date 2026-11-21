package repository

import (
	"database/sql"
	"encoding/json"
)

type seedProject struct {
	id             string
	title          string
	category       string
	tagline        string
	synopsis       string
	techStack      []string
	targetPersonas []string
	matchScore     int
	role           string
	metrics        string
	githubURL      string
	demoURL        string
	imageURL       string
	isFeatured     bool
	sortOrder      int
}

type seedMilestone struct {
	id            int64
	projectID     string
	episodeNumber int
	title         string
	duration      string
	summary       string
	codeSnippet   string
}

type seedSkill struct {
	rank        int
	name        string
	category    string
	description string
	imageURL    string
}

func SeedData(db *sql.DB) error {
	projects := []seedProject{
		{
			id:             "nexus-ai",
			title:          "Nexus AI Agent",
			category:       "ai",
			tagline:        "Autonomous Multi-Agent Workflow Engine for Resilient Cloud Automation",
			synopsis:       "An enterprise-grade autonomous agent runtime capable of dynamic workflow orchestration, self-healing diagnostics, and safe code synthesis in sandboxed cloud environments.",
			techStack:      []string{"Go", "Python", "gRPC", "Docker", "OpenAI API", "PostgreSQL"},
			targetPersonas: []string{"tech_lead", "architect", "engineering_manager"},
			matchScore:     99,
			role:           "Lead Systems Architect & Core Developer",
			metrics:        "Sub-100ms planning latency, 99.4% task completion rate across 50k automated runs",
			githubURL:      "https://github.com/iqbal/nexus-ai",
			demoURL:        "https://nexus-ai.demo.dev",
			imageURL:       "/static/images/projects/nexus-ai.jpg",
			isFeatured:     true,
			sortOrder:      1,
		},
		{
			id:             "storage-engine",
			title:          "Distributed Key-Value Engine",
			category:       "infra",
			tagline:        "LSM-Tree Storage Engine with Raft Consensus and Write-Ahead Logging",
			synopsis:       "High-performance distributed storage node featuring MemTable memory flushing, SSTable compaction, and Raft leader election.",
			techStack:      []string{"Go", "Raft", "LSM-Tree", "gRPC", "Linux"},
			targetPersonas: []string{"tech_lead", "recruiter", "architect"},
			matchScore:     97,
			role:           "Sole Author & Maintainer",
			metrics:        "250k write ops/sec, p99 write latency < 1.2ms on NVMe",
			githubURL:      "https://github.com/iqbal/storage-engine",
			demoURL:        "https://storage.demo.dev",
			imageURL:       "/static/images/projects/storage-engine.jpg",
			isFeatured:     false,
			sortOrder:      2,
		},
		{
			id:             "cost-orchestrator",
			title:          "Cloud Cost Optimizer",
			category:       "infra",
			tagline:        "Intelligent Multi-Cloud Spot Instance Arbiter and Workload Migrator",
			synopsis:       "Heuristic scheduler that predicts spot interruption probability and gracefully drains microservices across AWS and GCP regions.",
			techStack:      []string{"Go", "Kubernetes", "AWS", "Prometheus", "Terraform"},
			targetPersonas: []string{"engineering_manager", "tech_lead"},
			matchScore:     95,
			role:           "Cloud Infrastructure Engineer",
			metrics:        "Saved 42% on monthly Kubernetes compute bill with 99.99% availability",
			githubURL:      "https://github.com/iqbal/cost-orchestrator",
			demoURL:        "https://cost.demo.dev",
			imageURL:       "/static/images/projects/cost-orchestrator.jpg",
			isFeatured:     false,
			sortOrder:      3,
		},
		{
			id:             "stream-analytics",
			title:          "Real-Time Event Broker",
			category:       "backend",
			tagline:        "Sub-Millisecond Pub/Sub Event Pipeline with Partition Filtering",
			synopsis:       "Custom zero-allocation event routing daemon handling high-throughput telemetry streams with persistent ring buffers.",
			techStack:      []string{"Go", "Kafka", "Redis", "eBPF", "Prometheus"},
			targetPersonas: []string{"tech_lead", "architect"},
			matchScore:     96,
			role:           "Backend Systems Engineer",
			metrics:        "1.2M events/sec sustained throughput at 40MB total memory footprint",
			githubURL:      "https://github.com/iqbal/stream-analytics",
			demoURL:        "https://stream.demo.dev",
			imageURL:       "/static/images/projects/stream-analytics.jpg",
			isFeatured:     false,
			sortOrder:      4,
		},
		{
			id:             "grpc-gateway",
			title:          "High-Throughput API Gateway",
			category:       "backend",
			tagline:        "Edge Reverse Proxy with Dynamic Service Discovery and Token Bucket Rate Limiting",
			synopsis:       "Lightweight reverse proxy translating HTTP/2 JSON to upstream binary gRPC with distributed Redis rate limiting and OpenTelemetry tracing.",
			techStack:      []string{"Go", "gRPC", "Redis", "OpenTelemetry", "Docker"},
			targetPersonas: []string{"tech_lead", "recruiter"},
			matchScore:     94,
			role:           "Platform Engineer",
			metrics:        "Handles 180k req/sec with < 0.8ms added proxy overhead",
			githubURL:      "https://github.com/iqbal/grpc-gateway",
			demoURL:        "https://gateway.demo.dev",
			imageURL:       "/static/images/projects/grpc-gateway.jpg",
			isFeatured:     false,
			sortOrder:      5,
		},
		{
			id:             "fintech-ledger",
			title:          "Immutable Financial Ledger",
			category:       "backend",
			tagline:        "Double-Entry Bookkeeping Engine with Cryptographic Audit Trails",
			synopsis:       "ACID-compliant double-entry ledger implementing SHA-256 cryptographic chain proofs, multi-currency balance enforcement, and strict balance invariant verification.",
			techStack:      []string{"Go", "PostgreSQL", "Docker", "REST", "React"},
			targetPersonas: []string{"tech_lead", "engineering_manager", "recruiter"},
			matchScore:     98,
			role:           "Lead Backend Engineer",
			metrics:        "Zero ledger discrepancies across 10M+ test transactions, audited for SOC 2 Type II",
			githubURL:      "https://github.com/iqbal/fintech-ledger",
			demoURL:        "https://ledger.demo.dev",
			imageURL:       "/static/images/projects/fintech-ledger.jpg",
			isFeatured:     false,
			sortOrder:      6,
		},
	}

	milestones := []seedMilestone{
		{
			id:            1,
			projectID:     "nexus-ai",
			episodeNumber: 1,
			title:         "S1:E1 - Blueprint & Control Plane",
			duration:      "45m",
			summary:       "Architecting the distributed worker protocol, consensus-driven task state machine, and zero-trust sandbox execution runtime.",
			codeSnippet:   "type WorkerPool struct {\n\tmu sync.RWMutex\n\tworkers map[string]*AgentWorker\n\ttasks chan Task\n}",
		},
		{
			id:            2,
			projectID:     "nexus-ai",
			episodeNumber: 2,
			title:         "S1:E2 - Autonomous Graph Execution",
			duration:      "52m",
			summary:       "Implementing directed acyclic graph (DAG) dependency resolution with real-time backpressure monitoring and distributed failure recovery.",
			codeSnippet:   "func (e *Engine) ExecuteDAG(ctx context.Context, dag *WorkflowGraph) error {\n\treturn dag.WalkConcurrent(ctx, e.workerPool)\n}",
		},
		{
			id:            3,
			projectID:     "nexus-ai",
			episodeNumber: 3,
			title:         "S1:E3 - Hardening & Chaos Testing",
			duration:      "38m",
			summary:       "Subjecting the agent mesh to network partitions, simulated API throttling, and Byzantine worker faults with zero data loss.",
			codeSnippet:   "func TestByzantineFaultTolerance(t *testing.T) {\n\tcluster := spawnMockCluster(5)\n\tcluster.PartitionWorker(2)\n\tassert.NoError(t, cluster.VerifyQuorum())\n}",
		},
	}

	skills := []seedSkill{
		{
			rank:        1,
			name:        "Go",
			category:    "Language / Systems",
			description: "Concurrency, low-latency microservices, memory optimization, channels, stdlib internals",
			imageURL:    "/static/images/skills/go.svg",
		},
		{
			rank:        2,
			name:        "Rust",
			category:    "Language / Systems",
			description: "Memory safety without GC, async Tokio, low-level systems programming, high performance",
			imageURL:    "/static/images/skills/rust.svg",
		},
		{
			rank:        3,
			name:        "PostgreSQL",
			category:    "Database",
			description: "Schema design, index tuning, query plan analysis, ACID transactions, WAL replication",
			imageURL:    "/static/images/skills/postgres.svg",
		},
		{
			rank:        4,
			name:        "Kubernetes",
			category:    "Infrastructure",
			description: "Container orchestration, custom operators, Helm charts, ingress controllers, service mesh",
			imageURL:    "/static/images/skills/k8s.svg",
		},
		{
			rank:        5,
			name:        "Redis",
			category:    "Database / Cache",
			description: "In-memory caching, distributed locks (Redlock), pub/sub message brokers, sorted sets",
			imageURL:    "/static/images/skills/redis.svg",
		},
		{
			rank:        6,
			name:        "Docker",
			category:    "DevOps",
			description: "Multi-stage production container builds, minimal scratch images, container security",
			imageURL:    "/static/images/skills/docker.svg",
		},
		{
			rank:        7,
			name:        "gRPC",
			category:    "Networking",
			description: "Protobuf schema definition, HTTP/2 multiplexing, bi-directional streaming, interceptors",
			imageURL:    "/static/images/skills/grpc.svg",
		},
		{
			rank:        8,
			name:        "TypeScript",
			category:    "Frontend / Fullstack",
			description: "Strict type safety, modern web applications, Next.js, component architecture",
			imageURL:    "/static/images/skills/typescript.svg",
		},
		{
			rank:        9,
			name:        "Kafka",
			category:    "Distributed Systems",
			description: "Event-driven architecture, partition rebalancing, consumer group scaling, compaction",
			imageURL:    "/static/images/skills/kafka.svg",
		},
		{
			rank:        10,
			name:        "AWS",
			category:    "Cloud",
			description: "ECS, EKS, RDS, DynamoDB, S3, IAM security policies, multi-AZ high availability",
			imageURL:    "/static/images/skills/aws.svg",
		},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	projectStmt, err := tx.Prepare(`INSERT OR REPLACE INTO projects (
		id, title, category, tagline, synopsis, tech_stack, target_personas,
		match_score, role, metrics, github_url, demo_url, image_url, is_featured, sort_order
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer projectStmt.Close()

	for _, p := range projects {
		techJSON, _ := json.Marshal(p.techStack)
		personasJSON, _ := json.Marshal(p.targetPersonas)
		_, err := projectStmt.Exec(
			p.id, p.title, p.category, p.tagline, p.synopsis,
			string(techJSON), string(personasJSON), p.matchScore,
			p.role, p.metrics, p.githubURL, p.demoURL, p.imageURL,
			p.isFeatured, p.sortOrder,
		)
		if err != nil {
			return err
		}
	}

	milestoneStmt, err := tx.Prepare(`INSERT OR REPLACE INTO milestones (
		id, project_id, episode_number, title, duration, summary, code_snippet
	) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer milestoneStmt.Close()

	for _, m := range milestones {
		_, err := milestoneStmt.Exec(
			m.id, m.projectID, m.episodeNumber, m.title, m.duration, m.summary, m.codeSnippet,
		)
		if err != nil {
			return err
		}
	}

	skillStmt, err := tx.Prepare(`INSERT OR REPLACE INTO skills (
		rank, name, category, description, image_url
	) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer skillStmt.Close()

	for _, s := range skills {
		_, err := skillStmt.Exec(s.rank, s.name, s.category, s.description, s.imageURL)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
