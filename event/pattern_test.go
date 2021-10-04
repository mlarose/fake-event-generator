package event

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testPattern struct {
	called   uint
	counter  uint
	stubProp string
}

func (tp *testPattern) Name() string {
	return "test pattern"
}

func (tp *testPattern) Next() *Event {
	tp.called++
	if tp.counter > 0 {
		tp.counter--
		return &Event{
			Type:      tp.stubProp,
			TimeStamp: time.Now(),
			Level:     InfoLevel,
			ExtraProps: map[string]interface{}{
				"called":  tp.called,
				"counter": tp.counter,
			},
		}
	}
	return nil
}

func (tp *testPattern) Done() bool {
	return tp.counter == 0
}

func TestNewPatternFactory(t *testing.T) {
	ut := NewPatternFactory("jimboEventFactory", func() PatternInstance {
		return &testPattern{counter: 2, stubProp: "jimbo"}
	})

	assert.Equal(t, "jimboEventFactory", ut.Name())

	instance := ut.CreatePatternInstance()
	assert.NotNil(t, instance)

	evt1 := instance.Next()
	assert.NotNil(t, evt1)
	assert.Equal(t, "jimbo", evt1.Type)
	assert.False(t, instance.Done())

	evt2 := instance.Next()
	assert.NotNil(t, evt2)
	assert.True(t, instance.Done())

	evt3 := instance.Next()
	assert.Nil(t, evt3)
	assert.True(t, instance.Done())
}