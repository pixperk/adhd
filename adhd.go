package main

import "time"

//ADHD is similar to go's context package
type ADHD interface {
	Done() <-chan struct{}
	Err() error
	Value(key any) any
	Deadline() (deadline time.Time, ok bool)
}
