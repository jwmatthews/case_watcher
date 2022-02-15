package search

import (
	"context"
	"log"

	"github.com/jwmatthews/case_watcher/pkg/api"
)

// 	Search will run a keyword search to find relevant cases
//		url: Base URL of the remote server
//		username: Username to use for Basic Auth
//		password: Password to use for Basic Auth
//		searchQuery: A query string to pass in to a SOLR query
//		expression:  Related to how the results should be formatted and pagination
//
// 		Note:  I've seen bad requests returned due to lack of sufficient information in the 'expression'
//
func Search(url, username, password, searchQuery, expression string) (*api.ResponseCasesQueryBody, error) {
	log.Printf("Args are: url=`%s`, username=`%s`, password=`%s`, query=`%s`, expression=`%s`\n", url, username, "**REDACTED**", searchQuery, expression)
	client := api.NewClient(url, username, password)
	ctx := context.Background()

	resp, err := client.GetCases(ctx, searchQuery, expression)
	if err != nil {
		log.Fatalf("Error from client.GetCases():  err = '%v'", err)
	}
	log.Printf("client.GetCases() returnedm: \n%v\n", resp)
	log.Printf("Response:  Start = '%d', NumFound = '%d', len(Cases) = '%d'", resp.Start, resp.NumFound, len(resp.Cases))

	for index, theCase := range resp.Cases {
		log.Printf("Case [%d]: \n\t%v\n", index, theCase)
	}

	return resp, nil
}
