/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/niccoloCastelli/defiarb/config"
	"github.com/spf13/cobra"
	"os"
)

var (
	cfgFile    string
	walletAddr string
	conf       *config.Config
	Version    = "0.1.0"
	Revision   = ""
	Branch     = ""
)

func SetVersion(version string, revision string, branch string) {
	Version = version
	Revision = revision
	Branch = branch
	rootCmd.Version = fmt.Sprintf("%s (rev. %s) (branch: %s)", Version, Revision, Branch)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bscbot",
	Short: "Binance smart chain bot",
	Long:  ``,
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bscbot.yaml)")
	rootCmd.PersistentFlags().StringVarP(&walletAddr, "address", "a", "", "Wallet address")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	if cfgFile != "" {
		// Use config file from the flag.
		conf, err = config.ReadConfig(cfgFile)
	} else {
		conf, err = config.ReadConfig("config.json")
	}
	if err != nil {
		return
	}
}
