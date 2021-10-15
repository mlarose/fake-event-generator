package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	randomSeed int64
	rootCmd    = &cobra.Command{
		Use:   "event-gen",
		Short: "An event generator simulates security appliances",
		Long: `Generates events to imitate various security services or devices, such as firewalls, threat protection
and authentication systems.`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/event-gen.yaml)")
	rootCmd.PersistentFlags().Int64Var(&randomSeed, "seed", 0, "initialize random number generation with this seed")

	err := viper.BindPFlag("seed", rootCmd.PersistentFlags().Lookup("seed"))
	cobra.CheckErr(err)

	rootCmd.AddCommand(NewTcpCmd())
	rootCmd.AddCommand(NewHttpCmd())
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("event-gen")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		log.Infoln("Using config file: ", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("Terminated with unexpected error: %s", err)
		log.Exit(1)
	}
}
