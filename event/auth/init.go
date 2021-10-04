package auth

// Event types produced by the authentication system
const (
	AccountLockedEvent        = "AccountLockedEvent"
	FailedLoginAttemptEvent   = "FailedLoginAttempt"
	SuccessfulLoginEvent      = "SuccessfulLogin"
	PasswordResetRequestEvent = "PasswordResetRequest"
	PasswordChangedEvent      = "PasswordChanged"
)

/// Patterns emerging from different user interactions
const (
	AccountLockedPattern        = "AccountLockedPattern"
	ForeignLoginPattern         = "ForeignPasswordReset"
	ForeignPasswordResetPattern = "ForeignPasswordResetPattern"
	LegitimateLoginPattern      = "LegitimateLoginPattern"
)

var RestrictedForeignCountries = []string{
	"China", "Cuba", "Iran", "North Korea", "Russia", "Sudan", "Syria",
}

var UserDistributionPerCountry = map[string]float32{
	"Australia":      1,
	"Belgium":        2,
	"Canada":         20,
	"Denmark":        1,
	"France":         3,
	"Germany":        2,
	"Israel":         1,
	"Japan":          1,
	"Switzerland":    1,
	"Spain":          1,
	"United Kingdom": 5,
	"United States":  10,
}

var LegitimateCountries = make([]string, len(UserDistributionPerCountry))
var LegitimateCountriesWeight = make([]float32, len(UserDistributionPerCountry))
var legitimateCountriesAsInterfaceSlice = make([]interface{}, len(LegitimateCountries))

func init() {
	for _, country := range LegitimateCountries {
		legitimateCountriesAsInterfaceSlice = append(legitimateCountriesAsInterfaceSlice, country)
	}

	for k, v := range UserDistributionPerCountry {
		LegitimateCountries = append(LegitimateCountries, k)
		LegitimateCountriesWeight = append(LegitimateCountriesWeight, v)
	}
}
