package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"iqbal-portfolio/internal/database"
	"iqbal-portfolio/internal/models"
	"iqbal-portfolio/internal/repository"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}
	return db
}

func TestSeedDataIdempotency(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	if err := repository.SeedData(db); err != nil {
		t.Fatalf("first seed failed: %v", err)
	}

	// Calling seed again should succeed without duplicate key errors
	if err := repository.SeedData(db); err != nil {
		t.Fatalf("second seed failed: %v", err)
	}

	projectRepo := repository.NewProjectRepo(db)
	projects, err := projectRepo.ListProjects("")
	if err != nil {
		t.Fatalf("failed to list projects: %v", err)
	}
	if len(projects) < 6 {
		t.Fatalf("expected at least 6 projects, got %d", len(projects))
	}

	skillRepo := repository.NewSkillRepo(db)
	skills, err := skillRepo.ListSkills()
	if err != nil {
		t.Fatalf("failed to list skills: %v", err)
	}
	if len(skills) != 10 {
		t.Fatalf("expected 10 skills, got %d", len(skills))
	}
}

func TestListProjectsWithPersonaPriority(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	if err := repository.SeedData(db); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	projectRepo := repository.NewProjectRepo(db)
	projects, err := projectRepo.ListProjects("tech_lead")
	if err != nil {
		t.Fatalf("ListProjects error: %v", err)
	}
	if len(projects) == 0 {
		t.Fatalf("expected projects, got 0")
	}

	// First project should match tech_lead persona
	found := false
	for _, p := range projects[0].TargetPersonas {
		if p == "tech_lead" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected first project to target tech_lead, got personas: %v", projects[0].TargetPersonas)
	}
}

func TestGetProjectByIDWithMilestones(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	if err := repository.SeedData(db); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	projectRepo := repository.NewProjectRepo(db)
	project, err := projectRepo.GetProjectByID("nexus-ai")
	if err != nil {
		t.Fatalf("GetProjectByID error: %v", err)
	}
	if project == nil {
		t.Fatalf("expected project nexus-ai, got nil")
	}
	if project.ID != "nexus-ai" {
		t.Errorf("expected ID nexus-ai, got %s", project.ID)
	}
	if len(project.Milestones) != 3 {
		t.Errorf("expected 3 milestones, got %d", len(project.Milestones))
	}
	if project.Milestones[0].EpisodeNumber != 1 {
		t.Errorf("expected episode 1, got %d", project.Milestones[0].EpisodeNumber)
	}
}

func TestSearchProjects(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	if err := repository.SeedData(db); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	projectRepo := repository.NewProjectRepo(db)

	// Search by query and category
	results, err := projectRepo.SearchProjects("stream", "backend")
	if err != nil {
		t.Fatalf("SearchProjects error: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected search results for 'stream' in 'backend', got 0")
	}
	for _, p := range results {
		if p.Category != "backend" {
			t.Errorf("expected category backend, got %s", p.Category)
		}
	}

	// Search category all with query
	resultsAll, err := projectRepo.SearchProjects("nexus", "all")
	if err != nil {
		t.Fatalf("SearchProjects all category error: %v", err)
	}
	if len(resultsAll) == 0 {
		t.Fatalf("expected search results for 'nexus', got 0")
	}
}

func TestToggleShortlist(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	if err := repository.SeedData(db); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	shortlistRepo := repository.NewShortlistRepo(db)
	sessionID := "sess-1"
	projectID := "nexus-ai"

	// Toggle ON
	shortlisted, err := shortlistRepo.ToggleShortlist(sessionID, projectID)
	if err != nil {
		t.Fatalf("ToggleShortlist error: %v", err)
	}
	if !shortlisted {
		t.Errorf("expected shortlisted=true")
	}

	isShort, err := shortlistRepo.IsShortlisted(sessionID, projectID)
	if err != nil {
		t.Fatalf("IsShortlisted error: %v", err)
	}
	if !isShort {
		t.Errorf("expected IsShortlisted=true")
	}

	list, err := shortlistRepo.ListShortlist(sessionID)
	if err != nil {
		t.Fatalf("ListShortlist error: %v", err)
	}
	if len(list) != 1 || list[0].ID != projectID {
		t.Errorf("expected 1 shortlisted project %s, got %v", projectID, list)
	}

	// Toggle OFF
	shortlisted, err = shortlistRepo.ToggleShortlist(sessionID, projectID)
	if err != nil {
		t.Fatalf("ToggleShortlist toggle off error: %v", err)
	}
	if shortlisted {
		t.Errorf("expected shortlisted=false")
	}

	isShort, err = shortlistRepo.IsShortlisted(sessionID, projectID)
	if err != nil {
		t.Fatalf("IsShortlisted error: %v", err)
	}
	if isShort {
		t.Errorf("expected IsShortlisted=false")
	}
}

func TestCreateInquiry(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	inquiryRepo := repository.NewInquiryRepo(db)
	inq := &models.Inquiry{
		Name:           "Alex Rivera",
		Email:          "alex@example.com",
		EngagementType: "Full-time Staff Engineer",
		Message:        "Let's discuss distributed systems roles at our company.",
		CreatedAt:      time.Now().UTC(),
	}

	if err := inquiryRepo.CreateInquiry(inq); err != nil {
		t.Fatalf("CreateInquiry error: %v", err)
	}
	if inq.ID == 0 {
		t.Errorf("expected inquiry ID to be set, got 0")
	}
}
