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

	"github.com/gokapaya/cshelper/ulist"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print the list of users parsed from the CSV file",
	Long:  ``,
	Run:   runList,
}

var (
	flagLong bool
)

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&flagLong, "long", false, "print full details about each user")

	listCmd.PreRun = func(cmd *cobra.Command, args []string) {
		initUlist()
	}
}

func runList(cmd *cobra.Command, args []string) {
	ul := ulist.GetAllUsers()

	ul.Iter(func(_ int, u ulist.User) error {
		if flagLong {
			fmt.Println("======", u.String())
		} else {
			fmt.Println(u.Username)
		}
		return nil
	})
}
