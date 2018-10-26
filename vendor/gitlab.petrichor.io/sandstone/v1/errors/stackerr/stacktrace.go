package stackerr

import (
	"bytes"
	"fmt"
	"runtime"
)

type Traceable interface {
	Stacktrace() Frames
}

type Frame struct {
	File     string `bson:"file" json:"file"`
	Function string `bson:"function" json:"function"`
	Line     int    `bson:"line" json:"line"`
}

type Frames []Frame

type stack struct {
	frames []uintptr
	limit  int
}

func (w *wrap) Stacktrace() Frames {
	return w.base.(*base).Stacktrace()
}

func (b *base) Stacktrace() (fs Frames) {
	cfs := runtime.CallersFrames(b.stack.frames)

	var (
		f    runtime.Frame
		more bool = true
	)

	for i := 0; more && i < b.stack.limit; i++ {
		f, more = cfs.Next()

		if f.PC == 0 {
			break
		}

		fs = append(fs, Frame{
			File:     f.File,
			Function: f.Function,
			Line:     f.Line,
		})
	}

	return fs
}

func (f Frames) Output() string {
	b := &bytes.Buffer{}

	for _, v := range f {
		fmt.Fprintf(b, "%s() : %d\n\t%s\n", v.Function, v.Line, v.File)
	}

	return b.String()
}
