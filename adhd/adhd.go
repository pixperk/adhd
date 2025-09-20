package adhd

import (
	"errors"
	"time"
)

var ErrCanceled = errors.New("context canceled")
var ErrDeadlineExceeded = errors.New("context deadline exceeded")

type ADHD interface {
	Done() <-chan struct{}
	Err() error
	Value(key any) any
	Deadline() (deadline time.Time, ok bool)
}
