package updown

import (
	"fmt"
	"net/http"
)

// RecipientType represents the type of a recipient
type RecipientType string

const (
	RecipientTypeEmail    RecipientType = "email"
	RecipientTypeSMS      RecipientType = "sms"
	RecipientTypeTelegram RecipientType = "telegram"
	RecipientTypeSlack    RecipientType = "slack"
	RecipientTypeWebhook  RecipientType = "webhook"
	RecipientTypeZapier   RecipientType = "zapier"
)

// Recipient represents a recipient from the API
type Recipient struct {
	ID    string        `json:"id,omitempty"`
	Type  RecipientType `json:"type,omitempty"`
	Value string        `json:"value,omitempty"`
	Name  string        `json:"name,omitempty"`
}

// RecipientItem represents a recipient to create
type RecipientItem struct {
	Type  RecipientType `json:"type,omitempty"`
	Value string        `json:"value,omitempty"`
	Name  string        `json:"name,omitempty"`
}

// RecipientService interacts with the recipients API
type RecipientService struct {
	client *Client
}

// List lists all recipients
func (s *RecipientService) List() ([]Recipient, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "recipients", nil)
	if err != nil {
		return nil, nil, err
	}

	var res []Recipient
	resp, err := s.client.Do(req, &res)
	if err != nil {
		return nil, resp, err
	}

	return res, resp, err
}

// Add creates a new recipient
func (s *RecipientService) Add(data RecipientItem) (Recipient, *http.Response, error) {
	req, err := s.client.NewRequest("POST", "recipients", data)
	if err != nil {
		return Recipient{}, nil, err
	}

	var res Recipient
	resp, err := s.client.Do(req, &res)
	if err != nil {
		return Recipient{}, resp, err
	}

	return res, resp, err
}

// Remove deletes a recipient by ID
func (s *RecipientService) Remove(id string) (bool, *http.Response, error) {
	req, err := s.client.NewRequest("DELETE", fmt.Sprintf("recipients/%s", id), nil)
	if err != nil {
		return false, nil, err
	}

	var res struct {
		Deleted bool `json:"deleted"`
	}
	resp, err := s.client.Do(req, &res)
	if err != nil {
		return false, resp, err
	}

	return res.Deleted, resp, err
}
