//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Limiter struct {
	maxCount  int
	interval  time.Duration
	stopCh    chan struct{}
	acquireCh chan struct{}
	mu        sync.Mutex
	stopped   bool
}

var ErrStopped = errors.New("limiter stopped")

func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	return &Limiter{
		maxCount:  maxCount,
		interval:  interval,
		stopCh:    make(chan struct{}),
		acquireCh: make(chan struct{}, maxCount),
		stopped:   false,
	}
}

func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.stopCh:
		return ErrStopped
	default:
		l.mu.Lock()
		defer l.mu.Unlock()

		if l.stopped {
			return ErrStopped
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case l.acquireCh <- struct{}{}:
			time.AfterFunc(l.interval, func() {
				<-l.acquireCh
			})
			return nil
		}
	}
}

func (l *Limiter) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.stopped {
		l.stopped = true
		close(l.stopCh)
		close(l.acquireCh)
	}
}
