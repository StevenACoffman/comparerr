# comparerr - Comparison golang stacktrace error library output

Golang is great. I mostly love it. However, collecting an error with relevant context and
a nicely formatted stacktrace is kind of a mess of competing approaches. 
This [this feature comparison](https://github.com/cockroachdb/errors#features) is a good overview of the landscape.

I wanted to compare the output of a few different "nicely formatted" approaches with two different kinds of wrapped
errors, Sentinel errors and custom error types.

Some libraries capture the exact state of the stack when an error happens, including every function call. 

Some try to attach relevant contextual information (messages, variables) at strategic places along the call stack,
keeping stack traces compact and maximally useful.

Some eliminate stacktrace duplication from wrapped errors.

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
  /Users/steve/Documents/git/comparerr/emperror/main.go:23
  main.main
  /Users/steve/Documents/git/comparerr/emperror/main.go:31
  runtime.main
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:203
  runtime.goexit
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/asm_amd64.s:1373
  
  starting bar
  
  something went wrong
  got an error in bar
  main.bar
  /Users/steve/Documents/git/comparerr/emperror/main.go:27
  main.main
  /Users/steve/Documents/git/comparerr/emperror/main.go:39
  runtime.main
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:203
  runtime.goexit
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/asm_amd64.s:1373 
  ```
+ [cockroachdb/errors](https://github.com/cockroachdb/errors)
  ```
  Sentinel Something Went Wrong
  (1) attached stack trace
  -- stack trace:
  | main.foo
  | 	/Users/steve/Documents/git/comparerr/cockroachdb-errors/main.go:23
  | main.main
  | 	/Users/steve/Documents/git/comparerr/cockroachdb-errors/main.go:31
  | [...repeated from below...]
  Wraps: (2) attached stack trace
  -- stack trace:
  | main.init
  | 	/Users/steve/Documents/git/comparerr/cockroachdb-errors/main.go:10
  | runtime.doInit
  | 	/Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:5480
  | runtime.main
  | 	/Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:190
  | runtime.goexit
  | 	/Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/asm_amd64.s:1373
  Wraps: (3) Sentinel Something Went Wrong
  Error types: (1) *withstack.withStack (2) *withstack.withStack (3) *errutil.leafError
  
  
  starting bar
  
  got an error in bar: bar something went wrong
  (1) attached stack trace
  -- stack trace:
  | main.bar
  | 	/Users/steve/Documents/git/comparerr/cockroachdb-errors/main.go:27
  | main.main
  | 	/Users/steve/Documents/git/comparerr/cockroachdb-errors/main.go:39
  | runtime.main
  | 	/Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:203
  | runtime.goexit
  | 	/Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/asm_amd64.s:1373
  Wraps: (2) got an error in bar
  Wraps: (3) bar something went wrong
  Error types: (1) *withstack.withStack (2) *errutil.withPrefix (3) main.ErrMyError
  ```
+ [palantir/stacktrace](https://github.com/palantir/stacktrace)
  ```
  foo returned error
  --- at /Users/steve/Documents/git/comparerr/palantir-err/main.go:22 (foo) ---
  Caused by: Sentinel Something Went Wrong
  --- at /Users/steve/Documents/git/comparerr/palantir-err/main.go:9 (init) ---
  
  
  starting bar
  
  got an error in bar
  --- at /Users/steve/Documents/git/comparerr/palantir-err/main.go:26 (bar) ---
  Caused by: bar something went wrong
  ```
+ [pkg/errors](https://github.com/pkg/errors)
  ```
  Sentinel Something Went Wrong
  main.init
  /Users/steve/Documents/git/comparerr/pkg-errors/main.go:10
  runtime.doInit
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:5480
  runtime.main
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:190
  runtime.goexit
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/asm_amd64.s:1373
  main.foo
  /Users/steve/Documents/git/comparerr/pkg-errors/main.go:23
  main.main
  /Users/steve/Documents/git/comparerr/pkg-errors/main.go:31
  runtime.main
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:203
  runtime.goexit
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/asm_amd64.s:1373
  
  
  starting bar
  
  bar something went wrong
  got an error in bar
  main.bar
  /Users/steve/Documents/git/comparerr/pkg-errors/main.go:27
  main.main
  /Users/steve/Documents/git/comparerr/pkg-errors/main.go:39
  runtime.main
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/proc.go:203
  runtime.goexit
  /Users/steve/.asdf/installs/golang/1.14.15/go/src/runtime/asm_amd64.s:1373
  ```

You can clone this repo, cd into any of the directories and run `go run main.go` to see an example


+ [jba errfmt output example](https://play.golang.org/p/aYhdnfLSk8g)
+ [emperror output example](https://play.golang.org/p/OUrwpogR8_E)
+ cockroackdb/errors playground times out? Not sure what that is all about.
+ [pkg/errors](https://play.golang.org/p/TwKMNrVrqE8)
+ [palantir/stacktrace](https://play.golang.org/p/YCdTHCXEd0C)

