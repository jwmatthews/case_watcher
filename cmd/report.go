package cmd

import (
	"fmt"
	"github.com/jwmatthews/case_watcher/pkg/cache"
	"github.com/jwmatthews/case_watcher/pkg/report"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Will display a summary report of cached data to stdout",
	Long:  `Intended to help debug reports by looking at cached data and displaying summary data to stdout`,
	Run: func(cmd *cobra.Command, args []string) {
		VerifyParamsOrDie()
		// Parse configuration options
		var spreadsheetId = viper.GetString("spreadsheet")

		c, err := cache.Init(DBName)
		if err != nil {
			log.Fatalf("Error:  Unable to initialize cache: %s", err)
		}
		report := report.GetReport(&c, spreadsheetId)
		sinceLastWeek := time.Now().AddDate(0, 0, -7)
		activeCases, err := report.GetActiveCasesFrom(sinceLastWeek)
		if err != nil {
			fmt.Printf("Error from GetActiveCasesSince(): %s\n", err)
			os.Exit(1)
		}
		openCases, err := report.GetOpenCases()
		if err != nil {
			fmt.Printf("Error from GetOpenCases(): %s\n", err)
			os.Exit(1)
		}
		closedCases, err := report.GetClosedCases()
		if err != nil {
			fmt.Printf("Error from GetClosedCases(): %s\n", err)
			os.Exit(1)
		}
		uniqStatusValues, err := report.Cache.GetUniqueCaseStatusValues()
		if err != nil {
			fmt.Printf("Error from GetUniqueCaseStatusValues(): %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Spreadsheet URL: %s\n", report.GetSpreadsheetURL())
		fmt.Printf("Subject line: %s\n", report.GetSubjectLine())
		fmt.Printf("%d:  actives cases since last week (%s)\n", len(activeCases), sinceLastWeek.String())
		fmt.Printf("%d:  open cases\n", len(openCases))
		fmt.Printf("%d:  closed cases\n", len(closedCases))
		fmt.Printf("\n\n")
		fmt.Println("HTML Report")
		fmt.Println(report.ToHTML())
		fmt.Printf("\n\nDebug:\n")
		fmt.Printf("\t %d: Unique status values: %q\n", len(uniqStatusValues), uniqStatusValues)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
