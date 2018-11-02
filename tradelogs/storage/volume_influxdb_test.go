package storage

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/cq"
)

func doInfluxHTTPReq(client http.Client, cmd, endpoint, db string) error {
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("q", cmd)
	q.Add("db", db)
	req.URL.RawQuery = q.Encode()
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status code, expected: %d, got: %d", http.StatusOK, rsp.StatusCode)
	}
	return nil
}

func aggregationTestData(is *InfluxStorage) error {
	const (
		endpoint = "http://127.0.0.1:8086/"
	)
	cqs, err := tradelogcq.CreateAssetVolumeCqs(is.dbName)
	if err != nil {
		return err
	}
	for _, cq := range cqs {
		err = cq.Execute(is.influxClient, is.sugar)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestGetAssetVolume(t *testing.T) {
	const (
		dbName = "test_volume"
		// These params are expected to be change when export.dat changes.
		fromTime  = 1539248043000
		toTime    = 1539248666000
		ethAmount = 238.33849929550047
		freq      = "h"
		timeStamp = "2018-10-11T09:00:00Z"
	)

	is, err := newTestInfluxStorage(dbName)
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, is.tearDown())
	}()
	assert.NoError(t, loadTestData(dbName))
	assert.NoError(t, aggregationTestData(is))
	volume, err := is.GetAssetVolume(core.ETHToken, fromTime, toTime, freq)
	assert.NoError(t, err)

	t.Logf("Volume result %v", volume)

	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	result, ok := volume[timeUint]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}

	if result.USDAmount != ethAmount {
		t.Fatal(fmt.Errorf("Expect USD amount to be %.18f, got %.18f", ethAmount, result.USDAmount))
	}
}
