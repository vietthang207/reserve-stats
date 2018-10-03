package client

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Client is the the real implementation of core client interface.
type Client struct {
	url        string
	signingKey string

	client *http.Client
}

type commonResponse struct {
	Reason  string `json:"reason"`
	Success bool   `json:"success"`
}

// SortByKey sort all the params by key in string order
// This is required for the request to be signed correctly
func sortByKey(params map[string]string) map[string]string {
	newParams := make(map[string]string, len(params))
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		newParams[key] = params[key]
	}
	return newParams
}

func (c *Client) sign(msg string) (string, error) {
	mac := hmac.New(sha512.New, []byte(c.signingKey))
	if _, err := mac.Write([]byte(msg)); err != nil {
		return "", err
	}
	return common.Bytes2Hex(mac.Sum(nil)), nil
}

func GetTimepoint() uint64 {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return uint64(timestamp)
}

func (c *Client) newRequest(method, endpoint string, params map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.url, endpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	params = sortByKey(params)
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	nonce, ok := params["nonce"]
	if ok {
		signed, err := c.sign(q.Encode())
		if err != nil {
			return nil, err
		}
		req.Header.Add("nonce", nonce)
		req.Header.Add("signed", signed)
	}

	return req, nil
}

// NewClient creates a new core client instance.
func NewClient(url, signingKey string) (*Client, error) {
	const timeout = time.Minute
	c := &http.Client{
		Timeout: timeout,
	}
	return &Client{url: url, signingKey: signingKey, client: c}, nil
}
