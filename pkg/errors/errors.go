package errors

import (
	"errors"
	"fmt"
)

type Aggregate interface {
	error
	Errors() []error
	Is(error) bool
}

func NewAggregate(errlist []error) Aggregate {
	if len(errlist) == 0 {
		return nil
	}
	// In case of input error list contains nil
	var errs []error
	for _, e := range errlist {
		if e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return aggregate(errs)
}

type aggregate []error

// Error is part of the error interface.
func (agg aggregate) Error() string {
	if len(agg) == 0 {
		// This should never happen, really.
		return ""
	}
	if len(agg) == 1 {
		return agg[0].Error()
	}
	seenerrs := NewString()
	result := ""
	agg.visit(func(err error) bool {
		msg := err.Error()
		if seenerrs.Has(msg) {
			return false
		}
		seenerrs.Insert(msg)
		if len(seenerrs) > 1 {
			result += ", "
		}
		result += msg
		return false
	})
	if len(seenerrs) == 1 {
		return result
	}
	return "[" + result + "]"
}

func (agg aggregate) Is(target error) bool {
	return agg.visit(func(err error) bool {
		return errors.Is(err, target)
	})
}

func (agg aggregate) visit(f func(err error) bool) bool {
	for _, err := range agg {
		switch err := err.(type) {
		case aggregate:
			if match := err.visit(f); match {
				return match
			}
		case Aggregate:
			for _, nestedErr := range err.Errors() {
				if match := f(nestedErr); match {
					return match
				}
			}
		default:
			if match := f(err); match {
				return match
			}
		}
	}

	return false
}

// Errors is part of the Aggregate interface.
func (agg aggregate) Errors() []error {
	return []error(agg)
}

type withCode struct {
	err   error
	code  int
	cause error
	*stack
}

func WithCode(code int, format string, args ...interface{}) error {
	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  code,
		stack: callers(),
	}
}

func WrapC(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  code,
		cause: err,
		stack: callers(),
	}
}

// Error return the externally-safe error message.
func (w *withCode) Error() string { return fmt.Sprintf("%v", w) }

// Cause return the cause of the withCode error.
func (w *withCode) Cause() error { return w.cause }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withCode) Unwrap() error { return w.cause }
