package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"iqbal-portfolio/internal/handlers"
	"iqbal-portfolio/internal/repository"
	"iqbal-portfolio/web"
)

type Server struct {
	router     *http.ServeMux
	db         *sql.DB
	tmpl       *template.Template
	pages      map[string]*template.Template
	projects   *repository.ProjectRepo
	skills     *repository.SkillRepo
	shortlists *repository.ShortlistRepo
	inquiries  *repository.InquiryRepo
	handlers   *handlers.Handler

	Projects   *repository.ProjectRepo
	Skills     *repository.SkillRepo
	Shortlists *repository.ShortlistRepo
	Inquiries  *repository.InquiryRepo
}

func NewServer(db *sql.DB) (*Server, error) {
	tmpl, pages, err := parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	pRepo := repository.NewProjectRepo(db)
	sRepo := repository.NewSkillRepo(db)
	shRepo := repository.NewShortlistRepo(db)
	iRepo := repository.NewInquiryRepo(db)

	mux := http.NewServeMux()

	srv := &Server{
		router:     mux,
		db:         db,
		tmpl:       tmpl,
		pages:      pages,
		projects:   pRepo,
		skills:     sRepo,
		shortlists: shRepo,
		inquiries:  iRepo,
		Projects:   pRepo,
		Skills:     sRepo,
		Shortlists: shRepo,
		Inquiries:  iRepo,
	}

	srv.handlers = handlers.New(pRepo, sRepo, shRepo, iRepo, srv)

	if err := srv.registerRoutes(); err != nil {
		return nil, fmt.Errorf("failed to register routes: %w", err)
	}

	return srv, nil
}

func (s *Server) registerRoutes() error {
	staticSubFS, err := fs.Sub(web.Files, "static")
	if err != nil {
		return fmt.Errorf("failed to access static fs: %w", err)
	}

	fileServer := http.FileServer(http.FS(staticSubFS))
	s.router.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	s.router.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	s.router.HandleFunc("GET /{$}", s.HandleHome)
	s.router.HandleFunc("GET /gateway", s.HandleGateway)
	s.router.HandleFunc("POST /gateway/select", s.HandleGatewaySelect)
	s.router.HandleFunc("GET /projects/{id}/modal", s.HandleProjectModal)
	s.router.HandleFunc("GET /search", s.HandleSearch)
	s.router.HandleFunc("GET /api/search", s.HandleSearchAPI)
	s.router.HandleFunc("GET /projects/{id}/player", s.HandlePlayer)
	s.router.HandleFunc("GET /shortlist", s.HandleShortlist)
	s.router.HandleFunc("POST /shortlist/{id}", s.HandleShortlistToggle)
	s.router.HandleFunc("POST /inquiries", s.HandleInquiries)

	return nil
}

func (s *Server) HandlePlayer(w http.ResponseWriter, r *http.Request) {
	s.handlers.Player(w, r)
}

func (s *Server) HandleShortlist(w http.ResponseWriter, r *http.Request) {
	s.handlers.Shortlist(w, r)
}

func (s *Server) HandleShortlistToggle(w http.ResponseWriter, r *http.Request) {
	s.handlers.ShortlistToggle(w, r)
}

func (s *Server) HandleInquiries(w http.ResponseWriter, r *http.Request) {
	s.handlers.Inquiries(w, r)
}

func (s *Server) HandleProjectModal(w http.ResponseWriter, r *http.Request) {
	s.handlers.ProjectModal(w, r)
}

func (s *Server) HandleSearch(w http.ResponseWriter, r *http.Request) {
	s.handlers.Search(w, r)
}

func (s *Server) HandleSearchAPI(w http.ResponseWriter, r *http.Request) {
	s.handlers.SearchAPI(w, r)
}

func (s *Server) HandleGateway(w http.ResponseWriter, r *http.Request) {
	s.handlers.Gateway(w, r)
}

func (s *Server) HandleGatewaySelect(w http.ResponseWriter, r *http.Request) {
	s.handlers.GatewaySelect(w, r)
}

func (s *Server) HandleHome(w http.ResponseWriter, r *http.Request) {
	s.handlers.Home(w, r)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) DB() *sql.DB {
	return s.db
}

func (s *Server) Router() *http.ServeMux {
	return s.router
}

func (s *Server) Templates() *template.Template {
	return s.tmpl
}

func (s *Server) Render(w http.ResponseWriter, name string, data any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if pageTmpl, ok := s.pages[name]; ok {
		if pageTmpl.Lookup("base.html") != nil {
			return pageTmpl.ExecuteTemplate(w, "base.html", data)
		}
		return pageTmpl.ExecuteTemplate(w, name, data)
	}
	return s.tmpl.ExecuteTemplate(w, name, data)
}

func parseTemplates() (*template.Template, map[string]*template.Template, error) {
	funcMap := template.FuncMap{
		"sliceJoin": func(slice []string, sep string) string {
			return strings.Join(slice, sep)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"formatScore": func(score int) string {
			return fmt.Sprintf("%d%%", score)
		},
	}

	baseTmpl := template.New("base").Funcs(funcMap)
	allTmpl := template.New("all").Funcs(funcMap)
	pageFiles := make(map[string]string)

	err := fs.WalkDir(web.Files, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			content, err := web.Files.ReadFile(path)
			if err != nil {
				return err
			}
			name := filepath.Base(path)
			if _, err := allTmpl.New(name).Parse(string(content)); err != nil {
				return fmt.Errorf("failed parsing template %s into all: %w", path, err)
			}

			normPath := filepath.ToSlash(path)
			if strings.Contains(normPath, "templates/pages/") {
				pageFiles[name] = string(content)
			} else {
				if _, err := baseTmpl.New(name).Parse(string(content)); err != nil {
					return fmt.Errorf("failed parsing base/partial %s: %w", path, err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	pages := make(map[string]*template.Template)
	for pageName, content := range pageFiles {
		clone, err := baseTmpl.Clone()
		if err != nil {
			return nil, nil, fmt.Errorf("failed cloning base template for %s: %w", pageName, err)
		}
		if _, err := clone.New(pageName).Parse(content); err != nil {
			return nil, nil, fmt.Errorf("failed parsing page %s into clone: %w", pageName, err)
		}
		pages[pageName] = clone
	}

	return allTmpl, pages, nil
}
