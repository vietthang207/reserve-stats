package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	influxModel "github.com/influxdata/influxdb/models"
)

const (
	tokenVolumeField = "token_volume"
	ethVolumeField   = "eth_volume"
	fiatVolumeField  = "usd_volume"
)

// GetAssetVolume returns the volume of a specific assset(token) between a period and with desired frequency
func (is *InfluxStorage) GetAssetVolume(token core.Token, fromTime, toTime uint64, frequency string) (map[time.Time]common.VolumeStats, error) {
	cmd := fmt.Sprintf("SELECT SUM(token_volume) as %s, SUM(eth_volume) as %s, sum(fiat_volume) as %s FROM test_volume WHERE time >=%d%s AND time <= %d%s AND (dst_addr='%s' OR src_addr='%s') GROUP BY time(1%s)", tokenVolumeField, ethVolumeField, fiatVolumeField, fromTime, timePrecision, toTime, timePrecision, token.Address, token.Address, frequency)
	var (
		logger = is.sugar.With("asset Volume", token.Address, "from", fromTime, "to", toTime)
	)
	logger.Debugf("the query is %s", cmd)
	response, err := is.queryDB(is.influxClient, cmd)
	if err != nil {
		return nil, err
	}
	if len(response) == 0 || len(response[0].Series) == 0 {
		return nil, nil
	}
	return convertQueryResultToVolume(response[0].Series[0])
}

func convertQueryResultToVolume(row influxModel.Row) (map[time.Time]common.VolumeStats, error) {
	result := make(map[time.Time]common.VolumeStats)
	if len(row.Values) == 0 {
		return nil, nil
	}
	for _, v := range row.Values {
		ts, vol, err := convertRowValueToVolume(v)
		if err != nil {
			return nil, err
		}
		result[*ts] = *vol
	}
	return result, nil
}

var errCantConvert = errors.New("cannot convert response from influxDB to pre-defined struct")

func convertRowValueToVolume(v []interface{}) (*time.Time, *common.VolumeStats, error) {
	timestampString, ok := v[0].(string)
	if !ok {
		return nil, nil, errCantConvert
	}
	ts, err := time.Parse(time.RFC3339, timestampString)
	if err != nil {
		return nil, nil, err
	}
	volume, err := getFloat64FromInterface(v[1])
	if err != nil {
		return nil, nil, err
	}
	ethVolume, err := getFloat64FromInterface(v[2])
	if err != nil {
		return nil, nil, err
	}
	usdVolume, err := getFloat64FromInterface(v[3])
	if err != nil {
		return nil, nil, err
	}
	return &ts, &common.VolumeStats{
		Volume:    volume,
		ETHAmount: ethVolume,
		USDAmount: usdVolume,
	}, nil
}

func getFloat64FromInterface(v interface{}) (float64, error) {
	if v == nil {
		return 0, nil
	}
	number, convertible := v.(json.Number)
	if !convertible {
		return 0, errCantConvert
	}
	return number.Float64()
}