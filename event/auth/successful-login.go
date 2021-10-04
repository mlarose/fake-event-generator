package auth

import (
	"fmt"
	"homework-event-generator/event"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type LegitimateLogin struct {
	attempts int
	count    int
	email    string
	country  string
	ipv4     string
	ipv6     string
}

func (s *LegitimateLogin) Name() string {
	return LegitimateLoginPattern
}

func (s *LegitimateLogin) Next() *event.Event {
	// should generate only has many events as planned attempts
	if s.count >= s.attempts {
		return nil
	}
	s.count++

	ev := &event.Event{
		Type:      "FailedLoginAttempt",
		TimeStamp: time.Now(),
		Level:     event.WarningLevel,
		ExtraProps: map[string]interface{}{
			"Email":   s.email,
			"Country": s.country,
			"IPV4":    s.ipv4,
			"IPV6":    s.ipv6,
			"Reason":  gofakeit.RandomString([]string{"Wrong password", "MFA attempt failed"}),
		},
	}

	if s.count == s.attempts {
		ev.Type = "SuccessfulLogin"
		ev.Level = event.InfoLevel
		delete(ev.ExtraProps, "Reason")
	}

	return ev
}

func (s *LegitimateLogin) Done() bool {
	// should generate only has many events as planned attempts
	return s.count >= s.attempts
}

func LegitimateLoginFactory() event.PatternFactory {
	return event.NewPatternFactory(LegitimateLoginPattern, func() event.PatternInstance {
		ipv4 := gofakeit.IPv4Address()
		ipv6 := fmt.Sprintf("::FFFF:%s", ipv4)

		country, err := gofakeit.Weighted(legitimateCountriesAsInterfaceSlice, LegitimateCountriesWeight)
		if err != nil {
			country = "Canada"
		}

		return &LegitimateLogin{
			count:    0,
			attempts: gofakeit.RandomInt([]int{1, 2, 3, 4, 5}),
			email:    gofakeit.Email(),
			country:  country.(string),
			ipv4:     ipv4,
			ipv6:     ipv6,
		}
	})
}
