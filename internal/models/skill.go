package models

type Skill struct {
	Rank        int    `json:"rank"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}
