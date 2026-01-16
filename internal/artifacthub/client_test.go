package artifacthub

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// testServer creates a test HTTP server with the given handler and returns
// the server and a client configured to use it. The caller is responsible
// for closing the server.
func testServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	server := httptest.NewServer(handler)
	client := NewClientWithURL(server.URL)
	return server, client
}

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.baseURL != defaultAPIURL {
		t.Errorf("Expected baseURL to be %q, got %q", defaultAPIURL, client.baseURL)
	}
}

func TestNewClientWithURL(t *testing.T) {
	customURL := "https://custom.example.com/api"
	client := NewClientWithURL(customURL)
	if client == nil {
		t.Fatal("NewClientWithURL returned nil")
	}
	if client.baseURL != customURL {
		t.Errorf("Expected baseURL to be %q, got %q", customURL, client.baseURL)
	}
}

func TestSearchPackages(t *testing.T) {
	expectedResponse := SearchResponse{
		Packages: []Package{
			{
				PackageID:   "pkg-1",
				Name:        "test-task",
				DisplayName: "Test Task",
				Description: "A test task",
				Version:     "0.1.0",
				Repository: Repository{
					RepositoryID: "repo-1",
					Name:         "tekton-catalog",
					DisplayName:  "Tekton Catalog",
				},
			},
			{
				PackageID:   "pkg-2",
				Name:        "another-task",
				DisplayName: "Another Task",
				Description: "Another test task",
				Version:     "0.2.0",
				Repository: Repository{
					RepositoryID: "repo-1",
					Name:         "tekton-catalog",
					DisplayName:  "Tekton Catalog",
				},
			},
		},
	}

	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/packages/search" {
			t.Errorf("Expected path /packages/search, got %s", r.URL.Path)
		}

		// Verify query parameters
		query := r.URL.Query()
		if query.Get("ts_query_web") != "git" {
			t.Errorf("Expected ts_query_web=git, got %s", query.Get("ts_query_web"))
		}
		if query.Get("kind") != KindTektonTask {
			t.Errorf("Expected kind=%s, got %s", KindTektonTask, query.Get("kind"))
		}
		if query.Get("limit") != "10" {
			t.Errorf("Expected limit=10, got %s", query.Get("limit"))
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedResponse); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})
	defer server.Close()

	opts := SearchOptions{
		Text:  "git",
		Kinds: []string{KindTektonTask},
		Limit: 10,
	}

	resp, err := client.SearchPackages(t.Context(), opts)
	if err != nil {
		t.Fatalf("SearchPackages failed: %v", err)
	}

	if len(resp.Packages) != 2 {
		t.Errorf("Expected 2 packages, got %d", len(resp.Packages))
	}

	if diff := cmp.Diff(expectedResponse, *resp); diff != "" {
		t.Errorf("SearchPackages response mismatch (-want +got):\n%s", diff)
	}
}

func TestSearchPackages_Error(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	opts := SearchOptions{
		Text:  "git",
		Kinds: []string{KindTektonTask},
	}

	if _, err := client.SearchPackages(t.Context(), opts); err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestSearchTektonTasks(t *testing.T) {
	expectedResponse := SearchResponse{
		Packages: []Package{
			{
				PackageID:   "pkg-1",
				Name:        "git-clone",
				DisplayName: "Git Clone",
				Version:     "0.9.0",
			},
		},
	}

	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		// Verify it uses the correct kind for Tekton tasks
		if query.Get("kind") != KindTektonTask {
			t.Errorf("Expected kind=%s, got %s", KindTektonTask, query.Get("kind"))
		}
		if query.Get("sort") != "relevance" {
			t.Errorf("Expected sort=relevance, got %s", query.Get("sort"))
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedResponse); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})
	defer server.Close()

	resp, err := client.SearchTektonTasks(t.Context(), "git", 20)
	if err != nil {
		t.Fatalf("SearchTektonTasks failed: %v", err)
	}

	if len(resp.Packages) != 1 {
		t.Errorf("Expected 1 package, got %d", len(resp.Packages))
	}

	if resp.Packages[0].Name != "git-clone" {
		t.Errorf("Expected package name git-clone, got %s", resp.Packages[0].Name)
	}
}

func TestSearchTektonPipelines(t *testing.T) {
	expectedResponse := SearchResponse{
		Packages: []Package{
			{
				PackageID:   "pkg-1",
				Name:        "build-pipeline",
				DisplayName: "Build Pipeline",
				Version:     "1.0.0",
			},
		},
	}

	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		// Verify it uses the correct kind for Tekton pipelines
		if query.Get("kind") != KindTektonPipeline {
			t.Errorf("Expected kind=%s, got %s", KindTektonPipeline, query.Get("kind"))
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedResponse); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})
	defer server.Close()

	resp, err := client.SearchTektonPipelines(t.Context(), "build", 20)
	if err != nil {
		t.Fatalf("SearchTektonPipelines failed: %v", err)
	}

	if len(resp.Packages) != 1 {
		t.Errorf("Expected 1 package, got %d", len(resp.Packages))
	}
}

func TestGetPackage(t *testing.T) {
	expectedPackage := Package{
		PackageID:   "tekton-task/git-clone",
		Name:        "git-clone",
		DisplayName: "Git Clone",
		Description: "A task to clone a git repository",
		Version:     "0.9.0",
		ContentURL:  "https://raw.githubusercontent.com/tektoncd/catalog/main/task/git-clone/0.9/git-clone.yaml",
		Repository: Repository{
			RepositoryID: "repo-1",
			Name:         "tekton-catalog",
			DisplayName:  "Tekton Catalog",
		},
	}

	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/packages/tekton-task/git-clone" {
			t.Errorf("Expected path /packages/tekton-task/git-clone, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expectedPackage); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})
	defer server.Close()

	pkg, err := client.GetPackage(t.Context(), "tekton-task/git-clone")
	if err != nil {
		t.Fatalf("GetPackage failed: %v", err)
	}

	if pkg.PackageID != expectedPackage.PackageID {
		t.Errorf("Expected package ID %s, got %s", expectedPackage.PackageID, pkg.PackageID)
	}

	if pkg.ContentURL != expectedPackage.ContentURL {
		t.Errorf("Expected content URL %s, got %s", expectedPackage.ContentURL, pkg.ContentURL)
	}
}

func TestGetPackage_NotFound(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	defer server.Close()

	if _, err := client.GetPackage(t.Context(), "nonexistent-package"); err == nil {
		t.Fatal("Expected error for nonexistent package, got nil")
	}
}

func TestGetPackageContent_EmptyURL(t *testing.T) {
	client := NewClient()

	_, err := client.GetPackageContent(t.Context(), "")
	if err == nil {
		t.Fatal("Expected error for empty URL, got nil")
	}
}

func TestGetPackageContent_Error(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	if _, err := client.GetPackageContent(t.Context(), server.URL); err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestSearchOptions_AllParams(t *testing.T) {
	server, client := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// Check all parameters are passed correctly
		if query.Get("ts_query_web") != "test" {
			t.Errorf("Expected ts_query_web=test, got %s", query.Get("ts_query_web"))
		}
		if query.Get("sort") != "stars" {
			t.Errorf("Expected sort=stars, got %s", query.Get("sort"))
		}
		if query.Get("limit") != "50" {
			t.Errorf("Expected limit=50, got %s", query.Get("limit"))
		}
		if query.Get("offset") != "10" {
			t.Errorf("Expected offset=10, got %s", query.Get("offset"))
		}
		if query.Get("deprecated") != "false" {
			t.Errorf("Expected deprecated=false, got %s", query.Get("deprecated"))
		}
		if query.Get("verified_publisher") != "true" {
			t.Errorf("Expected verified_publisher=true, got %s", query.Get("verified_publisher"))
		}
		if query.Get("official") != "true" {
			t.Errorf("Expected official=true, got %s", query.Get("official"))
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(SearchResponse{}); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	})
	defer server.Close()

	deprecated := false
	verified := true
	official := true

	opts := SearchOptions{
		Text:       "test",
		Kinds:      []string{KindTektonTask},
		Sort:       "stars",
		Limit:      50,
		Offset:     10,
		Deprecated: &deprecated,
		Verified:   &verified,
		Official:   &official,
	}

	if _, err := client.SearchPackages(t.Context(), opts); err != nil {
		t.Fatalf("SearchPackages failed: %v", err)
	}
}

func TestConstants(t *testing.T) {
	// Verify constants are defined correctly
	if KindTektonTask != "7" {
		t.Errorf("Expected KindTektonTask to be '7', got %q", KindTektonTask)
	}
	if KindTektonPipeline != "11" {
		t.Errorf("Expected KindTektonPipeline to be '11', got %q", KindTektonPipeline)
	}
}
