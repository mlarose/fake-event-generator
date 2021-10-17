package cmd

import (
	"crypto/rand"
	"fake-event-generator/event"
	"fake-event-generator/event/auth"
	"math"
	"math/big"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func createEventGenerator() event.Generator {
	var err error

	jitter := event.NewJitterTicker(time.Millisecond*100, time.Second*4)

	gen := event.NewGenerator()

	if seed := viper.GetInt64("seed"); seed != 0 {
		err = gen.SetRandomSeed(seed)
		cobra.CheckErr(err)
	} else {
		bigSeed, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		cobra.CheckErr(err)

		int64Seed := bigSeed.Int64()
		err = gen.SetRandomSeed(int64Seed)
		cobra.CheckErr(err)
	}

	err = gen.RegisterPatternFactory(auth.NewLegitimateLoginFactory(jitter), 0.4, 3)
	cobra.CheckErr(err)

	err = gen.RegisterPatternFactory(auth.NewRestrictedCountryLoginFactory(jitter), 0.01, 2)
	cobra.CheckErr(err)

	err = gen.RegisterPatternFactory(auth.NewAccountLockedFactory(jitter), 0.05, 1)
	cobra.CheckErr(err)

	err = gen.RegisterPatternFactory(auth.NewRestrictedCountryPasswordResetFactory(jitter), 0.002, 1)
	cobra.CheckErr(err)

	return gen
}

func runEventGenerator(gen event.Generator) error {
	return gen.Run(
		event.WrapTimeTicker(time.NewTicker(10*time.Millisecond)),
		event.WrapTimeTicker(time.NewTicker(100*time.Millisecond)),
	)
}
