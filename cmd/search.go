package cmd

import (
	"github.com/jwmatthews/case_watcher/pkg/search"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Will perform a keyword search for relevant cases",
	Long: `Will perform a keyword search for relevant cases. 
	It is assumed that we may have some false positives in the list
	of cases, i.e. cases that are not directly related to our team 
	but matched from the keyword search.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse configuration options
		var url = viper.GetString("url")
		var username = viper.GetString("username")
		var password = viper.GetString("password")
		var searchQuery = viper.GetString("query")
		var expression = viper.GetString("expression")
		search.Search(url, username, password, searchQuery, expression)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
