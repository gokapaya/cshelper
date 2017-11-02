// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"

	"github.com/gokapaya/cshelper/match"
	"github.com/gokapaya/cshelper/ulist"
	"github.com/spf13/cobra"
)

// matchCmd represents the match command
var matchCmd = &cobra.Command{
	Use:   "match",
	Short: "Generate a list of pairings",
	Long:  ``,
	Run:   runMatch,
}

var (
	flagRegift bool
	flagOutput string
)

var (
	defaultFlagOutput = ".cshelper/pairs.csv"
)

func init() {
	RootCmd.AddCommand(matchCmd)

	matchCmd.Flags().BoolVarP(&flagRegift, "regift", "r", false, "Generate pairs for regifting")
	matchCmd.Flags().StringVarP(&flagOutput, "output", "o", defaultFlagOutput, "File to write the pairings")

	matchCmd.PreRun = func(cmd *cobra.Command, args []string) {
		initUlist()
	}
}

func runMatch(cmd *cobra.Command, args []string) {
	ul := ulist.GetAllUsers()
	if flagRegift {
		ul = ul.Filter(func(u ulist.User) bool {
			return u.Regift
		})
		Log.Info("filtered user list", "len", ul.Len())
	}

	Log.Info("matching users")
	p, err := match.Match(ul)
	if err != nil {
		Log.Error("matching failed", "err", err)
	}

	if err := match.Eval(p); err != nil {
		Log.Error("evaluating the pairings failed", "err", err)
	}
	Log.Info("evaluation successful")

	if err := match.SavePairings(defaultFlagOutput, p); err != nil {
		Log.Error("saving pairlist failed", "err", err)
		os.Exit(1)
	}
	Log.Info("pair csv saved", "file", defaultFlagOutput)
}
