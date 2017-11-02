package cmd

import (
	"os"

	"github.com/inconshreveable/log15"
	"github.com/spf13/viper"
)

type Config struct {
	Debug bool
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName("cshelper")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".cshelper")

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
	if cfg.Debug {
		maxLogLvl = log15.LvlDebug
	} else {
		maxLogLvl = log15.LvlWarn
	}
	Log.SetHandler(log15.LvlFilterHandler(maxLogLvl, Log.GetHandler()))
	Log.SetHandler(log15.CallerFileHandler(Log.GetHandler()))
}
