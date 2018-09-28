package client

import (
	"net/http"
	"time"
)

// Client is the the real implementation of core client interface.
type Client struct {
	url        string
	signingKey string

	c *http.Client
}

//func (sc *SettingClient) newRequest(method, url string, params map[string]string) (*http.Request, error) {
//	req, err := http.NewRequest(method, url, nil)
//	if err != nil {
//		return nil, err
//	}
//	// Add header
//	req.Header.Add("Accept", "application/json")
//	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//	// Create raw query
//	q := req.URL.Query()
//	for k, v := range params {
//		q.Add(k, v)
//	}
//	req.URL.RawQuery = q.Encode()
//	//sign
//	nonce, ok := params["nonce"]
//	if !ok {
//		log.Printf("there was no nonce")
//	} else {
//		sc.sign(req, q.Encode(), nonce)
//	}
//
//	return req, nil
//}

func (c *Client) newRequest(method, endpoint string, params map[string]string) (*http.Request, error) {

}

func (c *Client) Tokens() ([]Token, error) {
	const endpoint = "/setting/all-settings"
	//req, err := http.NewRequest(http.MethodGet, nil)
	return nil, nil
}

// NewClient creates a new core client instance.
func NewClient(url, signingKey string) (*Client, error) {
	const timeout = 5 * time.Second

	c := &http.Client{
		Timeout: timeout,
	}
	return &Client{url: url, signingKey: signingKey, c: c}, nil
}
