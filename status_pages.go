package updown

import (
	"fmt"
	"net/http"
)

// StatusPage represents a status page
type StatusPage struct {
	Token       string   `json:"token,omitempty"`
	URL         string   `json:"url,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Visibility  string   `json:"visibility,omitempty"`
	AccessKey   string   `json:"access_key,omitempty"`
	Checks      []string `json:"checks,omitempty"`
}

// StatusPageItem represents a status page to create or update
type StatusPageItem struct {
	// List of checks to show in the page (array of check tokens, order is respected)
	Checks []string `json:"checks,omitempty"`
	// Name of the status page
	Name string `json:"name,omitempty"`
	// Description text (displayed below the name, supports newlines and links)
	Description string `json:"description,omitempty"`
	// Page visibility: 'public', 'protected', or 'private'
	Visibility string `json:"visibility,omitempty"`
	// Access key for protected pages
	AccessKey string `json:"access_key,omitempty"`
}

// StatusPageService interacts with the status pages section of the API
type StatusPageService struct {
	client *Client
}

type removeStatusPageResponse struct {
	Deleted bool `json:"deleted,omitempty"`
}

// List lists all status pages
func (s *StatusPageService) List() ([]StatusPage, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "status_pages", nil)
	if err != nil {
		return nil, nil, err
	}

	var res []StatusPage
	resp, err := s.client.Do(req, &res)
	if err != nil {
		return nil, resp, err
	}

	return res, resp, err
}

// Get gets a single status page by its token
func (s *StatusPageService) Get(token string) (StatusPage, *http.Response, error) {
	req, err := s.client.NewRequest("GET", pathForStatusPageToken(token), nil)
	if err != nil {
		return StatusPage{}, nil, err
	}

	var res StatusPage
	resp, err := s.client.Do(req, &res)
	if err != nil {
		return StatusPage{}, resp, err
	}

	return res, resp, err
}

// Add creates a new status page
func (s *StatusPageService) Add(data StatusPageItem) (StatusPage, *http.Response, error) {
	req, err := s.client.NewRequest("POST", "status_pages", data)
	if err != nil {
		return StatusPage{}, nil, err
	}

	var res StatusPage
	resp, err := s.client.Do(req, &res)
	if err != nil {
		return StatusPage{}, resp, err
	}

	return res, resp, err
}

// Update updates a status page
func (s *StatusPageService) Update(token string, data StatusPageItem) (StatusPage, *http.Response, error) {
	req, err := s.client.NewRequest("PUT", pathForStatusPageToken(token), data)
	if err != nil {
		return StatusPage{}, nil, err
	}

	var res StatusPage
	resp, err := s.client.Do(req, &res)
	if err != nil {
		return StatusPage{}, resp, err
	}

	return res, resp, err
}

// Remove removes a status page by its token
func (s *StatusPageService) Remove(token string) (bool, *http.Response, error) {
	req, err := s.client.NewRequest("DELETE", pathForStatusPageToken(token), nil)
	if err != nil {
		return false, nil, err
	}

	var res removeStatusPageResponse
	resp, err := s.client.Do(req, &res)
	if err != nil {
		return false, resp, err
	}

	return res.Deleted, resp, err
}

func pathForStatusPageToken(token string) string {
	return fmt.Sprintf("status_pages/%s", token)
}
