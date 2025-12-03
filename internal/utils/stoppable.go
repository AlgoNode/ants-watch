package utils

import (
	"context"
	"time"
)

type Stoppable struct {
	// Used to signal when we are done
	ctx    context.Context
	cancel context.CancelFunc
}

func MakeStoppable(pCtx context.Context) *Stoppable {
	ctx, cancel := context.WithCancel(pCtx)
	return &Stoppable{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Stop manually triggers stop
func (s *Stoppable) Stop() {
	s.cancel()
}

// Chan returns a read only channel that is closed when the program should exit
func (s *Stoppable) Chan() <-chan struct{} {
	return s.ctx.Done()
}

// Context returns a context tied to the stop handler
func (s *Stoppable) Context() context.Context {
	return s.ctx
}

// Bool returns t/f if the stop handler has triggered
func (s *Stoppable) Stopped() bool {
	return s.ctx.Err() != nil
}

// Returns true if slept without interruption
func (s *Stoppable) Sleep(d time.Duration) bool {
	timer := time.NewTimer(d)
	select {
	case <-s.ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
	}
	return !s.Stopped()
}
