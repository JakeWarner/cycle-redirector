package stackerr

import (
	"encoding/gob"
	"strings"
)

type Accumulator struct {
	Errors []error `bson:"errors" json:"errors"`
}

func init() {
	gob.Register(Accumulator{})
}

func (a *Accumulator) Add(err error) {
	a.Errors = append(a.Errors, err)
}

func (a Accumulator) Error() string {
	es := make([]string, 0, len(a.Errors))

	for _, v := range a.Errors {
		es = append(es, v.Error())
	}

	return strings.Join(es, ", ")
}
