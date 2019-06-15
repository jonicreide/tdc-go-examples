package greatmail

//go:generate  mockgen -destination=mock_greatmail.go -package=greatmail main/greatmail HTTPProvider

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPProvider ...
type httpProvider interface {
	Do(*http.Request) (*http.Response, error)
}

// Client ...
type Client struct {
	http httpProvider
}

// Email ...
type Email struct {
	Message string   `json:"message"`
	Subject string   `json:"subject"`
	From    string   `json:"from"`
	To      []string `json:"to"`
}

// NewClient ...
func NewClient() Client {
	return Client{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendEmail ...
func (c Client) SendEmail(email Email) error {

	emailPayload, _ := json.Marshal(email)

	request, _ := http.NewRequest(
		"POST",
		"https://api.greatmail.com/send",
		ioutil.NopCloser(bytes.NewReader(emailPayload)),
	)

	_, err := c.http.Do(request)
	if err != nil {
		return err
	}

	return nil
}
