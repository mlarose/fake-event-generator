package event

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

type mockPattern struct {
	name            string
	nextEvent       <-chan *Event
	isDone          <-chan bool
	nextCallCounter *int
	done            bool
}

func (mp *mockPattern) Name() string {
	return mp.name
}

func (mp *mockPattern) Next() *Event {
	*mp.nextCallCounter++
	select {
	case evt := <-mp.nextEvent:
		return evt
	default:
		return nil
	}
}

func (mp *mockPattern) Done() bool {
	select {
	case val := <-mp.isDone:
		mp.done = val
	default:
	}
	return mp.done
}

func MockPatternFactory(nextEvent <-chan *Event, isDone <-chan bool, nextCallCounter *int) PatternFactory {
	return NewPatternFactory("mockPatternFactory", func() PatternInstance {
		return &mockPattern{
			name:            "mockPatternFactory",
			nextCallCounter: nextCallCounter,
			nextEvent:       nextEvent,
			isDone:          isDone,
		}
	})
}

func TestGenerator_EventGeneration(t *testing.T) {
	// arrange
	var nextEvent = make(chan *Event, 1)
	var isDone = make(chan bool, 1)
	var callCounter = 0
	ut := NewGenerator()
	err := ut.RegisterPatternFactory(MockPatternFactory(nextEvent, isDone, &callCounter), 1, 1)
	assert.Nil(t, err)

	var eventTicker = NewMockTicker()
	var patternTicker = NewMockTicker()

	wg := sync.WaitGroup{}

	// start async event generation
	go func() {
		wg.Add(1)
		_ = ut.SetRandomSeed(42134)
		err := ut.Run(eventTicker, patternTicker)
		assert.Nil(t, err)
		wg.Done()
	}()

	gen := ut.(*generator)

	assert.Len(t, gen.triggeredPatterns, 0, "triggered pattern list should be empty initially")
	patternTicker.SendTick()
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, gen.triggeredPatterns, 1, "pattern should be triggered only once")
	patternTicker.SendTick()
	time.Sleep(10 * time.Millisecond)
	assert.Len(t, gen.triggeredPatterns, 1, "pattern should be triggered only once")

	// Send a first event
	nextEvent <- &Event{Type: "foobar"}
	eventTicker.SendTick()
	evt1 := <-ut.Output()
	assert.NotNil(t, evt1)
	assert.Equal(t, "foobar", evt1.Type)
	assert.Equal(t, 1, callCounter)

	// Send a ticket with no produced event
	eventTicker.SendTick()
	select {
	case <-ut.Output():
		assert.Fail(t, "no event should have been produced here")
	case <-time.After(10 * time.Millisecond):
		assert.Equal(t, 2, callCounter)
	}

	// Produce and consume a second event
	nextEvent <- &Event{Type: "barfoo"}
	eventTicker.SendTick()
	evt2 := <-ut.Output()
	assert.NotNil(t, evt2)
	assert.Equal(t, "barfoo", evt2.Type)
	assert.Equal(t, 3, callCounter)

	// Mark event pattern as completed
	isDone <- true
	eventTicker.SendTick()
	select {
	case <-ut.Output():
		assert.Fail(t, "no extra events should have been inserted")
	case <-time.After(10 * time.Millisecond):
		assert.Equal(t, 4, callCounter)
		assert.Len(t, gen.triggeredPatterns, 0, "list of triggered patterns should be cleared")
	}

	ut.Terminate()
	wg.Wait()
}

func TestGenerator_RegisterPatternFactory(t *testing.T) {
	var nextEvent = make(chan *Event, 1)
	var isDone = make(chan bool, 1)
	var callCounter = 0
	ut := NewGenerator()
	mockPatternFactory := MockPatternFactory(nextEvent, isDone, &callCounter)
	err1 := ut.RegisterPatternFactory(mockPatternFactory, 0.0, 1)
	assert.Nil(t, err1, "should be able to register pattern factory")
	err2 := ut.RegisterPatternFactory(mockPatternFactory, 0.0, 1)
	assert.NotNil(t, err2, "should not be able to register same pattern twice")
}

func TestGenerator_TriggerPattern(t *testing.T) {
	// arrange
	var nextEvent = make(chan *Event, 1)
	var isDone = make(chan bool, 1)
	var callCounter = 0
	ut := NewGenerator()
	err := ut.RegisterPatternFactory(MockPatternFactory(nextEvent, isDone, &callCounter), 0.0, 1)
	assert.Nil(t, err)

	// act
	err1 := ut.TriggerPattern("bob")
	err2 := ut.TriggerPattern("mockPatternFactory")
	err3 := ut.TriggerPattern("mockPatternFactory")

	// assert
	assert.NotNil(t, err1, "unregistered pattern should not trigger")
	assert.Nil(t, err2, "registered pattern should trigger")
	assert.NotNil(t, err3, "pattern should not register when it has reached maximum concurrency")

	gen := ut.(*generator)
	assert.Len(t, gen.triggeredPatterns, 1)
}
