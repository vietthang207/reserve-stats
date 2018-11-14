package storage

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	influxModel "github.com/influxdata/influxdb/models"
)

const (
	tokenVolumeField = "token_volume"
	ethVolumeField   = "eth_volume"
	fiatVolumeField  = "usd_volume"
)

var (
	errCantConvert       = errors.New("cannot convert response from influxDB to pre-defined struct")
	assetMeasurementName = map[string]string{
		"h": "volume_hour",
		"d": "volume_day",
	}
	rsvMeasurementName = map[string]string{
		"h": "rsv_volume_hour",
		"d": "rsv_volume_day",
	}
)

// GetReserveVolume returns the volume of a specific asset(token) from a reserve
// between a period and with desired frequency
func (is *InfluxStorage) GetReserveVolume(rsvAddr ethereum.Address, token core.Token, fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error) {
	var (
		rsvAddrHex   = rsvAddr.Hex()
		tokenAddrHex = ethereum.HexToAddress(token.Address).Hex()
		logger       = is.sugar.With("reserve Address", rsvAddr.Hex(), "func", "tradelogs/storage/InfluxStorage.GetReserveVolume", "token Address", token.Address, "from", fromTime, "to", toTime)
	)
	mName, ok := rsvMeasurementName[strings.ToLower(frequency)]
	if !ok {
		return nil, fmt.Errorf("frequency %s is not supported", frequency)
	}

	addrFilter := fmt.Sprintf("((dst_addr='%s' OR src_addr='%s') AND (dst_rsv_addr='%s' OR src_rsv_addr='%s'))", tokenAddrHex, tokenAddrHex, rsvAddrHex, rsvAddrHex)
	timeFilter := fmt.Sprintf("(time >='%s' AND time <= '%s')", fromTime.UTC().Format(time.RFC3339), toTime.UTC().Format(time.RFC3339))
	cmd := fmt.Sprintf("SELECT SUM(token_volume) as %s, SUM(eth_volume) as %s, SUM(usd_volume) as %s FROM %s WHERE %s AND %s GROUP BY time(1%s) FILL(0)", tokenVolumeField, ethVolumeField, fiatVolumeField, mName, timeFilter, addrFilter, frequency)

	logger.Debugw("query rendered", "query", cmd)

	response, err := is.queryDB(is.influxClient, cmd)
	if err != nil {
		return nil, err
	}
	if len(response) == 0 || len(response[0].Series) == 0 {
		return nil, nil
	}
	return convertQueryResultToVolume(response[0].Series[0])
}

// GetAssetVolume returns the volume of a specific assset(token) between a period and with desired frequency
func (is *InfluxStorage) GetAssetVolume(token core.Token, fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error) {
	var (
		logger = is.sugar.With(
			"func", "tradelogs/storage/InfluxStorage.GetAssetVolume",
			"token", token.Address,
			"from", fromTime,
			"to", toTime,
		)
	)
	mName, ok := assetMeasurementName[strings.ToLower(frequency)]
	if !ok {
		return nil, fmt.Errorf("frequency %s is not supported", frequency)
	}
	var (
		tokenAddr  = ethereum.HexToAddress(token.Address).Hex()
		timeFilter = fmt.Sprintf("(time >='%s' AND time <= '%s')", fromTime.UTC().Format(time.RFC3339), toTime.UTC().Format(time.RFC3339))
		addrFilter = fmt.Sprintf("(dst_addr='%s' OR src_addr='%s')", tokenAddr, tokenAddr)
		cmd        = fmt.Sprintf("SELECT SUM(token_volume) as %s, SUM(eth_volume) as %s, sum(usd_volume) as %s FROM %s WHERE %s AND %s GROUP BY time(1%s) fill(0)", tokenVolumeField, ethVolumeField, fiatVolumeField, mName, timeFilter, addrFilter, frequency)
	)

	logger.Debugw("get asset volume query rendered", "query", cmd)
	response, err := is.queryDB(is.influxClient, cmd)

	if err != nil {
		return nil, err
	}

	logger.Debugw("got result for asset volume query", "response", response)

	if len(response) == 0 || len(response[0].Series) == 0 {
		return nil, nil
	}
	return convertQueryResultToVolume(response[0].Series[0])
}

func convertQueryResultToVolume(row influxModel.Row) (map[uint64]*common.VolumeStats, error) {
	result := make(map[uint64]*common.VolumeStats)
	if len(row.Values) == 0 {
		return nil, nil
	}
	for _, v := range row.Values {
		ts, vol, err := convertRowValueToVolume(v)
		if err != nil {
			return nil, err
		}
		result[ts] = vol
	}
	return result, nil
}

func convertRowValueToVolume(v []interface{}) (uint64, *common.VolumeStats, error) {
	// number of fields in record result
	// - time
	// - token_volume
	// - eth_volume
	// - usd_volume
	if len(v) != 4 {
		return 0, nil, errors.New("value fields is empty")
	}

	timestampString, ok := v[0].(string)
	if !ok {
		return 0, nil, errCantConvert
	}
	ts, err := time.Parse(time.RFC3339, timestampString)
	if err != nil {
		return 0, nil, err
	}
	tsUint64 := timeutil.TimeToTimestampMs(ts)
	volume, err := influxdb.GetFloat64FromInterface(v[1])
	if err != nil {
		return 0, nil, err
	}
	ethVolume, err := influxdb.GetFloat64FromInterface(v[2])
	if err != nil {
		return 0, nil, err
	}
	usdVolume, err := influxdb.GetFloat64FromInterface(v[3])
	if err != nil {
		return 0, nil, err
	}
	return tsUint64, &common.VolumeStats{
		Volume:    volume,
		ETHAmount: ethVolume,
		USDAmount: usdVolume,
	}, nil
}
