package errout

import (
	"gitlab.petrichor.io/sandstone/v1/log"
	"gitlab.petrichor.io/sandstone/v1/errors/stackerr"
	"sync"
)

type Wrapper interface {
	Error(error)
	Close() error
}

type errCh struct {
	ch        chan error
	closeOnce sync.Once
}

type errLog struct {
	l log.Logger
}

func Chan(ch chan error) Wrapper {
	return &errCh{
		ch: ch,
	}
}

func (o *errCh) Error(err error) {
	o.ch <- err
}

func (o *errCh) Close() error {
	o.closeOnce.Do(func() {
		close(o.ch)
	})

	return nil
}

func Log(l log.Logger) Wrapper {
	return &errLog{
		l: l,
	}
}

func (o *errLog) Error(err error) {
	switch err.(type) {
	case stackerr.Wrapper:
		o.l.Errorf(stackerr.FlattenAll(err).Error())
	default:
		o.l.Errorf(err.Error())
	}
}

func (o *errLog) Close() error {
	return nil
}
