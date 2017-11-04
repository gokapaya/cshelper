package cmd

import (
	"os"
	"path/filepath"

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

const defaultCsvListName = "ulist.csv"

var (
	Log = log15.New()

	cfg       Config
	maxLogLvl log15.Lvl

	flagDebug       bool
	flagCsvListPath string
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		Log.Error("executing root command failed", "err", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "print debug logs")
	RootCmd.PersistentFlags().StringVar(&flagCsvListPath, "csv-path", filepath.Join(defaultConfigPath, defaultCsvListName), "path to the CSV list with the form results")
	viper.BindPFlags(RootCmd.PersistentFlags())

	RootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		initLogging()
		initConfig(cmd)

		Log.Debug("dumping config", "ignore", flagIgnore, "csv-list", flagCsvListPath)
	}
}
