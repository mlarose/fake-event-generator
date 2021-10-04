package auth

import (
	"fmt"
	"homework-event-generator/event"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type ForeignPasswordReset struct {
	count         int
	attempts      int
	email         string
	recoveryEmail string
	country       string
	ipv4          string
	ipv6          string
}

func (s *ForeignPasswordReset) Name() string {
	return ForeignPasswordResetPattern
}

func (s *ForeignPasswordReset) Next() *event.Event {
	// should generate two events
	if s.count >= s.attempts {
		return nil
	}
	s.count++

	ev := &event.Event{
		Type:      PasswordResetRequestEvent,
		TimeStamp: time.Now(),
		Level:     event.InfoLevel,
		ExtraProps: map[string]interface{}{
			"Email":         s.email,
			"RecoveryEmail": s.recoveryEmail,
			"Country":       s.country,
			"IPV4":          s.ipv4,
			"IPV6":          s.ipv6,
		},
	}

	if s.count > 1 {
		ev.Type = PasswordChangedEvent
	}

	return ev
}

func (s *ForeignPasswordReset) Done() bool {
	return s.count >= s.attempts
}

func ForeignPasswordResetFactory() event.PatternFactory {
	return event.NewPatternFactory(ForeignPasswordResetPattern, func() event.PatternInstance {
		ipv4 := gofakeit.IPv4Address()
		ipv6 := fmt.Sprintf("::FFFF:%s", ipv4)
		return &ForeignPasswordReset{
			count:         0,
			attempts:      gofakeit.RandomInt([]int{1, 2}),
			email:         gofakeit.Email(),
			recoveryEmail: gofakeit.Email(),
			country:       gofakeit.RandomString(RestrictedForeignCountries),
			ipv4:          ipv4,
			ipv6:          ipv6,
		}
	})
}
