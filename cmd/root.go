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
	"log"
	"os"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

var IsInteractive bool

var rootCmd = &cobra.Command{
	Use:   "spotify_controller",
	Short: "Control your Spotify via CLI",
	Long:  `Control your Spotify via CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		IsInteractive = true
		commands := cmd.Commands()
		for true {
			idx, err := fuzzyfinder.Find(
				commands,
				func(i int) string {
					return commands[i].Use
				},
				fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
					if i == -1 {
						return ""
					}
					return fmt.Sprintf(
						"Args: %s\nDescription: %s\nExample: %s",
						commands[i].ValidArgs,
						commands[i].Short,
						commands[i].Example,
					)
				}))
			if err != nil {
				log.Fatal(err)
			}
			commands[idx].Run(cmd, []string{})

		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())

}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spotify_controller.yaml)")

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
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".spotify_controller" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".spotify_controller")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func RunBare(cmd *cobra.Command, args []string) {
	fmt.Println(rootCmd.Commands())
	for _, c := range rootCmd.Commands() {
		fmt.Println(c.Use)
	}
}
