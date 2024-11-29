package errs

import "errors"

var (
	ErrReservationNotFound  = errors.New("reservation not found")
	ErrTableAlreadyReserved = errors.New("table already reserved")
	ErrTableNotFound        = errors.New("table not found")
)
