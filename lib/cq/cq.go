package cq

import (
	"bytes"
	"strings"
	"text/template"
)

const cqTemplate = `CREATE CONTINUOUS QUERY "{{.Name}}" on "{{.Database}}" ` +
	`{{if or .ResampleEveryInterval .ResampleForInterval}}RESAMPLE {{if .ResampleEveryInterval}}EVERY {{.ResampleEveryInterval}} {{end}}{{if .ResampleForInterval}}FOR {{.ResampleForInterval}} {{end}}{{end}}` +
	`BEGIN {{.Query}}` +
	`{{if not .GroupByQuery}} GROUP BY {{else}}, {{end}}time({{.TimeInterval}}{{if .OffsetInterval}},{{.OffsetInterval}}{{end}}) END`

// ContinuousQuery represents an InfluxDB continous Query.
// By design ContinuousQuery doesn't try to be smart, it does not attempt to parse/validate any field,
// just act as a templating engine.
//
// Example:
// CREATE CONTINUOUS QUERY <cq_name> ON <database_name>
// RESAMPLE EVERY <interval> FOR <interval>
// BEGIN
// <cq_query>
// END
type ContinuousQuery struct {
	Name                  string
	Database              string
	ResampleEveryInterval string
	ResampleForInterval   string
	// the Query string without the GROUP BY time part which will be added by
	// examining TimeInterval and OffsetIntervals.
	Query           string
	TimeInterval    string
	OffsetIntervals []string
}

func (cq *ContinuousQuery) Queries() ([]string, error) {
	var queries []string

	if len(cq.OffsetIntervals) == 0 {
		cq.OffsetIntervals = []string{""}
	}

	for _, offsetInterval := range cq.OffsetIntervals {
		var query bytes.Buffer

		tmpl, err := template.New("cq.Queries").Parse(cqTemplate)
		if err != nil {
			return nil, err
		}

		err = tmpl.Execute(&query, struct {
			*ContinuousQuery
			GroupByQuery   bool // whether the query included GROUP BY statement
			OffsetInterval string
		}{
			ContinuousQuery: cq,
			GroupByQuery:    strings.Contains(cq.Query, "GROUP BY"),
			OffsetInterval:  offsetInterval,
		})
		if err != nil {
			return nil, err
		}

		queries = append(queries, query.String())
	}

	return queries, nil
}

// NewContinuousQuery creates new ContinousQuery instance.
func NewContinuousQuery(
	name, database, resampleEveryInterval, resampleForInterval, query,
	timeInterval string, offsetIntervals []string) *ContinuousQuery {
	return &ContinuousQuery{
		Name:                  name,
		Database:              database,
		ResampleEveryInterval: resampleEveryInterval,
		ResampleForInterval:   resampleForInterval,
		Query:                 query,
		TimeInterval:          timeInterval,
		OffsetIntervals:       offsetIntervals,
	}
}
