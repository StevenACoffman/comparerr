package main

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

// ErrSomethingWentWrong is a sentinel error which can be useful within a single API layer.
var ErrSomethingWentWrong = errors.New("Something went wrong")

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
	return errors.Wrap(ErrMyError{"something went wrong"}, "got an error in bar")
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
