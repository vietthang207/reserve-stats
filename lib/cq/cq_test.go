package cq

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueries(t *testing.T) {
	var tests = []struct {
		testName              string
		name                  string
		database              string
		resampleEveryInterval string
		resampleForInterval   string
		query                 string
		timeInterval          string
		offsetIntervals       []string
		queries               []string
	}{
		{
			testName:              "simple continuous query",
			name:                  "test_cq",
			database:              "test_db",
			resampleEveryInterval: "",
			resampleForInterval:   "",
			query:                 `SELECT * FROM super_database`,
			timeInterval:          "1h",
			offsetIntervals:       nil,
			queries:               []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" BEGIN SELECT * FROM super_database GROUP BY time(1h) END`},
		},
		{
			testName:              "continuous query with resample every interval",
			name:                  "test_cq",
			database:              "test_db",
			resampleEveryInterval: "2h",
			query:                 `SELECT * FROM super_database`,
			timeInterval:          "1h",
			queries:               []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 2h BEGIN SELECT * FROM super_database GROUP BY time(1h) END`},
		},
		{
			testName:            "continuous query with resample for interval",
			name:                "test_cq",
			database:            "test_db",
			resampleForInterval: "2h",
			query:               `SELECT * FROM super_database`,
			timeInterval:        "1h",
			queries:             []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE FOR 2h BEGIN SELECT * FROM super_database GROUP BY time(1h) END`},
		},
		{
			testName:              "continuous query with both resample every and for intervals",
			name:                  "test_cq",
			database:              "test_db",
			resampleEveryInterval: "1h",
			resampleForInterval:   "2h",
			query:                 `SELECT * FROM super_database`,
			timeInterval:          "1h",
			queries:               []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY time(1h) END`},
		},
		{
			testName:              "continuous query with group by in query clause",
			name:                  "test_cq",
			database:              "test_db",
			resampleEveryInterval: "1h",
			resampleForInterval:   "2h",
			query:                 `SELECT * FROM super_database GROUP BY "email"`,
			timeInterval:          "1h",
			queries:               []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h) END`},
		},
		{
			testName:              "continuous query with one offset interval",
			name:                  "test_cq",
			database:              "test_db",
			resampleEveryInterval: "1h",
			resampleForInterval:   "2h",
			query:                 `SELECT * FROM super_database GROUP BY "email"`,
			timeInterval:          "1h",
			offsetIntervals:       []string{"30m"},
			queries:               []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h,30m) END`},
		},
		{
			testName:              "continuous query with multiple offset intervals",
			name:                  "test_cq",
			database:              "test_db",
			resampleEveryInterval: "1h",
			resampleForInterval:   "2h",
			query:                 `SELECT * FROM super_database GROUP BY "email"`,
			timeInterval:          "1h",
			offsetIntervals:       []string{"10m", "20m"},
			queries: []string{
				`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h,10m) END`,
				`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h,20m) END`,
			},
		},
	}

	for _, tc := range tests {
		cq, err := NewContinuousQuery(
			tc.name,
			tc.database,
			tc.resampleEveryInterval,
			tc.resampleForInterval,
			tc.query,
			tc.timeInterval,
			tc.offsetIntervals,
		)
		require.NoError(t, err, tc.testName)
		assert.Equal(t, cq.queries, tc.queries, tc.testName)
	}
}
