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
	"fmt"
	"os"

	"github.com/gokapaya/cshelper/bot"
	"github.com/gokapaya/cshelper/tmpl"
	"github.com/gokapaya/cshelper/ulist"
	"github.com/spf13/cobra"
)

// pmCmd represents the pm command
var pmCmd = &cobra.Command{
	Use:   "pm [USER,...]",
	Short: "Send PMs to user(s)",
	Long: `Takes a comma seperated list of reddit usernames
that get send a PM. If the list is empty it sends to all users on the ulist.

You have to specify the template used for the PMs via the --template flag.

The app will ask for confirmation before sending anything.`,
	Run: runPm,
}

var (
	flagDryRun    bool
	flagTempl     string
)

func init() {
	RootCmd.AddCommand(pmCmd)
	pmCmd.Flags().BoolVarP(&flagDryRun, "dry-run", "n", false, "just print what will be done. Don't actually send messages")
	pmCmd.Flags().StringVarP(&flagTempl, "template", "t", "", "path to template for PMs")

	pmCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if flagTempl == "" {
			Log.Warn("no template found", "file", flagTempl)
			os.Exit(1)
		}

		initUlist()
		if !flagDryRun {
			initBot()
		}
	}
}

func runPm(cmd *cobra.Command, args []string) {
	users := args
	Log.Debug("found args", "", args)

	t := tmpl.Lookup(flagTempl)
	if t == nil {
		Log.Warn("no template found", "name", flagTempl)
		os.Exit(1)
	}

	var sendMsgTo *ulist.Ulist
	switch users[0] {
	case "all":
		// pm all users
		sendMsgTo = ulist.GetAllUsers()
	default:
		sendMsgTo = ulist.GetAllUsers().Filter(func(ulu ulist.User) bool {
			for _, u := range users {
				if ulist.CompareUsernames(ulu.Username, u) {
					return true
				}
			}
			return false
		})
	}

	Log.Info("sending messages to users", "total", sendMsgTo.Len(), "template", t.Name())
	if !ok() {
		os.Exit(2)
	}

	var num int
	err := sendMsgTo.Iter(func(i int, u ulist.User) error {
		num++
		Log.Info("==> " + u.Username)

		if !flagDryRun {
			// do stuff
			if err := bot.PmUserWithTemplate(u, t); err != nil {
				Log.Error("sending message failed", "err", err)
			}
		}
		return nil
	})
	if err != nil {
		Log.Error("failure when iterating over user list", "err", err)
		os.Exit(1)
	}
	Log.Info("sent all messages", "total", num)
}

func ok() bool {
	var yesno string
	fmt.Print("Are you sure? (y/n) ")
	_, err := fmt.Scanf("%s", &yesno)
	if err != nil {
		Log.Error("error getting input", "err", err)
		return false
	}
	switch yesno {
	case "y", "yes", "Y", "Yes":
		return true
	case "n", "no", "N", "No":
		return false
	}

	// always expect the dumbest user
	fmt.Println("unclear...?")
	return ok()
}
