package storage

import (
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

// Interface represent a storage for TradeLogs data
type Interface interface {
	SaveTradeLogs(logs []common.TradeLog, rates []tokenrate.ETHUSDRate) error
	LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error)
	GetAssetVolume(token core.Token, fromTime, toTime uint64, frequency string) (map[time.Time]*common.VolumeStats, error)
	GetReserveVolume(rsvAddr ethereum.Address, token core.Token, fromTime, toTime uint64, frequency string) (map[time.Time]common.VolumeStats, error)
}
