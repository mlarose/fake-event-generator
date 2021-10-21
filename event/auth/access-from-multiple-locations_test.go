package auth

import (
	"fake-event-generator/event"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewAccessFromMultipleLocationsFactory(t *testing.T) {
	ticker := event.NewMockTicker(10 * time.Millisecond)
	ut := NewAccessFromMultipleLocationsFactory(ticker)

	assert.NotNil(t, ut)
	assert.Equal(t, AccessFromMultipleLocationsPattern, ut.Name())

	instance := ut.CreatePatternInstance()
	assert.NotNil(t, instance)
	assert.Equal(t, AccessFromMultipleLocationsPattern, instance.Name())
	assert.False(t, instance.Done())
	assert.Nil(t, instance.Next())

	email := ""
	origins := map[string]bool{}

	for !instance.Done() {
		ticker.SendTick()
		ev := instance.Next()
		assert.NotNil(t, ev)
		assert.WithinDurationf(t, time.Now(), ev.TimeStamp, 50*time.Millisecond, "event should be received with tolerance")

		if email != "" {
			assert.Equal(t, email, ev.ExtraProps["Email"])
		} else {
			email = ev.ExtraProps["Email"]
		}

		if _, ok := origins[ev.ExtraProps["Country"]]; !ok {
			origins[ev.ExtraProps["Country"]] = true
		}
	}
	assert.Greater(t, len(origins), 3)

	assert.Nil(t, instance.Next())
	assert.True(t, instance.Done())

}
