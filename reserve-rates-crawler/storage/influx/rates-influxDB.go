package influx

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/KyberNetwork/reserve-stats/common"
	ethereum "github.com/ethereum/go-ethereum/common"
	influxClient "github.com/influxdata/influxdb/client/v2"
	influxModel "github.com/influxdata/influxdb/models"
)

const (
	//RateDBName is the name of influx database storing reserveRate
	RateDBName = "ReserveRate"
	//RateTableName is the name of influx table storing reserveRate
	RateTableName = "reserve_rate"
	//TimePrecision is the precision configured for influxDB
	TimePrecision = "ms"
)

var errCantConvert error = errors.New("cannot convert response from influxDB to pre-defined struct")

type InfluxRateStorage struct {
	client influxClient.Client
}

func NewRateInfluxDBStorage(url, uName, pwd string) (*InfluxRateStorage, error) {
	httpConf := influxClient.HTTPConfig{
		Addr:     url,
		Username: uName,
		Password: pwd,
	}
	client, err := influxClient.NewHTTPClient(httpConf)
	if err != nil {
		return nil, err
	}
	q := influxClient.NewQuery("CREATE DATABASE "+RateDBName, "", TimePrecision)
	response, err := client.Query(q)
	if err != nil {
		return nil, err
	}
	if response.Error() != nil {
		return nil, response.Error()
	}
	return &InfluxRateStorage{client: client}, nil
}

func (rs *InfluxRateStorage) UpdateRatesRecords(rateRecords map[string]common.ReserveRates) error {
	bp, err := influxClient.NewBatchPoints(
		influxClient.BatchPointsConfig{
			Database:  RateDBName,
			Precision: TimePrecision,
		},
	)
	if err != nil {
		return err
	}

	for rsvAddr, rateRecord := range rateRecords {

		for pair, rate := range rateRecord.Data {
			// InfluxDB get parsing error if the input is uint64. Must use int64
			tags := map[string]string{
				Reserve.String(): rsvAddr,
				Pair.String():    pair,
			}
			fields := map[string]interface{}{
				ReturnTime.String():     rateRecord.ReturnTime,
				BuyRate.String():        rate.BuyReserveRate,
				SellRate.String():       rate.SellReserveRate,
				BuySanityRate.String():  rate.BuySanityRate,
				SellSanityRate.String(): rate.SellSanityRate,
				BlockNumber.String():    rateRecord.BlockNumber,
			}
			pt, err := influxClient.NewPoint(RateTableName, tags, fields, common.TimepointMillisecToTime(rateRecord.Timestamp))
			if err != nil {
				return err
			}
			bp.AddPoint(pt)
		}
	}
	return rs.client.Write(bp)
}

func (rs *InfluxRateStorage) GetRatesByTimePoint(rsvAddr ethereum.Address, fromTime, toTime int64) ([]common.ReserveRates, error) {
	result := []common.ReserveRates{}
	command := fmt.Sprintf("SELECT * FROM %s WHERE time >= %d%s AND \"reserve\"='%s' AND time<= %d%s Order By time", RateTableName, fromTime, TimePrecision, rsvAddr.Hex(), toTime, TimePrecision)
	q := influxClient.NewQuery(command, RateDBName, TimePrecision)
	response, err := rs.client.Query(q)
	if err != nil {
		return result, err
	}
	if response.Error() != nil {
		return result, response.Error()
	}
	if len(response.Results) == 0 || len(response.Results[0].Series) == 0 {
		return []common.ReserveRates{}, nil
	}
	return convertQueryResultTorRate(response.Results[0].Series[0])
}

func getIndexOfFieldS(fieldNames []string) map[RateSchemaFieldName]int {
	result := make(map[RateSchemaFieldName]int)
	for idx, fieldNameStr := range fieldNames {
		fieldName, ok := RateSchemaFields[fieldNameStr]
		if ok {
			result[fieldName] = idx
		}
	}
	return result
}

func getint64FromInterface(v interface{}) (int64, error) {
	number, convertible := v.(json.Number)
	if !convertible {
		return 0, errCantConvert
	}
	return number.Int64()
}

func getFloat64FromInterface(v interface{}) (float64, error) {
	number, convertible := v.(json.Number)
	if !convertible {
		return 0, errCantConvert
	}
	return number.Float64()
}

func convertQueryResultTorRate(row influxModel.Row) ([]common.ReserveRates, error) {
	if len(row.Values) == 0 {
		return []common.ReserveRates{}, nil
	}
	idxs := getIndexOfFieldS(row.Columns)
	rateEntry := make(common.ReserveTokenRateEntry)
	rate := common.ReserveRates{
		Data: rateEntry,
	}
	rates := []common.ReserveRates{rate}
	firstRecordProcessed := false
	nRate := 0
	for _, v := range row.Values {
		// Get Time
		intNumber, err := getint64FromInterface(v[idxs[Time]])
		if err != nil {
			return nil, err
		}
		// New record with new Timestamp
		if rate.Timestamp != intNumber && firstRecordProcessed {
			rates = append(rates, rate)
			rate = common.ReserveRates{}
			rateEntry = make(common.ReserveTokenRateEntry)
			nRate++
		} else {
			rate = rates[nRate]
			rateEntry = rate.Data
		}
		rate.Timestamp = intNumber
		// get Return time
		intNumber, err = getint64FromInterface(v[idxs[ReturnTime]])
		if err != nil {
			return nil, err
		}
		rate.ReturnTime = intNumber
		// get Block number
		intNumber, err = getint64FromInterface(v[idxs[BlockNumber]])
		if err != nil {
			return nil, err
		}
		rate.BlockNumber = intNumber
		// get pair
		pairName, convertible := v[idxs[Pair]].(string)
		if !convertible {
			return nil, errCantConvert
		}
		buyRate, err := getFloat64FromInterface(v[idxs[BuyRate]])
		if err != nil {
			return nil, err
		}
		sellRate, err := getFloat64FromInterface(v[idxs[SellRate]])
		if err != nil {
			return nil, err
		}
		buySanityRate, err := getFloat64FromInterface(v[idxs[BuySanityRate]])
		if err != nil {
			return nil, err
		}
		sellSanityRate, err := getFloat64FromInterface(v[idxs[SellSanityRate]])
		if !convertible {
			return nil, errCantConvert
		}
		rateEntry[pairName] = common.ReserveRateEntry{
			BuyReserveRate:  buyRate,
			SellReserveRate: sellRate,
			BuySanityRate:   buySanityRate,
			SellSanityRate:  sellSanityRate,
		}
		rate.Data = rateEntry
		rates[nRate] = rate
		firstRecordProcessed = true
	}
	return rates, nil
}
