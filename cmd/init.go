package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	err := viper.BindPFlag("seed", rootCmd.PersistentFlags().Lookup("seed"))
	cobra.CheckErr(err)

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
