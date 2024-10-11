package fsq

import "sync/atomic"

type AtomicBool struct {
	value atomic.Value
}

func NewAtomicBool() *AtomicBool {
	return &AtomicBool{}
}

func (ab *AtomicBool) Set(value bool) {
	ab.value.Store(value)
}

func (ab *AtomicBool) Get() bool {
	return ab.value.Load().(bool)
}
