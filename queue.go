package fsq

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type Queue struct {
	shuttingDown IAtomic
	signal       chan os.Signal
	consumer     IQueue
}

func New(shuttingDown IAtomic, consumer IQueue) *Queue {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

	return &Queue{
		shuttingDown: shuttingDown,
		signal:       c,
		consumer:     consumer,
	}
}

func (q *Queue) Run(ct context.Context) error {
	ctx, cancel := context.WithCancel(ct)
	defer cancel()

	go func() {
		select {
		case <-q.signal:
			q.shuttingDown.Set(true)
			cancel()
		case <-ctx.Done():
			q.shuttingDown.Set(true)
		}
	}()

	if err := q.consumer.Consume(ctx); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}
