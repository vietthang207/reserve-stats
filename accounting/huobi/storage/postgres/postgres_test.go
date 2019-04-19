package postgres

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func TestSaveAndGetAccountingRates(t *testing.T) {
	var (
		testData = []huobi.TradeHistory{
			{
				ID:              15584072551,
				Symbol:          "cmtetsh",
				AccountID:       3375841,
				Amount:          "6000.000",
				Price:           "0.00045",
				CreatedAt:       1540793585678,
				Type:            "buy-limit",
				FieldAmount:     "6000.000",
				FieldCashAmount: "2.73336",
				FieldFees:       "12.00000",
				FinishedAt:      1540796135588,
				UserID:          0,
				Source:          "web",
				State:           "filled",
				CanceledAt:      0,
			},
		}
	)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	sugar := logger.Sugar()

	db, err := sqlx.Connect("postgres", "host=127.0.0.1 port=5432 user=reserve_stats password=reserve_stats dbname=reserve_stats sslmode=disable")
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		td := testData[0]
		td.ID++
		td.CreatedAt += 100
		testData = append(testData, td)
	}
	sugar.Debug(len(testData))

	hdb, err := NewDB(sugar, db)
	require.NoError(t, err)

	defer func() {
		err := hdb.TearDown()
		require.NoError(t, err)
		err = hdb.Close()
		require.NoError(t, err)
	}()

	ts, err := hdb.GetLastStoredTimestamp()
	require.NoError(t, err)
	assert.Equal(t, ts, time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC))

	err = hdb.UpdateTradeHistory(testData)
	require.NoError(t, err)

	lastestTimestamp, err := hdb.GetLastStoredTimestamp()
	require.NoError(t, err)
	assert.Equal(t, uint64(1540793585778), timeutil.TimeToTimestampMs(lastestTimestamp))
	sugar.Debugw("", "", timeutil.TimeToTimestampMs(lastestTimestamp))

	data, err := hdb.GetTradeHistory(timeutil.TimestampMsToTime(1540793585600), timeutil.TimestampMsToTime(1540793585699))
	require.NoError(t, err)
	assert.Equal(t, len(data), 1)
	assert.Equal(t, testData[0].FieldAmount, data[0].FieldAmount)

	// test database does not stored duplicated records(with the same id)
	data, err = hdb.GetTradeHistory(timeutil.TimestampMsToTime(1540793585679), timeutil.TimestampMsToTime(1540793586000))
	require.NoError(t, err)
	assert.Equal(t, len(data), 1)
}
