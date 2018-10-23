package storage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// Interface represent a storage for TradeLogs data
type Interface interface {
	LastBlock() (int64, error)
	SaveTradeLogs(logs []common.TradeLog) error
	LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error)
}
