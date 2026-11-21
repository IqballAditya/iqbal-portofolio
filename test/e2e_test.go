package test

import (
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"iqbal-portfolio/internal/database"
	"iqbal-portfolio/internal/repository"
	"iqbal-portfolio/internal/server"
)

type localRoundTripper struct {
	handler http.Handler
}

func (rt *localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	rt.handler.ServeHTTP(rec, req)
	resp := rec.Result()
	resp.Request = req
	return resp, nil
}

func TestVisitorJourney_E2E(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate failed: %v", err)
	}
	if err := repository.SeedData(db); err != nil {
		t.Fatalf("SeedData failed: %v", err)
	}

	srv, err := server.NewServer(db)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	ts := httptest.NewServer(srv)
	defer ts.Close()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar.New failed: %v", err)
	}

	var transport http.RoundTripper = http.DefaultTransport
	addr := strings.TrimPrefix(ts.URL, "http://")
	if conn, dialErr := net.DialTimeout("tcp", addr, 50*time.Millisecond); dialErr != nil {
		transport = &localRoundTripper{handler: srv}
	} else {
		_ = conn.Close()
	}

	client := &http.Client{
		Jar:       jar,
		Transport: transport,
	}

	t.Run("1_Home_Showcase", func(t *testing.T) {
		resp, err := client.Get(ts.URL + "/")
		if err != nil {
			t.Fatalf("GET / failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200 for GET /, got %d", resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed reading response body: %v", err)
		}
		bodyStr := string(body)
		if !strings.Contains(bodyStr, "Adith AI Agent") {
			t.Errorf("expected hero to feature 'Adith AI Agent'")
		}
		if !strings.Contains(bodyStr, "Top 10 Tech Stacks") && !strings.Contains(bodyStr, "TOP 10") {
			t.Errorf("expected home page to contain tech stack carousel")
		}
	})

	t.Run("2_Gateway_PersonaSelection", func(t *testing.T) {
		respGateway, err := client.Get(ts.URL + "/gateway")
		if err != nil {
			t.Fatalf("GET /gateway failed: %v", err)
		}
		defer respGateway.Body.Close()

		if respGateway.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200 for GET /gateway, got %d", respGateway.StatusCode)
		}

		noRedirectClient := &http.Client{
			Jar:       jar,
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		selectData := url.Values{"persona": {"tech_lead"}}
		respSelect, err := noRedirectClient.PostForm(ts.URL+"/gateway/select", selectData)
		if err != nil {
			t.Fatalf("POST /gateway/select failed: %v", err)
		}
		defer respSelect.Body.Close()

		if respSelect.StatusCode != http.StatusSeeOther {
			t.Fatalf("expected 303 See Other redirect, got %d", respSelect.StatusCode)
		}
		if loc := respSelect.Header.Get("Location"); loc != "/" {
			t.Errorf("expected redirect location '/', got %q", loc)
		}

		u, _ := url.Parse(ts.URL)
		var personaSet bool
		for _, c := range jar.Cookies(u) {
			if c.Name == "iqbal_persona" && c.Value == "tech_lead" {
				personaSet = true
				break
			}
		}
		if !personaSet {
			t.Errorf("expected cookie 'iqbal_persona=tech_lead' in jar")
		}
	})

	t.Run("3_CaseStudy_Modal", func(t *testing.T) {
		respModal, err := client.Get(ts.URL + "/projects/adith-ai/modal")
		if err != nil {
			t.Fatalf("GET /projects/adith-ai/modal failed: %v", err)
		}
		defer respModal.Body.Close()

		if respModal.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200 for modal, got %d", respModal.StatusCode)
		}
		modalBody, err := io.ReadAll(respModal.Body)
		if err != nil {
			t.Fatalf("failed reading modal body: %v", err)
		}
		modalStr := string(modalBody)
		if !strings.Contains(modalStr, "Adith AI Agent") {
			t.Errorf("expected modal to contain project title")
		}
		if !strings.Contains(modalStr, "Blueprint") {
			t.Errorf("expected modal to list milestones")
		}
		if !strings.Contains(modalStr, "WorkerPool") {
			t.Errorf("expected modal to show code preview")
		}
	})

	t.Run("4_Shortlist_Toggle", func(t *testing.T) {
		respToggle, err := client.Post(ts.URL+"/shortlist/adith-ai", "application/x-www-form-urlencoded", nil)
		if err != nil {
			t.Fatalf("POST /shortlist/adith-ai failed: %v", err)
		}
		defer respToggle.Body.Close()

		if respToggle.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200 for shortlist toggle, got %d", respToggle.StatusCode)
		}

		u, _ := url.Parse(ts.URL)
		var sessionSet bool
		for _, c := range jar.Cookies(u) {
			if c.Name == "iqbal_session" && c.Value != "" {
				sessionSet = true
				break
			}
		}
		if !sessionSet {
			t.Errorf("expected 'iqbal_session' cookie to be set")
		}
	})

	t.Run("5_Shortlist_Page", func(t *testing.T) {
		respShortlist, err := client.Get(ts.URL + "/shortlist")
		if err != nil {
			t.Fatalf("GET /shortlist failed: %v", err)
		}
		defer respShortlist.Body.Close()

		if respShortlist.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200 for shortlist page, got %d", respShortlist.StatusCode)
		}
		shortlistBody, err := io.ReadAll(respShortlist.Body)
		if err != nil {
			t.Fatalf("failed reading shortlist body: %v", err)
		}
		shortlistStr := string(shortlistBody)
		if !strings.Contains(shortlistStr, "Adith AI Agent") {
			t.Errorf("expected shortlisted project 'Adith AI Agent' on page")
		}

		tiers := []string{
			"Full-Time Staff / Lead",
			"Advisory / Architecture Sprint",
			"Contract / Fractional",
		}
		for _, tier := range tiers {
			if !strings.Contains(shortlistStr, tier) {
				t.Errorf("expected engagement tier %q on shortlist page", tier)
			}
		}
	})

	t.Run("6_HireMe_Inquiry", func(t *testing.T) {
		inquiryForm := url.Values{
			"name":            {"Sarah Connor"},
			"email":           {"sarah@netflix.com"},
			"engagement_type": {"sprint"},
			"message":         {"We need a 4-week architecture review for our distributed streaming cluster."},
		}
		respInquiry, err := client.PostForm(ts.URL+"/inquiries", inquiryForm)
		if err != nil {
			t.Fatalf("POST /inquiries failed: %v", err)
		}
		defer respInquiry.Body.Close()

		if respInquiry.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200 for valid inquiry, got %d", respInquiry.StatusCode)
		}
		inqBody, _ := io.ReadAll(respInquiry.Body)
		inqStr := string(inqBody)
		if !strings.Contains(inqStr, "Inquiry Received") && !strings.Contains(inqStr, "Thank you") {
			t.Errorf("expected inquiry confirmation response, got: %s", inqStr)
		}

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM inquiries WHERE email = ?", "sarah@netflix.com").Scan(&count)
		if err != nil {
			t.Fatalf("failed querying database for inquiry: %v", err)
		}
		if count != 1 {
			t.Errorf("expected 1 inquiry recorded in DB, got %d", count)
		}
	})

	t.Run("7_Interactive_Player", func(t *testing.T) {
		respPlayer, err := client.Get(ts.URL + "/projects/adith-ai/player")
		if err != nil {
			t.Fatalf("GET /projects/adith-ai/player failed: %v", err)
		}
		defer respPlayer.Body.Close()

		if respPlayer.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200 for player page, got %d", respPlayer.StatusCode)
		}
		playerBody, _ := io.ReadAll(respPlayer.Body)
		playerStr := string(playerBody)
		if !strings.Contains(playerStr, "Adith AI Agent") {
			t.Errorf("expected player page to include project title")
		}
		if !strings.Contains(playerStr, "terminal-player") {
			t.Errorf("expected player page to contain terminal player")
		}
		if !strings.Contains(playerStr, "S1:E1") {
			t.Errorf("expected player page to contain milestones")
		}
	})

	t.Run("8_Live_Search_API", func(t *testing.T) {
		respSearch, err := client.Get(ts.URL + "/api/search?q=AI")
		if err != nil {
			t.Fatalf("GET /api/search?q=AI failed: %v", err)
		}
		defer respSearch.Body.Close()

		if respSearch.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200 for live search API, got %d", respSearch.StatusCode)
		}
		searchBody, _ := io.ReadAll(respSearch.Body)
		searchStr := string(searchBody)
		if !strings.Contains(searchStr, "Adith AI Agent") {
			t.Errorf("expected search results to contain 'Adith AI Agent', got: %s", searchStr)
		}
	})
}
