package output_port

import (
	"errors"
)

var DatabaseError = errors.New("database error")

var StripeError = errors.New("stripe error")

func WrapDatabaseError(err *error) {
	if *err != nil {
		*err = errors.Join(DatabaseError, *err)
	}
}

func WrapStripeError(err *error) {
	if *err != nil {
		*err = errors.Join(StripeError, *err)
	}
}
