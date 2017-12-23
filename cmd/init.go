// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"path"
	"log"
	"github.com/TwistedSystemsLtd/lazenby-go/core"
	"github.com/TwistedSystemsLtd/lazenby-go/lazendata"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"strings"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new lazenfile",
	Long:  `Create a new lazenfile if it doesn't exist already.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		var lazenfile = cmd.Flag("file").Value.String()

		var lazenpath string
		if dir, err := os.Getwd(); err == nil {
			if path.IsAbs(lazenfile) {
				lazenpath = lazenfile
			} else {
				lazenpath = path.Join(dir, lazenfile)
			}
			if _, err := os.Stat(lazenpath); os.IsNotExist(err) {
				log.Print(fmt.Sprintf("No lazenfile exists, creating %s", lazenpath))
				createLazenfile(lazenpath)
			} else {
				log.Panic("Lazenfile already exists", lazenpath)
			}
		} else {
			log.Panic("Could not get current working directory")
		}
	},
}

func createLazenfile(lazenpath string) {
	lazenkey := core.GenerateLazenkey()
	publicKey, privateKey := core.GenerateUserKeys()
	lazenkeys := make(map[string]*lazendata.Keypair)
	encryptedLazenKey := core.EncryptWithUserKey(publicKey, privateKey, lazenkey[:])

	keypair := &lazendata.Keypair{PublicKey: publicKey[:], Lazenkey: *encryptedLazenKey}
	lazenkeys[core.ToHexString(publicKey[:])] = keypair

	lazenfile := &lazendata.Lazenfile{Lazenkeys: lazenkeys, Secrets: nil}
	lazenbytes, err := proto.Marshal(lazenfile)
	if err != nil {
		log.Panic("Error marshalling lazenfile", err)
	}

	hexString := core.ToHexString(lazenbytes)
	chunks := core.Chunk(hexString)
	body := []byte(strings.Join(chunks, "\n"))

	writeErr := ioutil.WriteFile(lazenpath, body, 0777)
	if writeErr != nil {
		log.Panic("Could not write lazenfile", writeErr)
	}
}

func toBlocks(slices []string) {
	for _, str := range slices {
		fmt.Printf("(%s) len: %d\n", str, len(str))
	}
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
