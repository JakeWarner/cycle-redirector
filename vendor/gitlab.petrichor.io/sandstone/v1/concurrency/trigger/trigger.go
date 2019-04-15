package trigger

import (
	"sync"
)

type Trigger struct {
	err       error
	ch        chan struct{}
	o         sync.Once
	listeners []chan error
}

func New() *Trigger {
	return &Trigger{
		ch: make(chan struct{}),
	}
}

func (t *Trigger) Triggered() <-chan struct{} {
	return t.ch
}

func (t *Trigger) TriggeredWithError(out **error) <-chan struct{} {
	if out != nil {
		*out = &t.err // set the address of the pointer to be the same address as the error
	}

	return t.ch
}

func (t *Trigger) Error() error {
	if t.err != nil {
		return t.err
	}

	return nil
}

func (t *Trigger) Trigger() {
	t.o.Do(func() {
		close(t.ch)
	})
}

func (t *Trigger) TriggerWithError(err error) {
	t.o.Do(func() {
		t.err = err
		close(t.ch)
	})
}

func NewClosed() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}

func NewClosedError() <-chan error {
	ch := make(chan error)
	close(ch)
	return ch
}
