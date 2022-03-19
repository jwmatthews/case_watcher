package cmd

import (
	"github.com/jwmatthews/case_watcher/pkg/cache"
	"github.com/jwmatthews/case_watcher/pkg/email"
	"github.com/jwmatthews/case_watcher/pkg/report"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var emailCmd = &cobra.Command{
	Use:   "email",
	Short: "Will email a summary report of relevant cases",
	Long:  `Will look at cached data and email a list of relevant cases.`,
	Run: func(cmd *cobra.Command, args []string) {
		VerifyParamsOrDie()
		// Parse configuration options
		var spreadsheetId = viper.GetString("spreadsheet")
		var sesRegion = viper.GetString("ses_region")
		var sesSender = viper.GetString("ses_sender")
		var reportEmailRecipients = viper.GetStringSlice("report_email_recipients")

		c, err := cache.Init(DBName)
		if err != nil {
			log.Fatalf("Error:  Unable to initialize cache: %s", err)
		}
		report := report.GetReport(&c, spreadsheetId)
		err = email.Send(report, sesSender, sesRegion, reportEmailRecipients)
		if err != nil {
			log.Fatalf("Error:  Unable to send report via email: %s", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(emailCmd)
}
