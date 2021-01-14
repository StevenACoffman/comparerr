package main

import (
	"fmt"

	"emperror.dev/errors"
)

// ErrSomethingWentWrong is a sentinel error which can be useful within a single API layer.
const ErrSomethingWentWrong = errors.Sentinel("something went wrong")

// ErrMyError is an error that can be returned from a public API.
type ErrMyError struct {
	Msg string
}

func (e ErrMyError) Error() string {
	return e.Msg
}

func foo() error {
	// Attach stack trace to the sentinel error.
	return errors.WithStack(ErrSomethingWentWrong)
}

func bar() error {
	return errors.Wrap(ErrMyError{"something went wrong"}, "error")
}

func main() {
	if err := foo(); err != nil {
		if errors.Cause(err) == ErrSomethingWentWrong { // or errors.Is(ErrSomethingWentWrong)
			fmt.Printf("%+v\n", err)
		}
	}

	if err := bar(); err != nil {
		if errors.As(err, &ErrMyError{}) {
			fmt.Printf("%+v\n", err)
		}
	}
}
