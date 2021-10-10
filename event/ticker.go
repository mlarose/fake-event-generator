package event

import (
	"math/rand"
	"time"
)

// Ticker allow to abstract the time.Ticker to allow alternate implementation for testing or other purposes.
type Ticker interface {
	Channel() <-chan time.Time
}

type timeTickerWrapper struct {
	*time.Ticker
}

// WrapTimeTicker provides the Ticker interface to the standard time.Ticker struct
func WrapTimeTicker(t *time.Ticker) Ticker {
	return &timeTickerWrapper{t}
}

// Channel to receives ticks
func (t *timeTickerWrapper) Channel() <-chan time.Time {
	return t.C
}

// A MockTicker implements the Ticker interface for use inside automated tests
// It exposes a SendTick() function to simulate passage of time.
type MockTicker struct {
	c              chan time.Time
	delayAfterTick time.Duration
}

// NewMockTicker creates an initialized and ready to use MockTicker
func NewMockTicker(delayAfterTick time.Duration) *MockTicker {
	return &MockTicker{make(chan time.Time, 1), delayAfterTick}
}

// Channel to receives ticks
func (m *MockTicker) Channel() <-chan time.Time {
	return m.c
}

// SendTick emits a tick through the channel
func (m *MockTicker) SendTick() {
	m.c <- time.Now()
	if m.delayAfterTick > 0 {
		time.Sleep(m.delayAfterTick)
	}
}

type JitterTicker struct {
	c      chan time.Time
	closed bool
}

func NewJitterTicker(min time.Duration, max time.Duration) *JitterTicker {
	channel := make(chan time.Time, 1)
	delta := max - min

	ticker := JitterTicker{c: channel}

	go func() {
		for !ticker.closed {
			delay := min + time.Duration(float64(delta)*rand.Float64())
			ts := <-time.After(delay)
			channel <- ts
		}
		close(channel)
	}()
	return &ticker
}

func (j *JitterTicker) Close() {
	j.closed = true
}

func (j *JitterTicker) Channel() <-chan time.Time {
	return j.c
}
