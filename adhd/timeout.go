package adhd

import (
	"errors"
	"time"
)

type deadlineCtx struct {
	*cancelCtx
	deadline time.Time
}

func WithDeadline(parent ADHD, deadline time.Time) (ADHD, func()) {
	c, cancel := WithCancel(parent)
	dc := &deadlineCtx{
		cancelCtx: c.(*cancelCtx),
		deadline:  deadline,
	}

	d := time.Until(deadline)
	if d <= 0 {
		dc.mu.Lock()
		dc.err = errors.New("deadline exceeded")
		close(dc.done)
		dc.mu.Unlock()
		return dc, cancel
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

func WithTimeout(parent ADHD, d time.Duration) (ADHD, func()) {
	return WithDeadline(parent, time.Now().Add(d))
}

func (dc *deadlineCtx) Deadline() (deadline time.Time, ok bool) {
	return dc.deadline, true
}
