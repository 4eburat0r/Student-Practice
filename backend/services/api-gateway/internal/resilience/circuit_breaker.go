package resilience

import (
	"time"

	"github.com/sony/gobreaker"
)

type CircuitBreaker struct {
	breaker *gobreaker.CircuitBreaker
}

func NewCircuitBreaker(name string, maxRequests uint32, failureRatio float64, timeout uint32) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: maxRequests,
		Interval:    time.Duration(timeout) * time.Second,
		Timeout:     time.Duration(timeout) * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			if counts.Requests >= maxRequests {
				failRate := float64(counts.TotalFailures) / float64(counts.Requests)
				return failRate >= failureRatio
			}
			return false
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
		},
	}

	return &CircuitBreaker{
		breaker: gobreaker.NewCircuitBreaker(settings),
	}
}

func (cb *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	return cb.breaker.Execute(fn)
}

func (cb *CircuitBreaker) State() gobreaker.State {
	return cb.breaker.State()
}
