package client

// Token is a ERC20 token allowed to trade in core.
type Token struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Address              string `json:"address"`
	Decimals             int64  `json:"decimals"`
	Active               bool   `json:"active"`
	Internal             bool   `json:"internal"`
	LastActivationChange uint64 `json:"last_activation_change"`
}
