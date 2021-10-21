package auth

import (
	"fake-event-generator/event"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func NewInvalidAccountFactory(ticker event.Ticker) event.PatternFactory {
	return event.NewPatternFactory(InvalidAccountPattern, func() event.PatternInstance {
		org := gofakeit.RandomString(GetLegitimateOrgs())
		firstName := gofakeit.FirstName()
		lastName := gofakeit.LastName()
		email := fmt.Sprintf("%s.%s@%s", firstName, lastName, org)
		country := gofakeit.Country()

		ipv4 := gofakeit.IPv4Address()
		ipv6 := fmt.Sprintf("::ffff:%s", ipv4)
		var events = []*event.Event{
			{
				Type:      FailedLoginAttemptEvent,
				TimeStamp: time.Now(),
				Level:     event.InfoLevel,
				ExtraProps: event.ExtraProps{
					"Email":   email,
					"Country": country,
					"IPV4":    ipv4,
					"IPV6":    ipv6,
					"Reason":  ReasonUnregisteredAccount,
				},
			},
		}

		return event.NewPatternInstance(InvalidAccountPattern, events, ticker)
	})
}
