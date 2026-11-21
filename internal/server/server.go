package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"iqbal-portfolio/internal/repository"
	"iqbal-portfolio/web"
)

type Server struct {
	router     *http.ServeMux
	db         *sql.DB
	tmpl       *template.Template
	projects   *repository.ProjectRepo
	skills     *repository.SkillRepo
	shortlists *repository.ShortlistRepo
	inquiries  *repository.InquiryRepo

	Projects   *repository.ProjectRepo
	Skills     *repository.SkillRepo
	Shortlists *repository.ShortlistRepo
	Inquiries  *repository.InquiryRepo
}

func NewServer(db *sql.DB) (*Server, error) {
	tmpl, err := parseTemplates()
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
		projects:   pRepo,
		skills:     sRepo,
		shortlists: shRepo,
		inquiries:  iRepo,
		Projects:   pRepo,
		Skills:     sRepo,
		Shortlists: shRepo,
		Inquiries:  iRepo,
	}

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

	return nil
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
	return s.tmpl.ExecuteTemplate(w, name, data)
}

func parseTemplates() (*template.Template, error) {
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

	tmpl := template.New("").Funcs(funcMap)
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
			if _, err := tmpl.New(name).Parse(string(content)); err != nil {
				return fmt.Errorf("failed parsing template %s: %w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}
