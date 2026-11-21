package repository

import (
	"database/sql"
	"encoding/json"
	"strings"

	"iqbal-portfolio/internal/models"
)

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) ListProjects(persona string) ([]models.Project, error) {
	var query string
	var args []any

	if persona != "" {
		query = `SELECT id, title, category, tagline, synopsis, tech_stack, target_personas, match_score, role, metrics, github_url, demo_url, image_url, is_featured, sort_order
			FROM projects
			ORDER BY CASE WHEN target_personas LIKE '%' || ? || '%' THEN 0 ELSE 1 END, sort_order ASC`
		args = append(args, persona)
	} else {
		query = `SELECT id, title, category, tagline, synopsis, tech_stack, target_personas, match_score, role, metrics, github_url, demo_url, image_url, is_featured, sort_order
			FROM projects
			ORDER BY sort_order ASC`
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProjects(rows)
}

func (r *ProjectRepo) GetProjectByID(id string) (*models.Project, error) {
	query := `SELECT id, title, category, tagline, synopsis, tech_stack, target_personas, match_score, role, metrics, github_url, demo_url, image_url, is_featured, sort_order
		FROM projects WHERE id = ?`
	row := r.db.QueryRow(query, id)

	p, err := scanProjectRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	milestones, err := r.ListMilestones(id)
	if err != nil {
		return nil, err
	}
	p.Milestones = milestones

	return p, nil
}

func (r *ProjectRepo) SearchProjects(query string, category string) ([]models.Project, error) {
	baseQuery := `SELECT id, title, category, tagline, synopsis, tech_stack, target_personas, match_score, role, metrics, github_url, demo_url, image_url, is_featured, sort_order
		FROM projects WHERE 1=1`
	var conditions []string
	var args []any

	if category != "" && category != "all" {
		conditions = append(conditions, "category = ?")
		args = append(args, category)
	}

	if query != "" {
		conditions = append(conditions, "(title LIKE ? OR tagline LIKE ? OR synopsis LIKE ? OR tech_stack LIKE ?)")
		pattern := "%" + query + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}
	baseQuery += " ORDER BY sort_order ASC"

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProjects(rows)
}

func (r *ProjectRepo) ListMilestones(projectID string) ([]models.Milestone, error) {
	query := `SELECT id, project_id, episode_number, title, duration, summary, code_snippet
		FROM milestones WHERE project_id = ? ORDER BY episode_number ASC`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var milestones []models.Milestone
	for rows.Next() {
		var m models.Milestone
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.EpisodeNumber, &m.Title, &m.Duration, &m.Summary, &m.CodeSnippet); err != nil {
			return nil, err
		}
		milestones = append(milestones, m)
	}
	return milestones, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanProjectRow(s rowScanner) (*models.Project, error) {
	var p models.Project
	var techStackRaw, personasRaw string
	var demoURL sql.NullString

	err := s.Scan(
		&p.ID,
		&p.Title,
		&p.Category,
		&p.Tagline,
		&p.Synopsis,
		&techStackRaw,
		&personasRaw,
		&p.MatchScore,
		&p.Role,
		&p.Metrics,
		&p.GithubURL,
		&demoURL,
		&p.ImageURL,
		&p.IsFeatured,
		&p.SortOrder,
	)
	if err != nil {
		return nil, err
	}
	p.DemoURL = demoURL.String

	parseStringSlice(techStackRaw, &p.TechStack)
	parseStringSlice(personasRaw, &p.TargetPersonas)
	return &p, nil
}

func scanProjects(rows *sql.Rows) ([]models.Project, error) {
	var projects []models.Project
	for rows.Next() {
		p, err := scanProjectRow(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *p)
	}
	return projects, rows.Err()
}

func parseStringSlice(raw string, target *[]string) {
	if raw == "" {
		*target = []string{}
		return
	}
	if strings.HasPrefix(raw, "[") {
		if err := json.Unmarshal([]byte(raw), target); err == nil {
			return
		}
	}
	parts := strings.Split(raw, ",")
	res := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}
	*target = res
}
