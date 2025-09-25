package artifacthub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	DefaultAPIURL      = "https://artifacthub.io/api/v1"
	TektonTaskKind     = "Tekton task"
	TektonPipelineKind = "Tekton pipeline"
)

// Client represents an Artifact Hub API client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Package represents a package from Artifact Hub
type Package struct {
	PackageID           string                 `json:"package_id"`
	Name                string                 `json:"name"`
	NormalizedName      string                 `json:"normalized_name"`
	DisplayName         string                 `json:"display_name"`
	Description         string                 `json:"description"`
	LogoImageID         string                 `json:"logo_image_id,omitempty"`
	Keywords            []string               `json:"keywords,omitempty"`
	HomeURL             string                 `json:"home_url,omitempty"`
	Readme              string                 `json:"readme,omitempty"`
	Version             string                 `json:"version"`
	AvailableVersions   []Version              `json:"available_versions,omitempty"`
	Deprecated          bool                   `json:"deprecated"`
	License             string                 `json:"license,omitempty"`
	Signed              bool                   `json:"signed"`
	ContentURL          string                 `json:"content_url,omitempty"`
	CreatedAt           time.Time              `json:"created_at"`
	Repository          Repository             `json:"repository"`
	Stats               Stats                  `json:"stats,omitempty"`
	ProductionOrgsCount int                    `json:"production_orgs_count,omitempty"`
	Recommendations     []Recommendation       `json:"recommendations,omitempty"`
	Data                map[string]interface{} `json:"data,omitempty"`
	Links               []Link                 `json:"links,omitempty"`
	Install             string                 `json:"install,omitempty"`
}

// Version represents a package version
type Version struct {
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	Digest    string    `json:"digest,omitempty"`
}

// Repository represents a repository
type Repository struct {
	RepositoryID            string `json:"repository_id"`
	Name                    string `json:"name"`
	DisplayName             string `json:"display_name,omitempty"`
	URL                     string `json:"url"`
	Kind                    int    `json:"kind"`
	UserAlias               string `json:"user_alias,omitempty"`
	OrganizationName        string `json:"organization_name,omitempty"`
	OrganizationDisplayName string `json:"organization_display_name,omitempty"`
}

// Stats represents package statistics
type Stats struct {
	Subscriptions int `json:"subscriptions"`
	Webhooks      int `json:"webhooks"`
}

// Recommendation represents a package recommendation
type Recommendation struct {
	URL string `json:"url"`
}

// Link represents a package link
type Link struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// SearchResponse represents the response from a package search
type SearchResponse struct {
	Packages []Package `json:"packages"`
	Facets   []Facet   `json:"facets,omitempty"`
}

// Facet represents search facets
type Facet struct {
	Title     string        `json:"title"`
	FilterKey string        `json:"filter_key"`
	Options   []FacetOption `json:"options"`
}

// FacetOption represents a facet option
type FacetOption struct {
	ID    interface{} `json:"id"`
	Name  string      `json:"name"`
	Total int         `json:"total"`
}

// SearchOptions represents search parameters
type SearchOptions struct {
	Text         string
	Kinds        []string
	Categories   []string
	Repositories []string
	Deprecated   *bool
	Operators    *bool
	Verified     *bool
	Official     *bool
	CNCF         *bool
	Sort         string
	Limit        int
	Offset       int
}

// NewClient creates a new Artifact Hub client
func NewClient() *Client {
	return &Client{
		baseURL: DefaultAPIURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClientWithURL creates a new Artifact Hub client with custom URL
func NewClientWithURL(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchPackages searches for packages on Artifact Hub
func (c *Client) SearchPackages(ctx context.Context, opts SearchOptions) (*SearchResponse, error) {
	endpoint := c.baseURL + "/packages/search"

	// Build query parameters
	params := url.Values{}
	if opts.Text != "" {
		params.Set("ts_query_web", opts.Text)
	}
	if len(opts.Kinds) > 0 {
		for _, kind := range opts.Kinds {
			params.Add("kind", kind)
		}
	}
	if len(opts.Categories) > 0 {
		for _, category := range opts.Categories {
			params.Add("category", category)
		}
	}
	if len(opts.Repositories) > 0 {
		for _, repo := range opts.Repositories {
			params.Add("repo", repo)
		}
	}
	if opts.Deprecated != nil {
		if *opts.Deprecated {
			params.Set("deprecated", "true")
		} else {
			params.Set("deprecated", "false")
		}
	}
	if opts.Operators != nil {
		if *opts.Operators {
			params.Set("operators", "true")
		}
	}
	if opts.Verified != nil {
		if *opts.Verified {
			params.Set("verified_publisher", "true")
		}
	}
	if opts.Official != nil {
		if *opts.Official {
			params.Set("official", "true")
		}
	}
	if opts.CNCF != nil {
		if *opts.CNCF {
			params.Set("cncf", "true")
		}
	}
	if opts.Sort != "" {
		params.Set("sort", opts.Sort)
	}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		params.Set("offset", strconv.Itoa(opts.Offset))
	}

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, err
	}

	return &searchResp, nil
}

// GetPackage retrieves a specific package by its ID
func (c *Client) GetPackage(ctx context.Context, packageID string) (*Package, error) {
	endpoint := fmt.Sprintf("%s/packages/%s", c.baseURL, packageID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pkg Package
	if err := json.Unmarshal(body, &pkg); err != nil {
		return nil, err
	}

	return &pkg, nil
}

// GetPackageContent retrieves the content/YAML definition of a package
func (c *Client) GetPackageContent(ctx context.Context, contentURL string) (string, error) {
	if contentURL == "" {
		return "", errors.New("content URL is empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, contentURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// SearchTektonTasks searches for Tekton tasks on Artifact Hub
func (c *Client) SearchTektonTasks(ctx context.Context, text string, limit int) (*SearchResponse, error) {
	opts := SearchOptions{
		Text:  text,
		Kinds: []string{"7"}, // Tekton task kind ID
		Limit: limit,
		Sort:  "relevance",
	}
	return c.SearchPackages(ctx, opts)
}

// SearchTektonPipelines searches for Tekton pipelines on Artifact Hub
func (c *Client) SearchTektonPipelines(ctx context.Context, text string, limit int) (*SearchResponse, error) {
	opts := SearchOptions{
		Text:  text,
		Kinds: []string{"7"}, // Tekton pipeline kind ID (same as tasks)
		Limit: limit,
		Sort:  "relevance",
	}
	return c.SearchPackages(ctx, opts)
}
