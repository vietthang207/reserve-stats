package storage

import "github.com/KyberNetwork/reserve-stats/users/common"

// Interface is the common interface of users persistent storage.
type Interface interface {
	CreateOrUpdate(common.UserData) error
	IsKYCed(string) (bool, error)
}
