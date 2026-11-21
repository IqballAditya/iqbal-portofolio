package handlers

import (
	"net/http"

	"iqbal-portfolio/internal/models"
)

type ModalData struct {
	Project  *models.Project
	NotFound bool
	Error    string
}

func (h *Handler) ProjectModal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		_ = h.renderer.Render(w, "modal_detail.html", ModalData{
			NotFound: true,
			Error:    "Project Not Found",
		})
		return
	}

	project, err := h.projects.GetProjectByID(id)
	if err != nil {
		http.Error(w, "Failed to retrieve project", http.StatusInternalServerError)
		return
	}

	if project == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = h.renderer.Render(w, "modal_detail.html", ModalData{
			NotFound: true,
			Error:    "Project Not Found",
		})
		return
	}

	if err := h.renderer.Render(w, "modal_detail.html", ModalData{
		Project: project,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
