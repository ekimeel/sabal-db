package services

import (
	"errors"
	"fmt"
)

const (
	errEntityAlreadyExists         = "entity already exists"
	errEntityCreateFailedWithCause = "failed to create entity: %s"
	errEntityMissingFieldWithCause = "missing required field: %s"
	errSqlExecWithCause            = "sql error: %s"
	errNoAffect                    = "sql error: no change"
	errRowScanWithCause            = "error scanning row: %s"
)

func NewError(message string, a ...any) error {
	return errors.New(fmt.Sprint(message, a))
}
