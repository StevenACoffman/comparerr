package main

import (
	"fmt"
	"github.com/palantir/stacktrace"
)

// ErrSomethingWentWrong is a sentinel error which can be useful within a single API layer.
var ErrSomethingWentWrong = stacktrace.NewError("Sentinel Something Went Wrong")

// ErrMyError is an error that can be returned from a public API.
type ErrMyError struct {
	Msg string
}

func (e ErrMyError) Error() string {
	return e.Msg
}

func foo() error {
	// Attach stack trace to the sentinel error.
	return stacktrace.Propagate(ErrSomethingWentWrong, "foo returned error")
}

func bar() error {
	return stacktrace.Propagate(ErrMyError{"bar something went wrong"}, "got an error in bar")
}

func main() {
	// check
	if err := foo(); err != nil {
			fmt.Printf("%+v\n", err)
	}

	fmt.Print("\n\nstarting bar\n\n")

	if err := bar(); err != nil {
			fmt.Printf("%+v\n", err)
	}
}
