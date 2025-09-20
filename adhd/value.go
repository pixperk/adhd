package adhd

import "time"

type valueCtx struct {
	parent   ADHD
	key, val any
}

func WithValue(parent ADHD, key, val any) ADHD {
	return &valueCtx{parent, key, val}
}

func (v *valueCtx) Done() <-chan struct{}       { return v.parent.Done() }
func (v *valueCtx) Err() error                  { return v.parent.Err() }
func (v *valueCtx) Deadline() (time.Time, bool) { return v.parent.Deadline() }
func (v *valueCtx) Value(key any) any {
	if key == v.key {
		return v.val
	}
	return v.parent.Value(key)
}
