package trigger

import (
	"sync"
)

type Trigger struct {
	err error
	ch chan struct{}
	o sync.Once
	listeners []chan error
}

func New() *Trigger {
	return &Trigger{
		ch: make(chan struct{}),
	}
}

func (t *Trigger) Triggered() (<-chan struct{}) {
	return t.ch
}

func (t *Trigger) TriggeredWithError() (<-chan error) {
	errCh := make(chan error, 1)

	select {
	case <-t.ch:
		if t.err != nil {
			errCh <- t.err
		}

		close(errCh)

		return errCh
	default:
		t.listeners = append(t.listeners, errCh)
	}

	return errCh
}

func (t *Trigger) Trigger() {
	t.o.Do(func() {
		close(t.ch)

		for _, v := range t.listeners {
			close(v)
		}
	})
}

func (t *Trigger) TriggerWithError(err error) {
	t.o.Do(func() {
		t.err = err
		close(t.ch)

		for _, v := range t.listeners {
			v <- err
			close(v)
		}
	})
}

func NewClosed() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}