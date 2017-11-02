package cmd

import (
	"os"

	"github.com/inconshreveable/log15"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "cshelper",
	Short: "Helper for r/ClosetSanta",
	Long: `Helper for r/ClosetSanta
	
Check out reddit.com/r/ClosetSanta for more.
We are active from around November to January.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var (
	Log = log15.New()

	cfg       Config
	maxLogLvl log15.Lvl
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		Log.Error("executing root command failed", "err", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().Bool("debug", false, "print debug logs")
	viper.BindPFlags(RootCmd.PersistentFlags())

	RootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		initConfig()
		initLogging()
	}
}
