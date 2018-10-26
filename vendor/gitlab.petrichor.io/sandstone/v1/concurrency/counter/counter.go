package counter

import (
	"sync/atomic"
)

type Counter struct {
	v uint32
}

func (c *Counter) Value() uint32 {
	return c.v
}

func (c *Counter) Zero() bool {
	if atomic.LoadUint32(&c.v) == 0 {
		return true
	}

	return false
}

func (c *Counter) Hit() uint32 {
	return atomic.AddUint32(&c.v, 1)
}

func (c *Counter) Reset() {
	atomic.StoreUint32(&c.v, 0)
}