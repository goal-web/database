package events

import "time"

type QueryExecuted struct {
	Sql        string
	Bindings   []interface{}
	Connection string
	Time       time.Duration
}

func (this *QueryExecuted) Event() string {
	return "QUERY_EXECUTED"
}

func (this *QueryExecuted) Sync() bool {
	return true
}
