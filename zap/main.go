package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

func another() {
	fmt.Println(TakeStacktrace(0))

	stack := callers(0)

	fmt.Printf("%+v", stack.StackTrace())
}
func something() {
	another()
}

func main() {
	something()
}

// stack represents a stack of program counters. This mirrors the
// (non-exported) type of the same name in github.com/pkg/errors.
type stack []uintptr

// StackTrace mirrors the code in github.com/pkg/errors.
func (s *stack) StackTrace() StackTrace {
	f := make([]StackFrame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = StackFrame((*s)[i])
	}

	return f
}

// callers mirrors the code in github.com/pkg/errors,
// but makes the depth customizable.
func callers(depth int) *stack {
	const numFrames = 32
	var pcs [numFrames]uintptr
	n := runtime.Callers(2+depth, pcs[:])
	var st stack = pcs[0:n]

	return &st
}

type stacktrace struct {
	pcs    []uintptr // program counters; always a subslice of storage
	frames *runtime.Frames

	// The size of pcs varies depending on requirements:
	// it will be one if the only the first frame was requested,
	// and otherwise it will reflect the depth of the call stack.
	//
	// storage decouples the slice we need (pcs) from the slice we pool.
	// We will always allocate a reasonably large storage, but we'll use
	// only as much of it as we need.
	storage []uintptr
}

// stacktraceDepth specifies how deep of a stack trace should be captured.
type stacktraceDepth int

const (
	// stacktraceFirst captures only the first frame.
	stacktraceFirst stacktraceDepth = iota

	// stacktraceFull captures the entire call stack, allocating more
	// storage for it if needed.
	stacktraceFull
)

// captureStacktrace captures a stack trace of the specified depth, skipping
// the provided number of frames. skip=0 identifies the caller of
// captureStacktrace.
//
// The caller must call Free on the returned stacktrace after using it.
func captureStacktrace(skip int, depth stacktraceDepth) *stacktrace {
	stack := &stacktrace{
		storage: make([]uintptr, 64),
	}

	switch depth {
	case stacktraceFirst:
		stack.pcs = stack.storage[:1]
	case stacktraceFull:
		stack.pcs = stack.storage
	}

	// Unlike other "skip"-based APIs, skip=0 identifies runtime.Callers
	// itself. +2 to skip captureStacktrace and runtime.Callers.
	numFrames := runtime.Callers(
		skip+2,
		stack.pcs,
	)

	// runtime.Callers truncates the recorded stacktrace if there is no
	// room in the provided slice. For the full stack trace, keep expanding
	// storage until there are fewer frames than there is room.
	if depth == stacktraceFull {
		pcs := stack.pcs
		for numFrames == len(pcs) {
			pcs = make([]uintptr, len(pcs)*2)
			numFrames = runtime.Callers(skip+2, pcs)
		}

		// Discard old storage instead of returning it to the pool.
		// This will adjust the pool size over time if stack traces are
		// consistently very deep.
		stack.storage = pcs
		stack.pcs = pcs[:numFrames]
	} else {
		stack.pcs = stack.pcs[:numFrames]
	}

	stack.frames = runtime.CallersFrames(stack.pcs)
	return stack
}

// Free releases resources associated with this stacktrace
// and returns it back to the pool.
func (st *stacktrace) Free() {
	st.frames = nil
	st.pcs = nil
}

// Count reports the total number of frames in this stacktrace.
// Count DOES NOT change as Next is called.
func (st *stacktrace) Count() int {
	return len(st.pcs)
}

// Next returns the next frame in the stack trace,
// and a boolean indicating whether there are more after it.
func (st *stacktrace) Next() (_ runtime.Frame, more bool) {
	return st.frames.Next()
}

func TakeStacktrace(skip int) string {
	stack := captureStacktrace(skip+1, stacktraceFull)
	defer stack.Free()

	buffer := &bytes.Buffer{}
	stackfmt := newStackFormatter(buffer)
	stackfmt.FormatStack(stack)
	return buffer.String()
}

// stackFormatter formats a stack trace into a readable string representation.
type stackFormatter struct {
	b        *bytes.Buffer
	nonEmpty bool // whehther we've written at least one frame already
}

// newStackFormatter builds a new stackFormatter.
func newStackFormatter(b *bytes.Buffer) stackFormatter {
	return stackFormatter{b: b}
}

// FormatStack formats all remaining frames in the provided stacktrace -- minus
// the final runtime.main/runtime.goexit frame.
func (sf *stackFormatter) FormatStack(stack *stacktrace) {
	// Note: On the last iteration, frames.Next() returns false, with a valid
	// frame, but we ignore this frame. The last frame is a runtime frame which
	// adds noise, since it's only either runtime.main or runtime.goexit.
	for frame, more := stack.Next(); more; frame, more = stack.Next() {
		sf.FormatFrame(frame)
	}
}

// FormatFrame formats the given frame.
func (sf *stackFormatter) FormatFrame(frame runtime.Frame) {
	if sf.nonEmpty {
		sf.b.WriteByte('\n')
	}
	sf.nonEmpty = true
	sf.b.WriteString(frame.Function)
	sf.b.WriteByte('\n')
	sf.b.WriteByte('\t')
	sf.b.WriteString(frame.File)
	sf.b.WriteByte(':')
	binary.Write(sf.b, binary.LittleEndian, int64(frame.Line))
}

// ElideSharedStackTraceSuffix removes the suffix of newStack that's already
// present in prevStack. The function returns true if some entries
// were elided.
func ElideSharedStackTraceSuffix(prevStack, newStack StackTrace) (StackTrace, bool) {
	if len(prevStack) == 0 {
		return newStack, false
	}
	if len(newStack) == 0 {
		return newStack, false
	}

	// Skip over the common suffix.
	var i, j int
	for i, j = len(newStack)-1, len(prevStack)-1; i > 0 && j > 0; i, j = i-1, j-1 {
		if newStack[i] != prevStack[j] {
			break
		}
	}
	if i == 0 {
		// Keep at least one entry.
		i = 1
	}

	return newStack[:i], i < len(newStack)-1
}

// StackTrace is the type of the data for a call stack.
// This mirrors the type of the same name in github.com/pkg/errors.
type StackTrace []StackFrame

// StackFrame is the type of a single call frame entry.
// This mirrors the type of the same name in github.com/pkg/errors.
type StackFrame uintptr

// StackTraceProvider is a provider of StackTraces.
// This is, intentionally, defined to be implemented by pkg/errors.stack.
type StackTraceProvider interface {
	StackTrace() StackTrace
}

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//    %s	lists source files for each Frame in the stack
//    %v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+v   Prints filename, function, and line number for each Frame in the stack.
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, f := range st {
				io.WriteString(s, "\n")
				f.Format(s, verb)
			}
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v", []StackFrame(st))
		default:
			st.formatSlice(s, verb)
		}
	case 's':
		st.formatSlice(s, verb)
	}
}

// Format mirrors the code in github.com/pkg/errors.
func (s *stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				f := StackFrame(pc)
				fmt.Fprintf(st, "\n%+v", f)
			}
		}
	}
}

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f StackFrame) pc() uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this StackFrame's pc.
func (f StackFrame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line returns the line number of source code of the
// function for this StackFrame's pc.
func (f StackFrame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// name returns the name of this function, if known.
func (f StackFrame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// Format formats the frame according to the fmt.Formatter interface.
//
//    %s    source file
//    %d    source line
//    %n    function name
//    %v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+s   function name and path of source file relative to the compile time
//          GOPATH separated by \n\t (<funcname>\n\t<path>)
//    %+v   equivalent to %+s:%d
func (f StackFrame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.name())
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.file())
		default:
			io.WriteString(s, path.Base(f.file()))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.line()))
	case 'n':
		io.WriteString(s, funcname(f.name()))
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

// formatSlice will format this StackTrace into the given buffer as a slice of
// Frame, only valid when called with '%s' or '%v'.
func (st StackTrace) formatSlice(s fmt.State, verb rune) {
	io.WriteString(s, "[")
	for i, f := range st {
		if i > 0 {
			io.WriteString(s, " ")
		}
		f.Format(s, verb)
	}
	io.WriteString(s, "]")
}
