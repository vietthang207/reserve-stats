package storage

import (
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// Interface represent a storage for TradeLogs data
type Interface interface {
	SaveTradeLogs(logs []common.TradeLog, rates []tokenrate.ETHUSDRate) error
	LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error)
}

// VolumeStorage define the functions require to get volume data from db
type VolumeStorage interface {
	GetAssetVolume(token core.Token, fromTime, toTime uint64, frequency string) (map[time.Time]common.VolumeStats, error)
}
