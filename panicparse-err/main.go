package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/kr/pretty"
	"io"
	"io/ioutil"
	"os"
	"runtime/debug"

	"github.com/maruel/panicparse/v2/stack"
)

func main() {
	parseStack := func(rawStack []byte) stack.Stack {
		s, _, err := stack.ScanSnapshot(bytes.NewReader(rawStack), ioutil.Discard, stack.DefaultOpts())
		if err != nil && err != io.EOF {
			panic(err)
		}

		if len(s.Goroutines) > 1 {
			panic(errors.New("provided stacktrace had more than one goroutine"))
		}
		f, err := os.Create("output.html")
		if err != nil {
			panic(err)
		}

		defer f.Close()


		s.ToHTML(f, `<p><em>Generated with style by Steve</em></p></body></html>`)
		return s.Goroutines[0].Signature.Stack
	}

	parsedStack := parseStack(debug.Stack())
	//fmt.Printf("parsedStack: %# v", parsedStack)
	fmt.Printf("parsedStack %# v", pretty.Formatter(parsedStack))
}

