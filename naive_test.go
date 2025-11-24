package statemachine

import (
	"sync"
	"testing"
)

type TrafficLightsSystem[T string] struct{ Naive[T] }

func (tl *TrafficLightsSystem[T]) event() T {
	c := make(chan T)
	tl.TransitionC <- func(state T) T {
		defer func() { c <- tl.State }()

		switch tl.State {
		case "": // initial or zero valued state
			return "RED"

		case "RED":
			return "YELLOW"

		case "YELLOW":
			return "GREEN"

		case "GREEN":
			return "RED"

		default: // handle undefined state
			panic("undefined state...")
		}
	}

	return <-c
}

func TestNaive(t *testing.T) {
	tls := TrafficLightsSystem[string]{
		Naive[string]{
			State:       "",
			TransitionC: make(chan func(string) string),
		},
	}
	go tls.Run(t.Context())

	wg := sync.WaitGroup{}
	for range 5 {
		wg.Go(func() { tls.event() })
	}
	wg.Wait()

	if got, want := tls.event(), "YELLOW"; got != want {
		t.Fatalf("expecting %s, got %s", want, got)
	}
}
