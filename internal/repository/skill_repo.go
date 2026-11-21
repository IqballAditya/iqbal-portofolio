package repository

import (
	"database/sql"

	"iqbal-portfolio/internal/models"
)

type SkillRepo struct {
	db *sql.DB
}

func NewSkillRepo(db *sql.DB) *SkillRepo {
	return &SkillRepo{db: db}
}

func (r *SkillRepo) ListSkills() ([]models.Skill, error) {
	query := `SELECT rank, name, category, description, image_url FROM skills ORDER BY rank ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.Skill
	for rows.Next() {
		var s models.Skill
		if err := rows.Scan(&s.Rank, &s.Name, &s.Category, &s.Description, &s.ImageURL); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}
	return skills, rows.Err()
}
