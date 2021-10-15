package auth

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
)

// Event types produced by the authentication system
const (
	AccountLockedEvent        = "AccountLockedEvent"
	FailedLoginAttemptEvent   = "FailedLoginAttempt"
	SuccessfulLoginEvent      = "SuccessfulLogin"
	PasswordResetRequestEvent = "PasswordResetRequest"
	PasswordChangedEvent      = "PasswordChanged"
)

// Reasons for failed login attempts
const (
	ReasonAccountLocked         = "User account locked"
	ReasonUnregisteredAccount   = "Unregistered account"
	ReasonWrongPassword         = "Wrong password"
	ReasonFailed2FA             = "Failed two factor authentication"
	ReasonTimeout               = "Attempt timed out"
	ReasonTooManyFailedAttempts = "Too many failed authentication attempts"
)

// Patterns emerging from different user interactions
const (
	AccountLockedPattern                  = "AccountLocked"
	RestrictedCountryLoginPattern         = "RestrictedCountryLogin"
	RestrictedCountryPasswordResetPattern = "RestrictedCountryPasswordReset"
	LegitimateLoginPattern                = "LegitimateLogin"
)

var RestrictedCountries = []string{
	"China", "Cuba", "Iran", "North Korea", "Russia", "Sudan", "Syria",
}

var UserDistributionPerCountry = map[string]float32{
	"Australia":      1,
	"Belgium":        2,
	"Canada":         30,
	"Denmark":        1,
	"France":         8,
	"Germany":        2,
	"Israel":         1,
	"Japan":          1,
	"Switzerland":    1,
	"Spain":          1,
	"United Kingdom": 5,
	"United States":  20,
}

var LegitimateCountries = make([]string, 0, len(UserDistributionPerCountry))
var LegitimateCountriesWeight = make([]float32, 0, len(UserDistributionPerCountry))
var legitimateCountriesAsInterfaceSlice = make([]interface{}, 0, len(LegitimateCountries))

var legitimateUsersCount = 1000
var legitimateUsers []string

var legitimateOrgsCount = 50
var legitimateOrgs []string

func init() {
	for _, country := range LegitimateCountries {
		legitimateCountriesAsInterfaceSlice = append(legitimateCountriesAsInterfaceSlice, country)
	}

	for k, v := range UserDistributionPerCountry {
		LegitimateCountries = append(LegitimateCountries, k)
		LegitimateCountriesWeight = append(LegitimateCountriesWeight, v)
	}
}

func GetLegitimateUsers() []string {
	if len(legitimateUsers) != legitimateUsersCount {
		legitimateUsers = make([]string, legitimateUsersCount)
		for i := 0; i < legitimateUsersCount; i++ {
			legitimateUsers[i] = fmt.Sprintf("%s.%s@%s",
				gofakeit.FirstName(),
				gofakeit.LastName(),
				gofakeit.RandomString(GetLegitimateOrgs()),
			)
		}
	}

	return legitimateUsers
}

func GetLegitimateOrgs() []string {
	if len(legitimateOrgs) != legitimateOrgsCount {
		legitimateOrgs = make([]string, legitimateOrgsCount)
		for i := 0; i < legitimateOrgsCount; i++ {
			legitimateOrgs[i] = gofakeit.DomainName()
		}
	}

	return legitimateOrgs
}
