package storage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// Interface represent a storage for TradeLogs data
type Interface interface {
	SaveTradeLogs(logs []common.TradeLog, rates []tokenrate.ETHUSDRate) error
	LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error)
	GetAggregatedWalletFee(reserveAddr, walletAddr, freq string,
		fromTime, toTime time.Time, timezone int64) (map[string]float64, error)
}
