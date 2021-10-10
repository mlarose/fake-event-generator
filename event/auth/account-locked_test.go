package auth

import (
	"github.com/stretchr/testify/assert"
	"homework-event-generator/event"
	"testing"
	"time"
)

func TestNewAccountLockedFactory(t *testing.T) {
	ticker := event.NewMockTicker(10 * time.Millisecond)
	ut := NewAccountLockedFactory(ticker)

	assert.NotNil(t, ut)
	assert.Equal(t, AccountLockedPattern, ut.Name())

	instance := ut.CreatePatternInstance()
	assert.NotNil(t, instance)
	assert.Equal(t, AccountLockedPattern, instance.Name())
	assert.False(t, instance.Done())
	assert.Nil(t, instance.Next())

	for i := 0; !instance.Done(); i++ {
		ticker.SendTick()
		ev := instance.Next()
		assert.NotNil(t, ev)
		assert.WithinDurationf(t, time.Now(), ev.TimeStamp, 50*time.Millisecond, "event should be received with tolerance")

		if i >= 5 {
			assert.Equal(t, event.WarningLevel, ev.Level)
			assert.Equal(t, AccountLockedEvent, ev.Type)
		} else {
			assert.Equal(t, event.InfoLevel, ev.Level)
			assert.Equal(t, FailedLoginAttemptEvent, ev.Type)
		}
	}

	assert.Nil(t, instance.Next())
	assert.True(t, instance.Done())
}
