package model

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

var (
	// ErrRecordNotFound indicates that specified record was not found.
	ErrRecordNotFound = errors.New("record not found")
	// ErrDuplicateEntry indicates that duplicate entry for this key exists.
	ErrDuplicateEntry = errors.New("record already exists")
	// ErrUnknown indicates an unknown error occurred at model.
	ErrUnknown = errors.New("unknown model error")
)

func parseError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%w: %s", ErrRecordNotFound, err.Error())
	}

	pqErr, typeOK := err.(*pq.Error)

	if !typeOK {
		return fmt.Errorf("%w: failed to cast error to pq.Error: %s", ErrUnknown, err.Error())
	}

	// reference: https://www.postgresql.org/docs/current/errcodes-appendix.html
	switch pqErr.Code {
	case "23505":
		return fmt.Errorf("%w: %s", ErrDuplicateEntry, err.Error())
	default:
		return fmt.Errorf("%w: undefined pq error: %s", ErrUnknown, err.Error())
	}
}
