package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var url string
var username string
var password string
var query string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "case_watcher",
	Short: "Intended to help get a better pulse of what is happening for cases our team may care about.",
	Long:  `Will perform a keyword search to find relevant cases and generate various reports`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is to look for '.case_watcher.yaml' in $HOME or present directory")
	rootCmd.PersistentFlags().StringVar(&url, "url", "", "URL for fetching case info")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "username for fetching case info")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "password for fetching case info")
	rootCmd.PersistentFlags().StringVar(&query, "query", "", "query for fetching case info")

	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("query", rootCmd.PersistentFlags().Lookup("query"))
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

		// Search config in home directory with name ".case_watcher" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".case_watcher")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Fprintln(os.Stderr, "Quitting early because no configuration file found.")
		os.Exit(1)
	}
}
