# comparerr - Compare golang error library output

Compare error output of different libraries

Golang is great. I mostly love it. However, collecting an error with relevant context and
a nicely formatted stacktrace is kind of a mess of competing approaches. 
This [this feature comparison](https://github.com/cockroachdb/errors#features) is good overview of the landscape.

I wanted to compare the output of a few different "nicely formatted" approaches.
+ [jba/errfmt](https://github.com/jba/errfmt)
  ```
  {Msg:This is a message Detail:Important detail Err:reading "file"
  	cmd/prog/reader.go:122
  parsing line 23
  	iff x > 3 {
  	cmd/prog/parser.go:85
  syntax error
  	cmd/prog/parser.go:214
  }
  ```
+ [emperror](https://github.com/emperror/emperror)
  ```
  something went wrong
  main.foo
  	/tmp/sandbox910738749/prog.go:23
  main.main
  	/tmp/sandbox910738749/prog.go:31
  runtime.main
  	/usr/local/go-faketime/src/runtime/proc.go:204
  runtime.goexit
  	/usr/local/go-faketime/src/runtime/asm_amd64.s:1374
  something went wrong
  error
  main.bar
  	/tmp/sandbox910738749/prog.go:27
  main.main
  	/tmp/sandbox910738749/prog.go:37
  runtime.main
  	/usr/local/go-faketime/src/runtime/proc.go:204
  runtime.goexit
  	/usr/local/go-faketime/src/runtime/asm_amd64.s:1374
  
  ```
+ [cockroachdb/errors](https://github.com/cockroachdb/errors)
  ```
  got an error in bar: something went wrong
  (1) attached stack trace
  -- stack trace:
  | main.bar
  | 	/Users/steve/Documents/git/comparerr/cockroachdb-errors/main.go:28
  | main.main
  | 	/Users/steve/Documents/git/comparerr/cockroachdb-errors/main.go:38
  | runtime.main
  | 	/Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:203
  | runtime.goexit
  | 	/Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/asm_amd64.s:1373
  Wraps: (2) got an error in bar
  Wraps: (3) something went wrong
  Error types: (1) *withstack.withStack (2) *errutil.withPrefix (3) main.ErrMyError
  ```
+ [pkg/errors](https://github.com/pkg/errors)
  ```
  Something went wrong
  main.init
  /Users/steve/Documents/git/comparerr/pkg-errors/main.go:10
  runtime.doInit
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:5480
  runtime.main
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:190
  runtime.goexit
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/asm_amd64.s:1373
  main.foo
  /Users/steve/Documents/git/comparerr/pkg-errors/main.go:24
  main.main
  /Users/steve/Documents/git/comparerr/pkg-errors/main.go:32
  runtime.main
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:203
  runtime.goexit
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/asm_amd64.s:1373
  something went wrong
  got an error in bar
  main.bar
  /Users/steve/Documents/git/comparerr/pkg-errors/main.go:28
  main.main
  /Users/steve/Documents/git/comparerr/pkg-errors/main.go:38
  runtime.main
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:203
  runtime.goexit
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/asm_amd64.s:1373
  ```

You can clone this repo, cd into any of the directories and run `go run main.go` to see an example


+ [jba errfmt output example](https://play.golang.org/p/aYhdnfLSk8g)
+ [emperror output example](https://play.golang.org/p/OUrwpogR8_E)
+ cockroackdb/errors playground times out? Not sure what that is all about.
+ [pkg/errors](https://play.golang.org/p/iHt_qG6mI2J)

