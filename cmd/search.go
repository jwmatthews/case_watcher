package cmd

import (
	"github.com/jwmatthews/case_watcher/pkg/cache"
	"github.com/jwmatthews/case_watcher/pkg/search"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
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
		VerifyParamsOrDie()
		// Parse configuration options
		var url = viper.GetString("url")
		var username = viper.GetString("username")
		var password = viper.GetString("password")
		var searchQuery = viper.GetString("query")
		var expression = viper.GetString("expression")
		//var spreadsheetId = viper.GetString("spreadsheet")
		//var email = viper.GetString("client_email")
		//var privkey = viper.GetString("private_key")
		//var privkeyId = viper.GetString("private_key_id")

		data, err := search.Search(url, username, password, searchQuery, expression)
		if err != nil {
			log.Fatalf("Error:  Failed to search for cases, error: %v\n", err)
		}
		c, err := cache.Init(DBName)
		if err != nil {
			log.Fatalf("Error:  Unable to initialize cache: %s", err)
		}
		err = c.StoreCases(data.Cases)
		if err != nil {
			log.Fatalf("Error updating cases in cache: %s\n", err)
		}
		accountIDs := c.GetMissingAccountIDs()
		if err != nil {
			log.Fatalf("Error processing updates to cache: %s\n", err)
		}
		log.Printf("Found missing account IDs: %s", accountIDs)
		/*
			cr := data.ToCaseReport()
			err = spreadsheet.Update(spreadsheetId, email, privkey, privkeyId, &cr)
			if err != nil {
				log.Fatalf("Error:  Unable to update spreadsheet, error: %v\n", err)
			}
		*/
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
