package lerrors

import "errors"

var (
	ErrDBNoRowsAffected        = errors.New("no rows affected")
	ErrDBMoreThan1RowsAffected = errors.New("more than one row affected")
)
