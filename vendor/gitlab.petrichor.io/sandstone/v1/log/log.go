package log

import (
	"bytes"
	"fmt"
	"gitlab.petrichor.io/sandstone/v1/builds"
	"gitlab.petrichor.io/sandstone/v1/errors/stackerr"
	"io"
	"os"
)

type Log struct {
	out        io.Writer
	errOut     io.Writer
	Caption    string
	RawOut     *Formatter
	NormalOut  *Formatter
	NoticeOut  *Formatter
	SuccessOut *Formatter
	DebugOut   *Formatter
	ErrorOut   *Formatter
	WarnOut    *Formatter
	Options    *Options
}

type Logger interface {
	Normalf(msg string, args ...interface{}) (err error)
	Noticef(msg string, args ...interface{}) (err error)
	Successf(msg string, args ...interface{}) (err error)
	Debugf(msg string, args ...interface{}) (err error)
	Warnf(msg string, args ...interface{}) (err error)
	Errorf(msg string, args ...interface{}) (err error)
	Write(in []byte) (written int, err error)
}

var Default = New("")

var DefaultOptions = Options{
	IncludeTime:    false,
	IncludeCaption: true,
	CaptionPadding: 30,
}

func New(caption string, args ...interface{}) (l *Log) {
	opts := DefaultOptions

	l = &Log{
		Caption: fmt.Sprintf(caption, args...),
		out:     os.Stdout,
		errOut:  os.Stderr,
		Options: &opts,
	}

	l.RawOut = NewFormatter(l, &l.out, rawFormatter)
	l.NormalOut = NewFormatter(l, &l.out, normalFormatter)
	l.NoticeOut = NewFormatter(l, &l.out, noticeFormatter)
	l.SuccessOut = NewFormatter(l, &l.out, successFormatter)
	l.DebugOut = NewFormatter(l, &l.out, debugFormatter)
	l.WarnOut = NewFormatter(l, &l.errOut, warnFormatter)
	l.ErrorOut = NewFormatter(l, &l.errOut, errorFormatter)

	return l
}

func (l *Log) Normalf(msg string, args ...interface{}) (err error) {
	if err := l.write(l.NormalOut, msg, args...); err != nil {
		return err
	}

	return nil
}

func (l *Log) Noticef(msg string, args ...interface{}) (err error) {
	if err := l.write(l.NoticeOut, msg, args...); err != nil {
		return err
	}

	return nil
}

func (l *Log) Successf(msg string, args ...interface{}) (err error) {
	if err := l.write(l.SuccessOut, msg, args...); err != nil {
		return err
	}

	return nil
}

func (l *Log) Debugf(msg string, args ...interface{}) (err error) {
	if builds.Debug {
		if err := l.write(l.DebugOut, msg, args...); err != nil {
			return err
		}
	}

	return nil
}

func (l *Log) Warnf(msg string, args ...interface{}) (err error) {
	if err := l.write(l.WarnOut, msg, args...); err != nil {
		return err
	}

	return nil
}

func (l *Log) Errorf(msg string, args ...interface{}) (err error) {
	if err := l.write(l.ErrorOut, msg, args...); err != nil {
		return err
	}

	return nil
}

func (l *Log) Write(in []byte) (written int, err error) {
	return l.NormalOut.Write(in)
}

func (l *Log) write(out io.Writer, msg string, args ...interface{}) (err error) {
	buf := bytes.NewBufferString(fmt.Sprintf(msg, args...))

	if _, err := buf.WriteTo(out); err != nil {
		return stackerr.Wrap(stackerr.Hide(err), "Cannot write to output")
	}

	return nil
}

func (l *Log) SetOutput(out, errOut io.Writer) {
	l.out = out
	l.errOut = errOut
}
