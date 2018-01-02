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
	"github.com/TwistedSystemsLtd/lazenby-go/lazendata"
	"github.com/spf13/cobra"
	"github.com/ugorji/go/codec"
	"golang.org/x/crypto/ripemd160"
	"log"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")

		key := args[0]
		value := args[1]

		lazenpath := cmd.Flag("file").Value.String()
		lazenfile := openLazenfile(lazenpath)

		publicKey, privateKey := core.ReadUserKeys(core.Lazenhome())
		lazenkey := core.DecryptLazenkey(publicKey, privateKey, lazenfile)

		revealedSecret := &lazendata.RevealedSecret{Name: key, Value: value}

		var revealedBytes []byte
		var ch codec.CborHandle
		enc := codec.NewEncoderBytes(&revealedBytes, &ch)

		err := enc.Encode(revealedSecret)

		if err != nil {
			log.Panic("Error marshalling secret", err)
		}

		encryptedSecret := core.EncryptWithLazenkey(lazenkey, revealedBytes)

		secret := core.EncodeString(encryptedSecret)

		hasher := ripemd160.New()
		hasher.Write([]byte(key))

		lazenfile.Secrets[core.EncodeString(hasher.Sum(nil))] = secret

		core.SaveLazenFile(lazenpath, lazenfile)
	},
}

func openLazenfile(lazenfile string) *lazendata.Lazenfile {
	lazenpath := core.GetLazenpath(lazenfile)
	return core.ReadLazenFile(lazenpath)
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
