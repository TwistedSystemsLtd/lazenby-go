// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lzb",
	Short: "Your secrets as a service",
	Long: `Lazenby is a service that keeps your secrets. Most projects have a bunch of shared private data
like database passwords, AWS/GCP tokens etc. Lazenby securely stores these in a file (called the 'lazenfile' by default) that is safe to commit
to your repo, including public repos. Only contributors you add can decrypt it. Lazenby can also emit the variables in ENV_VAR format suitable for
exporting to a shell.

On first use, lazenby will generate 3 files.
'$HOME/.lzb/privatekey': This is your private key (like an ssh key). Keep it safe!
'$HOME/.lzb/publickey': This is your public key. This is what you provide to have yourself added as a user of the lazenfile
'./lazenfile': This is the secrets store.

The lazenfile is a base64 encoded protocol buffer file containing the following:
1. A map of {publickey: encrypted_lazenkey}
2. A map of secrets in the form {name: (encrypted_value, [tags])}

When the lazenfile is first created it contains just a single entry consisting of the creator's public key and a randomly
generated lazenkey encrypted with the creator's public key. Only the matching private key can be used to decrypt the lazenkey.

Each secret is encrypted with the lazenkey using symmetric encryption.

Someone whose public key is already in the lazenfile can add new users, using 'lzb addkey'. This command will read their lazenkeys,
decrypt the lazenkey using their private key, encrypt the lazenkey with the supplied public key and update the lazenfile with the new
publickey: encrypted_lazenkey pair.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() { 
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lzb/config.yaml)")
	rootCmd.PersistentFlags().StringP("file", "f", "lazenfile", "Override the default lazenfile")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".lzb" (without extension).
		viper.AddConfigPath(home + "/.lzb")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
