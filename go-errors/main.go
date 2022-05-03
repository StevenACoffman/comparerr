package main

import (
	"fmt"
	"github.com/go-errors/errors"
)

// ErrSomethingWentWrong is a sentinel error which can be useful within a single API layer.
var ErrSomethingWentWrong = errors.Errorf("something went wrong %s", "...yup")

// ErrMyError is an error that can be returned from a public API.
type ErrMyError struct {
	Msg string
}

func (e ErrMyError) Error() string {
	return e.Msg
}

func foo() error {
	// Attach stack trace to the sentinel error.
	return errors.New(ErrSomethingWentWrong)
}

func bar() error {
	return errors.Wrap(errors.New(ErrMyError{"something went wrong"}), 0)
}

func main() {
	if err := foo(); err != nil {
		if errors.Is(err, ErrSomethingWentWrong) { // or errors.Is(ErrSomethingWentWrong)
			fmt.Printf("%+v\n", err.(*errors.Error).ErrorStack())
		}
	}

	fmt.Print("\n\nstarting bar\n\n")

	if err := bar(); err != nil {
		if errors.As(err, &ErrMyError{}) {
			fmt.Printf("%+v\n", err.(*errors.Error).ErrorStack())
		}
	}
}
