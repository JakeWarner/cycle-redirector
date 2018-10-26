package stackerr

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

type base struct {
	msg   string
	stack stack
}

type wrap struct {
	base    error
	baseErr error
	err     error
}

type causer interface {
	Cause() error
}

type Wrapper interface {
	Errors() []error
}

func (b *base) Error() string {
	return b.msg
}

func (w *wrap) Cause() error {
	if e, ok := w.err.(*wrap); ok {
		return e.Cause()
	}

	return w.err
}

func (w *wrap) Error() string {
	var parts []string

	if w.baseErr != nil {
		parts = []string{w.baseErr.Error()}
	} else {
		parts = []string{w.base.Error()}
	}

	if _, ok := w.err.(*hide); !ok {
		parts = append(parts, w.err.Error())
	}

	return strings.Join(parts, ": ")
}

func (w *wrap) Errors() (errs []error) {
	if w.baseErr != nil {
		errs = []error{w.baseErr}
	} else {
		errs = []error{w.base}
	}

	err := w.err

	for err != nil {
		errs = append(errs, err)

		w, ok := err.(*wrap)
		if !ok {
			break
		}

		err = w.err
	}

	return errs
}

func newBase(message string, skip int) error {
	err := &base{
		msg: message,
		stack: stack{
			frames: make([]uintptr, 5),
			limit:  5,
		},
	}

	runtime.Callers(2+skip, err.stack.frames)

	return err
}

func New(message string) error {
	return newBase(message, 1)
}

func Newf(message string, args ...interface{}) error {
	err := newBase(fmt.Sprintf(message, args...), 1)
	return err
}

func Wrap(err error, message string) error {
	return &wrap{
		base: newBase(message, 1),
		err:  err,
	}
}

func Wrapf(err error, message string, args ...interface{}) error {
	return &wrap{
		base: newBase(fmt.Sprintf(message, args...), 1),
		err:  err,
	}
}

func WrapErr(err error, top error) error {
	return &wrap{
		base:    newBase("", 1),
		baseErr: top,
		err:     err,
	}
}

func Match(err error, errs ...error) bool {
	inputErrs := []error{err}

	switch e := err.(type) {
	case *wrap:
		inputErrs = append(inputErrs, e.Errors()...)
	}

	for _, err := range inputErrs {
		errStr := err.Error()

		for _, v := range errs {
			if err == v {
				return true
			}

			vStr := v.Error()
			if vStr == errStr {
				return true
			}

			if strings.Contains(errStr, vStr) {
				return true
			}
		}
	}

	return false
}

func Cause(err error) error {
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

func Flatten(err error) error {
	return flatten(err, true)
}

func FlattenAll(err error) error {
	return flatten(err, false)
}

func flatten(err error, sanitize bool) error {
	if e, ok := err.(*Exported); ok {
		if sanitize {
			return New(e.Message)
		} else {
			return New(e.MessageInternal)
		}
	}

	b := &bytes.Buffer{}

loop:
	for err != nil {
		switch v := err.(type) {
		case *wrap:
			if v.baseErr != nil {
				b.WriteString(v.baseErr.Error())
			} else {
				b.WriteString(v.base.Error())
			}

			if sanitize {
				if _, ok := v.err.(*hide); ok {
					break loop
				} else {
					err = v.err
				}
			} else {
				if he, ok := v.err.(*hide); ok {
					err = he.error
				} else {
					err = v.err
				}

			}
			b.WriteString(": ")

			if err == nil {
				break loop
			}
			continue loop
		case error:
			b.WriteString(v.Error())
			break loop
		}
	}

	return newBase(b.String(), 1)
}
