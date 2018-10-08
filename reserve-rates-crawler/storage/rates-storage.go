package storage

import (
	"github.com/KyberNetwork/reserve-stats/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type ReserveRatesStorage interface {
	UpdateRatesRecords(rateRecords map[string]common.ReserveRates) error
	GetRatesByTimePoint(rsvAddr ethereum.Address, fromTime, toTime int64) ([]common.ReserveRates, error)
}
