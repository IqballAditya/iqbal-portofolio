package handlers

import (
	"net/http"

	"iqbal-portfolio/internal/models"
	"iqbal-portfolio/internal/repository"
)

type Renderer interface {
	Render(w http.ResponseWriter, name string, data any) error
}

type Handler struct {
	projects *repository.ProjectRepo
	skills   *repository.SkillRepo
	renderer Renderer
}

func New(projects *repository.ProjectRepo, skills *repository.SkillRepo, renderer Renderer) *Handler {
	return &Handler{
		projects: projects,
		skills:   skills,
		renderer: renderer,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	personaKey := "guest"
	if cookie, err := r.Cookie("iqbal_persona"); err == nil && cookie != nil {
		personaKey = ValidatePersona(cookie.Value)
	}
	persona := GetPersona(personaKey)

	projects, err := h.projects.ListProjects(persona.Key)
	if err != nil {
		http.Error(w, "Failed to load projects", http.StatusInternalServerError)
		return
	}

	skills, err := h.skills.ListSkills()
	if err != nil {
		http.Error(w, "Failed to load skills", http.StatusInternalServerError)
		return
	}

	var hero *models.Project
	for i := range projects {
		if projects[i].ID == "nexus-ai" {
			hero = &projects[i]
			break
		}
	}
	if hero == nil && len(projects) > 0 {
		hero = &projects[0]
	}

	var caseStudies []models.Project
	for _, p := range projects {
		if p.Metrics != "" {
			caseStudies = append(caseStudies, p)
		}
	}

	data := struct {
		Title          string
		HideNav        bool
		HideFooter     bool
		PersonaName    string
		PersonaColor   string
		PersonaInitial string
		ShortlistCount int
		Hero           *models.Project
		Projects       []models.Project
		TopSkills      []models.Skill
		CaseStudies    []models.Project
	}{
		Title:          "Home - Netflix Showcase",
		HideNav:        false,
		HideFooter:     false,
		PersonaName:    persona.Name,
		PersonaColor:   persona.Color,
		PersonaInitial: persona.Initial,
		ShortlistCount: 0,
		Hero:           hero,
		Projects:       projects,
		TopSkills:      skills,
		CaseStudies:    caseStudies,
	}

	if err := h.renderer.Render(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
