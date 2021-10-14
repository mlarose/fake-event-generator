package auth

import (
	"homework-event-generator/event"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewForeignLoginFactory(t *testing.T) {
	ticker := event.NewMockTicker(10 * time.Millisecond)
	ut := NewForeignLoginFactory(ticker)

	assert.NotNil(t, ut)
	assert.Equal(t, ForeignLoginPattern, ut.Name())

	instance := ut.CreatePatternInstance()
	assert.NotNil(t, instance)
	assert.Equal(t, ForeignLoginPattern, instance.Name())
	assert.False(t, instance.Done())
	assert.Nil(t, instance.Next())

	ticker.SendTick()
	ev := instance.Next()
	assert.NotNil(t, ev)
	assert.WithinDurationf(t, time.Now(), ev.TimeStamp, 50*time.Millisecond, "event should be received with tolerance")
	assert.Equal(t, SuccessfulLoginEvent, ev.Type)
	assert.Equal(t, event.InfoLevel, ev.Level)
	assert.Containsf(t, RestrictedForeignCountries, ev.ExtraProps["Country"], "should be a restricted country")
}
