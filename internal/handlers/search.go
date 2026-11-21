package handlers

import (
	"net/http"
	"strings"

	"iqbal-portfolio/internal/models"
)

type SearchPageData struct {
	Title          string
	HideNav        bool
	HideFooter     bool
	PersonaName    string
	PersonaColor   string
	PersonaInitial string
	ShortlistCount int
	Query          string
	Category       string
	Projects       []models.Project
}

type SearchResultsData struct {
	Projects []models.Project
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	personaKey := "guest"
	if cookie, err := r.Cookie("iqbal_persona"); err == nil && cookie != nil {
		personaKey = ValidatePersona(cookie.Value)
	}
	persona := GetPersona(personaKey)

	q := strings.TrimSpace(r.URL.Query().Get("q"))
	category := strings.TrimSpace(r.URL.Query().Get("category"))

	projects, err := h.projects.SearchProjects(q, category)
	if err != nil {
		http.Error(w, "Failed to search projects", http.StatusInternalServerError)
		return
	}

	data := SearchPageData{
		Title:          "Search Projects - IQBAL",
		HideNav:        false,
		HideFooter:     false,
		PersonaName:    persona.Name,
		PersonaColor:   persona.Color,
		PersonaInitial: persona.Initial,
		ShortlistCount: 0,
		Query:          q,
		Category:       category,
		Projects:       projects,
	}

	if err := h.renderer.Render(w, "search.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) SearchAPI(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	category := strings.TrimSpace(r.URL.Query().Get("category"))

	projects, err := h.projects.SearchProjects(q, category)
	if err != nil {
		http.Error(w, "Failed to execute search", http.StatusInternalServerError)
		return
	}

	data := SearchResultsData{
		Projects: projects,
	}

	if err := h.renderer.Render(w, "search_results.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
