// Copyright © 2017 Twisted Systems Ltd
//
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
	"github.com/spf13/cobra"
	"os"
	"log"
	"github.com/TwistedSystemsLtd/lazenby-go/core"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new lazenfile",
	Long:  `Create a new lazenfile if it doesn't exist already.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		var lazenfile = cmd.Flag("file").Value.String()
		lazenpath := core.GetLazenpath(lazenfile)

		_, statErr := os.Stat(lazenpath)

		if os.IsNotExist(statErr) {
			log.Print(fmt.Sprintf("No lazenfile exists, creating %s", lazenpath))
			core.CreateLazenFile(lazenpath)
		} else {
			log.Panic("Lazenfile already exists", lazenpath)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
