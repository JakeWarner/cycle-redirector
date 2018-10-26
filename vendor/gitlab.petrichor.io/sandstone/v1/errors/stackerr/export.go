package stackerr

import (
	"encoding/gob"
	"encoding/json"
	"reflect"
	"strings"
)

type Exported struct {
	Message         string `bson:"message" json:"message"`
	MessageInternal string `bson:"message_internal" json:"message_internal,omitempty"`
	Stack           Frames `bson:"stack" json:"stack,omitempty"`
}

var SanitizeJSON bool

func (e *Exported) MarshalJSON() ([]byte, error) {
	ne := *e

	if SanitizeJSON {
		ne.MessageInternal = ""
		ne.Stack = nil
	}

	return json.Marshal(ne)
}

func init() {
	gob.Register(Exported{})
}

func Exportable(err error) bool {
	switch e := err.(type) {
	case *base:
		return true
	case *wrap:
		return true
	case *hide:
		return Exportable(e.error)
	default:
		if strings.Contains(reflect.TypeOf(err).String(), "errors.errorString") {
			return true
		}
		return false
	}
}

func Export(err error) *Exported {
	if err == nil {
		panic("Error is nil")
	}

	switch e := err.(type) {
	case Exported:
		return &e
	case *Exported:
		return e
	case *base:
		return &Exported{
			Message:         Flatten(e).Error(),
			MessageInternal: flatten(e, false).Error(),
			Stack:           e.Stacktrace(),
		}
	case *wrap:
		return &Exported{
			Message:         Flatten(e).Error(),
			MessageInternal: flatten(e, false).Error(),
			Stack:           e.base.(*base).Stacktrace(),
		}
	case *hide:
		return Export(e.error)
	default:
		return &Exported{
			Message:         e.Error(),
			MessageInternal: e.Error(),
		}
	}
}

func (e Exported) Error() string {
	return e.Message
}
