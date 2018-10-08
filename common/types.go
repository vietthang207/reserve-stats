package common

// Token holds all necessary data regard a token.
type Token struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Address              string `json:"address"`
	Decimals             int64  `json:"decimals"`
	Active               bool   `json:"active"`
	Internal             bool   `json:"internal"`
	LastActivationChange int64  `json:"last_activation_change"`
}

// ReserveTokenRateEntry is a  map[ETH-tokenID]ratesEntry
type ReserveTokenRateEntry map[string]ReserveRateEntry

// ReserveRateEntry hold 4 float number represent necessary data for a rate entry
type ReserveRateEntry struct {
	BuyReserveRate  float64
	BuySanityRate   float64
	SellReserveRate float64
	SellSanityRate  float64
}

// ReserveRates hold all the pairs's rate for a particular reserve and metadata
type ReserveRates struct {
	Timestamp   int64
	ReturnTime  int64
	BlockNumber int64
	Data        ReserveTokenRateEntry
}

// NewToken creates a new Token.
func NewToken(id, name, address string, decimal int64, active, internal bool, timepoint int64) Token {
	return Token{
		ID:                   id,
		Name:                 name,
		Address:              address,
		Decimals:             decimal,
		Active:               active,
		Internal:             internal,
		LastActivationChange: timepoint,
	}
}
