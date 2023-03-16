package mapstructure

import (
	"fmt"
	"sort"
	"strings"
)

type Error struct {
	Err  error
	Name string
}

func (err *Error) Error() string {
	return fmt.Sprintf("'%s': %s", err.Name, err.Err)
}

func NewError(name string, format any, vals ...any) *Error {
	var err error

	switch format := format.(type) {
	case string:
		err = fmt.Errorf(format, vals...)
	case error:
		err = format
	}

	return &Error{
		Err:  err,
		Name: name,
	}
}

// Errors implements the error interface and can represents multiple
// errors that occur in the course of a single decode.
type Errors []*Error

func (errs *Errors) Error() string {
	if errs == nil {
		return ""
	}

	points := make([]string, len(*errs))
	for i, err := range *errs {
		points[i] += fmt.Sprintf("* %s", err)
	}

	sort.Strings(points)
	return fmt.Sprintf("%d error(s) decoding:\n\t%s", len(*errs), strings.Join(points, "\n\t"))
}

// WrappedErrors implements the errwrap.Wrapper interface to make this
// return value more useful with the errwrap and go-multierror libraries.
func (errs *Errors) WrappedErrors() []error {
	if errs == nil {
		return nil
	}

	result := make([]error, len(*errs))
	for i, err := range *errs {
		result[i] = err
	}

	return result
}

func (errs *Errors) ErrorOrNil() error {
	if errs == nil || len(errs.WrappedErrors()) == 0 {
		return nil
	}

	return errs
}

func appendErrors(errs *Errors, err error) *Errors {
	if errs == nil {
		errs = &Errors{}
	}

	switch err := err.(type) {
	case *Error:
		*errs = append(*errs, err)
	default:
		*errs = append(*errs, &Error{Err: err})
	}

	return errs
}
