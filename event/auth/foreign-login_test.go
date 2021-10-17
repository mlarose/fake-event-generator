package auth

import (
	"fake-event-generator/event"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewForeignLoginFactory(t *testing.T) {
	ticker := event.NewMockTicker(10 * time.Millisecond)
	ut := NewRestrictedCountryLoginFactory(ticker)

	assert.NotNil(t, ut)
	assert.Equal(t, RestrictedCountryLoginPattern, ut.Name())

	instance := ut.CreatePatternInstance()
	assert.NotNil(t, instance)
	assert.Equal(t, RestrictedCountryLoginPattern, instance.Name())
	assert.False(t, instance.Done())
	assert.Nil(t, instance.Next())

	ticker.SendTick()
	ev := instance.Next()
	assert.NotNil(t, ev)
	assert.WithinDurationf(t, time.Now(), ev.TimeStamp, 50*time.Millisecond, "event should be received with tolerance")
	assert.Equal(t, SuccessfulLoginEvent, ev.Type)
	assert.Equal(t, event.InfoLevel, ev.Level)
	assert.Containsf(t, RestrictedCountries, ev.ExtraProps["Country"], "should be a restricted country")
}
