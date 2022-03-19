package report

import (
	"fmt"
	"github.com/jwmatthews/case_watcher/pkg/cache"
	"time"
)

type Report struct {
	Cache         *cache.Cache
	SpreadsheetID string
}

func (r Report) GetSpreadsheetURL() string {
	return fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", r.SpreadsheetID)
}

func (r Report) toString() {

}

func (r Report) GetSubjectLine() string {
	currentDate := time.Now().Format("2006-01-02")
	return fmt.Sprintf("Case Report for %s", currentDate)
}

func (r Report) ToHTML() string {
	currentDate := time.Now().Format("2006-01-02")
	openCases, err := r.GetOpenCases()
	if err != nil {
		return "<h1>Error processing report</h1>"
	}
	closedCases, err := r.GetClosedCases()
	if err != nil {
		return "<h1>Error processing report</h1>"
	}
	activeCases, err := r.GetActiveCasesFrom(time.Now().AddDate(0, 0, -7))
	if err != nil {
		return "<h1>Error processing report</h1>"
	}
	html := fmt.Sprintf("<h1>Department Case Report %s</h1>"+
		"<p>This email was sent with "+
		"<a href='https://github.com/jwmatthews/case_watcher'>Case Watcher</a></p>"+
		"<p>%d Open Cases</p>"+
		"<p>%d Active Cases updated in past week</p>"+
		"<p>%d Closed Cases</p>"+
		"<p>For more details visit the <a href='%s'>spreadsheet here</a></p>",
		currentDate, len(openCases), len(activeCases), len(closedCases), r.GetSpreadsheetURL())
	return html
}

func (r Report) GetOpenCases() ([]cache.Case, error) {
	return r.Cache.GetOpenCases()
}

func (r Report) GetClosedCases() ([]cache.Case, error) {
	return r.Cache.GetClosedCases()
}

func (r Report) GetActiveCasesFrom(since time.Time) ([]cache.Case, error) {
	return r.Cache.GetCasesActiveFrom(since)
}

func GetReport(myCache *cache.Cache, spreadsheetId string) Report {
	return Report{Cache: myCache, SpreadsheetID: spreadsheetId}
}
