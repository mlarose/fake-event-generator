package auth

import (
	"fake-event-generator/event"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand"
	"time"
)

func NewAccessFromMultipleLocationsFactory(ticker event.Ticker) event.PatternFactory {
	return event.NewPatternFactory(AccessFromMultipleLocationsPattern, func() event.PatternInstance {
		type origin struct {
			Country string
			IPV4    string
			IPV6    string
		}

		email := gofakeit.RandomString(GetLegitimateUsers())

		originCount := rand.Intn(5) + 5
		origins := make([]origin, originCount)
		originVisited := make([]bool, originCount)
		for i := 0; i < originCount; i++ {
			ipv4 := gofakeit.IPv4Address()
			ipv6 := fmt.Sprintf("::ffff:%s", ipv4)
			origins[i] = origin{
				Country: gofakeit.Country(),
				IPV4:    ipv4,
				IPV6:    ipv6,
			}
		}

		events := make([]*event.Event, 0, originCount*2)

		done := false
		for done != true {
			typ := gofakeit.RandomString([]string{FailedLoginAttemptEvent, SuccessfulLoginEvent, PasswordResetRequestEvent})

			originIndex := rand.Intn(len(origins))
			origin := origins[originIndex]

			ev := &event.Event{
				Type:      typ,
				TimeStamp: time.Now(),
				Level:     event.InfoLevel,
				ExtraProps: event.ExtraProps{
					"Email":   email,
					"Country": origin.Country,
					"IPV4":    origin.IPV4,
					"IPV6":    origin.IPV6,
				},
			}

			switch typ {
			case FailedLoginAttemptEvent:
				ev.ExtraProps["Reason"] = gofakeit.RandomString([]string{ReasonWrongPassword, ReasonTimeout, ReasonFailed2FA})
				ev.Level = event.WarningLevel

			case PasswordChangedEvent:
				fallthrough
			case PasswordResetRequestEvent:
				ev.ExtraProps["RecoveryEmail"] = gofakeit.Email()
			}

			events = append(events, ev)

			// are we done yet ?
			allset := true
			originVisited[originIndex] = true
			for _, v := range originVisited {
				allset = allset && v
			}
			done = allset
		}

		return event.NewPatternInstance(AccessFromMultipleLocationsPattern, events, ticker)
	})
}
