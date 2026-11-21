package models

import "time"

type Inquiry struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	EngagementType string    `json:"engagement_type"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
}
