package log

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

type FormatFn func(f *Formatter, in []byte) (out []byte)

type Formatter struct {
	log           *Log
	out           *io.Writer
	formatFn      FormatFn
	cachedCaption string
}

var buffers = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 500))
	},
}

func NewFormatter(l *Log, out *io.Writer, formatter FormatFn) *Formatter {
	return &Formatter{
		log:      l,
		out:      out,
		formatFn: formatter,
	}
}

var rawFormatter = func(f *Formatter, in []byte) (out []byte) {
	buf := buffers.Get().(*bytes.Buffer)

	defer func() {
		buf.Reset()
		buffers.Put(buf)
	}()

	f.prefix(buf)

	buf.Write(in)

	return buf.Bytes()
}

var normalFormatter = func(f *Formatter, in []byte) (out []byte) {
	buf := buffers.Get().(*bytes.Buffer)

	defer func() {
		buf.Reset()
		buffers.Put(buf)
	}()

	f.prefix(buf)
	buf.Write(in)
	buf.WriteString("\n")

	return buf.Bytes()
}

var noticeFormatter = func(f *Formatter, in []byte) (out []byte) {
	buf := buffers.Get().(*bytes.Buffer)

	defer func() {
		buf.Reset()
		buffers.Put(buf)
	}()

	f.prefix(buf)
	buf.WriteString(TextColor(string(in), "lightblue"))
	buf.WriteString("\n")

	return buf.Bytes()
}

var successFormatter = func(f *Formatter, in []byte) (out []byte) {
	buf := buffers.Get().(*bytes.Buffer)

	defer func() {
		buf.Reset()
		buffers.Put(buf)
	}()

	f.prefix(buf)
	buf.WriteString(TextColor(string(in), "green"))
	buf.WriteString("\n")

	return buf.Bytes()
}

var debugFormatter = func(f *Formatter, in []byte) (out []byte) {
	buf := buffers.Get().(*bytes.Buffer)

	defer func() {
		buf.Reset()
		buffers.Put(buf)
	}()

	f.prefix(buf)
	buf.WriteString(TextColor(string(in), "gray"))
	buf.WriteString("\n")

	return buf.Bytes()
}

var warnFormatter = func(f *Formatter, in []byte) (out []byte) {
	buf := buffers.Get().(*bytes.Buffer)

	defer func() {
		buf.Reset()
		buffers.Put(buf)
	}()

	f.prefix(buf)
	buf.WriteString(TextColor(string(in), "yellow"))
	buf.WriteString("\n")

	return buf.Bytes()
}

var errorFormatter = func(f *Formatter, in []byte) (out []byte) {
	buf := buffers.Get().(*bytes.Buffer)

	defer func() {
		buf.Reset()
		buffers.Put(buf)
	}()

	f.prefix(buf)
	buf.WriteString(TextColor(string(in), "red"))
	buf.WriteString("\n")

	return buf.Bytes()
}

func (f *Formatter) Write(in []byte) (written int, err error) {
	written = len(in)

	s := strings.Replace(string(in), "\x00", "", -1)
	buf := bytes.NewBufferString(s)
	scanner := bufio.NewScanner(buf)

	o := *f.out

	for scanner.Scan() {
		out := f.formatFn(f, scanner.Bytes())

		if len(out) > 0 {
			o.Write(out)
		}
	}

	return written, nil
}

func (f *Formatter) WriteString(in string) (written int, err error) {
	return f.Write([]byte(in))
}

func (f *Formatter) prefix(buf *bytes.Buffer) {
	if f.log.Options != nil {
		if f.log.Options.IncludeTime {
			buf.WriteString(TextColor(fmt.Sprintf("[%s] ", time.Now().Format(time.StampMilli)), "gray"))
		}

		if f.log.Options.IncludeCaption {
			buf.WriteString(f.caption())
		}
	}
}

func (f *Formatter) caption() string {
	if f.cachedCaption != "" {
		return f.cachedCaption
	}

	if f.log.Options != nil {
		if f.log.Options.IncludeCaption {
			var caption string

			if f.log.Options.CaptionPadding > 0 {
				caption = fmt.Sprintf("%"+fmt.Sprintf("%d", f.log.Options.CaptionPadding)+"s", f.log.Caption)
			} else {
				caption = f.log.Caption
			}

			f.cachedCaption = TextColor(fmt.Sprintf("[%s] ", caption), "lightgray")
		}
	}

	return f.cachedCaption
}
