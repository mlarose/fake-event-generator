package auth

import (
	"fake-event-generator/event"
	"fmt"
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
		email := gofakeit.RandomString(GetLegitimateUsers())
		ipv4 := gofakeit.IPv4Address()
		ipv6 := fmt.Sprintf("::ffff:%s", ipv4)
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
				ExtraProps: event.ExtraProps{
					"Email":   email,
					"Country": country.(string),
					"IPV4":    ipv4,
					"IPV6":    ipv6,
				},
			}

			if i < allowedAuthAttempts {
				ev.Type = FailedLoginAttemptEvent
				ev.Level = event.InfoLevel
				ev.ExtraProps["Reason"] = gofakeit.RandomString([]string{ReasonWrongPassword, ReasonFailed2FA, ReasonTimeout})
			} else if i == allowedAuthAttempts {
				ev.Type = AccountLockedEvent
				ev.Level = event.WarningLevel
				ev.ExtraProps["Reason"] = ReasonTooManyFailedAttempts
			} else {
				ev.Type = FailedLoginAttemptEvent
				ev.Level = event.InfoLevel
				ev.ExtraProps["Reason"] = ReasonAccountLocked
			}
			events[i] = ev
		}
		return event.NewPatternInstance(AccountLockedPattern, events, ticker)
	})
}
