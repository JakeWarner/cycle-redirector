package events

import (
	"gitlab.petrichor.io/sandstone/v1/builds"
	"gitlab.petrichor.io/sandstone/v1/errors/errout"
	"gitlab.petrichor.io/sandstone/v1/log"
	"reflect"
	"runtime"
	"sync"
)

type Dispatcher struct {
	ls []*listener
}

var DefaultDispatcher = NewDispatcher()

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Call(fn BindingFn) (l *listener) {
	l = &listener{
		binding: &binding{
			fn: fn,
			d:  d,
		},
	}

	d.ls = append(d.ls, l)

	return l
}

func (d *Dispatcher) listeners(trigger string) (bs []*binding) {
	for _, l := range d.ls {
		for _, v := range l.triggers {
			if v == trigger {
				bs = append(bs, l.binding)
			}
		}
	}

	return bs
}

func (d *Dispatcher) Fire(trigger string, data interface{}) error {
	var lg *log.Log

	if builds.Debug {
		lg = log.New("Event Synchronous Fire (%s)", trigger)
	}

	for _, v := range d.listeners(trigger) {
		if builds.Debug {
			fn := runtime.FuncForPC(reflect.ValueOf(v.fn).Pointer()).Name()
			lg.Debugf("Calling %s", fn)
		}

		if err := v.trigger(trigger, data); err != nil {
			return err
		}
	}

	return nil
}

func (d *Dispatcher) FireAsync(trigger string, data interface{}, ew errout.Wrapper) (doneCh <-chan bool) {
	var lg *log.Log

	if builds.Debug {
		lg = log.New("Event Async Fire (%s)", trigger)
	}

	bs := d.listeners(trigger)

	var wg sync.WaitGroup
	wg.Add(len(bs))

	doneChOut := make(chan bool)

	for _, v := range d.listeners(trigger) {
		if builds.Debug {
			fn := runtime.FuncForPC(reflect.ValueOf(v.fn).Pointer()).Name()
			lg.Debugf("Calling %s", fn)
		}

		go func(b *binding) {
			defer wg.Done()

			if err := b.trigger(trigger, data); err != nil && ew != nil {
				ew.Error(err)
			}
		}(v)
	}

	go func() {
		wg.Wait()
		close(doneChOut)

		if ew != nil {
			ew.Close()
		}
	}()

	return doneChOut
}
