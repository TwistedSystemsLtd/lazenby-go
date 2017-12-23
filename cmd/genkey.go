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
	"github.com/TwistedSystemsLtd/lazenby-go/core"
	"os"
	"os/user"
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

		fmt.Println("PUBLIC KEY", core.ToHexString(publicKey[:]))
		fmt.Println("PRIVATE KEY", core.ToHexString(privateKey[:]))

		user, err := user.Current()
		if err != nil {
			log.Panic(err)
		}

		home := user.HomeDir
		lazenhome := path.Join(home, ".lzb")

		if _, err := os.Stat(lazenhome); os.IsNotExist(err) {
			log.Print(fmt.Sprintf("No lazenhome exists, creating %s", lazenhome))
			mkdirErr := os.Mkdir(lazenhome, os.ModeDir | 0700)
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
			writeKeyFile(publicKeyFile, core.ToHexString(publicKey[:]))
			writeKeyFile(privateKeyFile, core.ToHexString(privateKey[:]))

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
