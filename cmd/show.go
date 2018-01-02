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
	"bytes"
	"fmt"
	"github.com/TwistedSystemsLtd/lazenby-go/core"
	"github.com/TwistedSystemsLtd/lazenby-go/lazendata"
	"github.com/spf13/cobra"
	"github.com/ugorji/go/codec"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("show called")
		var lazenfilename = cmd.Flag("file").Value.String()
		lazenpath := core.GetLazenpath(lazenfilename)

		lazenfile := core.ReadLazenFile(lazenpath)

		publicKey, privateKey := core.ReadUserKeys(core.Lazenhome())

		lazenkey := core.DecryptLazenkey(publicKey, privateKey, lazenfile)

		parsedSecret := &lazendata.RevealedSecret{}

		for _, secret := range lazenfile.Secrets {
			secretBytes := core.DecodeString(secret)

			revealedBytes := bytes.NewBuffer(core.DecryptWithLazenkey(lazenkey, secretBytes))

			var ch codec.CborHandle
			dec := codec.NewDecoderBytes(revealedBytes.Bytes(), &ch)
			dec.Decode(&parsedSecret)

			fmt.Println(parsedSecret.Name, "=", parsedSecret.Value)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
