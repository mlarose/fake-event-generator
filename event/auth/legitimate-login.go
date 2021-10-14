package auth

import (
	"fmt"
	"homework-event-generator/event"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func NewLegitimateLoginFactory(ticker event.Ticker) event.PatternFactory {
	return event.NewPatternFactory(LegitimateLoginPattern, func() event.PatternInstance {
		email := gofakeit.Email()
		country, err := gofakeit.Weighted(legitimateCountriesAsInterfaceSlice, LegitimateCountriesWeight)
		if err != nil {
			country = "Canada"
		}

		ipv4 := gofakeit.IPv4Address()
		ipv6 := fmt.Sprintf("::FFFF:%s", ipv4)
		unsuccessfulAttempts := rand.Intn(4)

		events := make([]*event.Event, unsuccessfulAttempts)
		for i := 0; i < unsuccessfulAttempts; i++ {
			events[i] = &event.Event{
				Type:      FailedLoginAttemptEvent,
				TimeStamp: time.Now(),
				Level:     event.WarningLevel,
				ExtraProps: event.ExtraProps{
					"Email":   email,
					"Country": country.(string),
					"IPV4":    ipv4,
					"IPV6":    ipv6,
					"Reason":  gofakeit.RandomString([]string{ReasonWrongPassword, ReasonTimeout, ReasonFailed2FA}),
				},
			}
		}

		events = append(events, &event.Event{
			Type:      SuccessfulLoginEvent,
			TimeStamp: time.Now(),
			Level:     event.InfoLevel,
			ExtraProps: event.ExtraProps{
				"Email":   email,
				"Country": country.(string),
				"IPV4":    ipv4,
				"IPV6":    ipv6,
			},
		})

		return event.NewPatternInstance(LegitimateLoginPattern, events, ticker)
	})
}
