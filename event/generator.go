package event

import (
	"fmt"
	"math/rand"
)

const defaultOutputChannelBufferSize = 16

var outputChannelBufferSize = defaultOutputChannelBufferSize

func SetOutputChannelBufferSize(size int) {
	if size < 1 {
		panic("output channel buffer size must be larger than 0")
	}
	outputChannelBufferSize = size
}

type Generator interface {
	RegisterPatternFactory(pattern PatternFactory, probability float64, maxConcurrency int) error

	TriggerPattern(patternName string) error

	Output() <-chan *Event

	Run(eventTicker Ticker, patternTicker Ticker) error

	SetRandomSeed(seed int64) error

	Terminate()
}

type patternFactoryRegistration struct {
	patternFactory PatternFactory
	probability    float64
	maxConcurrency int
}

type generator struct {
	patternFactories  map[string]patternFactoryRegistration
	triggeredPatterns []PatternInstance
	output            chan *Event
	started           bool
	terminated        bool
	termination       chan interface{}
	seed              int64
}

func NewGenerator() Generator {
	return &generator{
		patternFactories:  make(map[string]patternFactoryRegistration, 0),
		triggeredPatterns: make([]PatternInstance, 0),
		output:            make(chan *Event, outputChannelBufferSize),
		started:           false,
		terminated:        false,
		termination:       make(chan interface{}, 1),
		seed:              0,
	}
}

func (g *generator) RegisterPatternFactory(pattern PatternFactory, probability float64, maxConcurrency int) error {
	_, ok := g.patternFactories[pattern.Name()]
	if ok {
		return fmt.Errorf("pattern factory is already registered")
	}

	g.patternFactories[pattern.Name()] = patternFactoryRegistration{
		patternFactory: pattern,
		probability:    probability,
		maxConcurrency: maxConcurrency,
	}

	return nil
}

func (g *generator) TriggerPattern(patternName string) error {
	patternRegistration, ok := g.patternFactories[patternName]
	if !ok {
		return fmt.Errorf("pattern '%s' is not registered", patternName)
	}

	concurrency := 0
	for _, triggeredPattern := range g.triggeredPatterns {
		if triggeredPattern.Name() == patternName {
			concurrency++
		}
	}

	if concurrency >= patternRegistration.maxConcurrency {
		return fmt.Errorf("pattern '%s' has reached max concurrency", patternName)
	}

	g.triggeredPatterns = append(g.triggeredPatterns, patternRegistration.patternFactory.CreatePatternInstance())
	return nil
}

func (g *generator) Output() <-chan *Event {
	return g.output
}

func (g *generator) Run(eventTicker Ticker, patternTicker Ticker) error {
	if g.terminated {
		return fmt.Errorf("this event generator has been terminated")
	}

	g.started = true

	var randomSource rand.Source
	if g.seed > 0 {
		randomSource = rand.NewSource(g.seed)
	} else {
		randomSource = rand.NewSource(rand.Int63())
	}
	random := rand.New(randomSource)

	for !g.terminated {
		select {
		case <-g.termination:
			// nothing to do

		case <-patternTicker.Channel():
			// attempt to trigger new pattern instances
			for _, patternFactoryRegistration := range g.patternFactories {
				roll := random.Float64()
				if patternFactoryRegistration.probability >= roll {
					_ = g.TriggerPattern(patternFactoryRegistration.patternFactory.Name())
				}
			}

		case <-eventTicker.Channel():
			// for each triggered patterns
			activePatterns := make([]PatternInstance, 0)
			for _, triggeredPattern := range g.triggeredPatterns {
				evt := triggeredPattern.Next()
				if evt != nil {
					g.output <- evt
				}

				if !triggeredPattern.Done() {
					activePatterns = append(activePatterns, triggeredPattern)
				} else {
					fmt.Printf("pattern %s has completed\n", triggeredPattern.Name())
				}
			}

			// filter completed patterns out
			g.triggeredPatterns = activePatterns
		}

	}

	close(g.output)
	return nil
}

func (g *generator) SetRandomSeed(seed int64) error {
	if seed <= 0 {
		return fmt.Errorf("random seed must be a non-negative integer, %d given", seed)
	}

	if g.started {
		return fmt.Errorf("random seed not applied as generator was already started")
	}

	g.seed = seed
	return nil
}

func (g *generator) Terminate() {
	if !g.started {
		panic("generator was not started")
	}
	g.terminated = true
	g.termination <- true
}
