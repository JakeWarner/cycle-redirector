package events

import (
	"gitlab.petrichor.io/sandstone/v1/builds"
	"gitlab.petrichor.io/sandstone/v1/errors/stackerr"
)

type BindingFn func(data interface{}) (err error)

type binding struct {
	fn BindingFn
	d  *Dispatcher
}

func (b *binding) trigger(trigger string, data interface{}) (err error) {
	if b.fn != nil {
		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(error); ok {
					err = stackerr.Wrapf(stackerr.Hide(r.(error)), "Panic while firing event: %s", trigger)
				} else {
					err = stackerr.Newf("Panic while firing event: %s (%v)", trigger, r)
				}

				builds.DebugStack(false)
			}
		}()

		return b.fn(data)
	}

	return nil
}
