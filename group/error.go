package group

import (
	"fmt"
	"strings"
)

// Error ...
type Error struct {
	errs []error
}

// Error ...
func (e *Error) Error() string {
	if len(e.errs) == 1 {
		return fmt.Sprintf("1 error occurred:\n\t* %s\n\n", e.errs[0])
	}

	ee := make([]string, 0, len(e.errs))
	for i, err := range e.errs {
		ee[i] = fmt.Sprintf("* %s", err)
	}

	return fmt.Sprintf(
		"%d errors occurred:\n\t%s\n\n",
		len(e.errs), strings.Join(ee, "\n\t"))
}

// Append ...
//
//nolint:errorlint
func Append(err error, errs ...error) *Error {
	switch err := err.(type) {
	case *Error:
		if err == nil {
			err = new(Error)
		}

		// Go through each error and flatten
		for _, e := range errs {
			switch e := e.(type) {
			case *Error:
				if e != nil {
					err.errs = append(err.errs, e.errs...)
				}
			default:
				if e != nil {
					err.errs = append(err.errs, e)
				}
			}
		}

		return err
	default:
		errs := make([]error, 0, len(errs)+1)
		if err != nil {
			errs = append(errs, err)
		}
		errs = append(errs, errs...)

		return Append(&Error{}, errs...)
	}
}
