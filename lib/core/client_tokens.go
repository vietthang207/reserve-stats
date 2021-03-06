package core

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

// generateNonce returns nonce header required to use Core API,
// which is current timestamp in milliseconds.
func generateNonce() string {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(now, 10)
}

type allSettingsResponse struct {
	commonResponse
	Data tokensData `json:"data"`
}

type tokensData struct {
	Tokens *tokenList `json:"tokens"`
}

type tokenList struct {
	Tokens []Token `json:"tokens"`
}

// Tokens returns all configured tokens.
// Example response JSON:
// {
//  "data": {
//    "tokens": {
//      "tokens": [
//        {
//          "id": "ABT",
//          "name": "ArcBlock",
//          "address": "0xb98d4c97425d9908e66e53a6fdf673acca0be986",
//          "decimals": 18,
//          "active": true,
//          "internal": true,
//          "last_activation_change": 1535021910190
//        },
//        {
//          "id": "ADX",
//          "name": "AdEx",
//          "address": "0x4470BB87d77b963A013DB939BE332f927f2b992e",
//          "decimals": 4,
//          "active": true,
//          "internal": false,
//          "last_activation_change": 1535021910195
//        }
//      ]
//    }
//  }
//}
func (c *Client) Tokens() ([]Token, error) {
	const endpoint = "/setting/all-settings"
	var params = make(map[string]string)

	params["nonce"] = generateNonce()

	req, err := c.newRequest(http.MethodGet, endpoint, params)
	if err != nil {
		return nil, err
	}

	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected return code: %d", rsp.StatusCode)
	}

	var settingsResponse = allSettingsResponse{}
	if err = json.NewDecoder(rsp.Body).Decode(&settingsResponse); err != nil {
		return nil, err
	}

	if settingsResponse.Success != true {
		return nil, fmt.Errorf("got an error from server: %s", settingsResponse.Reason)
	}

	return settingsResponse.Data.Tokens.Tokens, nil
}

// tokensReply is the struct to contain core's reply
type tokensReply struct {
	commonResponse
	Data []Token `json:"data"`
}

func (c *Client) getTokens(endpoint string) ([]Token, error) {
	var params = make(map[string]string)
	params["nonce"] = generateNonce()
	req, err := c.newRequest(http.MethodGet, endpoint, params)
	if err != nil {
		return nil, err
	}
	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected return code: %d", rsp.StatusCode)
	}
	var tokenReply = tokensReply{}
	if err = json.NewDecoder(rsp.Body).Decode(&tokenReply); err != nil {
		return nil, err
	}
	return tokenReply.Data, nil
}

// FromWei formats the given amount in wei to human friendly
// number with preconfigured token decimals.
func (c *Client) FromWei(address common.Address, amount *big.Int) (float64, error) {
	tokens, err := c.Tokens()
	if err != nil {
		return 0, err
	}

	for _, token := range tokens {
		if common.HexToAddress(token.Address) == address {
			return token.FromWei(amount), nil
		}
	}

	return 0, fmt.Errorf("token not found: %s", address)
}

// ToWei return the given human friendly number to wei unit.
func (c *Client) ToWei(address common.Address, amount float64) (*big.Int, error) {
	tokens, err := c.Tokens()
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		if common.HexToAddress(token.Address) == address {
			return token.ToWei(amount), nil
		}
	}

	return nil, fmt.Errorf("token not found: %s", address)
}

// GetInternalTokens return list of internal token from Kyber reserve
func (c *Client) GetInternalTokens() ([]Token, error) {
	const endpoint = "setting/internal-tokens"
	return c.getTokens(endpoint)
}

// GetActiveTokens return list of active token from external reserve
func (c *Client) GetActiveTokens() ([]Token, error) {
	const endpoint = "setting/active-tokens"
	return c.getTokens(endpoint)
}
