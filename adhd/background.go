package adhd

import "time"

type emptyCtx struct{}

func (e *emptyCtx) Done() <-chan struct{}       { return nil }
func (e *emptyCtx) Err() error                  { return nil }
func (e *emptyCtx) Value(key any) any           { return nil }
func (e *emptyCtx) Deadline() (time.Time, bool) { return time.Time{}, false }

func Background() ADHD {
	return &emptyCtx{}
}

func TODO() ADHD {
	return &emptyCtx{}
}
