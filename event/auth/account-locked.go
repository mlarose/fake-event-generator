package auth

import (
	"fmt"
	"homework-event-generator/event"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type AccountLocked struct {
	attempts int
	count    int
	email    string
	country  string
	ipv4     string
	ipv6     string
}

func (s *AccountLocked) Name() string {
	return LegitimateLoginPattern
}

func (s *AccountLocked) Next() *event.Event {
	// should generate five events
	if s.count >= s.attempts {
		return nil
	}
	s.count++

	ev := &event.Event{
		Type:      FailedLoginAttemptEvent,
		TimeStamp: time.Now(),
		Level:     event.WarningLevel,
		ExtraProps: map[string]interface{}{
			"Email":   s.email,
			"Country": s.country,
			"IPV4":    s.ipv4,
			"IPV6":    s.ipv6,
			"Reason":  "Wrong password",
		},
	}

	// after 5 attempts the account is locked
	if s.count > 5 {
		ev.Type = AccountLockedEvent
		ev.Level = event.ErrorLevel
	}

	return ev
}

func (s *AccountLocked) Done() bool {
	// should generate only has many events as planned attempts
	return s.count >= s.attempts
}

func AccountLockedFactory() event.PatternFactory {
	return event.NewPatternFactory(AccountLockedPattern, func() event.PatternInstance {
		ipv4 := gofakeit.IPv4Address()
		ipv6 := fmt.Sprintf("::FFFF:%s", ipv4)

		country, err := gofakeit.Weighted(legitimateCountriesAsInterfaceSlice, LegitimateCountriesWeight)
		if err != nil {
			country = "Canada"
		}

		return &AccountLocked{
			count:    0,
			attempts: gofakeit.RandomInt([]int{6, 7, 8, 9}),
			email:    gofakeit.Email(),
			country:  country.(string),
			ipv4:     ipv4,
			ipv6:     ipv6,
		}
	})
}
