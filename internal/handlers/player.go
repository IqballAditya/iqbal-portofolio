package handlers

import (
	"encoding/json"
	"net/http"

	"iqbal-portfolio/internal/models"
)

type PlayerData struct {
	Title             string
	HideNav           bool
	HideFooter        bool
	HideGatewayHeader bool
	Project           *models.Project
	MilestonesJSON    string
}

func (h *Handler) Player(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	project, err := h.projects.GetProjectByID(id)
	if err != nil {
		http.Error(w, "Failed to load project", http.StatusInternalServerError)
		return
	}
	if project == nil {
		http.NotFound(w, r)
		return
	}

	milestonesJSON, err := json.Marshal(project.Milestones)
	if err != nil {
		milestonesJSON = []byte("[]")
	}

	data := PlayerData{
		Title:             project.Title + " - Interactive Demo Player",
		HideNav:           true,
		HideFooter:        true,
		HideGatewayHeader: true,
		Project:           project,
		MilestonesJSON:    string(milestonesJSON),
	}

	if err := h.renderer.Render(w, "player.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
