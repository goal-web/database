package events

import "time"

type QueryExecuted struct {
	Sql        string
	Bindings   []any
	Connection string
	Time       time.Duration
	Error      error
}

func (event *QueryExecuted) Event() string {
	return "QUERY_EXECUTED"
}

func (event *QueryExecuted) Sync() bool {
	return true
}
