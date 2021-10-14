package auth

import (
	"fmt"
	"homework-event-generator/event"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

const (
	allowedAuthAttempts = 5
	minExtraAttempts    = 1
	maxExtraAttempts    = 3
)

func NewAccountLockedFactory(ticker event.Ticker) event.PatternFactory {
	return event.NewPatternFactory(AccountLockedPattern, func() event.PatternInstance {
		email := gofakeit.Email()
		ipv4 := gofakeit.IPv4Address()
		ipv6 := fmt.Sprintf("::FFFF:%s", ipv4)
		count := allowedAuthAttempts + rand.Intn(maxExtraAttempts) + minExtraAttempts
		country, err := gofakeit.Weighted(legitimateCountriesAsInterfaceSlice, LegitimateCountriesWeight)
		if err != nil {
			country = "Canada"
		}

		events := make([]*event.Event, count)
		for i := 0; i < count; i++ {
			ev := &event.Event{
				Type:      FailedLoginAttemptEvent,
				TimeStamp: time.Now(),
				Level:     event.InfoLevel,
				ExtraProps: event.ExtraProps{
					"Email":   email,
					"Country": country.(string),
					"IPV4":    ipv4,
					"IPV6":    ipv6,
					"Reason":  ReasonWrongPassword,
				},
			}

			if i >= allowedAuthAttempts {
				ev.Type = AccountLockedEvent
				ev.Level = event.WarningLevel
				ev.ExtraProps["Reason"] = "Too many failed authentication attempts"
			}
			events[i] = ev
		}
		return event.NewPatternInstance(AccountLockedPattern, events, ticker)
	})
}
