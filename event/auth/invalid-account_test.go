package auth

import (
	"fake-event-generator/event"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewInvalidAccountFactory(t *testing.T) {
	ticker := event.NewMockTicker(20 * time.Millisecond)
	ut := NewInvalidAccountFactory(ticker)

	assert.NotNil(t, ut)
	assert.Equal(t, InvalidAccountPattern, ut.Name())

	instance := ut.CreatePatternInstance()
	assert.NotNil(t, instance)
	assert.Equal(t, InvalidAccountPattern, instance.Name())
	assert.False(t, instance.Done())
	assert.Nil(t, instance.Next())

	for !instance.Done() {
		ticker.SendTick()
		ev := instance.Next()
		assert.NotNil(t, ev)
		assert.WithinDurationf(t, time.Now(), ev.TimeStamp, 50*time.Millisecond, "event should be received with tolerance")
		assert.Equal(t, FailedLoginAttemptEvent, ev.Type)
		assert.Equal(t, event.InfoLevel, ev.Level)
		assert.Equal(t, ReasonUnregisteredAccount, ev.ExtraProps["Reason"])
	}

	assert.Nil(t, instance.Next())
	assert.True(t, instance.Done())
}
