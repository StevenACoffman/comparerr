package main

import (
	"fmt"

	"github.com/StevenACoffman/errfmt"
)

func main() {

	err := &errfmt.DetailError{
		Msg:    `reading "file"`,
		Detail: `cmd/prog/reader.go:122`,
		Err: &errfmt.DetailError{
			Msg:    "parsing line 23",
			Detail: "iff x > 3 {\n\tcmd/prog/parser.go:85",
			Err: &errfmt.DetailError{
				Msg:    "syntax error",
				Detail: "cmd/prog/parser.go:214",
			},
		},
	}

	detailerr := errfmt.DetailError{
		Msg:    "This is a message",
		Detail: "Important detail",
		Err:    err,
	}
	fmt.Printf("%+v\n", detailerr)

}
