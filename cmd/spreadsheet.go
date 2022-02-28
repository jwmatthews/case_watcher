package cmd

import (
	"github.com/jwmatthews/case_watcher/pkg/api"
	"github.com/jwmatthews/case_watcher/pkg/spreadsheet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// spreadsheetCmd represents the spreadsheet command
var spreadsheetCmd = &cobra.Command{
	Use:   "spreadsheet",
	Short: "Will update a google spreadsheet",
	Long:  `Updates a google spreadsheet with cached data`,
	Run: func(cmd *cobra.Command, args []string) {
		VerifyParamsOrDie()
		// Parse configuration options
		var spreadsheetId = viper.GetString("spreadsheet")
		var email = viper.GetString("client_email")
		var privkey = viper.GetString("private_key")
		var privkeyId = viper.GetString("private_key_id")

		// At the moment this command is to help test the spreadsheet integration
		// We create an empty report for now.
		// In future I expect this will read cached data and update the spreadsheet from it
		data := api.CaseReport{}
		err := spreadsheet.Update(spreadsheetId, email, privkey, privkeyId, &data)
		if err != nil {
			log.Fatalf("Error:  Unable to update spreadsheet, error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(spreadsheetCmd)
}
