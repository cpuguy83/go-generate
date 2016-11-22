package deepcopy

import (
	"errors"
	"fmt"
)

// errors used by this package
var (
	ErrUnexportedType  = errors.New("use of unexported type from another package")
	ErrUnsettableField = errors.New("use of imported type with an unexported field")
	ErrUnsupportedType = errors.New("unsupported type")
)

// typeError implements causer to integrate with the github.com/pkg/errors API
// without having to import the package.
type typeError struct {
	err error
	msg string
}

func (e *typeError) Cause() error {
	return e.err
}

func (e *typeError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func wrapErr(err error, msg string) error {
	return &typeError{err: err, msg: msg}
}

type causer interface {
	Cause() error
}

func cause(err error) error {
	causer, ok := err.(causer)
	if ok {
		return causer.Cause()
	}
	return err
}
