package event

import "time"

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
	c             chan time.Time
	TickSendDelay time.Duration
}

// NewMockTicker creates an initialized and ready to use MockTicker
func NewMockTicker() *MockTicker {
	return &MockTicker{make(chan time.Time, 1), 10 * time.Millisecond}
}

// Channel to receives ticks
func (m *MockTicker) Channel() <-chan time.Time {
	return m.c
}

// SendTick emits a tick through the channel
func (m *MockTicker) SendTick() {
	m.c <- time.Now()
}
