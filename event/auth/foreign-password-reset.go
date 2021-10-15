package auth

import (
	"fmt"
	"homework-event-generator/event"

	"github.com/brianvoe/gofakeit/v6"
)

func NewForeignPasswordResetFactory(ticker event.Ticker) event.PatternFactory {
	return event.NewPatternFactory(ForeignPasswordResetPattern, func() event.PatternInstance {
		email := gofakeit.RandomString(GetLegitimateUsers())
		recoveryEmail := gofakeit.Email()
		country := gofakeit.RandomString(RestrictedForeignCountries)

		ipv4 := gofakeit.IPv4Address()
		ipv6 := fmt.Sprintf("::FFFF:%s", ipv4)

		events := []*event.Event{
			{
				Type:  PasswordResetRequestEvent,
				Level: event.InfoLevel,
				ExtraProps: event.ExtraProps{
					"Email":         email,
					"RecoveryEmail": recoveryEmail,
					"Country":       country,
					"IPV4":          ipv4,
					"IPV6":          ipv6,
				},
			},
			{
				Type:  PasswordChangedEvent,
				Level: event.InfoLevel,
				ExtraProps: event.ExtraProps{
					"Email":         email,
					"RecoveryEmail": recoveryEmail,
					"Country":       country,
					"IPV4":          ipv4,
					"IPV6":          ipv6,
				},
			},
		}

		return event.NewPatternInstance(ForeignPasswordResetPattern, events, ticker)
	})
}
