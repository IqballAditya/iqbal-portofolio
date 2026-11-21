package models

type Project struct {
	ID             string      `json:"id"`
	Title          string      `json:"title"`
	Category       string      `json:"category"`
	Tagline        string      `json:"tagline"`
	Synopsis       string      `json:"synopsis"`
	TechStack      []string    `json:"tech_stack"`
	TargetPersonas []string    `json:"target_personas"`
	MatchScore     int         `json:"match_score"`
	Role           string      `json:"role"`
	Metrics        string      `json:"metrics"`
	GithubURL      string      `json:"github_url"`
	DemoURL        string      `json:"demo_url"`
	ImageURL       string      `json:"image_url"`
	IsFeatured     bool        `json:"is_featured"`
	SortOrder      int         `json:"sort_order"`
	Milestones     []Milestone `json:"milestones,omitempty"`
}

type Milestone struct {
	ID            int64  `json:"id"`
	ProjectID     string `json:"project_id"`
	EpisodeNumber int    `json:"episode_number"`
	Title         string `json:"title"`
	Duration      string `json:"duration"`
	Summary       string `json:"summary"`
	CodeSnippet   string `json:"code_snippet"`
}
