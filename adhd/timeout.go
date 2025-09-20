package adhd

import (
	"errors"
	"time"
)

type deadlineCtx struct {
	*cancelCtx
	deadline time.Time
}

func WithTimeout(parent ADHD, d time.Duration) (ADHD, func()) {
	deadline := time.Now().Add(d)
	c, cancel := WithCancel(parent)
	dc := &deadlineCtx{
		cancelCtx: c.(*cancelCtx),
		deadline:  deadline,
	}

	timer := time.AfterFunc(d, func() {
		dc.mu.Lock()
		defer dc.mu.Unlock()
		if dc.err == nil {
			dc.err = errors.New("deadline exceeded")
			close(dc.done)
		}
	})

	return dc, func() {
		timer.Stop()
		cancel()
	}
}

func (dc *deadlineCtx) Deadline() (deadline time.Time, ok bool) {
	return dc.deadline, true
}
