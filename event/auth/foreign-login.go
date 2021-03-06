package auth

import (
	"fake-event-generator/event"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func NewRestrictedCountryLoginFactory(ticker event.Ticker) event.PatternFactory {
	return event.NewPatternFactory(RestrictedCountryLoginPattern, func() event.PatternInstance {
		ipv4 := gofakeit.IPv4Address()
		ipv6 := fmt.Sprintf("::ffff:%s", ipv4)

		events := []*event.Event{
			{
				Type:      SuccessfulLoginEvent,
				TimeStamp: time.Now(),
				Level:     event.InfoLevel,
				ExtraProps: event.ExtraProps{
					"Email":   gofakeit.RandomString(GetLegitimateUsers()),
					"Country": gofakeit.RandomString(RestrictedCountries),
					"IPV4":    ipv4,
					"IPV6":    ipv6,
				},
			},
		}

		return event.NewPatternInstance(RestrictedCountryLoginPattern, events, ticker)
	})
}
