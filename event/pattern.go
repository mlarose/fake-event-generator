package event

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
