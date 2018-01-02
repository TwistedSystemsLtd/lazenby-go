// Lazenby - Your secrets as a service
// Copyright © 2018 Twisted Systems Ltd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.


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
