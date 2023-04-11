package events

import "time"

type QueryExecuted struct {
	Sql        string
	Bindings   []any
	Connection string
	Time       time.Duration
	Error      error
}

func (this *QueryExecuted) Event() string {
	return "QUERY_EXECUTED"
}

func (this *QueryExecuted) Sync() bool {
	return true
}
