package fsq

import "context"

// IAtomic is an interface that defines the methods for an atomic boolean.
// It is used for shuttingDown value on the queue. Use atomics.go or implement your own.
type IAtomic interface {
	Set(value bool)
	Get() bool
}

// IQueue is an interface that defines the methods for a queue.
// Use rabbit.go or implement your own.
type IQueue interface {
	SendToQueue(to string, subject string, body string) error
	Consume(ctx context.Context) error
}

// ISender is an interface that defines the methods for an email sender service.
// Use smtp.go or implement your own.
type ISender interface {
	SendMail(to string, subject string, body string) error
}
