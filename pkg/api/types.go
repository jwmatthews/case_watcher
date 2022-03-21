package api

import (
	"fmt"
	"time"
)

type CasesQuery struct {
	Q             string `json:"q"`
	Start         int    `json:"start"`
	Rows          int    `json:"rows"`
	PartnerSearch bool   `json:"partnerSearch"`
	Expression    string `json:"expression"`
}

type ResponseCasesQuery struct {
	Response ResponseCasesQueryBody `json:"response"`
}

type ResponseCasesQueryBody struct {
	NumFound int    `json:"numFound"`
	Start    int    `json:"start"`
	Cases    []Case `json:"docs"`
}

// ToCellValues returns an array we can use to update a spreadsheet
func (r ResponseCasesQueryBody) ToCellValues() [][]interface{} {
	values := make([][]interface{}, len(r.Cases)+1)
	values[0] = CaseCellHeaderValues()
	startingIndex := 1
	for index, c := range r.Cases {
		values[startingIndex+index] = c.ToCellValues()
	}
	return values
}

func (r ResponseCasesQueryBody) ToCaseReport() CaseReport {
	cr := CaseReport{}
	for _, c := range r.Cases {
		if c.Status == "Closed" {
			cr.ClosedCases = append(cr.ClosedCases, c)
		} else {
			cr.OpenCases = append(cr.OpenCases, c)
		}
	}
	return cr
}

type CaseReport struct {
	OpenCases   []Case
	ClosedCases []Case
}

func (cr CaseReport) ToCellValues() ([][]interface{}, [][]interface{}) {
	openCaseValues := make([][]interface{}, len(cr.OpenCases)+1)
	openCaseValues[0] = CaseCellHeaderValues()
	startingIndex := 1
	for index, c := range cr.OpenCases {
		openCaseValues[startingIndex+index] = c.ToCellValues()
	}

	closedCaseValues := make([][]interface{}, len(cr.ClosedCases)+1)
	closedCaseValues[0] = CaseCellHeaderValues()
	startingIndex = 1
	for index, c := range cr.ClosedCases {
		closedCaseValues[startingIndex+index] = c.ToCellValues()
	}
	return openCaseValues, closedCaseValues
}

type Case struct {
	AccountNumber        string    `json:"case_accountNumber"`
	CaseNumber           string    `json:"case_number"`
	ContactName          string    `json:"case_contactName"`
	CreatedByName        string    `json:"case_createdByName"`
	CreatedDate          time.Time `json:"case_createdDate"`
	CustomerEscalation   bool      `json:"case_customer_escalation"`
	Id                   string    `json:"id"`
	LastModifiedByName   string    `json:"case_lastModifiedByName"`
	LastModifiedDate     time.Time `json:"case_lastModifiedDate"`
	LastPublicUpdateBy   string    `json:"case_last_public_update_by"`
	LastPublicUpdateDate time.Time `json:"case_last_public_update_date"`
	Owner                string    `json:"case_owner"`
	Products             []string  `json:"case_product"`
	Severity             string    `json:"case_severity"`
	Summary              string    `json:"case_summary"`
	Status               string    `json:"case_status"`
	Type                 string    `json:"case_type"`
	Uri                  string    `json:"uri"`
	Version              string    `json:"case_version"`
}

func (c Case) String() string {
	return fmt.Sprintf(
		"Case[%s] %s, Created %s\n"+
			"\tSeverity: %s, Status: %s\n"+
			"\tLastPublicUpdateDate: %s\n"+
			"\tCreatedByName; %s, ContactName: %s\n"+
			"\tURI: %s\n"+
			"\tId: %s\n",
		c.CaseNumber,
		c.Summary,
		c.CreatedDate,
		c.Severity,
		c.Status,
		c.LastPublicUpdateDate,
		c.CreatedByName,
		c.ContactName,
		c.Uri,
		c.Id)
}

func CaseCellHeaderValues() []interface{} {
	return []interface{}{"Uri",
		"Severity",
		"Id",
		"Status",
		"Summary",
		"CreatedByName",
		"CreatedDate",
		"LastModifiedData",
	}
}
func (c Case) ToCellValues() []interface{} {
	values := make([]interface{}, 0)
	values = append(values,
		fmt.Sprintf("=hyperlink(\"%s\")", c.Uri),
		c.Severity,
		c.Id,
		c.Status,
		c.Summary,
		c.CreatedByName,
		c.CreatedDate.String(),
		c.LastModifiedDate.String(),
	)
	return values
}

type Account struct {
	AccountNumber  string `json:"accountNumber"`
	GSCSMSegment   string `json:"gscsmSegment"`
	Name           string `json:"name"`
	CSMUserID      string `json:"csmUserId"`
	CSMUserName    string `json:"csmUserName"`
	CSMUserSSOName string `json:"csmUserSSOName"`
	Strategic      bool   `json:"strategic"`
	HasEnhancedSLA bool   `json:"hasEnhancedSLA"`
	HasSRM         bool   `json:"hasSRM"`
	HasTAM         bool   `json:"hasTAM"`
}
