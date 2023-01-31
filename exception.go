package database

import (
	"github.com/goal-web/contracts"
)

type ConnectionErrorCode int

const (
	DbDriverDontExist ConnectionErrorCode = iota
	DbConnectionDontExist
)

type DBConnectionException struct {
	Err        error
	Connection string
	Code       ConnectionErrorCode
	Config     contracts.Fields
	previous   contracts.Exception
}

func (D DBConnectionException) Error() string {
	return D.Err.Error()
}

func (D DBConnectionException) GetPrevious() contracts.Exception {
	return D.previous
}
