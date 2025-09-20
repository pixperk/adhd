package adhd

import (
	"errors"
	"sync"
	"time"
)

type cancelCtx struct {
	mu     sync.Mutex
	done   chan struct{}
	err    error
	parent ADHD
}

func WithCancel(parent ADHD) (ADHD, func()) {
	c := &cancelCtx{
		parent: parent,
		done:   make(chan struct{}),
	}

	cancel := func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.err != nil {
			c.err = errors.New("context canceled")
			close(c.done)
		}
	}
	return c, cancel
}

func (c *cancelCtx) Done() <-chan struct{} {
	return c.done
}

func (c *cancelCtx) Err() error {
	return c.err
}

func (c *cancelCtx) Value(key any) any {
	return c.parent.Value(key)
}

func (c *cancelCtx) Deadline() (deadline time.Time, ok bool) {
	return c.parent.Deadline()
}
