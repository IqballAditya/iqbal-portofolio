package handlers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"

	"iqbal-portfolio/internal/models"
)

type ShortlistData struct {
	Title          string
	HideNav        bool
	HideFooter     bool
	PersonaName    string
	PersonaColor   string
	PersonaInitial string
	ShortlistCount int
	Shortlist      []models.Project
}

type ShortlistItemData struct {
	Project       *models.Project
	IsShortlisted bool
	Count         int
}

type InquiryFeedbackData struct {
	Success bool
	Name    string
	Error   string
}

func getOrCreateSessionID(w http.ResponseWriter, r *http.Request) string {
	if cookie, err := r.Cookie("iqbal_session"); err == nil && cookie != nil && strings.TrimSpace(cookie.Value) != "" {
		return cookie.Value
	}

	b := make([]byte, 16)
	_, _ = rand.Read(b)
	sessionID := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])

	http.SetCookie(w, &http.Cookie{
		Name:     "iqbal_session",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   86400 * 365,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return sessionID
}

func (h *Handler) Shortlist(w http.ResponseWriter, r *http.Request) {
	sessionID := getOrCreateSessionID(w, r)

	personaKey := "guest"
	if cookie, err := r.Cookie("iqbal_persona"); err == nil && cookie != nil {
		personaKey = ValidatePersona(cookie.Value)
	}
	persona := GetPersona(personaKey)

	items, err := h.shortlists.ListShortlist(sessionID)
	if err != nil {
		http.Error(w, "Failed to load shortlist", http.StatusInternalServerError)
		return
	}

	data := ShortlistData{
		Title:          "My Shortlist & Hire Me - IQBAL",
		HideNav:        false,
		HideFooter:     false,
		PersonaName:    persona.Name,
		PersonaColor:   persona.Color,
		PersonaInitial: persona.Initial,
		ShortlistCount: len(items),
		Shortlist:      items,
	}

	if err := h.renderer.Render(w, "hire.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) ShortlistToggle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Project ID required", http.StatusBadRequest)
		return
	}

	sessionID := getOrCreateSessionID(w, r)

	isShort, err := h.shortlists.ToggleShortlist(sessionID, id)
	if err != nil {
		http.Error(w, "Failed to toggle shortlist", http.StatusInternalServerError)
		return
	}

	list, _ := h.shortlists.ListShortlist(sessionID)
	count := len(list)

	if r.Header.Get("HX-Target") == "shortlist-counter" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(fmt.Sprintf("%d", count)))
		return
	}

	project, _ := h.projects.GetProjectByID(id)
	if project == nil {
		project = &models.Project{ID: id}
	}

	itemData := ShortlistItemData{
		Project:       project,
		IsShortlisted: isShort,
		Count:         count,
	}

	if err := h.renderer.Render(w, "shortlist_item.html", itemData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Inquiries(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	engagementType := strings.TrimSpace(r.FormValue("engagement_type"))
	message := strings.TrimSpace(r.FormValue("message"))

	if name == "" || email == "" || !strings.Contains(email, "@") || !strings.Contains(email, ".") || message == "" {
		w.WriteHeader(http.StatusBadRequest)
		feedback := InquiryFeedbackData{
			Success: false,
			Error:   "Please provide a valid name, email address (with @ and .), and a message outlining your requirements.",
		}
		_ = h.renderer.Render(w, "inquiry_feedback.html", feedback)
		return
	}

	switch engagementType {
	case "staff", "sprint", "fractional":
	default:
		engagementType = "staff"
	}

	inq := models.Inquiry{
		Name:           name,
		Email:          email,
		EngagementType: engagementType,
		Message:        message,
	}

	if err := h.inquiries.CreateInquiry(&inq); err != nil {
		http.Error(w, "Failed to save inquiry", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	feedback := InquiryFeedbackData{
		Success: true,
		Name:    name,
	}
	_ = h.renderer.Render(w, "inquiry_feedback.html", feedback)
}
