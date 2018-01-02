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
	"github.com/TwistedSystemsLtd/lazenby-go/core"
	"os"
	"log"
	"path"
	"io/ioutil"
)

// genkeyCmd represents the genkey command
var genkeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("genkey called")
		publicKey, privateKey := core.GenerateUserKeys()

		fmt.Println("PUBLIC KEY", core.EncodeString(publicKey[:]))
		fmt.Println("PRIVATE KEY", core.EncodeString(privateKey[:]))


		lazenhome := core.Lazenhome()

		if _, err := os.Stat(lazenhome); os.IsNotExist(err) {
			log.Print(fmt.Sprintf("No lazenhome exists, creating %s", lazenhome))
			mkdirErr := os.Mkdir(lazenhome, os.ModeDir|0700)
			if mkdirErr != nil {
				log.Panic("Error creating lazenhome", mkdirErr)
			}
		} else {
			log.Print("lazenhome exists", lazenhome)
		}

		publicKeyFile := path.Join(lazenhome, "publickey")
		privateKeyFile := path.Join(lazenhome, "privatekey")

		_, pubKeyErr := os.Stat(publicKeyFile)
		_, privKeyErr := os.Stat(privateKeyFile)

		if os.IsNotExist(pubKeyErr) && os.IsNotExist(privKeyErr) {
			log.Print(fmt.Sprintf("Creating keyfiles: %s / %s", publicKeyFile, privateKeyFile))
			writeKeyFile(publicKeyFile, core.EncodeString(publicKey[:]))
			writeKeyFile(privateKeyFile, core.EncodeString(privateKey[:]))

		} else {
			log.Panic("Keyfile(s) already present, aborting", pubKeyErr, privKeyErr)
		}
	},
}

func writeKeyFile(keyPath string, keyData string) {
	keyErr := ioutil.WriteFile(keyPath, []byte(keyData), 0600)
	if keyErr != nil {
		log.Panic("Error writing key data", keyPath, keyErr)
	}
}

func init() {
	rootCmd.AddCommand(genkeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genkeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genkeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
