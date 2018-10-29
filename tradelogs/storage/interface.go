package storage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

// Interface represent a storage for TradeLogs data
type Interface interface {
	SaveTradeLogs(logs []common.TradeLog, rates []tokenrate.ETHUSDRate) error
	LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error)
	GetAggregatedBurnFee(from, to time.Time, freq string, reserveAddrs []ethereum.Address) (map[ethereum.Address]map[string]float64, error)
	GetAssetVolume(token core.Token, fromTime, toTime uint64, frequency string) (map[uint64]*common.VolumeStats, error)
	GetCountryStats(countryCode string, timezone int64, fromTime, toTime time.Time) (map[uint64]float64, error)
}
