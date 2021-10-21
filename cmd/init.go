package cmd

import (
	crand "crypto/rand"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"math"
	"math/big"
	"time"
)

const timeoutDelay = time.Millisecond * 100

var (
	randomSeed int64
	rootCmd    = &cobra.Command{
		Use:   "event-gen",
		Short: "An event generator simulates security appliances",
		Long: `Generates events to imitate various security services or devices, such as firewalls, threat protection
and authentication systems.`,
	}
)

func init() {
	rootCmd.PersistentFlags().Int64Var(&randomSeed, "seed", 0, "initialize random number generation with this seed")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// if no seed is provided, initialize a seed using crypto/rand
		if randomSeed == 0 {
			bigSeed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
			cobra.CheckErr(err)
			randomSeed = bigSeed.Int64()
		}
	}

	rootCmd.AddCommand(NewHttpCmd(timeoutDelay))
	rootCmd.AddCommand(NewTcpCmd(timeoutDelay))
	rootCmd.AddCommand(NewStdoutCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("Terminated with unexpected error: %s", err)
		log.Exit(1)
	}
}
