package handlers

import (
	"net/http"
)

type Persona struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Initial  string `json:"initial"`
	Subtitle string `json:"subtitle"`
}

func ValidatePersona(p string) string {
	switch p {
	case "recruiter", "tech_lead", "founder", "guest":
		return p
	default:
		return "guest"
	}
}

func GetPersona(key string) Persona {
	switch ValidatePersona(key) {
	case "recruiter":
		return Persona{
			Key:      "recruiter",
			Name:     "Recruiter",
			Color:    "#0071EB",
			Initial:  "R",
			Subtitle: "Impact metrics, resume, hire details",
		}
	case "tech_lead":
		return Persona{
			Key:      "tech_lead",
			Name:     "Tech Lead",
			Color:    "#46D369",
			Initial:  "T",
			Subtitle: "Architectures, throughput, benchmarks",
		}
	case "founder":
		return Persona{
			Key:      "founder",
			Name:     "Founder",
			Color:    "#8B5CF6",
			Initial:  "F",
			Subtitle: "Business ROI, scaling speed, case studies",
		}
	default:
		return Persona{
			Key:      "guest",
			Name:     "Guest",
			Color:    "#E50914",
			Initial:  "G",
			Subtitle: "Standard portfolio showcase",
		}
	}
}

func AllPersonas() []Persona {
	return []Persona{
		GetPersona("recruiter"),
		GetPersona("tech_lead"),
		GetPersona("founder"),
		GetPersona("guest"),
	}
}

func (h *Handler) Gateway(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title      string
		HideNav    bool
		HideFooter bool
		Personas   []Persona
	}{
		Title:      "Who's Watching?",
		HideNav:    true,
		HideFooter: true,
		Personas:   AllPersonas(),
	}

	if err := h.renderer.Render(w, "gateway.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GatewaySelect(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	selected := ValidatePersona(r.FormValue("persona"))
	http.SetCookie(w, &http.Cookie{
		Name:     "iqbal_persona",
		Value:    selected,
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
