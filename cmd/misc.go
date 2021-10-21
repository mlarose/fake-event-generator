package cmd

import (
	"fake-event-generator/event"
	"fake-event-generator/event/auth"
	"github.com/brianvoe/gofakeit/v6"
	mrand "math/rand"
	"time"

	"github.com/spf13/cobra"
)

const (
	generatorEventPeriod   = time.Millisecond * 10
	generatorPatternPeriod = time.Millisecond * 100
	factoryTickerBase      = time.Millisecond * 20
	factoryTickerJitter    = time.Millisecond * 2000
)

func createEventGenerator() event.Generator {
	var err error

	mrand.Seed(randomSeed)
	generatorSeed := mrand.Int63()
	fakeItSeed := mrand.Int63()
	gofakeit.Seed(fakeItSeed)

	gen := event.NewGenerator()
	err = gen.SetRandomSeed(generatorSeed)
	cobra.CheckErr(err)

	jitter := event.NewJitterTicker(factoryTickerBase, factoryTickerJitter)

	err = gen.RegisterPatternFactory(auth.NewLegitimateLoginFactory(jitter), 0.4, 3)
	cobra.CheckErr(err)

	err = gen.RegisterPatternFactory(auth.NewRestrictedCountryLoginFactory(jitter), 0.01, 2)
	cobra.CheckErr(err)

	err = gen.RegisterPatternFactory(auth.NewAccountLockedFactory(jitter), 0.05, 1)
	cobra.CheckErr(err)

	err = gen.RegisterPatternFactory(auth.NewRestrictedCountryPasswordResetFactory(jitter), 0.002, 1)
	cobra.CheckErr(err)

	err = gen.RegisterPatternFactory(auth.NewInvalidAccountFactory(jitter), 0.05, 1)
	cobra.CheckErr(err)

	return gen
}

func runEventGenerator(gen event.Generator) error {
	return gen.Run(
		event.WrapTimeTicker(time.NewTicker(generatorEventPeriod)),
		event.WrapTimeTicker(time.NewTicker(generatorPatternPeriod)),
	)
}
