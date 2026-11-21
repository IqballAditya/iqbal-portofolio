CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    title TEXT,
    category TEXT,
    tagline TEXT,
    synopsis TEXT,
    tech_stack TEXT,
    target_personas TEXT,
    match_score INTEGER,
    role TEXT,
    metrics TEXT,
    github_url TEXT,
    demo_url TEXT,
    image_url TEXT,
    is_featured BOOLEAN,
    sort_order INTEGER
);

CREATE TABLE IF NOT EXISTS milestones (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id TEXT REFERENCES projects(id),
    episode_number INTEGER,
    title TEXT,
    duration TEXT,
    summary TEXT,
    code_snippet TEXT
);

CREATE TABLE IF NOT EXISTS skills (
    rank INTEGER PRIMARY KEY,
    name TEXT,
    category TEXT,
    description TEXT,
    image_url TEXT
);

CREATE TABLE IF NOT EXISTS shortlists (
    session_id TEXT,
    project_id TEXT,
    created_at DATETIME,
    PRIMARY KEY (session_id, project_id)
);

CREATE TABLE IF NOT EXISTS inquiries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    email TEXT,
    engagement_type TEXT,
    message TEXT,
    created_at DATETIME
);
