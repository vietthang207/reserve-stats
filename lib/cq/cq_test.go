package cq

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueries(t *testing.T) {
	var tests = []struct {
		name    string
		cq      *ContinuousQuery
		queries []string
	}{
		{
			name: "simple continous query",
			cq: NewContinuousQuery(
				"test_cq",                      // Name
				"test_db",                      // Database
				"",                             // ResampleEveryInterval
				"",                             // ResampleForInterval
				`SELECT * FROM super_database`, // Query
				"1h",                           // TimeInterval
				nil,                            // OffsetIntervals
			),
			queries: []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" BEGIN SELECT * FROM super_database GROUP BY time(1h) END`},
		},
		{
			name: "continous query with resample every interval",
			cq: NewContinuousQuery(
				"test_cq",                      // Name
				"test_db",                      // Database
				"2h",                           // ResampleEveryInterval
				"",                             // ResampleForInterval
				`SELECT * FROM super_database`, // Query
				"1h",                           // TimeInterval
				nil,                            // OffsetIntervals
			),
			queries: []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 2h BEGIN SELECT * FROM super_database GROUP BY time(1h) END`},
		},
		{
			name: "continous query with resample for interval",
			cq: NewContinuousQuery(
				"test_cq",                      // Name
				"test_db",                      // Database
				"",                             // ResampleEveryInterval
				"2h",                           // ResampleForInterval
				`SELECT * FROM super_database`, // Query
				"1h",                           // TimeInterval
				nil,                            // OffsetIntervals
			),
			queries: []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE FOR 2h BEGIN SELECT * FROM super_database GROUP BY time(1h) END`},
		},
		{
			name: "continous query with both resample every and for intervals",
			cq: NewContinuousQuery(
				"test_cq",                      // Name
				"test_db",                      // Database
				"1h",                           // ResampleEveryInterval
				"2h",                           // ResampleForInterval
				`SELECT * FROM super_database`, // Query
				"1h",                           // TimeInterval
				nil,                            // OffsetIntervals
			),
			queries: []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY time(1h) END`},
		},
		{
			name: "continous query with group by in query clause",
			cq: NewContinuousQuery(
				"test_cq", // Name
				"test_db", // Database
				"1h",      // ResampleEveryInterval
				"2h",      // ResampleForInterval
				`SELECT * FROM super_database GROUP BY "email"`, // Query
				"1h", // TimeInterval
				nil,  // OffsetIntervals
			),
			queries: []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h) END`},
		},
		{
			name: "continous query with one offset interval",
			cq: NewContinuousQuery(
				"test_cq", // Name
				"test_db", // Database
				"1h",      // ResampleEveryInterval
				"2h",      // ResampleForInterval
				`SELECT * FROM super_database GROUP BY "email"`, // Query
				"1h",            // TimeInterval
				[]string{"30m"}, // OffsetIntervals
			),
			queries: []string{`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h,30m) END`},
		},
		{
			name: "continous query with multiple offset intervals",
			cq: NewContinuousQuery(
				"test_cq", // Name
				"test_db", // Database
				"1h",      // ResampleEveryInterval
				"2h",      // ResampleForInterval
				`SELECT * FROM super_database GROUP BY "email"`, // Query
				"1h",                   // TimeInterval
				[]string{"10m", "20m"}, // OffsetIntervals
			),
			queries: []string{
				`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h,10m) END`,
				`CREATE CONTINUOUS QUERY "test_cq" on "test_db" RESAMPLE EVERY 1h FOR 2h BEGIN SELECT * FROM super_database GROUP BY "email", time(1h,20m) END`,
			},
		},
	}

	for _, tc := range tests {
		queries, err := tc.cq.Queries()
		require.NoError(t, err, tc.name)
		assert.Equal(t, queries, tc.queries, tc.name)
	}
}
