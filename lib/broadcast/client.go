package broadcast

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KyberNetwork/httpsign-utils/sign"
	"go.uber.org/zap"
)

// errNotFoundMsg is the error returned by broadcast API
// if the given transaction address not found.
var errNotFoundMsg = "Can not find the transaction. Check Tx again"

// Client is the the real implementation of broadcast client interface
type Client struct {
	host   string
	sugar  *zap.SugaredLogger
	client *http.Client

	readKeyID     string
	readSecretKey string
}

type tradeLogGeoInfoResp struct {
	Success bool   `json:"success"`
	Err     string `json:"err"`
	Data    struct {
		IP      string `json:"IP"`
		Country string `json:"Country"`
	} `json:"data"`
}

const timeout = time.Minute * 5

// NewClient creates a new broadcast client instance.
func NewClient(sugar *zap.SugaredLogger, host, readKeyID, readSecretKey string) (*Client, error) {
	return &Client{
		host:          host,
		sugar:         sugar,
		client:        &http.Client{Timeout: timeout},
		readKeyID:     readKeyID,
		readSecretKey: readSecretKey,
	}, nil
}

// GetTxInfo get ip, country info of a tx
func (c *Client) GetTxInfo(tx string) (ip string, country string, err error) {
	url := fmt.Sprintf("%s/get-tx-info/%s", c.host, tx)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", "", err
	}
	signedReq, err := sign.Sign(req, c.readKeyID, c.readSecretKey)
	if err != nil {
		return "", "", err
	}
	resp, err := c.client.Do(signedReq)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			c.sugar.Errorw("failed to close body", "err", cErr.Error())
		}
	}()
	response := tradeLogGeoInfoResp{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", "", err
	}
	if !response.Success {
		if response.Err == errNotFoundMsg {
			c.sugar.Debugw("transaction not found", "tx", tx, "err", response.Err)
			return "", "", nil
		}
		c.sugar.Errorw("server returns unknown error")
		return "", "", err
	}
	return response.Data.IP, response.Data.Country, nil
}
