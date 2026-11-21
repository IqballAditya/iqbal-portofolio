package repository

import (
	"database/sql"
	"time"

	"iqbal-portfolio/internal/models"
)

type InquiryRepo struct {
	db *sql.DB
}

func NewInquiryRepo(db *sql.DB) *InquiryRepo {
	return &InquiryRepo{db: db}
}

func (r *InquiryRepo) CreateInquiry(inq *models.Inquiry) error {
	if inq.CreatedAt.IsZero() {
		inq.CreatedAt = time.Now().UTC()
	}

	query := `INSERT INTO inquiries (name, email, engagement_type, message, created_at) VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.Exec(query, inq.Name, inq.Email, inq.EngagementType, inq.Message, inq.CreatedAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	inq.ID = id
	return nil
}
