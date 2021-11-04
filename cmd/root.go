/*
Copyright © 2021 Loïc Yavercovski

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
	exporter "github.com/arsiesys/jenkins-node-state-exporter/pkg/exporter"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jenkins-node-state-exporter",
	Short: "Generate prometheus metrics for jenkins nodes states",
	Long: `Generate prometheus metrics for jenkins nodes states`,
	Run: func(cmd *cobra.Command, args []string) {
		exporter.Entrypoint()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringP("address","a","http://localhost/jenkins","address of the jenkins server")
	if err := viper.BindPFlag("address", rootCmd.Flags().Lookup("address")); err != nil {
		log.Fatal(err)
	}
	rootCmd.Flags().Int("port", 9827,"port to listen on")
	if err := viper.BindPFlag("port", rootCmd.Flags().Lookup("port")); err != nil {
		log.Fatal(err)
	}
	rootCmd.Flags().Duration("fetch-interval",30*time.Second,"fetch-interval in seconds")
	if err := viper.BindPFlag("fetch-interval", rootCmd.Flags().Lookup("fetch-interval")); err != nil {
		log.Fatal(err)
	}
	rootCmd.Flags().Bool("disable-authentication",false,"disable authentication")
	if err := viper.BindPFlag("disable-authentication", rootCmd.Flags().Lookup("disable-authentication")); err != nil {
		log.Fatal(err)
	}
	rootCmd.Flags().StringP("username","u","admin","username of the jenkins user account")
	if err := viper.BindPFlag("username", rootCmd.Flags().Lookup("username")); err != nil {
		log.Fatal(err)
	}
	rootCmd.Flags().StringP("password","p","admin","password of the jenkins user account")
	if err := viper.BindPFlag("password", rootCmd.Flags().Lookup("password")); err != nil {
		log.Fatal(err)
	}
	rootCmd.Flags().StringP("labelrole","r","role=","prefix of the label to parse a role associated with the node")
	if err := viper.BindPFlag("labelrole", rootCmd.Flags().Lookup("labelrole")); err != nil {
		log.Fatal(err)
	}
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
			log.Fatal(err)
		}

		// Search config in home directory with name ".jenkins-node-state-exporter" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".jenkins-node-state-exporter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	log.Printf("Using jenkins api: %s", viper.GetString("address"))
	if viper.GetBool("disable-authentication") {
		log.Printf("Authentication disabled")
	} else {
		log.Printf("Using username: %s", viper.GetString("username"))
	}
	log.Printf("Listening on port: %d",viper.GetInt("port"))
	log.Printf("Parsing interval: %fs",viper.GetDuration("fetch-interval").Seconds())
}
