package auth

import (
	"fake-event-generator/event"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLegitimateLoginFactory(t *testing.T) {
	ticker := event.NewMockTicker(10 * time.Millisecond)
	ut := NewLegitimateLoginFactory(ticker)

	assert.NotNil(t, ut)
	assert.Equal(t, LegitimateLoginPattern, ut.Name())

	instance := ut.CreatePatternInstance()
	assert.NotNil(t, instance)
	assert.Equal(t, LegitimateLoginPattern, instance.Name())
	assert.False(t, instance.Done())
	assert.Nil(t, instance.Next())

	for !instance.Done() {
		ticker.SendTick()
		ev := instance.Next()
		assert.NotNil(t, ev)
		assert.WithinDurationf(t, time.Now(), ev.TimeStamp, 50*time.Millisecond, "event should be received with tolerance")
		assert.NotContains(t, RestrictedCountries, ev.ExtraProps["Country"], "should be a restricted country")
		assert.Contains(t, LegitimateCountries, ev.ExtraProps["Country"])
		assert.Contains(t, []string{SuccessfulLoginEvent, FailedLoginAttemptEvent}, ev.Type)

		if ev.Type == FailedLoginAttemptEvent {
			assert.Equal(t, event.WarningLevel, ev.Level)
			assert.Contains(t, []string{ReasonWrongPassword, ReasonTimeout, ReasonFailed2FA}, ev.ExtraProps["Reason"])
		} else {
			assert.Equal(t, event.InfoLevel, ev.Level)
		}
	}

	assert.Nil(t, instance.Next())
	assert.True(t, instance.Done())
}
