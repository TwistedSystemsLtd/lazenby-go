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
	"github.com/TwistedSystemsLtd/lazenby-go/core"
	"github.com/spf13/cobra"
)

// adduserCmd represents the adduser command
var adduserCmd = &cobra.Command{
	Use:   "adduser",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("adduser called")
		newUserPublicKey := core.DecodeString(args[0])

		lazenpath := cmd.Flag("file").Value.String()
		lazenfile := openLazenfile(lazenpath)

		authorisedPublicKey, authorisedPrivateKey := core.ReadUserKeys(core.Lazenhome())

		lazenkey := core.DecryptLazenkey(authorisedPublicKey, authorisedPrivateKey, lazenfile)

		var newPublicKeyBytes [32]byte
		copy(newPublicKeyBytes[:], newUserPublicKey)
		newUserLazenkey := core.EncryptWithUserKey(&newPublicKeyBytes, authorisedPrivateKey, authorisedPublicKey, lazenkey[:])

		lazenfile.Lazenkeys[core.EncodeString(newUserPublicKey)] = core.EncodeString(newUserLazenkey)
		core.SaveLazenFile(lazenpath, lazenfile)
	},
}

func init() {
	rootCmd.AddCommand(adduserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adduserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adduserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
