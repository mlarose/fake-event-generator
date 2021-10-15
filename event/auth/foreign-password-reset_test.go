package auth

import (
	"github.com/stretchr/testify/assert"
	"homework-event-generator/event"
	"testing"
	"time"
)

func TestNewForeignPasswordResetFactory(t *testing.T) {
	ticker := event.NewMockTicker(10 * time.Millisecond)
	ut := NewForeignPasswordResetFactory(ticker)

	assert.NotNil(t, ut)
	assert.Equal(t, ForeignPasswordResetPattern, ut.Name())

	instance := ut.CreatePatternInstance()
	assert.NotNil(t, instance)
	assert.Equal(t, ForeignPasswordResetPattern, instance.Name())
	assert.False(t, instance.Done())
	assert.Nil(t, instance.Next())

	ticker.SendTick()
	ev1 := instance.Next()
	assert.NotNil(t, ev1)
	assert.WithinDurationf(t, time.Now(), ev1.TimeStamp, 50*time.Millisecond, "event should be received with tolerance")
	assert.Equal(t, PasswordResetRequestEvent, ev1.Type)
	assert.Equal(t, event.InfoLevel, ev1.Level)
	assert.Containsf(t, RestrictedForeignCountries, ev1.ExtraProps["Country"], "should be a restricted country")

	ticker.SendTick()
	ev2 := instance.Next()
	assert.NotNil(t, ev2)
	assert.WithinDurationf(t, time.Now(), ev2.TimeStamp, 50*time.Millisecond, "event should be received with tolerance")
	assert.Equal(t, PasswordChangedEvent, ev2.Type)
	assert.Equal(t, event.InfoLevel, ev2.Level)
	assert.Containsf(t, RestrictedForeignCountries, ev2.ExtraProps["Country"], "should be a restricted country")

	assert.True(t, ev2.TimeStamp.After(ev1.TimeStamp))
	assert.EqualValues(t, ev1.ExtraProps, ev2.ExtraProps)
}
