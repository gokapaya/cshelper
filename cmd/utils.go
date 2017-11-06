package cmd

import (
	"os"

	"github.com/gokapaya/cshelper/bot"
	"github.com/gokapaya/cshelper/match"
	"github.com/gokapaya/cshelper/tmpl"
	"github.com/gokapaya/cshelper/ulist"
	"github.com/inconshreveable/log15"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultConfigPath = ".cshelper"

// Config values used at runtime
type Config struct {
	Bot bot.Config
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(defaultConfigPath)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		Log.Error("unable to read config file", "err", err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		Log.Error("decoding config into struct failed", "err", err)
		os.Exit(1)
	}

	Log.Debug("found config file", "cfg", viper.ConfigFileUsed())
}

func initLogging() {
	if flagDebug {
		maxLogLvl = log15.LvlDebug
	} else {
		maxLogLvl = log15.LvlInfo
	}
	Log.SetHandler(log15.StderrHandler)
	Log.SetHandler(log15.LvlFilterHandler(maxLogLvl, Log.GetHandler()))
	Log.SetHandler(log15.CallerFileHandler(Log.GetHandler()))

	ulist.Log = Log.New("pkg", "ulist")
	bot.Log = Log.New("pkg", "bot")
	tmpl.Log = Log.New("pkg", "tmpl")
	match.Log = Log.New("pkg", "match")
}

func initUlist() {
	if err := ulist.Init(flagCsvListPath, flagIgnore); err != nil {
		Log.Error("initializing user list failed", "err", err)
		os.Exit(1)
	}
	Log.Debug("user list loaded", "len", ulist.GetAllUsers().Len())
}

func initBot() {
	if err := bot.Init(&cfg.Bot); err != nil {
		Log.Error("initializing bot failed", "err", err)
		os.Exit(1)
	}
}
