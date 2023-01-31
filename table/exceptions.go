package table

import "github.com/goal-web/contracts"

type Exception struct {
	Sql      string
	Bindings []interface{}
	Err      error
	previous contracts.Exception
}

func (e *Exception) Error() string {
	return e.Err.Error()
}

func (e *Exception) GetPrevious() contracts.Exception {
	return e.previous
}

type CreateException = Exception

type InsertException = Exception

type UpdateException = Exception

type DeleteException = Exception

type SelectException = Exception

type NotFoundException = Exception
