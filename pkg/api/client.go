package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type Client struct {
	BaseURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

func NewClient(url string, username string, password string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		BaseURL:  url,
		Username: username,
		Password: password,
		HTTPClient: &http.Client{
			Timeout:   time.Minute,
			Transport: tr,
		},
	}
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	log.Printf("sendRequest() invoked\n")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.SetBasicAuth(c.Username, c.Password)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Warning: unable to read response body, received error: %v\n", err)
		return err
	}
	log.Printf("Received a status code of %v\n", res.StatusCode)
	log.Printf("Received a response body of:\n%s\n", body)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	var objmap map[string]json.RawMessage
	err = json.Unmarshal(body, &objmap)
	if err != nil {
		log.Printf("Warning: unable to unmarshall string received error: %v\n", err)
		log.Printf("Raw response is:\n#{body}\n")
		return err
	}
	log.Printf("ObjMap is:\n#{objmap}\n")
	if _, ok := objmap["response"]; !ok {
		log.Printf("Warning: Didn't receive a 'response' key as expected in response below:\n%v\n", objmap)
		return errors.New("Missing 'response' key")
	}
	err = json.Unmarshal(objmap["response"], &v)
	if err != nil {
		log.Printf("Warning: unable to unmarshall 'response', received error: %v\n", err)
		return err
	}
	//json.NewDecoder(res.Body).Decode(&)

	//var users []Users
	//err := json.Unmarshal(body, &users)
	//if err != nil {
	//	log.Fatal(err)
	//}
	/*fullResponse := successResponse{
		Data: v,
	}
	if err = json.NewDecoder(res.Body).Decode(&fullResponse); err != nil {
		return err
	}
	*/
	return nil
}

func (c *Client) GetCases(ctx context.Context, searchQuery string, expression string) (*ResponseCasesQueryBody, error) {
	log.Printf("GetCases(ctx, '%v') invoked\n", searchQuery)
	var url = fmt.Sprintf("%s/search/v2/cases", c.BaseURL)
	var q = CasesQuery{
		Q:             searchQuery,
		Start:         0,
		Rows:          100,
		PartnerSearch: false,
		Expression:    expression,
	}
	var jsonData, err = json.Marshal(q)
	if err != nil {
		log.Fatalf("Failed to parse json of\n%v\n", q)
		return nil, err
	}
	log.Printf("Body of searchQuery is:\n%s\n", jsonData)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to form http.NewRequest. err = %v\n", err)
		return nil, err
	}
	req = req.WithContext(ctx)

	res := ResponseCasesQueryBody{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
