package auth

import (
	"fmt"
	"homework-event-generator/event"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type ForeignLogin struct {
	count int
}

func (s *ForeignLogin) Name() string {
	return ForeignLoginPattern
}

func (s *ForeignLogin) Next() *event.Event {
	// should generate a single event
	if s.count > 0 {
		return nil
	}
	s.count++

	ipv4 := gofakeit.IPv4Address()
	ipv6 := fmt.Sprintf("::FFFF:%s", ipv4)

	return &event.Event{
		Type:      SuccessfulLoginEvent,
		TimeStamp: time.Now(),
		Level:     event.InfoLevel,
		ExtraProps: map[string]interface{}{
			"Email":   gofakeit.Email(),
			"Country": gofakeit.RandomString(RestrictedForeignCountries),
			"IPV4":    ipv4,
			"IPV6":    ipv6,
		},
	}
}

func (s *ForeignLogin) Done() bool {
	// should generate a single event
	return s.count > 0
}

func ForeignLoginFactory() event.PatternFactory {
	return event.NewPatternFactory(ForeignLoginPattern, func() event.PatternInstance {
		return &ForeignLogin{count: 0}
	})
}
