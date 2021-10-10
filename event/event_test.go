package event

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTestEvent(timestamp time.Time) Event {
	return Event{
		Type:      "test",
		Level:     InfoLevel,
		TimeStamp: timestamp,
		ExtraProps: ExtraProps{
			"foo":  "bar",
			"some": "one",
		},
	}
}

func TestEventExtraAttrs(t *testing.T) {
	now := time.Now()

	evt := createTestEvent(now)

	assert.Equal(t, "bar", evt.ExtraProps["foo"])
	assert.Equal(t, "one", evt.ExtraProps["some"])
}

func TestEventMarshalJson(t *testing.T) {
	now := time.Now()

	evt := createTestEvent(now)

	buf, err := json.Marshal(evt)
	assert.Nil(t, err)

	assert.Equal(t, fmt.Sprintf(`{"type":"test","timestamp":"%s","level":"info","extraProps":{"foo":"bar","some":"one"}}`, now.Format(time.RFC3339Nano)), string(buf))

	evt.ExtraProps = make(ExtraProps)

	buf, err = json.Marshal(evt)
	assert.Nil(t, err)

	assert.Equal(t, fmt.Sprintf(`{"type":"test","timestamp":"%s","level":"info"}`, now.Format(time.RFC3339Nano)), string(buf))
}

func TestEventUnmarshalJson(t *testing.T) {
	now := time.Now()
	buf := []byte(fmt.Sprintf(`{"type":"test","timestamp":"%s","level":"info","extraProps":{"foo":"bar","some":"one"}}`, now.Format(time.RFC3339Nano)))

	var evt Event
	err := json.Unmarshal(buf, &evt)
	assert.Nil(t, err)

	assert.WithinDurationf(t, now, evt.TimeStamp, time.Millisecond, "within 1 ms")
	assert.Equal(t, "bar", evt.ExtraProps["foo"])
	assert.Equal(t, "one", evt.ExtraProps["some"])
}
