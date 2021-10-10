package event

import "time"

type PatternFactory interface {
	Name() string
	CreatePatternInstance() PatternInstance
}

type PatternInstance interface {
	Name() string
	Next() *Event
	Done() bool
}

type CreateFunc func() PatternInstance

type patternFactory struct {
	name    string
	creator CreateFunc
}

func NewPatternFactory(name string, creator CreateFunc) PatternFactory {
	return &patternFactory{
		name:    name,
		creator: creator,
	}
}

func (pf *patternFactory) Name() string {
	return pf.name
}

func (pf *patternFactory) CreatePatternInstance() PatternInstance {
	return pf.creator()
}

type patternInstance struct {
	name   string
	events chan *Event
	done   bool
}

func NewPatternInstance(name string, events []*Event, ticker Ticker) PatternInstance {
	c := make(chan *Event, len(events))
	p := patternInstance{
		name:   name,
		events: c,
	}

	go func() {
		for len(events) > 0 {
			select {
			case <-ticker.Channel():
				ev := events[0]
				if ev != nil {
					ev.TimeStamp = time.Now()
					p.events <- ev
				}
				events = events[1:]
			}
		}

		p.done = true
		close(p.events)
	}()

	return &p
}

func (p *patternInstance) Next() *Event {
	select {
	case evt := <-p.events:
		return evt
	default:
		return nil
	}
}

func (p *patternInstance) Name() string {
	return p.name
}

func (p *patternInstance) Done() bool {
	return p.done
}
