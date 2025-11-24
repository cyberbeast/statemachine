package statemachine

import "context"

type Naive[T any] struct {
	State       T
	TransitionC chan func(T) T
}

func (sm *Naive[T]) Run(ctx context.Context) {
	for {
		select {
		case fn := <-sm.TransitionC:
			sm.State = fn(sm.State)
		case <-ctx.Done():
			return
		}
	}
}
