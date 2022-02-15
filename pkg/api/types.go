package api

import "time"

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

type Case struct {
	Id                   string    `json:"id"`
	Uri                  string    `json:"uri"`
	CreatedByName        string    `json:"case_createdByName"`
	ContactName          string    `json:"case_contactName"`
	Version              string    `json:"case_version"`
	Product              []string  `json:"case_product"`
	Number               string    `json:"case_number"`
	LastPublicUpdateBy   string    `json:"case_last_public_update_by"`
	Severity             string    `json:"case_severity"`
	Owner                string    `json:"case_owner"`
	LastPublicUpdateDate time.Time `json:"case_last_public_update_date"`
	CreatedDate          time.Time `json:"case_createdDate"`
	Summary              string    `json:"case_summary"`
	LastModifiedDate     time.Time `json:"case_lastModifiedDate"`
	AccountNumber        string    `json:"case_accountNumber"`
	Type                 string    `json:"case_type"`
	LastModifiedByName   string    `json:"case_lastModifiedByName"`
	CustomerEscalation   bool      `json:"case_customer_escalation"`
	Status               string    `json:"case_status"`
}
