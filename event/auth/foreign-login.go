package auth

import (
	"fmt"
	"homework-event-generator/event"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func NewForeignLoginFactory(ticker event.Ticker) event.PatternFactory {
	return event.NewPatternFactory(ForeignLoginPattern, func() event.PatternInstance {
		ipv4 := gofakeit.IPv4Address()
		ipv6 := fmt.Sprintf("::FFFF:%s", ipv4)

		events := []*event.Event{
			{
				Type:      SuccessfulLoginEvent,
				TimeStamp: time.Now(),
				Level:     event.InfoLevel,
				ExtraProps: event.ExtraProps{
					"Email":   gofakeit.Email(),
					"Country": gofakeit.RandomString(RestrictedForeignCountries),
					"IPV4":    ipv4,
					"IPV6":    ipv6,
				},
			},
		}

		return event.NewPatternInstance(ForeignLoginPattern, events, ticker)
	})
}
