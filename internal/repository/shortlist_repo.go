package repository

import (
	"database/sql"
	"time"

	"iqbal-portfolio/internal/models"
)

type ShortlistRepo struct {
	db *sql.DB
}

func NewShortlistRepo(db *sql.DB) *ShortlistRepo {
	return &ShortlistRepo{db: db}
}

func (r *ShortlistRepo) ToggleShortlist(sessionID, projectID string) (bool, error) {
	isShort, err := r.IsShortlisted(sessionID, projectID)
	if err != nil {
		return false, err
	}

	if isShort {
		_, err := r.db.Exec("DELETE FROM shortlists WHERE session_id = ? AND project_id = ?", sessionID, projectID)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	_, err = r.db.Exec("INSERT INTO shortlists (session_id, project_id, created_at) VALUES (?, ?, ?)", sessionID, projectID, time.Now().UTC())
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ShortlistRepo) ListShortlist(sessionID string) ([]models.Project, error) {
	query := `SELECT p.id, p.title, p.category, p.tagline, p.synopsis, p.tech_stack, p.target_personas, p.match_score, p.role, p.metrics, p.github_url, p.demo_url, p.image_url, p.is_featured, p.sort_order
		FROM projects p
		INNER JOIN shortlists s ON p.id = s.project_id
		WHERE s.session_id = ?
		ORDER BY s.created_at DESC`

	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProjects(rows)
}

func (r *ShortlistRepo) IsShortlisted(sessionID, projectID string) (bool, error) {
	var exists int
	err := r.db.QueryRow("SELECT 1 FROM shortlists WHERE session_id = ? AND project_id = ? LIMIT 1", sessionID, projectID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
