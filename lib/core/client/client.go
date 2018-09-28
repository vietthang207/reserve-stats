package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"time"
)

// Client is the the real implementation of core client interface.
type Client struct {
	url        string
	signingKey string

	c *http.Client
}

func (c *Client) newRequest(method, endpoint string, params map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, path.Join(c.url, endpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	return req, nil
}

type allSettingsResponse struct {
	Data allSettings `json:"data"`
}

type allSettings struct {
	Tokens *response `json:"tokens"`
}

type response struct {
	Tokens []Token `json:"tokens"`
}

func (c *Client) Tokens() ([]Token, error) {
	const endpoint = "/setting/all-settings"
	req, err := c.newRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	rsp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected return code: %d", rsp.StatusCode)
	}

	var settingsResponse = &allSettingsResponse{}
	if err = json.NewDecoder(rsp.Body).Decode(&settingsResponse); err != nil {
		return nil, err
	}

	return settingsResponse.Data.Tokens.Tokens, nil
}

// NewClient creates a new core client instance.
func NewClient(url, signingKey string) (*Client, error) {
	const timeout = 5 * time.Second

	c := &http.Client{
		Timeout: timeout,
	}
	return &Client{url: url, signingKey: signingKey, c: c}, nil
}
