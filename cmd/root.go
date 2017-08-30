package cmd

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"code.cloudfoundry.org/uaa-cli/uaa"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetSavedConfig() uaa.Config {
	c := config.ReadConfig()
	c.Trace = trace
	c.ZoneSubdomain = zoneSubdomain
	return c
}

var cfgFile string
var trace bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = cobra.Command{
	Use:   "uaa",
	Short: "A cli for interacting with UAAs",
	Long:  help.Root(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().BoolVarP(&trace, "trace", "t", false, "See additional info on HTTP requests")
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

		// Search config in home directory with name ".uaa" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".uaa")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
